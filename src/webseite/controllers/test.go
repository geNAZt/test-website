package controllers

import (
	"github.com/astaxie/beego"
	//"github.com/astaxie/beego/orm"
	//"webseite/models"
)

type TestController struct {
	beego.Controller
}

func (c *TestController) Get() {
	/*// ORM
	o := orm.NewOrm()
	o.Using("default")

	// Create new User
	user := new(models.User)
	user.Name = "test"

	// Insert new Object
	o.Insert(user)

	// Get all Users
	var users []models.User

	// Build up the Query
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("user.name").
		From("user").
		OrderBy("name").Desc().
		Limit(10)

	// Get the SQL Statement and execute it
	sql := qb.String()
	o.Raw(sql).QueryRows(&users)

	var counter int32
	counter = 0

	v := c.GetSession("counter")
	if v == nil {
		c.SetSession("counter", counter)
	} else {
		counter = v.(int32)
	}

	counter++

	c.SetSession("counter", counter)*/

	// Put the data into the template
	c.TplNames = "test.tpl"
	c.Data["host"] = "127.0.0.1:8080"
}
