package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/context"
	"github.com/shuaiming/openid"
)

const (
	// OpenIDContextKey context store key for OpenID
	OpenIDContextKey string = "contextopenid"
	// OpenIDURLReturnKey URL variable key for redirection after verified
	OpenIDURLReturnKey string = "openidreturn"
)

// OpenID OpenID
type OpenID struct {
	opEndpoint string
	urlPrefix  string
	realm      string
	openid     *openid.OpenID
}

// NewOpenID make new OpenID
func NewOpenID(endpoint, realm, prefix string) *OpenID {

	return &OpenID{
		opEndpoint: endpoint,
		urlPrefix:  prefix,
		realm:      realm,
		openid:     openid.New(realm),
	}
}

// ServeHTTP make mung middleware
func (o *OpenID) ServeHTTP(
	rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	// verifyURL is url for OpenID Server back redirecetion
	verifyURL := fmt.Sprintf("%s/verify", o.urlPrefix)

	session := GetSession(r)
	if session == nil {
		fmt.Fprintln(os.Stderr, "can not enable openid without session")
		next(rw, r)
		return
	}

	// Put user to requesting context for later usage
	if user, ok := session.Values[OpenIDContextKey]; ok {
		context.Set(r, OpenIDContextKey, user)
		defer context.Delete(r, OpenIDContextKey)
	}

	// muopenid only http Method GET and HEAD supported
	if r.Method != "GET" && r.Method != "HEAD" {
		next(rw, r)
		return
	}

	if !strings.HasPrefix(r.URL.Path, o.urlPrefix) {
		next(rw, r)
		return
	}

	switch r.URL.Path {

	case fmt.Sprintf("%s/login", o.urlPrefix):

		// returnRUL is url will return back to after login finished
		// We will store it to "session" for later usage
		if returnURL := r.URL.Query().Get(OpenIDURLReturnKey); returnURL != "" {
			session.Values[OpenIDURLReturnKey] = returnURL
			if err := session.Save(r, rw); err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
			}
		}

		// redirectURL is url redirect to OpenID Server
		if redirectURL, err := o.openid.CheckIDSetup(
			o.opEndpoint, verifyURL); err == nil {
			http.Redirect(rw, r, redirectURL, http.StatusFound)
		} else {
			fmt.Fprintln(os.Stderr, err.Error())
		}

	case verifyURL:

		user, err := o.openid.IDRes(r)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			break
		}

		session.Values[OpenIDContextKey] = user
		if err := session.Save(r, rw); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}

		if returnURL, ok := session.Values[OpenIDURLReturnKey]; ok {
			http.Redirect(rw, r, returnURL.(string), http.StatusFound)
		} else {
			http.Redirect(rw, r, o.realm, http.StatusFound)
		}

	default:
		next(rw, r)
	}
}

// GetOpenIDUser return an map of openid user info
func GetOpenIDUser(r *http.Request) map[string]string {
	id := context.Get(r, OpenIDContextKey)
	if id != nil {
		return id.(map[string]string)
	}
	return nil
}
