package json

import (
	gojson "encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	status "github.com/geNAZt/minecraft-status/data"
	"strconv"
	"time"
	"webseite/cache"
	"webseite/models"
	"webseite/websocket"
	"fmt"
)

const createdFormat = "2006-01-02 15:04:05"

type Server struct {
	Id         int32
	Name       string
	Website    string
	IP         string
	Online     int32
	MaxPlayers int32
	Record     int64
	Average    int32
	Favicon    string
	Ping24     int32
	Uptime     float32
	UptimeLast float32
	Favicons   []status.Favicon `json:"-"`
}

type PlayerUpdate struct {
	Id      int32
	Online  int32
	Time    int64
	Ping    int32
	Ping24  int32
	Record  int64
	Average int32
}

type UptimeUpdate struct {
	Id         int32
	Uptime     float32
	UptimeLast float32
}

type JSONUpdatePlayerResponse struct {
	Ident string
	Value PlayerUpdate
}

type dbUptimeResponse struct {
	Uptime float32
}

var Servers map[int32]Server

func init() {
	Servers = make(map[int32]Server)
}

func SendAvailableViews(c *websocket.Connection) {
	// Get the default View
	// ORM
	o := orm.NewOrm()

	// Build up the Query
	qb, _ := orm.NewQueryBuilder("mysql")
	qb = qb.Select("*").
	From("`view`").
	Where("`owner_id` = ?")

	systemUserId, _ := beego.AppConfig.Int("SystemUserID")

	// Check if User is logged in, if so include his views
	var rawSeter orm.RawSeter
	if c.Session.Get("userId") != nil && c.Session.Get("userId").(int32) != -1 {
		qb.Or("`owner_id` = ?")
		rawSeter = o.Raw(qb.String(), int32(systemUserId), c.Session.Get("userId").(int32))
	} else {
		rawSeter = o.Raw(qb.String(), int32(systemUserId))
	}

	views := []models.View{}
	rawSeter.QueryRows(&views)

	// Remap for JSON
	jsonResponse := &JSONResponse{
		Ident: "views",
		Value: make(map[string]int32, len(views)),
	}

	for viewI := range views {
		view := views[viewI]
		jsonResponse.Value.(map[string]int32)[view.Name] = view.Id
	}

	// Send to client
	jsonResponse.Send(c)
}

func SendView(c *websocket.Connection) {
	// Get the default View
	// ORM
	o := orm.NewOrm()

	// Build up the Query
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("*").
	From("`view`").
	Where("`id` = ?")

	// Get the SQL Statement and execute it
	sql := qb.String()
	view := &models.View{}
	o.Raw(sql, c.Session.Get("view").(int32)).QueryRow(&view)
	o.LoadRelated(view, "Servers")
	o.LoadRelated(view, "Owner")

	// Check if user is owner of this View or if its a system view
	systemUserId, _ := beego.AppConfig.Int("SystemUserID")
	if view.Owner.Id == int32(systemUserId) || view.Owner.Id == c.Session.Get("userId").(int32) {
		// Send the user all servers which belong to this view
		jsonResponse := JSONResponse{
			Ident: "servers",
			Value: []Server{},
		}

		serverIds := make(map[int32]bool, len(view.Servers))
		for serverI := range view.Servers {
			server := GetServer(view.Servers[serverI].Id)
			if server.Id != -1 {
				serverIds[server.Id] = true
				jsonResponse.Value = append(jsonResponse.Value.([]Server), server)
			}
		}

		c.Session.Set("servers", serverIds)
		jsonResponse.Send(c)
	}
}

func ReloadServers(servers []models.Server) {
	// Prepare ORM (database)
	o := orm.NewOrm()

	// Get the current time
	_, offset := time.Now().Zone()

	// Iterate over all Servers to calc 24h Pings
	for serverI := range servers {
		fmt.Printf("Reloading server %d...\n", serverI)

		sqlServer := servers[serverI]

		// Check if there is a old entry
		jsonServer := Server{
			Online: 0,
			Uptime: 0,
		}
		if tempJsonServer, ok := Servers[sqlServer.Id]; ok {
			jsonServer = tempJsonServer
		}

		// Update basic informations
		jsonServer.Id = sqlServer.Id
		jsonServer.IP = sqlServer.Ip
		jsonServer.Name = sqlServer.Name
		jsonServer.Website = sqlServer.Website

		// Check for 24h Ping
		past24Hours := time.Unix((time.Now().Add(time.Duration(-24 * 60) * time.Minute).Unix()) - int64(offset), 0).Format(createdFormat)
		past24HoursAnd2Minutes := time.Unix((time.Now().Add(time.Duration((-24 * 60) + 2) * time.Minute).Unix()) - int64(offset), 0).Format(createdFormat)

		// Build up the Query
		qb, _ := orm.NewQueryBuilder("mysql")
		qb.Select("*").
		From("`ping`").
		Where("`server_id` = ?").
		And("`time` > ?").And("`time` < ?").
		OrderBy("`time`").
		Desc().
		Limit(1)

		// Ask the Database for 24h Ping
		sql := qb.String()
		ping := models.Ping{}

		err := o.Raw(sql, strconv.FormatInt(int64(jsonServer.Id), 10), past24Hours, past24HoursAnd2Minutes).QueryRow(&ping)
		if err == nil {
			jsonServer.Ping24 = ping.Online
		}

		// Get the Favicons for this Server entities
		if ent, ok := cache.Favicons.Get(jsonServer.Name); ok {
			jsonServer.Favicons = ent.(cache.StoredFavicon).Favicons
			jsonServer.Favicon = ent.(cache.StoredFavicon).Favicon
		}

		// Recalc Average and record counters
		jsonServer.RecalcAverage()
		jsonServer.RecalcRecord()
		jsonServer.RecalcUptime()
		Servers[jsonServer.Id] = jsonServer
	}
}

func GetServer(id int32) Server {
	if server, ok := Servers[id]; ok {
		return server
	}

	return Server{
		Id: -1,
	}
}

func UpdateStatus(id int32, status *status.Status) {
	_, offset := time.Now().Zone()

	online := int32(status.Players.Online)
	max := int32(status.Players.Max)

	server, ok := Servers[id]
	if !ok {
		return
	}

	server.Online = online

	if status.Favicon != "" {
		server.Favicon = status.Favicon
		server.Favicons = status.Favicons

		storedFavicon := cache.StoredFavicon{
			Favicon:  server.Favicon,
			Favicons: server.Favicons,
		}

		cache.Favicons.Add(server.Name, storedFavicon)
	}

	jsonPlayerUpdate := JSONUpdatePlayerResponse{
		Ident: "updatePlayer",
		Value: PlayerUpdate{
			Id:      server.Id,
			Online:  online,
			Time:    time.Now().Unix() - int64(offset),
			Average: server.Average,
			Record:  server.Record,
		},
	}

	if server.MaxPlayers != max {
		jsonMaxPlayer := &JSONMaxPlayerResponse{
			Id:         server.Id,
			MaxPlayers: max,
		}

		jsonMaxPlayer.Broadcast()
		server.MaxPlayers = max
	}

	jsonPlayerUpdate.Value.Ping24 = server.Ping24
	Servers[server.Id] = server

	jsonBytes, err := gojson.Marshal(jsonPlayerUpdate)
	if err != nil {
		beego.BeeLogger.Warn("Could not convert to json: %v", err)
		return
	}

	for c := range websocket.Hub.Connections {
		allowedServers := c.Session.Get("servers").(map[int32]bool)
		if val, ok := allowedServers[server.Id]; !ok || !val {
			continue
		}

		select {
		case c.Send <- jsonBytes:
		default:
			c.CloseCustomChannels()
			close(c.Send)
			delete(websocket.Hub.Connections, c)
		}
	}
}

func (s *Server) RecalcRecord() {
	fmt.Printf("Recalc record for server %d...\n", s.Id)

	// ORM
	o := orm.NewOrm()

	// Build up the Query
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("MAX(`online`) AS `MaxOnline`").
	From("ping").
	Where("server_id = " + strconv.FormatInt(int64(s.Id), 10));

	// Get the SQL Statement and execute it
	sql := qb.String()

	var maps []orm.Params
	_, err := o.Raw(sql).Values(&maps)
	if err != nil {
		fmt.Printf( "%v", err );
	}

	// Set the record
	if len(maps) > 0 {
		temp := maps[0]["MaxOnline"];
		if temp == nil {
			return
		}

		s.Record, err = strconv.ParseInt( maps[0]["MaxOnline"].(string), 10, 32 )
		if err != nil {
			fmt.Printf( "%v", err );
		}
	}
}

func (s *Server) RecalcAverage() {
	fmt.Printf("Recalc average for server %d...\n", s.Id)

	// ORM
	o := orm.NewOrm()

	// Build up the Query
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("*").
	From("ping").
	Where("server_id = " + strconv.FormatInt(int64(s.Id), 10)).
	OrderBy("time").
	Desc().
	Limit(60 * 24)

	// Get the SQL Statement and execute it
	sql := qb.String()
	pings := []models.Ping{}
	o.Raw(sql).QueryRows(&pings)

	// Calc the average
	overall := int32(0)
	for pingI := range pings {
		overall = overall + pings[pingI].Online
	}

	len := int32(len(pings))
	if len > 0 {
		s.Average = overall / len
	}
}

func (s *Server) RecalcUptime() {
	fmt.Printf("Recalc uptime for server %d...\n", s.Id)

	// ORM
	o := orm.NewOrm()

	// Get the current time
	_, offset := time.Now().Zone()
	pastMonth := time.Unix((time.Now().Add(time.Duration(-30 * 24) * time.Hour).Unix()) - int64(offset), 0).Format(createdFormat)

	// Get the SQL Statement and execute it
	sql := "SELECT 100-((SELECT COUNT(`id`) FROM `ping` WHERE `server_id` = ? AND `online` = 0 AND `time` > ?) / COUNT(`id`)) AS `uptime` FROM `ping` WHERE `server_id` = ? AND `time` > ?;"
	uptime := []dbUptimeResponse{}
	o.Raw(sql, s.Id, pastMonth, s.Id, pastMonth).QueryRows(&uptime)

	past2Month := time.Unix((time.Now().Add(time.Duration(-60 * 24) * time.Hour).Unix()) - int64(offset), 0).Format(createdFormat)
	sql = "SELECT 100-((SELECT COUNT(`id`) FROM `ping` WHERE `server_id` = ? AND `online` = 0 AND `time` > ? AND `time` < ?) / COUNT(`id`)) AS `uptime` FROM `ping` WHERE `server_id` = ? AND `time` > ? AND `time` < ?;"
	uptimeLast := []dbUptimeResponse{}
	o.Raw(sql, s.Id, past2Month, pastMonth, s.Id, past2Month, pastMonth).QueryRows(&uptimeLast)

	// Remap
	s.Uptime = uptime[0].Uptime
	s.UptimeLast = uptimeLast[0].Uptime

	// Send update to all
	jsonMessage := &JSONResponse{
		Ident: "uptime",
		Value: &UptimeUpdate{
			Id: s.Id,
			Uptime: s.Uptime,
			UptimeLast: s.UptimeLast,
		},
	}

	jsonMessage.BroadcastToServerID(s.Id)
}
