package main

type Blog struct {
	Title       string
	Description string
	Posts       []Post
}

type Post struct {
	ID   int
	Name string
	Body string
}

var simpleBlog = Blog{
	Title:       "Название блога",
	Description: "Описание блога",
	Posts: []Post{
		Post{1, "Первый пост", "Содержание первый пост"},
		Post{2, "Второй пост", "Содержание второй пост"},
		Post{3, "Третий пост", "Содержание третий пост"},
	},
}
