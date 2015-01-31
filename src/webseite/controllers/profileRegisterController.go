package controllers

import (
	"crypto/sha512"
	"encoding/base64"
	valid "github.com/asaskevich/govalidator"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"webseite/models"
	util "webseite/util"
)

type ProfileRegisterController struct {
	beego.Controller
}

func (c *ProfileRegisterController) Get() {
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
	c.Data["lastEmail"] = c.GetSession("profile.register.lastEmail")
	c.TplNames = "profile/register.tpl"

	if c.GetSession("profile.register.lastEmail") != nil {
		c.DelSession("profile.register.lastEmail")
	}
}

func (c *ProfileRegisterController) Post() {
	// ORM
	o := orm.NewOrm()
	o.Using("default")

	email := c.GetString("email")
	pass := c.GetString("password")

	c.SetSession("profile.register.lastEmail", email)

	// Check for a valid E-Mail
	if email != "" {
		if !valid.IsEmail(email) {
			c.SetSession("profile.emailError", "E-Mail is not valid")
			c.Redirect("/profile/register/", 302)
			return
		} else {
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
				c.SetSession("profile.emailError", "E-Mail is already registered")
				c.Redirect("/profile/register/", 302)
				return
			}
		}
	} else {
		c.SetSession("profile.emailError", "E-Mail is empty")
		c.Redirect("/profile/register/", 302)
		return
	}

	// Check for a valid Password
	if pass == "" {
		c.SetSession("profile.passwordError", "Password is empty")
		c.Redirect("/profile/register/", 302)
		return
	}

	// Generate random salt
	salt := util.RandomString(32)
	hasher := sha512.New()
	hasher.Write([]byte(salt))
	hasher.Write([]byte(pass))

	// New user
	user := &models.User{
		Email:    email,
		Salt:     salt,
		PassHash: base64.URLEncoding.EncodeToString(hasher.Sum(nil)),
		Avatar:   "default",
	}
	o.Insert(user)

	// Flash for MainPage
	c.SetSession("profile.registerComplete", "Your registration has been completed")

	// Redirect back to the Mainpage
	c.Redirect("/", 302)
}
