package main

import (
	"os"

	_ "go-web-dev/lesson5/beeGo_habr/routers"

	"github.com/astaxie/beego"
)

func main() {
	beego.Run("localhost", os.Getenv("httpport"))
}
