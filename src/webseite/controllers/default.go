package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["host"] = "127.0.0.1:8080"
	c.TplNames = "index.tpl"
}
