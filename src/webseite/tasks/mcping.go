package tasks

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/toolbox"
	statusresolver "github.com/geNAZt/minecraft-status"
	statusdata "github.com/geNAZt/minecraft-status/data"
	"time"
	"webseite/models"
	"webseite/models/json"
	"webseite/util"
	"sync"
	"webseite/cache"
	"fmt"
)

var servers []models.Server
var queue *util.Queue
var queueMutex = &sync.Mutex{}
var serverMutex = &sync.Mutex{}

func InitTasks() {
	// ORM
	o := orm.NewOrm()

	// Build up the Query
	fmt.Printf("Building up server caches...\n")
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("*").
	From("`server`")

	// Get the SQL Statement and execute it
	sql := qb.String()
	servers = []models.Server{}
	o.Raw(sql).QueryRows(&servers)
	fmt.Printf("Done building up server caches...\n")

	// Reload the JSON side
	json.ReloadServers(servers)

	// Prepare the queue and let the pinger roll
	queue = &util.Queue{Nodes: make([]*models.Ping, 100)}

	mcping := toolbox.NewTask("mcping", "0 * * * * *", func() error {
		// Ping all da servers
		serverMutex.Lock();
		for serverId := range servers {
			go ping(&servers[serverId])
		}
		serverMutex.Unlock();

		return nil
	})

	batchInserter := toolbox.NewTask("batchInserter", "0/10 * * * * *", func() error {
		queueMutex.Lock()
		defer queueMutex.Unlock()

		if queue.Size() == 0 {
			return nil
		}

		// Remap into bulk array
		bulk := make([]*models.Ping, queue.Size())
		count := 0
		for {
			ele := queue.Pop()
			if ele == nil {
				break
			}

			bulk[count] = ele
			count++
		}

		fmt.Printf("Inserting %d new pings...\n", len(bulk))

		// Insert with max 20 in a Query
		o.InsertMulti(20, bulk)
		queue = &util.Queue{Nodes: make([]*models.Ping, 100)}

		// Reload servers
		serverMutex.Lock()
		servers = []models.Server{}
		o.Raw(sql).QueryRows(&servers)
		serverMutex.Unlock()

		// Reload the JSON side
		json.ReloadServers(servers)

		return nil
	})

	toolbox.AddTask("mcping", mcping)
	toolbox.AddTask("batchInserter", batchInserter)

	// Start the tasks
	toolbox.StartTask()
}

func ping(server *models.Server) {
	fmt.Printf("Pinging server %s for new data...\n", server.Name)

	// Ask the JSON side if we have a animated Favicon
	fetchFavicon := true
	fetchAnimated := server.DownloadAnimatedFavicon
	if _, ok := cache.Favicons.Get(server.Name); ok {
		fetchFavicon = false
		fetchAnimated = false
	}

	// Make ping
	status, err := statusresolver.GetStatus(server.Ip, fetchAnimated)
	if err != nil {
		beego.BeeLogger.Warn("Error while pinging: %v", err)

		// Create "fake" ping
		status = &statusdata.Status{
			Players: &statusdata.MCPlayers{
				Online: 0,
				Max:    0,
			},
		}
	}

	// "NULL" the favicon if needed
	if !fetchFavicon {
		status.Favicon = ""
	}

	// Save ping
	ping := &models.Ping{
		Server: server,
		Online: int32(status.Players.Online),
		Time:   time.Now(),
	}

	queueMutex.Lock()
	queue.Push(ping)
	queueMutex.Unlock()

	// Notify the JSON side
	json.UpdateStatus(server.Id, status)
}
