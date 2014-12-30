package controllers

import (
	"github.com/astaxie/beego"
	"time"
	"webseite/models/json"
	"webseite/websocket"
)

type WSController struct {
	beego.Controller
}

func (w *WSController) Get() {
	w.EnableRender = false

	conn := websocket.Upgrade(w.Controller)

	// Send this Client all known Servers
	json.SendAllServers(conn)
}

func write(c *websocket.Connection) chan struct{} {
	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				c.Send <- []byte("time:" + time.Now().UTC().String())
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
	return quit
}
