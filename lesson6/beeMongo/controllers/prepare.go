package controllers

import (
	"github.com/astaxie/beego"
)

type Prepare struct {
	beego.Controller
}

func (c *Prepare) Get() {
	c.TplName = "edit.tpl"
}
