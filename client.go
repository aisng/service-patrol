package main

import (
	"net/http"
	"time"
)

type HttpClient *http.Client

func NewHttpClient(timeout uint) HttpClient {
	return HttpClient(&http.Client{
		Timeout: time.Second * time.Duration(timeout),
	})
}
