package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// the first version of search()
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