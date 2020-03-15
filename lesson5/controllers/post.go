package controllers

import (
	"database/sql"
	"encoding/json"
	"io"
	"io/ioutil"
	"strconv"

	"not-for-work/GeekBrainsWebinars/current-lessons/5/99_hw/habr-posts/models"

	"github.com/astaxie/beego"
)

type SinglePost struct {
	beego.Controller
	Db *sql.DB
}

func (c *SinglePost) Get() {
	post, err := getPost(c.Db, c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(500)
		_, _ = c.Ctx.ResponseWriter.Write([]byte(err.Error()))
		return
	}

	c.Data["Post"] = post
	c.TplName = "single.tpl"
}

func (c *SinglePost) Post() {
	id, err := strconv.Atoi(c.Ctx.Request.FormValue("id"))
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(500)
		_, _ = c.Ctx.ResponseWriter.Write([]byte(err.Error()))
		return
	}

	post := models.Post{
		Id:      id,
		Title:   c.Ctx.Request.FormValue("title"),
		Date:    c.Ctx.Request.FormValue("date"),
		Link:    c.Ctx.Request.FormValue("link"),
		Comment: c.Ctx.Request.FormValue("comment"),
	}

	if err := addPost(c.Db, post); err != nil {
		c.Ctx.ResponseWriter.WriteHeader(500)
		_, _ = c.Ctx.ResponseWriter.Write([]byte(err.Error()))
		return
	}

	c.Redirect("/", 301)
}

/*
	curl -vX PUT -H "Content-Type: application/json"  -d'{"date":"date","link":"link","comment":"comment",
"title":"NewTitle"}' http://localhost:8090/post/
*/

func (c *SinglePost) Put() {
	id := c.Ctx.Input.Param(":id")
	if len(id) == 0 {
		c.Ctx.ResponseWriter.WriteHeader(500)
		_, _ = c.Ctx.ResponseWriter.Write([]byte("empty id"))
		return
	}

	post := new(models.Post)
	if err := readAndUnmarshall(post, c.Ctx.Request.Body); err != nil {
		c.Ctx.ResponseWriter.WriteHeader(500)
		_, _ = c.Ctx.ResponseWriter.Write([]byte(err.Error()))
		return
	}

	if err := editPost(c.Db, post, id); err != nil {
		c.Ctx.ResponseWriter.WriteHeader(500)
		_, _ = c.Ctx.ResponseWriter.Write([]byte(err.Error()))
		return
	}

	c.Redirect("/", 301)
}

func readAndUnmarshall(resp interface{}, body io.ReadCloser) error {
	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bytes, resp); err != nil {
		return err
	}

	return nil
}

/*
	curl -vX DELETE  http://localhost:8080/post/:id
*/

func (c *SinglePost) Delete() {
	id := c.Ctx.Input.Param(":id")
	if len(id) == 0 {
		c.Ctx.ResponseWriter.WriteHeader(500)
		_, _ = c.Ctx.ResponseWriter.Write([]byte("empty id"))
		return
	}

	if err := deletePost(c.Db, id); err != nil {
		c.Ctx.ResponseWriter.WriteHeader(500)
		_, _ = c.Ctx.ResponseWriter.Write([]byte(err.Error()))
		return
	}

	c.Redirect("/", 301)
}
