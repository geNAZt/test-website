package json

import (
	"strconv"
	"time"
	"github.com/astaxie/beego/orm"
	"webseite/models"
)

type JSONPingResponse struct {
	Id      int32
	Players map[string]int32
}

func (j *JSONPingResponse) FillPings(days int32) {
	// Construct pasttime and the map
	_, offset := time.Now().Zone()
	j.Players = make(map[string]int32)

	// ORM
	o := orm.NewOrm()
	o.Using("default")

	// Check for 24h Ping
	past24Hours := time.Unix( (time.Now().Add(time.Duration(-days*24*60) * time.Minute).Unix()) - int64(offset), 0 ).Format( createdFormat )

	// Build up the Query
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("*").
		From("`ping`").
		Where("`server_id` = ?").
		And("`time` > ?").
		OrderBy("`time`").
		Asc()

	// Ask the Database for 24h Ping
	sql := qb.String()
	pings := []models.Ping{}

	_, err := o.Raw(sql, strconv.FormatInt(int64(j.Id), 10), past24Hours).QueryRows(&pings)
	if err != nil {
		// Select the pings we need to fill in
		for pingI := range pings {
			sqlPing := pings[pingI]
			j.Players[strconv.FormatInt(int64(pingI), 10)] = sqlPing.Online
		}

		// Cap to a maximum of 300 data pointers
		length := len(j.Players)
		skip := 0

		// Calc which we should skip
		if length > 3000 {
			skip = (length - 3000) / 3000

			// Remap if we need to
			tempMap := make(map[string]int32)
			counter := 0
			for playerI := range j.Players {
				if skip > counter {
					counter++
					continue
				}

				counter = 0
				tempMap[playerI] = j.Players[playerI]
			}

			j.Players = tempMap
		}
	}
}
