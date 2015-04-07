package websocket

import (
	"strconv"
	"strings"
	"webseite/websocket"
)

func ParseInt(m websocket.Message) int32 {
	message := string(m.Message)

	if !strings.Contains(message, ":") {
		return -1
	}

	serverIdStr := strings.Split(message, ":")[1]
	serverId, errParse := strconv.ParseInt(serverIdStr, 10, 32)
	if errParse != nil {
		return -1
	}

	return int32(serverId)
}

func ParseInts(m websocket.Message) []int32 {
	message := string(m.Message)

	if !strings.Contains(message, ":") {
		return []int32{-1}
	}

	splits := strings.Split(message, ":")
	ints := make([]int32, len(splits) - 1)
	counter := 0

	for key := range splits {
		serverId, errParse := strconv.ParseInt(splits[key], 10, 32)
		if errParse != nil {
			ints[counter] = serverId
			counter++;
		}
	}

	return ints
}
