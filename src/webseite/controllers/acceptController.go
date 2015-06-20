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
    // Generate new Flash Data
    flash := beego.NewFlash()

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
    users := []models.User{}
    o.Raw(sql, token).QueryRows(&users)

    // Check found users
    if (len(users) < 1) {
        flash.Error("No user with this accept Token could not be found")
    } else {
        // Reset the token and update
        user := users[0]
        user.AcceptToken = ""
        o.Update(&user)

        // Flash for MainPage
        flash.Success("Your account has been activated. You can login now")
    }

    // Redirect back to the Mainpage
    flash.Store(&c.Controller)
    c.Redirect("/", 302)
}