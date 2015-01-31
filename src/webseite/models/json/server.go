package json

import (
	gojson "encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	status "github.com/geNAZt/minecraft-status/data"
	"strconv"
	"sync"
	"time"
	"webseite/cache"
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

var Servers JSONResponse
var lock sync.RWMutex
var Favicons *cache.TimeoutCache

func init() {
	tempCache, err := cache.NewFaviconCache()
	if err != nil {
		panic("Could not init favicon cache")
	}

	Favicons = tempCache
}

func ReloadServers(servers []models.Server) {
	lock.Lock()
	defer lock.Unlock()

	Servers = JSONResponse{
		Ident: "servers",
		Value: []Server{},
	}

	for serverI := range servers {
		sqlServer := servers[serverI]

		pings := sqlServer.Pings

		for pingI := range pings {
			ping := pings[pingI]
			AddPing(sqlServer.Id, ping.Time.Unix(), ping.Online)
		}

		var ping24 *models.Ping
		if len(pings) > 0 {
			ping24 = pings[0]
		}

		jsonServer := Server{
			Id:      sqlServer.Id,
			IP:      sqlServer.Ip,
			Name:    sqlServer.Name,
			Website: sqlServer.Website,
			Online:  sqlServer.Pings[len(sqlServer.Pings)-1].Online,
		}

		if ping24 != nil {
			jsonServer.Ping24 = ping24.Online
		}

		if ent, ok := Favicons.Get(jsonServer.Name); ok {
			jsonServer.Favicons = ent.(StoredFavicon).Favicons
			jsonServer.Favicon = ent.(StoredFavicon).Favicon
		}

		jsonServer.RecalcAverage()
		jsonServer.RecalcRecord()
		Servers.Value = append(Servers.Value.([]Server), jsonServer)
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

func GetServer(id int32) *Server {
	for serverI := range Servers.Value.([]Server) {
		server := &Servers.Value.([]Server)[serverI]
		if server.Id == id {
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

	for serverI := range Servers.Value.([]Server) {
		server := &Servers.Value.([]Server)[serverI]

		if server.Id == id {
			server.RecalcRecord()
			server.RecalcAverage()

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

			server.Ping = int32(status.Ping)

			AddPing(server.Id, time.Now().Unix()-int64(offset), online)

			jsonPlayerUpdate := JSONUpdatePlayerResponse{
				Ident: "updatePlayer",
				Value: PlayerUpdate{
					Id:      server.Id,
					Online:  online,
					Time:    time.Now().Unix() - int64(offset),
					Ping:    server.Ping,
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
