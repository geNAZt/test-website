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

	jsonServers := json.GetPingResponse(serverId, int32(m.Connection.Session.Get("days").(int)))
	sendResponse := make([]*json.JSONPingResponse, len(jsonServers))

	for sId := range jsonServers {
		sendResponse = append(sendResponse, &json.JSONPingResponse{
			Id: jsonServers[sId].Id,
			Players: jsonServers[sId].Players,
		})
	}

	jsonResponse := &json.JSONResponse{
		Ident: "pings",
		Value: sendResponse,
	}

	jsonResponse.Send(m.Connection)

	elapsed := time.Since(start)
	json.SendLog(m.Connection, "Sending Pings took " + fmt.Sprintf("%s", elapsed) )
}
