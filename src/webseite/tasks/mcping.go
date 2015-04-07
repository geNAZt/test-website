package tasks

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/toolbox"
	status "github.com/geNAZt/minecraft-status"
	statusdata "github.com/geNAZt/minecraft-status/data"
	"time"
	"webseite/models"
	"webseite/models/json"
	"webseite/util"
	"sync"
)

var servers []models.Server
var queue *util.Queue
var queueMutex = &sync.Mutex{}

func InitTasks() {
	// ORM
	o := orm.NewOrm()
	o.Using("default")

	// Build up the Query
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("*").
		From("`server`")

	// Get the SQL Statement and execute it
	sql := qb.String()
	servers = []models.Server{}
	o.Raw(sql).QueryRows(&servers)

	// Reload the JSON side
	json.ReloadServers(servers)

	// Prepare the queue and let the pinger roll
	queue = &util.Queue{nodes: make([]*util.Node, 100)}

	mcping := toolbox.NewTask("mcping", "0 * * * * *", func() error {
		// Reload servers
		servers = []models.Server{}
		o.Raw(sql).QueryRows(&servers)

		// Reload the JSON side
		json.ReloadServers(servers)

		// Ping all da servers
		for serverId := range servers {
			go ping(&servers[serverId])
		}

		return nil
	})

	batchInserter := toolbox.NewTask("batchInserter", "0 */5 * * * *", func() error {
		queueMutex.Lock()
		defer queueMutex.Unlock()

		bulk := make([]*models.Ping, queue.Size())
		count := 0
		for {
			ele := queue.Pop().Value
			if ele == nil {
				break
			}

			bulk[count] = ele
			count++
		}

		o.InsertMulti(20, bulk)
		queue = &util.Queue{nodes: make([]*util.Node, 100)}
		return nil
	})

	toolbox.AddTask("mcping", mcping)
	toolbox.AddTask("batchInserter", batchInserter)

	// Start the tasks
	toolbox.StartTask()
}

func ping(server *models.Server) {
	// Ask the JSON side if we have a animated Favicon
	fetchFavicon := true
	fetchAnimated := server.DownloadAnimatedFavicon
	if _, ok := json.Favicons.Get(server.Name); ok {
		fetchFavicon = false
		fetchAnimated = false
	}

	// Make ping
	status, err := status.GetStatus(server.Ip, fetchAnimated)
	if err != nil {
		beego.BeeLogger.Warn("Error while pinging: %v", err)

		// Create "fake" ping
		status = &statusdata.Status{
			Players: &statusdata.MCPlayers{
				Online: 0,
				Max:    0,
			},
			Ping: time.Duration(30 * time.Second),
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
	queue.Push(&util.Node{ping})
	queueMutex.Unlock()

	// Notify the JSON side
	json.UpdateStatus(server.Id, status)
}
