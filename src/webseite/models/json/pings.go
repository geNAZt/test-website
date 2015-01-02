package json

import (
	"github.com/astaxie/beego/orm"
	"strconv"
	"time"
	"webseite/models"
)

type JSONPingResponse struct {
	Id      int32
	Players map[string]int32
}

func (j *JSONPingResponse) FillPings(days int32) {
	// Get the database
	o := orm.NewOrm()
	o.Using("default")

	// Load the 24 hour before ping
	// Build up the Query
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("*").
		From("ping").
		Where("server_id = " + strconv.FormatInt(int64(j.Id), 10)).
		OrderBy("time").
		Desc().
		Limit(int(days * 24 * 60))

	// Get the SQL Statement and execute it
	sql := qb.String()
	pings := []models.Ping{}
	o.Raw(sql).QueryRows(&pings)

	// Construct pasttime and the map
	j.Players = make(map[string]int32)
	pastTime := time.Now().Add(time.Duration(-days*24*60) * time.Minute)

	// Select the pings we need to fill in
	for pingI := range pings {
		sqlPing := pings[(len(pings)-1)-pingI]
		if sqlPing.Time.Before(pastTime) {
			continue
		}

		j.Players[strconv.FormatInt(sqlPing.Time.Unix(), 10)] = sqlPing.Online
	}
}
