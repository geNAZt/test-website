package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"webseite/models"
	"webseite/models/json"
	"webseite/websocket"
	"time"
)

type WSController struct {
	beego.Controller
}

func (w *WSController) Get() {
	w.EnableRender = false
	w.SetSession("days", 2)

	conn := websocket.Upgrade(w.Controller)

	start := time.Now()

	// Get the default View
	// ORM
	o := orm.NewOrm()
	o.Using("default")

	// Build up the Query
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("*").
		From("`view`").
		Where("`id` = ?")

	// Get the SQL Statement and execute it
	sql := qb.String()
	view := &models.View{}
	o.Raw(sql, 1).QueryRow(&view)

	o.LoadRelated(view, "Servers")

	// Send this Client all known Servers
	json.SendAllServers(conn, view)

	elapsed := time.Since(start)
	json.SendLog(conn, "Sending the View with servers took " + elapsed )
}
