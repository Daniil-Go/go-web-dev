package controllers

import (
	"context"

	"not-for-work/GeekBrainsWebinars/current-lessons/6/99_hw/habr-mongo/models"

	"github.com/pkg/errors"

	"go.mongodb.org/mongo-driver/bson"
)

func (e Explorer) getPosts() ([]models.Post, error) {
	c := e.Db.Database(e.DbName).Collection(e.DbCollection)

	cur, err := c.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, errors.Wrap(err, "Find")
	}

	res := make([]models.Post, 0, 1)

	if err := cur.All(context.Background(), &res); err != nil {
		return nil, errors.Wrap(err, "All")
	}

	return res, nil
}

func (e Explorer) getPost(id string) (models.Post, error) {
	c := e.Db.Database(e.DbName).Collection(e.DbCollection)

	filter := bson.D{{Key: "id", Value: id}}

	res := c.FindOne(context.Background(), filter)

	post := new(models.Post)
	if err := res.Decode(post); err != nil {
		return models.Post{}, errors.Wrap(err, "decode")
	}

	return *post, nil
}

func (e Explorer) addPost(post models.Post) error {
	c := e.Db.Database(e.DbName).Collection(e.DbCollection)
	_, err := c.InsertOne(context.Background(), post)

	return err
}

func (e Explorer) editPost(post *models.Post, id string) error {
	filter := bson.D{{Key: "id", Value: id}}

	update := createUpdates(*post)

	c := e.Db.Database(e.DbName).Collection(e.DbCollection)
	_, err := c.UpdateOne(context.Background(), filter, update)

	return err
}

func createUpdates(post models.Post) bson.D {
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

func createPosts() []models.Post {
	return []models.Post{
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
