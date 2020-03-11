package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/pkg/errors"
)



type SearchParams struct {
	SearchString string   `json:"search"`
	Sites        []string `json:"sites"`
}

type ResponseParams struct {
	Sites  []string `json:"sites"`
	Errors int      `json:"errors,omitempty"`

	ServerError string `json:"error,omitempty"`
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/sites", searchHandler)

	port := "8080"
	log.Printf("start server on port %v", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

/*
 curl -vX POST -H "Content-Type: application/json" -d'
{"string": "html", "sites": ["https://yandex.ru", "https://golang.org"]}' http://localhost:8080/sites
try to send curl with bad json
*/

func searchHandler(wr http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		sendError(wr, http.StatusBadRequest, errors.New("it is not POST request"))
		return
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		sendError(wr, http.StatusInternalServerError, errors.Wrap(err, "read request body"))
		return
	}

	searchParams := new(SearchParams)
	if err := json.Unmarshal(body, searchParams); err != nil {
		sendError(wr, http.StatusInternalServerError, errors.Wrap(err, "unmarshal request"))
		return
	}

	result, errs := search(searchParams.SearchString, searchParams.Sites)
	response := ResponseParams{
		Sites:  result,
		Errors: errs,
	}

	bytes, err := json.Marshal(response)
	if err != nil {
		sendError(wr, http.StatusInternalServerError, errors.Wrap(err, "marshal response"))
		return
	}

	_, err = wr.Write(bytes)
	if err != nil {
		sendError(wr, http.StatusInternalServerError, errors.Wrap(err, "write to response"))
		return
	}
}

func sendError(wr http.ResponseWriter, statusCode int, err error) {
	wr.WriteHeader(statusCode)

	response := ResponseParams{
		ServerError: err.Error(),
	}

	bytes, err := json.Marshal(response)
	if err != nil {
		log.Print(errors.Wrap(err, "during marshal response"))
		return
	}

	_, err = wr.Write(bytes)
	if err != nil {
		log.Print(errors.Wrap(err, "during write response"))
		return
	}
}
