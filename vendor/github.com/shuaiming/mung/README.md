# mung


## description

mung is tiny web framwork with middlewares supported. the ideas mainly come from expressjs and negroni. It is much more smaller and simpler also very easy to use and extend. Any struct/type can be made to a middlewarea only by adding a function ServeHTTP(rw, r, next), then add the middleware to the mung's call statck.


## usage exapmple

	middleware1 := ...  // only ServeHTTP(rw, r, next) needed
	middleware2 := ...
	middleware3 := ...
	app := mung.New()
	app.Use(middleware1)
	app.Use(middleware2)
	app.Use(middleware3)
	app.Run(":8888")

Mung.handle is the entry of middlewares stack.


## about middlewares

the order of middleware's ServeHTTP() execution is:

	middleware1			// before next() in middleware1's ServeHTTP()
	  middleware2		// before next() in middleware2's ServeHTTP()
	    middleware3		// middleware3's ServeHTTP()
	  middleware2		// after next() in middleware2's ServeHTTP()
	middleware1			// after next() in middleware1's ServeHTTP()

Any functionality should be made to a reusable middleware inluding staics serving , router, sessions, logger, database accessing, openid/oauth and so on.
