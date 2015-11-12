package controllers

import (
	"github.com/astaxie/beego"
)

type ProfileLogoutController struct {
	beego.Controller
}

func (c *ProfileLogoutController) Get() {
	// Logout
	c.DelSession("user")

	// Set flash
	flash := beego.NewFlash();
	flash.Success("You have been logged out")
	flash.Store(&c.Controller)

	// Redirect back to the Mainpage
	c.Redirect("/", 302)
}
