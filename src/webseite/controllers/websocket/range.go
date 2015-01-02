package websocket

import (
	"bytes"
	"webseite/models/json"
	"webseite/websocket"
)

func init() {
	go listenRange()
}

func listenRange() {
	c := websocket.Hub.Listen(func(message websocket.Message) bool {
		return bytes.Index(message.Message, []byte("range:")) != -1
	})

	for {
		select {
		case m := <-c:
			go sendRangePings(m)
		}
	}
}

func sendRangePings(m websocket.Message) {
	days := ParseServerId(m)
	if days == -1 {
		return
	}

	for serverI := range json.Servers.Value.([]json.Server) {
		jsonServer := json.Servers.Value.([]json.Server)[serverI]

		jsonPings := &json.JSONPingResponse{
			Id: jsonServer.Id,
		}

		jsonPings.FillPings(days)

		jsonResponse := &json.JSONResponse{
			Ident: "pings",
			Value: jsonPings,
		}

		jsonResponse.Send(m.Connection)
	}
}