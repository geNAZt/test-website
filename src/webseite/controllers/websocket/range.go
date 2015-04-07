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
	days := ParseInt(m)
	if days == -1 {
		return
	}

	m.Connection.Session.Set("days", int(days))
}
