package websocket

import (
	"bytes"
	"webseite/models/json"
	"webseite/websocket"
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
	serverId := ParseServerId(m)
	if serverId == -1 {
		return
	}

	server := json.GetServer(int32(serverId))
	if server.Id != -1 {
		jsonPings := &json.JSONPingResponse{
			Id: server.Id,
		}

		jsonPings.FillPings(int32(2))

		jsonResponse := &json.JSONResponse{
			Ident: "pings",
			Value: jsonPings,
		}

		jsonResponse.Send(m.Connection)
	}
}
