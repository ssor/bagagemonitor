package routers

import (
	"bagageMonitor/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/weixin", &controllers.MainController{}, "get:Index")
	beego.Router("/weixin", &controllers.MainController{}, "post:ReceiveMsg")
}
