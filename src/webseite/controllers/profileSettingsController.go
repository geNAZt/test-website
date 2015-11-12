package controllers

import (
	"github.com/astaxie/beego"
	"fmt"
	"webseite/models"
	"github.com/astaxie/beego/orm"
)

type ProfileSettingsController struct {
	beego.Controller
}

func (c *ProfileSettingsController) Get() {
	c.Data["user"] = c.GetSession("user")
	c.TplNames = "profile/settings.tpl"
}

func (c *ProfileSettingsController) Post() {
	o := orm.NewOrm()

	_, meta, err := c.GetFile("avatarfile")
	if err != nil {
		// Check if Gravatar is enabled
		if ok, err := c.GetBool("avatargravatar"); ok && err == nil {
			var user models.User = c.GetSession("user").(models.User)
			user.Avatar = "gravatar:" + user.Email

			// Build up the Query
			qb, _ := orm.NewQueryBuilder("mysql")
			qb.Update("user").Set("avatar='" + user.Avatar + "'").Where("`email` = ?")
			o.Raw(qb.String(), user.Email).Exec()
		} else {
			c.SetSession("settings.avatar", "No Avatar selected. Upload one or tick Gravatar")
		}

		// Redirect to Settings page
		c.Redirect("/profile/settings/", 302)
		return
	}

	fmt.Printf("%v", meta.Filename)

	c.Redirect("/profile/settings/", 302)
}