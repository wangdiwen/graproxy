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

var flagGrafana = flag.String("grafana", "localhost", "grafana host and port")
var flagOpenIDEndpoint = flag.String("endpoint", "", "openid endpoint")
var flagListen = flag.String("l", ":8080", "proxy server listen address")
var flagServer = flag.String("n", "localhost", "proxy domain name openid will return to")

func main() {

	flag.Parse()

	al := middlewares.NewAccessLog(os.Stdout)

	// Security bytes lenght must be 32 or 64
	// Cookie has a size limit of 4096?
	store := sessions.NewCookieStore([]byte("KwGPH3acQfOXscHHMMxCRb0HqO01+GTh"))
	store.MaxAge(172800)

	sessionMgr := middlewares.NewSessions(store)
	openid := middlewares.NewOpenID(
		"https://login.netease.com/openid",
		fmt.Sprintf("http://%s:%s", *flagServer, strings.Split(*flagListen, ":")[1]),
		"/openid",
	)

	proxy := middlewares.NewHandler(NewProxy(*flagGrafana))

	// use of middleware, please tack care the order!
	app := mung.New()
	app.Use(al)
	app.Use(sessionMgr)
	app.Use(openid)
	app.Use(proxy)

	// start the app
	app.Run(*flagListen)
	// app.RunTLS(*flagListen, "ssl/server.crt", "ssl/server.key")
}
