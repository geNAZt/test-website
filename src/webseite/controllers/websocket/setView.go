package websocket

import (
	"bytes"
	"webseite/models/json"
	"webseite/websocket"
	"time"
	"fmt"
)

func init() {
	go listenSetView()
}

func listenSetView() {
	c := websocket.Hub.Listen(func(message websocket.Message) bool {
		return bytes.Index(message.Message, []byte("setview:")) != -1
	})

	for {
		select {
		case m := <-c:
			go setView(m)
		}
	}
}

func setView(m websocket.Message) {
	viewId := ParseInt(m)
	if viewId == -1 {
		return
	}

	m.Connection.Session.Set("view", viewId)
}
