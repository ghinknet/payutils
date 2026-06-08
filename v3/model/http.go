package model

import "net/http"

// HttpInstance is the unified routing surface an http driver must expose.
//
// Every method registers a standard net/http handler on the underlying
// framework router for the corresponding HTTP verb. payutils itself currently
// only registers callbacks via Post, but the full verb set is part of the
// contract so http drivers stay future-proof and usable for arbitrary routes.
type HttpInstance interface {
	Get(path string, handler http.HandlerFunc)
	Post(path string, handler http.HandlerFunc)
	Put(path string, handler http.HandlerFunc)
	Patch(path string, handler http.HandlerFunc)
	Delete(path string, handler http.HandlerFunc)
	Head(path string, handler http.HandlerFunc)
	Options(path string, handler http.HandlerFunc)
	// Any registers the handler for all standard HTTP verbs.
	Any(path string, handler http.HandlerFunc)
}
