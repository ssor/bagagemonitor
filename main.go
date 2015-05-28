package main

import (
	_ "bagageMonitor/routers"
	"github.com/astaxie/beego"
)

func main() {
	beego.SetStaticPath("/images", "static/img")

	beego.Run()
}
