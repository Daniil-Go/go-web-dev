package main

import (
	"os"

	_ "not-for-work/GeekBrainsWebinars/current-lessons/5/99_hw/habr-posts/routers"

	"github.com/astaxie/beego"
)

func main() {
	beego.Run("localhost", os.Getenv("httpport"))
}
