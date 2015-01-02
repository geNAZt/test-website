package controllers

import (
	"github.com/astaxie/beego"
	"webseite/models/json"
	"webseite/websocket"
)

type WSController struct {
	beego.Controller
}

func (w *WSController) Get() {
	w.EnableRender = false
	w.SetSession("days", 2)

	conn := websocket.Upgrade(w.Controller)

	// Send this Client all known Servers
	json.SendAllServers(conn)
}
