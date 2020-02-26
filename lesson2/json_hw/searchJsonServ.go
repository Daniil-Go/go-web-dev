package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type SearchJSON struct {
	Search  string   `json:"search"`
	Sites   []string `json:"sites"`
	Results []string `json:"results"`
}

var searchStruct = SearchJSON{
	Search:  "",
	Sites:   []string{"https://yandex.ru", "https://golang.org", "https://google.com", "https://github.com"},
	Results: []string{""},
}

func main() {

	router := http.NewServeMux()
	router.HandleFunc("/", startHandle)
	router.HandleFunc("/search", searchHandle)
	router.HandleFunc("/results", resultsHandle)
	router.HandleFunc("/cookie", setCookieHandle)

	port := "8081"
	log.Printf("start listen on port %v", port)

	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Printf("%v", err)
	}

}

func startHandle(wr http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintf(wr, "Welcome to URL search! Type \"/search?param=YourSearchRequest\" in URL string!\n\t"+
		"To set cookie with your name type \"/cookie?name=YourName\" in URL string!")
}

func searchHandle(wr http.ResponseWriter, req *http.Request) {
	searchReq := req.URL.Query().Get("param")
	searchStruct.Search = searchReq
	_, _ = fmt.Fprintf(wr, "Your search request is \"%s\" \n\t"+
		"Type \"/results\" in URL string for results ", searchStruct.Search)
}

func resultsHandle(wr http.ResponseWriter, req *http.Request) {
	searchResults, _ := search(searchStruct.Search, searchStruct.Sites)
	searchStruct.Results = searchResults
	_, _ = fmt.Fprintf(wr, "The results of searching for \"%s\" \n\t"+
		"are: %s \n\tAlso available in %s JSON file", searchStruct.Search, searchStruct.Results, searchStruct.Search)
	log.Print(marshallJSON(searchStruct.Search, searchStruct))
}

func setCookieHandle(wr http.ResponseWriter, req *http.Request) {
	getName := req.URL.Query().Get("name")
	http.SetCookie(wr, &http.Cookie{
		Name:    getName,
		Value:   "Go-Web-Dev:Lesson2",
		Expires: time.Now().Add(time.Minute * 5),
	})

	_, _ = fmt.Fprintf(wr, "Set cookie with name \"%s\"", getName)

}

func search(str string, sites []string) ([]string, int) {
	out := make([]string, 0, 1)
	errs := 0

	for _, site := range sites {
		res, err := getReq(site)
		if err != nil {
			errs++
			log.Print(err)
			continue
		}

		if strings.Contains(string(res), str) {
			out = append(out, site)
		}
	}

	return out, errs
}

func getReq(reqURL string) ([]byte, error) {
	resp, err := http.Get(reqURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func marshallJSON(filename string, data interface{}) error {
	obj, _ := json.Marshal(data)
	err := ioutil.WriteFile(filename, obj, 0755)
	if err != nil {
		return err
	}
	return nil
}
