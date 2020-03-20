package main

import (
	"context"
	"log"
	"reflect"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Case struct {
	Method dbMethod
	ID     string
	Post   Post

	ExpectedPosts []Post
	ExpectedPost  Post
}

type dbMethod string

const (
	readAll dbMethod = "readAll"
	readOne dbMethod = "readOne"
	create  dbMethod = "create"
	update  dbMethod = "update"
	delete  dbMethod = "delete"
)

func TestDb(t *testing.T) {
	e, err := initDb()
	if err != nil {
		t.Error(err)
		return
	}

	defer func() {
		_ = e.Truncate()
		_ = e.Db.Disconnect(context.Background())
	}()

	for i, c := range createCases() {
		var (
			posts = make([]Post, 0, 1)
			post  = Post{}
			err   error
		)

		switch c.Method {
		case readAll:
			posts, err = e.getPosts()
		case readOne:
			post, err = e.getPost(c.ID)
		case create:
			err = e.addPost(c.Post)
		case update:
			err = e.editPost(&c.Post, c.ID)
		case delete:
			err = e.deletePost(c.ID)
		default:
			t.Error("unknown method")
			continue
		}

		if err != nil {
			t.Error(err)
		}

		if c.Method == readAll {
			if !reflect.DeepEqual(posts, c.ExpectedPosts) {
				t.Errorf("[%d] Expected: %v; Result: %v", i, c.ExpectedPosts, posts)
				break
			}
		} else if c.Method == readOne {
			if !reflect.DeepEqual(post, c.ExpectedPost) {
				t.Errorf("[%d] Expected: %v; Result: %v", i, c.ExpectedPost, post)
				break
			}
		}
	}
}

func createCases() []Case {
	initialPost := createPosts()
	return []Case{
		// show all
		{
			Method:        readAll,
			ExpectedPosts: initialPost, // the first our state
		},
		// show one
		{
			Method: readOne,
			ID:     "1",
			ExpectedPost: Post{
				Id:      "1",
				Title:   "Hello world! Or Habr in English, v1.0",
				Date:    "15.01.2019 in 14:15",
				Link:    "https://habr.com/ru/company/habr/blog/435764/",
				Comment: "Nice one!",
			},
		},
		// add one
		{
			Method: create,
			Post: Post{
				Id:      "107",
				Title:   "Test Title",
				Date:    "test date",
				Link:    "test link",
				Comment: "test comment",
			},
		},
		// show all (check previous case)
		{
			Method: readAll,
			ExpectedPosts: append(initialPost, Post{
				Id:      "107",
				Title:   "Test Title",
				Date:    "test date",
				Link:    "test link",
				Comment: "test comment",
			}),
		},
		// edit one
		{
			Method: update,
			ID:     "107",
			Post: Post{
				Title:   "Test update Title",
				Date:    "test update date",
				Link:    "test update link",
				Comment: "test update comment",
			},
		},
		// show all (check previous case)
		{
			Method: readAll,
			ID:     "107",
			ExpectedPosts: append(initialPost, Post{
				Id:      "107",
				Title:   "Test update Title",
				Date:    "test update date",
				Link:    "test update link",
				Comment: "test update comment",
			}),
		},
		// delete one
		{
			Method: delete,
			ID:     "107",
		},
		// show all (check previous case)
		{
			Method:        readAll,
			ID:            "107",
			ExpectedPosts: initialPost,
		},
	}
}

func initDb() (Explorer, error) {
	db, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return Explorer{}, err
	}

	if err = db.Connect(context.Background()); err != nil {
		return Explorer{}, err
	}
	log.Print("mongo connected")

	e := Explorer{
		Db:           db,
		DbName:       "habr",
		DbCollection: "test_collection",
	}

	if err := e.InsertDefault(); err != nil {
		return Explorer{}, err
	}

	return e, nil
}

func (e Explorer) Truncate() error {
	c := e.Db.Database(e.DbName).Collection(e.DbCollection)
	_, err := c.DeleteMany(context.Background(), bson.D{})

	return err
}

func (e Explorer) InsertDefault() error {
	for _, post := range createPosts() {
		if err := e.addPost(post); err != nil {
			return err
		}
	}

	return nil
}

func createPosts() []Post {
	return []Post{
		{
			Id:      "1",
			Title:   "Hello world! Or Habr in English, v1.0",
			Date:    "15.01.2019 in 14:15",
			Link:    "https://habr.com/ru/company/habr/blog/435764/",
			Comment: "Nice one!",
		},
		{
			Id:      "2",
			Title:   "Common Errors in English Usage",
			Date:    "29.06.2010 in 19:51",
			Link:    "https://habr.com/ru/post/97778/",
			Comment: "hmst",
		},
	}
}
