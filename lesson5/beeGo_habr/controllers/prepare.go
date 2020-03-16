package controllers

import (
	"database/sql"

	"github.com/astaxie/beego"
)

type Prepare struct {
	beego.Controller
	Db *sql.DB
}

func (c *Prepare) Get() {
	c.TplName = "edit.tpl"
}
