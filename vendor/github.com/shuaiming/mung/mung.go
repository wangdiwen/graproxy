package mung

import (
	"container/list"
	"log"
	"net/http"
)

// Middleware interface
// http.HandleFunc is the entry of next middleware.
type Middleware interface {
	ServeHTTP(http.ResponseWriter, *http.Request, http.HandlerFunc)
}

// Mung is tiny web framwork with middlewares supported.
type Mung struct {
	middlewares *list.List
	handle      http.HandlerFunc
}

// New mung instance
func New() *Mung {
	m := &Mung{list.New(), nil}
	m.rebuild()
	return m
}

// rebuild middlewares stack when any middleware added.
func (mung *Mung) rebuild() {
	mung.handle = chainMiddleware(mung.middlewares.Front())
}

// ServeHTTP make a http.Handler.
func (mung *Mung) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	mung.handle(rw, r)
}

// Use add one middleware to middlewares stack.
func (mung *Mung) Use(middleware Middleware) {
	mung.middlewares.PushBack(middleware)
	// rebuild the middlewares stack.
	mung.rebuild()
}

// Run bind address and serve in http.
func (mung *Mung) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, mung))
}

// RunTLS bind address and serve in https.
func (mung *Mung) RunTLS(addr, certFile, keyFile string) {
	log.Fatal(http.ListenAndServeTLS(addr, certFile, keyFile, mung))
}

// chainMiddleware build middlewares to stack with recursion
func chainMiddleware(el *list.Element) http.HandlerFunc {
	// nil is the last element of the middlewares chain,
	// and the bottom of middleware's stack.
	if el == nil {
		return func(rw http.ResponseWriter, r *http.Request) {}
	}
	return func(rw http.ResponseWriter, r *http.Request) {
		el.Value.(Middleware).ServeHTTP(rw, r, chainMiddleware(el.Next()))
	}
}
