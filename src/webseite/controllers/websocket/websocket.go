package websocket

import (
	"strconv"
	"strings"
	"webseite/websocket"
)

func ParseServerId(m websocket.Message) int32 {
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
