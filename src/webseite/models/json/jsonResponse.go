package json

import (
	gojson "encoding/json"
	"github.com/astaxie/beego"
	"webseite/websocket"
)

type JSONResponse struct {
	Ident string
	Value interface{}
}

func (j *JSONResponse) marshal() []byte {
	jsonBytes, err := gojson.Marshal(j)
	if err != nil {
		beego.BeeLogger.Warn("Could not convert to json: %v", err)
		return make([]byte, 0)
	}

	return jsonBytes
}

func (j *JSONResponse) BroadcastToServerID(serverID int32) {
	jsonBytes := j.marshal()
	if len(jsonBytes) > 0 {
		for c := range websocket.Hub.Connections {
			allowedServers := c.Session.Get("servers").(map[int32]bool)
			if val, ok := allowedServers[serverID]; !ok || !val {
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

func (j *JSONResponse) Send(c *websocket.Connection) {
	defer func() {
		recover()
	}()

	by := j.marshal()
	if len(by) != 0 {
		c.Send <- by
	}
}

func (j *JSONResponse) Broadcast() {
	by := j.marshal()
	if len(by) != 0 {
		websocket.Hub.Broadcast <- by
	}
}
