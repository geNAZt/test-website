package json

import (
	"webseite/websocket"
)

type JSONMaxPlayerResponse struct {
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

	jsonBytes := jsonResponse.marshal()
	if len(jsonBytes) > 0 {
		for c := range websocket.Hub.Connections {
			allowedServers := c.Session.Get("servers").(map[int32]bool)
			if val, ok := allowedServers[j.Id]; !ok || !val {
				continue
			}

			select {
			case c.Send <- jsonBytes:
			default:
				c.CloseCustomChannels()
				close(c.Send)
				delete(websocket.Hub.Connections, c)
			}
		}
	}
}
