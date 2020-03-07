package main

import (
	"html/template"
	"log"
	"net/http"
)


func main() {
	port := "8080"

	router := http.NewServeMux()

	router.HandleFunc("/", viewBlog)
	router.HandleFunc("/post", viewPost)

	log.Printf("start listen on port %v", port)

	log.Fatal(http.ListenAndServe(":"+port, router))
}

func viewBlog(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.New("blog").ParseFiles("template.html"))
	if err := tmpl.ExecuteTemplate(w, "blog", simpleBlog); err != nil {
		log.Println(err)
	}
}

func viewPost(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.New("post").ParseFiles("post.html"))
	if err := tmpl.ExecuteTemplate(w, "post", onePost); err != nil {
		log.Println(err)
	}
}
