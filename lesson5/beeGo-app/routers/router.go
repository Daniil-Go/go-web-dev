package routers

import (
	"go-web-dev/lesson5/beeGo-app/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
}
