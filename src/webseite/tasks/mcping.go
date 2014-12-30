package tasks

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/toolbox"
	status "github.com/geNAZt/minecraft-status"
	statusdata "github.com/geNAZt/minecraft-status/data"
	"strconv"
	"time"
	"webseite/models"
	"webseite/models/json"
)

var servers []models.Server

func InitTasks() {
	// ORM
	o := orm.NewOrm()
	o.Using("default")

	// Build up the Query
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("*").
		From("server")

	// Get the SQL Statement and execute it
	sql := qb.String()
	servers = []models.Server{}
	o.Raw(sql).QueryRows(&servers)

	// Load
	for serverI := range servers {
		o.LoadRelated(&servers[serverI], "Pings", 0, 2*24*60, 0, "Time")
	}

	// Reload the JSON side
	json.ReloadServers(servers)

	mcping := toolbox.NewTask("mcping", "0 * * * * *", func() error {
		// Reload servers
		servers = []models.Server{}
		o.Raw(sql).QueryRows(&servers)

		// Load
		for serverI := range servers {
			o.LoadRelated(&servers[serverI], "Pings", 0, 2*24*60, 0, "Time")
		}

		// Reload the JSON side
		json.ReloadServers(servers)

		// Ping all da servers
		for serverId := range servers {
			go ping(&servers[serverId])
		}

		return nil
	})

	toolbox.AddTask("mcping", mcping)

	// Start the tasks
	toolbox.StartTask()
}

func ping(server *models.Server) {
	// Get the database
	o := orm.NewOrm()
	o.Using("default")

	// Make ping
	status, err := status.GetStatus(server.Ip)
	if err != nil {
		beego.BeeLogger.Warn("Error while pinging: %v", err)

		// Create "fake" ping
		status = &statusdata.Status{
			Players: &statusdata.MCPlayers{
				Online: 0,
				Max:    0,
			},
			Ping: time.Duration(30 * time.Nanosecond),
		}
	}

	// Save ping
	ping := &models.Ping{
		Server: server,
		Online: int32(status.Players.Online),
		Ping:   int64(status.Ping),
		Time:   time.Now(),
	}

	o.Insert(ping)

	// Update record if needed
	if int32(status.Players.Online) > server.Record {
		server.Record = int32(status.Players.Online)
		o.Update(server)
	}

	// Load the 24 hour before ping
	// Build up the Query
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("*").
		From("ping").
		Where("server_id = " + strconv.FormatInt(int64(server.Id), 10)).
		Limit(1).
		Offset(24 * 60)

	// Get the SQL Statement and execute it
	sql := qb.String()
	pings := []models.Ping{}
	o.Raw(sql).QueryRows(&pings)

	var ping24 models.Ping
	if len(pings) > 0 {
		ping24 = pings[0]
	}

	// Notify the JSON side
	json.UpdateStatus(server.Id, status, &ping24)
}
