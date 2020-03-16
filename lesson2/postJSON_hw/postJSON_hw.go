package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	resp, err := http.Post(
		"http://localhost:8080/post",
		"application/json",
		strings.NewReader(`{"Search":"git","Sites":["hhttps://yandex.ru", "https://golang.org", "https://google.com", 
"https://github.com"]}`),
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body)
}
