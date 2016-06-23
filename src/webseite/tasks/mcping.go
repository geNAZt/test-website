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
	"sync"
	"webseite/cache"
	"fmt"
)

var servers []models.Server
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

	mcping := toolbox.NewTask("mcping", "0 * * * * *", func() error {
		// Ping all da servers
		serverMutex.Lock();
		for serverId := range servers {
			go ping(&servers[serverId])
		}
		serverMutex.Unlock();

		return nil
	})

	batchInserter := toolbox.NewTask("batchInserter", "0 30 * * * *", func() error {
		queueMutex.Lock()
		defer queueMutex.Unlock()

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
	o := orm.NewOrm()

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

	o.Insert(ping)

	// Notify the JSON side
	json.UpdateStatus(server.Id, status)
}
