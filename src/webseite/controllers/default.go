package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["host"] = "localhost:8080"
	c.TplNames = "index.tpl"
}
