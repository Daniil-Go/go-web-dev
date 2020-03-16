package main

import (
	"log"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
)

func main() {
	router := http.NewServeMux()

	router.HandleFunc("/setCookie", setCookie)
	router.HandleFunc("/showCookie", showCookie)

	port := "8080"
	log.Printf("start server at port: %v", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

const cookieName = "hw-2-cookie"

func setCookie(wr http.ResponseWriter, req *http.Request) {
	http.SetCookie(wr, &http.Cookie{
		Name:    cookieName,
		Value:   uuid.NewV4().String(),
		Path:    "/",
		Expires: time.Now().Add(time.Minute * 5),
	})

	http.Redirect(wr, req, "http://localhost:8080/showCookie", http.StatusFound)
}

func showCookie(wr http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie(cookieName)
	if err != nil {
		log.Print(err)
		return
	}

	_, _ = wr.Write([]byte(cookie.Value))
}
