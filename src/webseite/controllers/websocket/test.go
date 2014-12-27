package websocket

import (
	"bytes"
	"webseite/websocket"
)

func init() {
	go listen()
}

func listen() {
	websocket.Hub.Listen(func(message websocket.Message) bool {
		return bytes.Index(message.Message, []byte("time:")) != -1
	})

	/*for {
		select {
		case m := <-c:
			websocket.Hub.Broadcast <- m.Message
		}
	}*/
}
