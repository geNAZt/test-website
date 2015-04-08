package json

import (
	"strconv"
	"time"
	"github.com/astaxie/beego/orm"
)

type JSONPingResponse struct {
	Id      int32
	Players map[string]int32
}

type TempPingRow struct {
	Time time.Time
	ServerId int32
	Online int32
}

func GetPingResponse(serverIds []int32, days int32) map[int32]*JSONPingResponse {
	// Prepare the map
	sqlIds := make([]string, len(serverIds))
	returnMap := make(map[int32]*JSONPingResponse)
	skip := make(map[int32]int)
	for sId := range serverIds {
		returnMap[serverIds[sId]] = &JSONPingResponse{
			Id: serverIds[sId],
			Players: make(map[string]int32),
		}

		skip[serverIds[sId]] = 0
		sqlIds[sId] = strconv.FormatInt(int64(serverIds[sId]),10)
	}

	// Construct pasttime and the map
	_, offset := time.Now().Zone()

	// ORM
	o := orm.NewOrm()
	o.Using("default")

	// Check for 24h Ping
	past24Hours := time.Unix( (time.Now().Add(time.Duration(-days*24*60) * time.Minute).Unix()) - int64(offset), 0 ).Format( createdFormat )

	// Build up the Query
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("`time`, `server_id`, `online`").
		From("`ping`").
		Where("`server_id`").
		In(sqlIds).
		And("`time` > ?").
		OrderBy("`time`").
		Asc()

	// Ask the Database for 24h Ping
	sql := qb.String()

	var pings []TempPingRow
	_, err := o.Raw(sql, past24Hours).QueryRows(&pings)
	if err == nil {
		length := len(pings) / len(serverIds)
		shouldSkip := 0

		if length > 3000 {
			shouldSkip = (length - 3000) / 3000
		}

		// Select the pings we need to fill in
		for pingI := range pings {
			sqlPing := pings[pingI]

			if shouldSkip > 0 {
				if shouldSkip > skip[sqlPing.ServerId] {
					skip[sqlPing.ServerId]++
					continue
				}

				skip[sqlPing.ServerId] = 0
				returnMap[sqlPing.ServerId].Players[strconv.FormatInt(sqlPing.Time.Unix(), 10)] = sqlPing.Online
			} else {
				returnMap[sqlPing.ServerId].Players[strconv.FormatInt(sqlPing.Time.Unix(), 10)] = sqlPing.Online
			}
		}
	}

	return returnMap
}