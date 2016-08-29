package middlewares

import (
	"fmt"
	"log"
	"log/syslog"
	"net/http"

	"github.com/gorilla/context"
)

const (
	// SyslogContextKey context store key for SQL
	SyslogContextKey string = "contextsyslog"
)

// Syslog write to syslog
type Syslog struct {
	logger *syslog.Writer
}

// NewSyslog make new Syslog
func NewSyslog(priority syslog.Priority, tag string) *Syslog {
	logger, err := syslog.New(priority, tag)
	if err != nil {
		log.Fatal(err)
	}
	logger.Warning(fmt.Sprintf("syslog enabled %d %s", priority, tag))
	return &Syslog{logger: logger}
}

// ServeHTTP make mung middleware
func (sl *Syslog) ServeHTTP(
	rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	context.Set(r, SyslogContextKey, sl.logger)
	defer context.Delete(r, SyslogContextKey)
	next(rw, r)

}

// GetSyslogWriter return syslog.Writer
func GetSyslogWriter(r *http.Request) *syslog.Writer {
	l := context.Get(r, SyslogContextKey)
	if l != nil {
		return l.(*syslog.Writer)
	}
	return nil
}
