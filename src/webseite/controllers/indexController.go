package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	// Check if last email is set from
	if c.GetSession("profile.register.lastEmail") != nil {
		c.DelSession("profile.register.lastEmail")
	}

    // Read flash
    flash := beego.ReadFromRequest(&c.Controller)
    flashes := flash.Data

	// Check for login id
	if c.GetSession("userId") == nil {
		c.SetSession("userId", int32(-1))
	}

	c.Data["flash"] = flashes
	c.Data["host"] = "minecrafttracker.net"
	c.TplNames = "index.tpl"
}
