package json

import (
	"webseite/websocket"
)

type JSONMaxPlayerResponse struct {
	JSONResponse
	Id         int32
	MaxPlayers int32
}

func (j *JSONMaxPlayerResponse) Send(c *websocket.Connection) {
	jsonResponse := JSONResponse{
		Ident: "maxPlayer",
		Value: j,
	}

	jsonResponse.Send(c)
}

func (j *JSONMaxPlayerResponse) Broadcast() {
	jsonResponse := JSONResponse{
		Ident: "maxPlayer",
		Value: j,
	}

	jsonResponse.BroadcastToServerID(j.Id)
}
