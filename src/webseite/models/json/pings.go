package json

import (
	"strconv"
	"time"
	"github.com/astaxie/beego/orm"
	"webseite/cache"
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

var timestampCache *cache.TimeoutCache
var parseCache *cache.TimeoutCache
var datetimeCache *cache.TimeoutCache

func init() {
	timestampCache, _ = cache.NewTimeoutCache(int64(24 * 31) * int64(time.Hour))
	parseCache, _ = cache.NewTimeoutCache(int64(24 * 31) * int64(time.Hour))
	datetimeCache, _ = cache.NewTimeoutCache(int64(24 * 31) * int64(time.Hour))
}

func getParsedTime(representation string) time.Time {
	// Check cache
	if val, ok := datetimeCache.Get(representation); ok {
		return val.(time.Time)
	}

	// Calc new time
	str, _ := time.ParseInLocation(createdFormat, representation, time.Local)
	datetimeCache.Add(representation, str)
	return str
}

func getParsedInt(representation string) int32 {
	// Check cache
	if val, ok := parseCache.Get(representation); ok {
		return val.(int64)
	}

	// Calc new integer
	str, _ := strconv.ParseInt(representation, 10, 32);
	parseCache.Add(representation, str)
	return int32(str)
}

func getStringRepresentation(unix int64) string {
	// Check cache
	if val, ok := timestampCache.Get(unix); ok {
		return val.(string)
	}

	// Calc new string
	str := strconv.FormatInt(unix, 10);
	timestampCache.Add(unix, str)
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

			serverId, _ := getParsedInt(sqlPing["server_id"].(string))
			time, _ := getParsedTime(sqlPing["time"].(string))
			online, _ := getParsedInt(sqlPing["online"].(string))

			if shouldSkip > 0 {
				if shouldSkip > skip[int32(serverId)] {
					skip[int32(serverId)]++
					continue
				}

				skip[int32(serverId)] = 0
				returnMap[int32(serverId)].Players[getStringRepresentation(time.Unix())] = int32(online)
			} else {
				returnMap[int32(serverId)].Players[getStringRepresentation(time.Unix())] = int32(online)
			}
		}
	}

	return returnMap
}