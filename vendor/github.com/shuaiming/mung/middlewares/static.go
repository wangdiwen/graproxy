package middlewares

import (
	"net/http"
	"path"
	"strings"
)

// Static serve static files
type Static struct {
	Dir       http.FileSystem
	Prefix    string
	IndexFile string
}

// NewStatic make new Static
func NewStatic(prefix string, fs http.FileSystem, index string) *Static {
	return &Static{
		Prefix:    prefix, // statics url prefix
		Dir:       fs,
		IndexFile: index,
	}
}

// ServeHTTP make mung middleware
func (s *Static) ServeHTTP(
	rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	file := r.URL.Path

	// only surpport GET and HEAD method for statics
	if r.Method != "GET" && r.Method != "HEAD" {
		next(rw, r)
		return
	}

	// if we have a prefix, filter requests by stripping the prefix
	if s.Prefix != "" {
		// bypass the requests not matching the prefix
		if !strings.HasPrefix(file, s.Prefix) {
			next(rw, r)
			return
		}
		// remove the prefix from request url
		file = file[len(s.Prefix):]
		if file != "" && file[0] != '/' {
			next(rw, r)
			return
		}
	}

	f, err := s.Dir.Open(file)
	if err != nil {
		// discard the error?
		next(rw, r)
		return
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		next(rw, r)
		return
	}

	// try to serve index file
	if fi.IsDir() {
		// redirect if missing trailing slash
		if !strings.HasSuffix(r.URL.Path, "/") {
			http.Redirect(rw, r, r.URL.Path+"/", http.StatusFound)
			return
		}

		file = path.Join(file, s.IndexFile)
		f, err = s.Dir.Open(file)
		if err != nil {
			next(rw, r)
			return
		}
		defer f.Close()

		fi, err = f.Stat()
		if err != nil || fi.IsDir() {
			next(rw, r)
			return
		}
	}

	// write the content
	http.ServeContent(rw, r, file, fi.ModTime(), f)
}
