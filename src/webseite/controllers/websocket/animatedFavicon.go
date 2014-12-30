package websocket

import (
	"bytes"
	"strings"
	"time"
	"webseite/models/json"
	"webseite/websocket"
)

func init() {
	go listen()
}

func listen() {
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
	servername := strings.Split(string(m.Message), ":")[1]
	server := json.GetServer(servername)
	if server != nil && len(server.Favicons) > 1 {
		for faviconI := range server.Favicons {
			favicon := server.Favicons[faviconI]

			if favicon.Icon == "" {
				continue
			}

			json.SendFavicon(m.Connection, servername, favicon.Icon)
			time.Sleep(time.Duration(favicon.DisplayTime) * time.Millisecond)
		}
	}
}
