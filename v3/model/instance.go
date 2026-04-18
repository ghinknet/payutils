package model

import "net/http"

type HttpInstance interface {
	Post(path string, handler http.HandlerFunc)
}
