package main

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson"
)

type Explorer struct {
	Db           *mongo.Client
	DbName       string
	DbCollection string
}

type Post struct {
	Id      string `json:"id"`
	Title   string `json:"title"`
	Date    string `json:"date"`
	Link    string `json:"link"`
	Comment string `json:"comment"`
}

func (e Explorer) getPosts() ([]Post, error) {
	c := e.Db.Database(e.DbName).Collection(e.DbCollection)

	cur, err := c.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, errors.Wrap(err, "Find")
	}

	res := make([]Post, 0, 1)
	if err := cur.All(context.Background(), &res); err != nil {
		return nil, errors.Wrap(err, "All")
	}

	return res, nil
}

func (e Explorer) getPost(id string) (Post, error) {
	c := e.Db.Database(e.DbName).Collection(e.DbCollection)

	filter := bson.D{{Key: "id", Value: id}}

	res := c.FindOne(context.Background(), filter)

	post := new(Post)
	if err := res.Decode(post); err != nil {
		return Post{}, errors.Wrap(err, "decode")
	}

	return *post, nil
}

func (e Explorer) addPost(post Post) error {
	c := e.Db.Database(e.DbName).Collection(e.DbCollection)
	_, err := c.InsertOne(context.Background(), post)

	return err
}

func (e Explorer) editPost(post *Post, id string) error {
	filter := bson.D{{Key: "id", Value: id}}

	update := createUpdates(*post)

	c := e.Db.Database(e.DbName).Collection(e.DbCollection)
	_, err := c.UpdateOne(context.Background(), filter, update)

	return err
}

func createUpdates(post Post) bson.D {
	update := bson.D{}
	if len(post.Title) != 0 {
		update = append(update, bson.E{Key: "title", Value: post.Title})
	}

	if len(post.Date) != 0 {
		update = append(update, bson.E{Key: "date", Value: post.Date})
	}

	if len(post.Link) != 0 {
		update = append(update, bson.E{Key: "link", Value: post.Link})
	}

	if len(post.Comment) != 0 {
		update = append(update, bson.E{Key: "comment", Value: post.Comment})
	}

	return bson.D{{"$set", update}}

}

func (e Explorer) deletePost(id string) error {
	filter := bson.D{{Key: "id", Value: id}}

	c := e.Db.Database(e.DbName).Collection(e.DbCollection)
	_, err := c.DeleteOne(context.Background(), filter)

	return err
}
