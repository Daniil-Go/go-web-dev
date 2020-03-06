package main

type Blog struct {
	Title       string
	Description string
	Posts       []Post
}

type Post struct {
	ID      int
	Name    string
	Body    string
	Comment string
}

var simpleBlog = Blog{
	Title:       "Название блога",
	Description: "Описание блога",
	Posts: []Post{
		Post{ID: 1, Name: "Первый пост", Body: "Содержание первый пост", Comment: "Комментарий 1"},
		Post{ID: 2, Name: "Второй пост", Body: "Содержание второй пост", Comment: "Комментарий 2"},
		Post{ID: 3, Name: "Третий пост", Body: "Содержание третий пост", Comment: "Комментарий 3"},
	},
}

var onePost = Blog{
	Title:       "Название блога",
	Description: "Описание блога",
	Posts: []Post{
		Post{ID: 2, Name: "Второй пост", Body: "Содержание второй пост", Comment: "Комментарий 2"},
	},
}
