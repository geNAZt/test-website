package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"time"
	"webseite/websocket"
)

type WSController struct {
	beego.Controller
}

func (w *WSController) Get() {
	w.EnableRender = false

	ws, err := websocket.Upgrade(w.Ctx.ResponseWriter, w.Ctx.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	conn := websocket.CreateConnection(ws)
	conn.AppendChannel(write(conn))
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
