package routers

import (
	"github.com/astaxie/beego"
	"webseite/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/ws", &controllers.WSController{})

	// Register routes
	beego.Router("/profile/register/", &controllers.ProfileRegisterController{})
	beego.Router("/profile/login/", &controllers.ProfileLoginController{})
	beego.Router("/profile/logout/", &controllers.ProfileLogoutController{})
	beego.Router("/profile/settings/", &controllers.ProfileSettingsController{})
	beego.Router("/profile/accept/:token([A-Za-z]{64})", &controllers.AcceptController{})
}
