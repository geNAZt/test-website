package controllers
import (
    "github.com/astaxie/beego"
    "github.com/astaxie/beego/orm"
    "webseite/models"
)

type AcceptController struct {
    beego.Controller
}

func (c *AcceptController) Get() {
    // Get the token out of the routing
    token := c.Ctx.Input.Param(":token")

    // ORM
    o := orm.NewOrm()

    // Build up the Query
    qb, _ := orm.NewQueryBuilder("mysql")
    qb.Select("*").
        From("user").
        Where("`accept_token` = ?")

    // Get the SQL Statement and execute it
    sql := qb.String()
    user := []models.User{}
    o.Raw(sql, token).QueryRows(&user)

    // Flash for MainPage
    c.SetSession("profile.registerComplete", "We found your accept Token for EMail " + user[0].Email)

    // Redirect back to the Mainpage
    c.Redirect("/", 302)
}