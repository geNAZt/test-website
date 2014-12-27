package routers

import (
	"github.com/astaxie/beego"
	"webseite/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/test", &controllers.TestController{})
	beego.Router("/ws", &controllers.WSController{})
}
