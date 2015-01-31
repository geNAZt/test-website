package json

import (
	"strconv"
	"time"
)

type JSONPingResponse struct {
	Id      int32
	Players map[string]int32
}

type PingCache struct {
	Players map[int32]map[int64]int32
}

var pingCache = PingCache{
	Players: make(map[int32]map[int64]int32),
}

func AddPing(id int32, time int64, players int32) {
	if _, ok := pingCache.Players[id]; !ok {
		pingCache.Players[id] = make(map[int64]int32)
	}

	pingCache.Players[id][time] = players
}

func (j *JSONPingResponse) FillPings(days int32) {
	pings := pingCache.Players[j.Id]

	// Construct pasttime and the map
	j.Players = make(map[string]int32)
	pastTime := time.Now().Add(time.Duration(-days*24*60) * time.Minute).Unix()

	// Select the pings we need to fill in
	for pingI := range pings {
		sqlPing := pings[int64(len(pings)-1)-pingI]
		if pingI < pastTime {
			continue
		}

		j.Players[strconv.FormatInt(pingI, 10)] = sqlPing
	}
}
