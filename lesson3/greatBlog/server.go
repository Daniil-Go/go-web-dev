package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
)

type Posts map[int]Post

type Server struct {
	Title string
	Posts Posts

	Templates map[templateName]*template.Template
}

type Post struct {
	Id      int
	Title   string
	Date    string
	Link    string
	Comment string
}

func main() {
	router := http.NewServeMux()

	server := Server{
		Title:     "Posts from Habr",
		Posts:     createPosts(),
		Templates: createTemplates(),
	}

	router.HandleFunc("/", server.handlePostsList)
	router.HandleFunc("/post/", server.handleSinglePost)
	router.HandleFunc("/edit/", server.handleEditPost)

	port := "8080"
	log.Printf("start server on port: %v", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func (server *Server) handlePostsList(wr http.ResponseWriter, req *http.Request) {
	tmpl := getTemplate(server.Templates, List)
	if tmpl == nil {
		err := errors.New("empty template")
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	idVal := req.FormValue("id")
	if len(idVal) != 0 {
		id, err := strconv.Atoi(idVal)
		if err != nil {
			err := errors.Wrapf(err, "id from form value: %v", idVal)
			http.Error(wr, err.Error(), http.StatusInternalServerError)
			log.Print(err)
			return
		}

		post := Post{
			Id:      id,
			Title:   req.FormValue("title"),
			Date:    req.FormValue("date"),
			Link:    req.FormValue("link"),
			Comment: req.FormValue("comment"),
		}
		mu.Lock()
		server.Posts[id] = post
		mu.Unlock()
	}

	if err := tmpl.ExecuteTemplate(wr, "page", server); err != nil {
		err = errors.Wrap(err, "execute template")
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}
}

func (server *Server) handleSinglePost(wr http.ResponseWriter, req *http.Request) {
	tmpl := getTemplate(server.Templates, Single)
	if tmpl == nil {
		err := errors.New("empty template")
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	post, err := server.Posts.getID(req.URL.Query().Get("id"))
	if err != nil {
		err := errors.Wrap(err, "empty post")
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	if err := tmpl.ExecuteTemplate(wr, "page", post); err != nil {
		err := errors.Wrap(err, "execute template")
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}
}

func (server *Server) handleEditPost(wr http.ResponseWriter, req *http.Request) {
	tmpl := getTemplate(server.Templates, Edit)
	if tmpl == nil {
		err := errors.New("empty template")
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	post, err := server.Posts.getID(req.URL.Query().Get("id"))
	if err != nil {
		err := errors.Wrap(err, "empty post")
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	if err := tmpl.ExecuteTemplate(wr, "page", post); err != nil {
		err = errors.Wrap(err, "execute template")
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}
}
