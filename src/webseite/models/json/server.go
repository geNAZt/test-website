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
)

const createdFormat = "2006-01-02 15:04:05"

type Server struct {
	Id         int32
	Name       string
	Website    string
	IP         string
	Online     int32
	MaxPlayers int32
	Record     int32
	Average    int32
	Favicon    string
	Ping24     int32
	Favicons   []status.Favicon `json:"-"`
}

type PlayerUpdate struct {
	Id      int32
	Online  int32
	Time    int64
	Ping    int32
	Ping24  int32
	Record  int32
	Average int32
}

type JSONUpdatePlayerResponse struct {
	Ident string
	Value PlayerUpdate
}

type JSONUpdateFaviconResponse struct {
	Ident string
	Value ServerFavicon
}

type ServerFavicon struct {
	Id   int32
	Icon string
}

type StoredFavicon struct {
	Favicon  string
	Favicons []status.Favicon
}

var Servers map[int32]Server
var Favicons *cache.TimeoutCache

func init() {
	tempCache, err := cache.NewFaviconCache()
	if err != nil {
		panic("Could not init favicon cache")
	}

	Favicons = tempCache
	Servers = make(map[int32]Server)
}

func ReloadServers(servers []models.Server) {
	// Prepare ORM (database)
	o := orm.NewOrm()
	o.Using("default")

	// Get the current time
	_, offset := time.Now().Zone()

	// Iterate over all Servers to calc 24h Pings
	for serverI := range servers {
		sqlServer := servers[serverI]

		// Check if there is a old entry
		jsonServer := Server{
			Online: 0,
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
		past24Hours := time.Unix( (time.Now().Add(time.Duration(-24*60) * time.Minute).Unix()) - int64(offset), 0 ).Format( createdFormat )
		past24HoursAnd2Minutes := time.Unix( (time.Now().Add(time.Duration((-24*60)+2) * time.Minute).Unix()) - int64(offset), 0 ).Format( createdFormat )

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
		if ent, ok := Favicons.Get(jsonServer.Name); ok {
			jsonServer.Favicons = ent.(StoredFavicon).Favicons
			jsonServer.Favicon = ent.(StoredFavicon).Favicon
		}

		// Recalc Average and record counters
		jsonServer.RecalcAverage()
		jsonServer.RecalcRecord()
		Servers[jsonServer.Id] = jsonServer
	}
}

func SendLog(c *websocket.Connection, message string) {
	defer func() {
		recover()
	}()

	jsonResponse := JSONResponse{
		Ident: "log",
		Value: message,
	}

	jsonBytes, err := gojson.Marshal(jsonResponse)
	if err != nil {
		beego.BeeLogger.Warn("Could not convert to json: %v", err)
		return
	}

	c.Send <- jsonBytes
}

func SendAllServers(c *websocket.Connection, view *models.View) {
	defer func() {
		recover()
	}()

	jsonResponse := JSONResponse{
		Ident: "servers",
		Value: []Server{},
	}

	for serverI := range view.Servers {
		server := GetServer(view.Servers[serverI].Id)
		if server.Id != -1 {
			jsonResponse.Value = append(jsonResponse.Value.([]Server), server)
		}
	}

	jsonBytes, err := gojson.Marshal(jsonResponse)
	if err != nil {
		beego.BeeLogger.Warn("Could not convert to json: %v", err)
		return
	}

	c.Send <- jsonBytes
}

func SendFavicon(c *websocket.Connection, serverId int32, favicon string) {
	defer func() {
		recover()
	}()

	fav := JSONUpdateFaviconResponse{
		Ident: "favicon",
		Value: ServerFavicon{
			Icon: favicon,
			Id:   serverId,
		},
	}

	jsonBytes, err := gojson.Marshal(fav)
	if err != nil {
		beego.BeeLogger.Warn("Could not convert to json: %v", err)
		return
	}

	c.Send <- jsonBytes
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

		storedFavicon := StoredFavicon{
			Favicon:  server.Favicon,
			Favicons: server.Favicons,
		}

		Favicons.Add(server.Name, storedFavicon)
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

	websocket.Hub.Broadcast <- jsonBytes
}

func (s *Server) RecalcRecord() {
	// ORM
	o := orm.NewOrm()
	o.Using("default")

	// Build up the Query
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("*").
		From("ping").
		Where("server_id = " + strconv.FormatInt(int64(s.Id), 10)).
		OrderBy("online").
		Desc().
		Limit(1)

	// Get the SQL Statement and execute it
	sql := qb.String()
	pings := []models.Ping{}
	o.Raw(sql).QueryRows(&pings)

	// Set the record
	if len(pings) > 0 {
		s.Record = pings[0].Online
	}
}

func (s *Server) RecalcAverage() {
	// ORM
	o := orm.NewOrm()
	o.Using("default")

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
