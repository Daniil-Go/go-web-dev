package controllers

import (
	"database/sql"

	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
	Db *sql.DB
}

func (c *MainController) Get() {
	posts, err := getPosts(c.Db)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(500)
		_, _ = c.Ctx.ResponseWriter.Write([]byte(err.Error()))
		return
	}

	c.Data["Posts"] = posts
	c.TplName = "index.tpl"
}
