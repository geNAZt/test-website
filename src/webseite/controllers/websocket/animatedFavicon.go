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
	serverId := ParseInt(m)
	if serverId == -1 {
		return
	}

	server := json.GetServer(serverId)
	if server.Id != -1 && len(server.Favicons) > 1 {
		for faviconI := range server.Favicons {
			favicon := server.Favicons[faviconI]

			if favicon.Icon == "" {
				continue
			}

			json.SendFavicon(m.Connection, serverId, favicon.Icon)
			time.Sleep(time.Duration(favicon.DisplayTime) * time.Millisecond)
		}
	}
}
