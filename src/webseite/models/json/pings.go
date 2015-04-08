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

var timestampCache map[int64]string

func init() {
	timestampCache = make(map[int64]string)
}

func getStringRepresentation(unix int64) string {
	// Check cache
	if val, ok := timestampCache[unix]; ok {
		return val
	}

	// Calc new string
	str := strconv.FormatInt(unix, 10);
	timestampCache[unix] = str
	return str
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
		In(sqlIds...).
		And("`time` > ?").
		OrderBy("`time`").
		Asc()

	// Ask the Database for 24h Ping
	sql := qb.String()

	var pings []orm.Params
	_, err := o.Raw(sql, past24Hours).Values(&pings)
	if err == nil {
		length := len(pings) / len(serverIds)
		shouldSkip := 0

		if length > 3000 {
			shouldSkip = (length - 3000) / 3000
		}

		// Select the pings we need to fill in
		for pingI := range pings {
			sqlPing := pings[pingI]

			serverId, _ := strconv.ParseInt(sqlPing["server_id"].(string), 10, 32);
			time, _ := time.ParseInLocation(createdFormat, sqlPing["time"].(string), time.Local)
			online, _ := strconv.ParseInt(sqlPing["online"].(string), 10, 32);

			serverId = int32(serverId)
			online = int32(online)

			if shouldSkip > 0 {
				if shouldSkip > skip[serverId] {
					skip[serverId]++
					continue
				}

				skip[serverId] = 0
				returnMap[serverId].Players[getStringRepresentation(time.Unix())] = online
			} else {
				returnMap[serverId].Players[getStringRepresentation(time.Unix())] = online
			}
		}
	}

	return returnMap
}