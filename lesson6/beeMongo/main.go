package main

import (
	_ "not-for-work/GeekBrainsWebinars/current-lessons/6/99_hw/habr-mongo/routers"
	"os"

	"github.com/astaxie/beego"
)

func main() {
	beego.Run("localhost", os.Getenv("httpport"))
}
