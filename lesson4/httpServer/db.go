package main

import (
	"database/sql"
	"fmt"
	"log"
)

// Читаем посты из БД
func getPosts(db *sql.DB) ([]Post, error) {
	// Делаем слайс под результат запроса в БД
	res := make([]Post, 0, 1)

	// Делаем запрос в БД
	rows, err := db.Query("select * from posts.habr_posts")
	if err != nil {
		return res, err
	}
	defer rows.Close()

	// Проходим по всем строчкам таблицы
	for rows.Next() {
		// Готовим структуру
		post := Post{}

		// Наполняем структуру post
		if err := rows.Scan(&post.Id, &post.Title, &post.Date, &post.Link, &post.Comment); err != nil {
			log.Println(err)
			continue
		}

		// Добавляем в слайс res post1, post2...
		res = append(res, post)
	}

	return res, nil
}

// Получение поста по id
func getPost(db *sql.DB, id string) (Post, error) {
	// Делаем запрос
	row := db.QueryRow(fmt.Sprintf("select * from posts.habr_posts WHERE id = %v", id))

	// Готовим структуру
	post := Post{}
	// Наполняем стуктуру
	if err := row.Scan(&post.Id, &post.Title, &post.Date, &post.Link, &post.Comment); err != nil {
		return Post{}, err
	}

	return post, nil
}

// Редактировани поста
func editPost(db *sql.DB, post Post, id string) error {
	// Обновляем данные в таблице
	query := fmt.Sprintf(`UPDATE posts.habr_posts SET title="%s", date="%s", link="%s", comment="%s"  where id=?;`,
		post.Title, post.Date, post.Link, post.Comment)
	_, err := db.Exec(query, id)

	return err
}
