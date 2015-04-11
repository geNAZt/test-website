package websocket

import (
	"bytes"
	"time"
	"webseite/models/json"
	"webseite/websocket"
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
		server := json.GetServer(serverId)
		if server.Id != -1 && len(server.Favicons) > 1 {
			abort := false
			go func() {
				for faviconI := range server.Favicons {
					if abort {
						break;
					}

					favicon := server.Favicons[faviconI]

					if favicon.Icon == "" {
						continue
					}

					serverFavicon := &json.ServerFavicon{
						Id: serverId,
						Icon: favicon.Icon,
					}

					abort = !serverFavicon.Send(m.Connection)
					time.Sleep(time.Duration(favicon.DisplayTime) * time.Millisecond)
				}
			}()
		}
	}
}
