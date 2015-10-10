package controllers

import (
	"crypto/sha512"
	"encoding/base64"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"webseite/models"
)

type ProfileLoginController struct {
	beego.Controller
}

func (c *ProfileLoginController) Get() {
	errEmail := c.GetSession("profile.emailError")
	errPass := c.GetSession("profile.passwordError")

	c.DelSession("profile.emailError")
	c.DelSession("profile.passwordError")

	errors := make(map[string]interface{})
	if errEmail != nil {
		errors["email"] = errEmail
	}

	if errPass != nil {
		errors["password"] = errPass
	}

	c.Data["errors"] = errors
	c.Data["lastEmail"] = c.GetSession("profile.login.lastEmail")
	c.TplNames = "profile/login.tpl"

	if c.GetSession("profile.login.lastEmail") != nil {
		c.DelSession("profile.login.lastEmail")
	}
}

func (c *ProfileLoginController) Post() {
	// ORM
	o := orm.NewOrm()

	email := c.GetString("email")
	pass := c.GetString("password")

	c.SetSession("profile.login.lastEmail", email)
	error := false

    // Build up the Query
    qb, _ := orm.NewQueryBuilder("mysql")
    qb.Select("*").
    From("user").
    Where("`email` = ?")

    // Get the SQL Statement and execute it
    sql := qb.String()
    user := []models.User{}
    o.Raw(sql, email).QueryRows(&user)

    if len(user) > 0 {
        // Generate the pw hash from the input
        salt := user[0].Salt
        hasher := sha512.New()
        hasher.Write([]byte(salt))
        hasher.Write([]byte(pass))

        // Check if password is correct
        if base64.URLEncoding.EncodeToString(hasher.Sum(nil)) == user[0].PassHash {
            c.SetSession("user", user[0]);

            flash := beego.NewFlash();
            flash.Success("You have been logged in")
            flash.Store(&c.Controller)
        } else {
            c.SetSession("profile.passwordError", "The entered Pasword is wrong")
            error = true
        }
    } else {
        c.SetSession("profile.emailError", "E-Mail could not be found")
        error = true
    }

	if error {
		c.Redirect("/profile/login/", 302)
		return
	}

	// Redirect back to the Mainpage
	c.Redirect("/", 302)
}
