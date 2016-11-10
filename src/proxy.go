package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/shuaiming/mung/middlewares"
)

// Proxy wapper ReverseProxy
type Proxy struct {
	proxy *httputil.ReverseProxy
}

// NewProxy New proxy
func NewProxy(grafana string) *Proxy {
	backend, err := url.Parse(fmt.Sprintf("%s", grafana))
	if err != nil {
		log.Fatal(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(backend)
	return &Proxy{proxy: proxy}
}

func (p *Proxy) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	sess := middlewares.GetSession(r)
	openid := middlewares.GetOpenIDUser(r)

	// overwrite grafana's login
	if r.URL.Path == "/login" {
		http.Redirect(rw, r, "/openid/login", http.StatusFound)
		return
	}

	// overwrite grafana's logout
	if r.URL.Path == "/logout" {
		delete(sess.Values, middlewares.OpenIDContextKey)
		if err := sess.Save(r, rw); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}

		p.proxy.ServeHTTP(rw, r)
		return
	}

	// redirect to login url if openid not login
	email, ok := openid["sreg.email"]
	if !ok {
		http.Redirect(rw, r, "/openid/login", http.StatusFound)
		return
	}

	// overwirte X-WEBAUTH-USER with openid email name
	r.Header.Set("X-WEBAUTH-USER", email)
	p.proxy.ServeHTTP(rw, r)
}
