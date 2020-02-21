package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	if err := reqGet(); err != nil{
		log.Fatal(err)
	}
}

const (
	url = "https://golang.org/"
)

func reqGet() error {
	resp, err := http.Get(url)
	if err !=nil{
		return err
	}
	defer resp.Body.Close()

	log.Printf("bodyOfRequest: %v", resp.Header)

	_, err = io.Copy(os.Stdout, resp.Body)
	if err != nil {
		return err
	}

	return nil
}