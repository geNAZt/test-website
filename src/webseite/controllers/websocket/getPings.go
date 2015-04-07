package websocket

import (
	"bytes"
	"webseite/models/json"
	"webseite/websocket"
	"time"
	"fmt"
)

func init() {
	go listenGetPings()
}

func listenGetPings() {
	c := websocket.Hub.Listen(func(message websocket.Message) bool {
		return bytes.Index(message.Message, []byte("pings:")) != -1
	})

	for {
		select {
		case m := <-c:
			go sendPings(m)
		}
	}
}

func sendPings(m websocket.Message) {
	start := time.Now()

	serverId := ParseInts(m)
	if serverId[0] == -1 {
		return
	}

	jsonServers := json.GetPingResponse(serverId, int32(2))
	for sId := range jsonServers {
		server := json.GetServer(int32(jsonServers[sId].Id))
		if server.Id != -1 {
			jsonResponse := &json.JSONResponse{
				Ident: "pings",
				Value: &json.JSONPingResponse{
					Id: sId,
					Players: jsonServers[sId].Players,
				},
			}

			jsonResponse.Send(m.Connection)
		}
	}

	elapsed := time.Since(start)
	json.SendLog(m.Connection, "Sending Pings took " + fmt.Sprintf("%s", elapsed) )
}
