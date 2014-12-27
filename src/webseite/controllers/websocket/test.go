package websocket

import (
	"bytes"
	"fmt"
	"webseite/websocket"
)

func init() {
	go listen()
}

func listen() {
	c := websocket.Hub.Listen(func(message *websocket.Message) bool {
		return bytes.Index(message.Message, []byte("time:")) != -1
	})

	for {
		select {
		case m := <-c:
			fmt.Println("got new time message: ")
			fmt.Printf("%v", m)
		}
	}
}
