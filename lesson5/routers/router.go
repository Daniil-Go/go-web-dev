package routers

import (
	"database/sql"
	"fmt"
	"log"

	"not-for-work/GeekBrainsWebinars/current-lessons/5/99_hw/habr-posts/controllers"
	"not-for-work/GeekBrainsWebinars/current-lessons/5/99_hw/habr-posts/models"

	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
)

const DSN = "root:1234@tcp(localhost:3306)/posts?charset=utf8"

func init() {
	db, err := sql.Open("mysql", DSN)
	if err != nil {
		log.Fatal(err)
	}
	// defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	truncate(db)
	insertDefault(db)

	beego.Router("/", &controllers.MainController{
		Db: db,
	})

	beego.Router("/post/:id", &controllers.SinglePost{
		Db: db,
	})

	beego.Router("/prepare/:id", &controllers.Prepare{
		Db: db,
	})
}

const (
	databaseName = "posts"
	tableName    = "habr_posts"
)

func truncate(db *sql.DB) {
	query := fmt.Sprintf("truncate %s.%s;", databaseName, tableName)
	_, err := db.Exec(query)
	if err != nil {
		log.Println(err)
	}
}

func insertDefault(db *sql.DB) {
	for _, post := range createPosts() {
		query := fmt.Sprintf(
			`insert into %s.%s (id,title,date,link,comment) values (?,?,?,?,?);`, databaseName, tableName)
		_, err := db.Exec(query, post.Id, post.Title, post.Date, post.Link, post.Comment)
		if err != nil {
			log.Println(err)
			return
		}
	}
	log.Print("inserted default")
}

func createPosts() []models.Post {
	return []models.Post{
		{
			Id:      1,
			Title:   "Hello world! Or Habr in English, v1.0",
			Date:    "15.01.2019 in 14:15",
			Link:    "https://habr.com/ru/company/habr/blog/435764/",
			Comment: "Nice one!",
		},
		{
			Id:      2,
			Title:   "Common Errors in English Usage",
			Date:    "29.06.2010 in 19:51",
			Link:    "https://habr.com/ru/post/97778/",
			Comment: "hmst",
		},
	}
}
