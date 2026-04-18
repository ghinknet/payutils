package model

import "net/http"

type PayClient interface {
	Create(http.ResponseWriter, *http.Request)
	Callback(http.ResponseWriter, *http.Request)
}

type Client struct {
	PayClient map[string]PayClient
}
