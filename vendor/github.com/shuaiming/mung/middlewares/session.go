package middlewares

import (
	"encoding/gob"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/context"
	"github.com/shuaiming/mung/middlewares/sessions"
)

const (
	// SessionContextKey context store key for *sessions.Session
	SessionContextKey string = "contextsession"
)

// Sessions manager sessions
type Sessions struct {
	store sessions.Store
}

// NewSessions make new Sessions
func NewSessions(store sessions.Store) *Sessions {
	// securecookie use gob to [de]serialize
	// regiester to gob before you can sessions.CookieStore save to.
	gob.Register(map[string]string{})

	return &Sessions{store: store}
}

// ServeHTTP make mung middleware
func (ss *Sessions) ServeHTTP(
	rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	if s, err := ss.store.Get(r, SessionContextKey); err == nil {

		context.Set(r, SessionContextKey, s)
		defer context.Delete(r, SessionContextKey)
	} else {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	next(rw, r)

}

// GetSession return Session
func GetSession(r *http.Request) *sessions.Session {
	s := context.Get(r, SessionContextKey)
	if s != nil {
		return s.(*sessions.Session)
	}
	return nil
}
