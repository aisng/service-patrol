package main

import (
	"net/http"
	"time"
)

func NewHttpClient(timeout uint) *http.Client {
	return &http.Client{
		Timeout: time.Second * time.Duration(timeout),
	}
}
