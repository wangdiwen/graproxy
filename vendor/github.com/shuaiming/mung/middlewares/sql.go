package middlewares

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/context"
)

const (
	// SQLContextKey context store key for SQL
	SQLContextKey string = "contextsql"
)

type stmt struct {
	stmt *sql.Stmt
}

func (s stmt) Exec(args ...interface{}) (sql.Result, error) {
	return s.stmt.Exec(args...)
}

func (s stmt) Query(args ...interface{}) (*sql.Rows, error) {
	return s.stmt.Query(args...)
}

func (s stmt) QueryRow(args ...interface{}) *sql.Row {
	return s.stmt.QueryRow(args...)
}

// SQL sql like database surpporting
type SQL struct {
	db    *sql.DB
	stmts map[string]Stmt
}

// Stmt wrap up the sql.Stmt without method close()
type Stmt interface {
	Exec(...interface{}) (sql.Result, error)
	Query(...interface{}) (*sql.Rows, error)
	QueryRow(...interface{}) *sql.Row
}

// NewSQL make new SQL
func NewSQL(driverName, dataSourceName string) *SQL {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	stmts := map[string]Stmt{}
	return &SQL{db: db, stmts: stmts}
}

// ServeHTTP make mung middleware
func (msql *SQL) ServeHTTP(
	rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	context.Set(r, SQLContextKey, msql.stmts)
	defer context.Delete(r, SQLContextKey)
	next(rw, r)
}

// InitDB initialize database structs
func (msql *SQL) InitDB(dbInitSQL string) error {
	_, err := msql.db.Exec(dbInitSQL)
	return err
}

// StmtPrepare push Stmt to MuSQL.stmts for later usage
func (msql *SQL) StmtPrepare(name, sql string) error {
	s, err := msql.db.Prepare(sql)
	if err != nil {
		return err
	}
	msql.stmts[name] = &stmt{stmt: s}
	return nil
}

// GetStmts return an map of musql.Stmt
func GetStmts(r *http.Request) map[string]Stmt {
	l := context.Get(r, SQLContextKey)
	if l != nil {
		return l.(map[string]Stmt)
	}
	return nil
}
