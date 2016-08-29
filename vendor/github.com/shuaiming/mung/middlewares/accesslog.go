package middlewares

import (
	"io"
	"log"
	"net/http"
	"time"
)

// responseWriter warp http.ResponseWriter.
// There is no way to get http status code and written size by
// using default http.ResponseWriter.
type responseWriter struct {
	http.ResponseWriter
	wroteHeader bool
	status      int
	size        int
}

// WriteHeader warp http.ResponseWriter.WriteHeader
func (rw *responseWriter) WriteHeader(s int) {
	if !rw.wroteHeader {
		rw.wroteHeader = true
		rw.status = s
	}

	rw.ResponseWriter.WriteHeader(s)
}

// Write warp http.ResponseWriter.Write
func (rw *responseWriter) Write(b []byte) (int, error) {
	//	look at http.ResponseWriter.WriteHeader() implementation
	//	if !rw.wroteHeader {
	//		rw.WriteHeader(http.StatusOK)
	//	}
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

// Size of http server written bytes
func (rw *responseWriter) Size() int {
	return rw.size
}

// Status return http server status code
func (rw *responseWriter) Status() int {
	return rw.status
}

// AccessLog write access log to io.Writer
type AccessLog struct {
	log.Logger
}

// NewAccessLog make new AccessLog
func NewAccessLog(out io.Writer) *AccessLog {
	l := log.New(out, "", log.LstdFlags|log.Lmicroseconds)
	return &AccessLog{*l}
}

// ServeHTTP make mung middleware
func (al *AccessLog) ServeHTTP(
	rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	timeStart := time.Now()

	newrw := &responseWriter{rw, false, 0, 0}
	httpMethod := r.Method
	urlPath := r.URL.String()

	next(newrw, r)

	timeEnd := time.Now()
	du := timeEnd.Sub(timeStart)

	al.Printf(
		"%s %s %s %s %d %d",
		r.RemoteAddr, httpMethod, urlPath, du, newrw.Size(), newrw.Status())
}
