package controllers

import (
	"github.com/astaxie/beego"
	"go.mongodb.org/mongo-driver/mongo"
)

type Explorer struct {
	Db           *mongo.Client
	DbName       string
	DbCollection string
}

type MainController struct {
	beego.Controller
	Explorer Explorer
}

func (c *MainController) Get() {
	posts, err := c.Explorer.getPosts()
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(500)
		_, _ = c.Ctx.ResponseWriter.Write([]byte(err.Error()))
		return
	}

	c.Data["Posts"] = posts
	c.TplName = "index.tpl"
}
