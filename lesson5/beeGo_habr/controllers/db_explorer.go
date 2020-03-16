package controllers

import (
	"database/sql"
	"fmt"
	"log"

	"not-for-work/GeekBrainsWebinars/current-lessons/5/99_hw/habr-posts/models"
)

func getPosts(db *sql.DB) ([]models.Post, error) {
	res := make([]models.Post, 0, 1)

	rows, err := db.Query("select * from posts.habr_posts")
	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		post := models.Post{}

		if err := rows.Scan(&post.Id, &post.Title, &post.Date, &post.Link, &post.Comment); err != nil {
			log.Println(err)
			continue
		}

		res = append(res, post)
	}

	return res, nil
}

func getPost(db *sql.DB, id string) (models.Post, error) {
	row := db.QueryRow(fmt.Sprintf("select * from posts.habr_posts WHERE id = %v", id))

	post := models.Post{}
	if err := row.Scan(&post.Id, &post.Title, &post.Date, &post.Link, &post.Comment); err != nil {
		return models.Post{}, err
	}

	return post, nil
}

func addPost(db *sql.DB, post models.Post) error {
	_, err := db.Exec(`INSERT into posts.habr_posts (id,title,date,link,comment) values (?,?,?,?,?);`,
		post.Id, post.Title, post.Date, post.Link, post.Comment)

	return err
}

func editPost(db *sql.DB, post *models.Post, id string) error {
	query := fmt.Sprintf(`UPDATE posts.habr_posts SET title="%s", date="%s", link="%s", comment="%s"  where id=?;`,
		post.Title, post.Date, post.Link, post.Comment)
	_, err := db.Exec(query, id)

	return err
}

func deletePost(db *sql.DB, id string) error {
	_, err := db.Exec(`DELETE FROM posts.habr_posts where id=?;`, id)

	return err
}
