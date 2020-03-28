package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func (c *Configuration) init() {
	bytes, err := ioutil.ReadFile("app.json")
	if err != nil {
		log.Fatal(err)
	}

	if err = json.Unmarshal(bytes, &c); err != nil {
		log.Fatal(err)
	}

	f, err := os.OpenFile(c.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	// Не понимаю, как правильно обработать (в каком месте)
	//defer f.Close()
	log.SetOutput(f)

	//log.Printf("%+v", Config)
}

func main() {

	var Config Configuration
	Config.init()

	/*
		const (
			endpointURL = "https://api.exchangeratesapi.io/latest?base=USD"
			port        = "8080"
		)

	*/

	var rates Cache

	go func() {
		for {
			rates.getRates(Config.EndpointUrl)
			log.Printf("Rates successfully updated")
			time.Sleep(60 * time.Second)
		}
	}()

	router := http.NewServeMux()

	router.HandleFunc("/rates", rates.ratesHandler)
	router.HandleFunc("/search", rates.searchHandler)

	log.Printf("Start server on port %v", Config.Port)
	log.Fatal(http.ListenAndServe(":"+Config.Port, router))

	/*
		curl -v http://localhost:8080/rates
		http://localhost:8080/search?base=USD&quote=RUB

	*/
}

func (rates *Cache) getRates(reqURL string) {
	resp, err := http.Get(reqURL)
	if err != nil {
		log.Println(fmt.Errorf("%v during GET request", err))
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(fmt.Errorf("%v during reading body of ratest", err))
		return
	}
	rates.mu.Lock()
	defer rates.mu.Unlock()
	if err = json.Unmarshal(body, &rates); err != nil {
		log.Println(fmt.Errorf("%v during Unmarshal", err))
		return
	}
}

func (rates *Cache) ratesHandler(wr http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		sendError(wr, http.StatusBadRequest, fmt.Errorf("this is not GET request"))
		return
	}
	bytes, err := json.Marshal(rates.Rates)
	if err != nil {
		sendError(wr, http.StatusInternalServerError, fmt.Errorf("%v marshal response", err))
		return
	}
	_, err = wr.Write(bytes)
	if err != nil {
		sendError(wr, http.StatusInternalServerError, fmt.Errorf("%v write to response", err))
		return
	}

}

func (rates *Cache) searchHandler(wr http.ResponseWriter, req *http.Request) {
	searchParams := new(SearchParams)
	searchParams.Base = req.URL.Query().Get("base")
	searchParams.Quote = req.URL.Query().Get("quote")

	out := new(Response)

	rates.mu.Lock()
	defer rates.mu.Unlock()
	base, ok := rates.Rates[searchParams.Base]
	if !ok {
		fmt.Fprintf(wr, "The base \"%s\" is not correct", searchParams.Base)
		return
	}
	quote, ok := rates.Rates[searchParams.Quote]
	if !ok {
		fmt.Fprintf(wr, "The quote \"%s\" is not correct", searchParams.Quote)
		return
	}

	out.BaseName = searchParams.Base
	out.BaseRate = base
	out.QuoteName = searchParams.Quote
	out.QuoteRate = quote
	out.exchange()

	fmt.Fprintf(wr, "Base is \"%s\" \nQuote is \"%s\" \nRate is \"%v\"",
		out.BaseName, out.QuoteName, out.Rate)
}

func (c *Response) exchange() {
	c.Rate = 1 / c.BaseRate * c.QuoteRate
	return
}

func sendError(wr http.ResponseWriter, statusCode int, err error) {
	wr.WriteHeader(statusCode)

	response := Response{
		Error: err.Error(),
	}

	bytes, err := json.Marshal(&response)
	if err != nil {
		log.Print(fmt.Errorf("%v during marshal response", err))
		return
	}

	_, err = wr.Write(bytes)
	if err != nil {
		log.Print(fmt.Errorf("%v during write response", err))
		return
	}
}
