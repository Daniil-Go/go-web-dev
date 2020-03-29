package main

import (
	"os"
	"sync"
)

type Configuration struct {
	EndpointUrl string `json:"endpoint_url" check:"required"`
	Port        string `json:"port" check:"required"`
	LogFile     string `json:"log_file" check:"required"`
	outputFile  *os.File
}

type Cache struct {
	Rates map[string]float64 `json:"rates"`
	Base  string             `json:"base"`
	Date  string             `json:"date"`
	mu    sync.Mutex
}

type Response struct {
	BaseName  string
	BaseRate  float64
	QuoteName string
	QuoteRate float64
	Rate      float64
	Error     string
}

type SearchParams struct {
	Base  string
	Quote string
}
