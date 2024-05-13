package main

import (
	"net/http"
	"time"
)

type Client *http.Client

func NewClient(timeout uint) Client {
	return Client(&http.Client{
		Timeout: time.Second * time.Duration(timeout),
	})
}
