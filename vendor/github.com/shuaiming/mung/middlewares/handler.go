package middlewares

import "net/http"

// Handler turn any http.Handler to mung middleware
type Handler struct {
	handler http.Handler
}

// NewHandler make new Handler
func NewHandler(handler http.Handler) *Handler {
	return &Handler{handler: handler}
}

// ServeHTTP make mung middlewares
func (h *Handler) ServeHTTP(
	rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	h.handler.ServeHTTP(rw, r)
	next(rw, r)
}
