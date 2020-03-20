package routers

import (
	"context"
	"log"

	"not-for-work/GeekBrainsWebinars/current-lessons/6/99_hw/habr-mongo/controllers"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/astaxie/beego"
)

const (
	dbName         = "habr"
	collectionName = "habr_posts"
)

func init() {
	db, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Connect(context.Background()); err != nil {
		log.Fatal(err)
	}
	// defer db.Disconnect(context.Background())

	log.Print("connected")

	e := controllers.Explorer{
		Db:           db,
		DbName:       dbName,
		DbCollection: collectionName,
	}

	if err := e.Truncate(); err != nil {
		log.Fatal(err)
	}
	log.Print("truncated")

	if err := e.InsertDefault(); err != nil {
		log.Fatal(err)
	}
	log.Print("inserted default")

	beego.Router("/", &controllers.MainController{
		Explorer: e,
	})

	beego.Router("/post/", &controllers.SinglePost{
		Explorer: e,
	})

	beego.Router("/prepare/:id", &controllers.Prepare{})
}
