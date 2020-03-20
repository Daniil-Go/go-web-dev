package models

type Post struct {
	Id      string `json:"id"`
	Title   string `json:"title"`
	Date    string `json:"date"`
	Link    string `json:"link"`
	Comment string `json:"comment"`
}
