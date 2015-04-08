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
	defaultView, _ := beego.AppConfig.Int("DefaultView")

	w.EnableRender = false
	w.SetSession("days", 2)
	w.SetSession("view", int32(defaultView))
	json.SendView(websocket.Upgrade(w.Controller))
}
