package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/shuaiming/mung"
	"github.com/shuaiming/mung/middlewares"
	"github.com/shuaiming/mung/middlewares/sessions"
)

var flagGrafana = flag.String("grafana", "http://127.0.0.1:3000", "grafana host and port")
var flagOpenIDEndpoint = flag.String("endpoint", "", "openid endpoint")
var flagListen = flag.String("l", ":8080", "proxy server listen address")
var flagServer = flag.String("n", "localhost", "proxy domain name openid will return to")

var flagSSL = flag.Bool("ssl", false, "enable https")
var flagSSLCertFile = flag.String("cert", "ssl/server.crt", "ssl server certificate file")
var flagSSLKeyFile = flag.String("key", "ssl/server.key", "ssl server key file")

func main() {

	flag.Parse()

	proto := "http"
	if *flagSSL {
		proto = "https"
	}

	al := middlewares.NewAccessLog(os.Stdout)

	// Security bytes lenght must be 32 or 64
	// Cookie has a size limit of 4096?
	store := sessions.NewCookieStore([]byte("KwGPH3acQfOXscHHMMxCRb0HqO01+GTh"))
	store.MaxAge(172800)

	sessionMgr := middlewares.NewSessions(store)
	openid := middlewares.NewOpenID(
		*flagOpenIDEndpoint,
		fmt.Sprintf("%s://%s:%s", proto, *flagServer, strings.Split(*flagListen, ":")[1]),
		"/openid",
	)

	proxy := middlewares.NewHandler(NewProxy(*flagGrafana))

	// use of middleware, please tack care the order!
	app := mung.New()
	app.Use(al)
	app.Use(sessionMgr)
	app.Use(openid)
	app.Use(proxy)

	// start the server
	if *flagSSL {
		app.RunTLS(*flagListen, *flagSSLCertFile, *flagSSLKeyFile)
	} else {
		app.Run(*flagListen)
	}
}
