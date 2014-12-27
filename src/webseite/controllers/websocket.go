package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
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
	c := &connection{send: make(chan []byte, 256), ws: ws}
	quit := write(c)
	h.register <- c
	go c.writePump(quit)
	c.readPump()
}

func write(c *connection) chan struct{} {
	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				c.send <- []byte(time.Now().UTC().Format("dd.mm.yyyy"))
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
	return quit
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
