package websocket

import (
	"bytes"
	"time"
	"webseite/models/json"
	"webseite/websocket"
	"fmt"
)

func init() {
	go listenAnimatedFavicon()
}

func listenAnimatedFavicon() {
	c := websocket.Hub.Listen(func(message websocket.Message) bool {
		return bytes.Index(message.Message, []byte("animated:")) != -1
	})

	for {
		select {
		case m := <-c:
			go displayAnimatedFavicon(m)
		}
	}
}

func displayAnimatedFavicon(m websocket.Message) {
	serverIds := ParseInts(m)
	if serverIds[0] == -1 {
		return
	}

	for serverI := range serverIds {
		serverId := serverIds[serverI]

		fmt.Printf("Looking up favicons for Server: %v", serverId)

		server := json.GetServer(serverId)
		fmt.Printf("Server has %v favicons cached", len(server.Favicons))

		if server.Id != -1 && len(server.Favicons) > 1 {
			for faviconI := range server.Favicons {
				favicon := server.Favicons[faviconI]

				if favicon.Icon == "" {
					continue
				}

				serverFavicon := &json.ServerFavicon{
					Id: serverId,
					Icon: favicon.Icon,
				}

				serverFavicon.Send(m.Connection)
				time.Sleep(time.Duration(favicon.DisplayTime) * time.Millisecond)
			}
		}
	}
}
