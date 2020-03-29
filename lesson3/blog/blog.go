package main

import (
	"html/template"
	"log"
	"net/http"
)

var tmpl = template.Must(template.New("MyTemplate").ParseFiles("template.html"))

func main() {
	port := "8080"

	router := http.NewServeMux()

	router.HandleFunc("/", viewBlog)

	log.Printf("start listen on port %v", port)

	log.Fatal(http.ListenAndServe(":"+port, router))
}

func viewBlog(w http.ResponseWriter, r *http.Request) {
	if err := tmpl.ExecuteTemplate(w, "blog", simpleBlog); err != nil {
		log.Println(err)
	}
}
