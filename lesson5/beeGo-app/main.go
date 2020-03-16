package main

import (
	_ "go-web-dev/lesson5/beeGo-app/routers"
	"github.com/astaxie/beego"
)

func main() {
	beego.Run("localhost")
}

