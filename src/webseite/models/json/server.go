package json

import (
	gojson "encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	status "github.com/geNAZt/minecraft-status/data"
	"strconv"
	"sync"
	"time"
	"webseite/models"
	"webseite/websocket"
)

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
	Ping       int32
	Ping24     int32
	Players    []Ping
	Favicons   []status.Favicon `json:"-"`
}

type Ping struct {
	Online int32
	Time   int64
}

type JSONServerResponse struct {
	Ident string
	Value []Server
}

type PlayerUpdate struct {
	Name       string
	Online     int32
	MaxPlayers int32
	Time       int64
	Ping       int32
	Ping24     int32
	Record     int32
	Average    int32
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
	Server string
	Icon   string
}

var Servers JSONServerResponse
var lock sync.RWMutex

func ReloadServers(servers []models.Server) {
	lock.Lock()
	defer lock.Unlock()

	Servers = JSONServerResponse{
		Ident: "servers",
	}

	Servers.Value = []Server{}

	for serverI := range servers {
		sqlServer := servers[serverI]

		jsonPings := []Ping{}
		var jsonPing Ping

		for pingI := range sqlServer.Pings {
			sqlPing := sqlServer.Pings[pingI]

			jsonPing = Ping{
				Online: sqlPing.Online,
				Time:   sqlPing.Time.Unix(),
			}

			jsonPings = append(jsonPings, jsonPing)
		}

		// Get the database
		o := orm.NewOrm()
		o.Using("default")

		// Load the 24 hour before ping
		// Build up the Query
		qb, _ := orm.NewQueryBuilder("mysql")
		qb.Select("*").
			From("ping").
			Where("server_id = " + strconv.FormatInt(int64(sqlServer.Id), 10)).
			Limit(1).
			Offset(24 * 60)

		// Get the SQL Statement and execute it
		sql := qb.String()
		pings := []models.Ping{}
		o.Raw(sql).QueryRows(&pings)

		var ping24 *models.Ping
		if len(pings) > 0 {
			ping24 = &pings[0]
		}

		jsonServer := Server{
			Id:      sqlServer.Id,
			IP:      sqlServer.Ip,
			Name:    sqlServer.Name,
			Website: sqlServer.Website,
			Online:  jsonPing.Online,
			Players: jsonPings,
		}

		if ping24 != nil {
			jsonServer.Ping24 = ping24.Online
		}

		jsonServer.RecalcAverage()
		jsonServer.RecalcRecord()
		Servers.Value = append(Servers.Value, jsonServer)
	}
}

func SendAllServers(c *websocket.Connection) {
	lock.RLock()
	defer lock.RUnlock()

	jsonBytes, err := gojson.Marshal(Servers)
	if err != nil {
		beego.BeeLogger.Warn("Could not convert to json: %v", err)
		return
	}

	c.Send <- jsonBytes
}

func SendFavicon(c *websocket.Connection, servername, favicon string) {
	defer func() {
		recover()
	}()

	fav := JSONUpdateFaviconResponse{
		Ident: "favicon",
		Value: ServerFavicon{
			Icon:   favicon,
			Server: servername,
		},
	}

	jsonBytes, err := gojson.Marshal(fav)
	if err != nil {
		beego.BeeLogger.Warn("Could not convert to json: %v", err)
		return
	}

	c.Send <- jsonBytes
}

func GetServer(name string) *Server {
	for serverI := range Servers.Value {
		server := &Servers.Value[serverI]
		if server.Name == name {
			return server
		}
	}

	return nil
}

func UpdateStatus(id int32, status *status.Status, ping24 *models.Ping) {
	lock.RLock()
	defer lock.RUnlock()

	_, offset := time.Now().Zone()

	online := int32(status.Players.Online)
	max := int32(status.Players.Max)

	for serverI := range Servers.Value {
		server := &Servers.Value[serverI]

		if server.Id == id {
			server.RecalcRecord()
			server.RecalcAverage()

			server.Online = online
			server.MaxPlayers = max
			if status.Favicon != "" {
				server.Favicon = status.Favicon
				server.Favicons = status.Favicons
			}
			server.Ping = int32(status.Ping)

			jsonPlayerUpdate := JSONUpdatePlayerResponse{
				Ident: "updatePlayer",
				Value: PlayerUpdate{
					Name:       server.Name,
					Online:     online,
					MaxPlayers: max,
					Time:       time.Now().Unix() - int64(offset),
					Ping:       server.Ping,
					Average:    server.Average,
					Record:     server.Record,
				},
			}

			if ping24 != nil {
				jsonPlayerUpdate.Value.Ping24 = ping24.Online
			}

			jsonBytes, err := gojson.Marshal(jsonPlayerUpdate)
			if err != nil {
				beego.BeeLogger.Warn("Could not convert to json: %v", err)
				return
			}

			websocket.Hub.Broadcast <- jsonBytes
			return
		}
	}
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
