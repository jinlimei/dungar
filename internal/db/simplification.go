package db

import (
	"database/sql"
	"errors"
	"log"

	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
)

// ErrNoDBConnection returns when zero DB connections have been made.
var ErrNoDBConnection = errors.New("no connection to database")

// QueryAble is our interface for working with either a db connection or a transaction
type QueryAble interface {
	QueryRow(query string, data ...interface{}) *sql.Row
	Query(query string, data ...interface{}) (*sql.Rows, error)
	Exec(query string, data ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
}

func isQueryAbleNil (inp QueryAble) bool {
	switch it := inp.(type) {
	case *sql.DB:
		return it == nil
	case *sql.Tx:
		return it == nil
	case nil:
		return true
	default:
		log.Printf("Unable to determine type for QueryAble: %+v\n", inp)
		return true
	}
}

func useGlobalIfNil(incoming QueryAble) QueryAble {
	if isQueryAbleNil(incoming) {
		return GetDatabase()
	}

	return incoming
}

// ConQueryRow is QueryRow ran on the standard connection
func ConQueryRow(query string, data ...interface{}) *sql.Row {
	return QueryRow(GetDatabase(), query, data...)
}

// ConQuery is Query ran on the standard connection
func ConQuery(query string, data ...interface{}) (*sql.Rows, error) {
	return Query(GetDatabase(), query, data...)
}

// ConMustQuery is ConQuery which *must* execute (error will halt)
func ConMustQuery(query string, data ...interface{}) *sql.Rows {
	return MustQuery(GetDatabase(), query, data...)
}

// QueryRow runs QueryRow on intf or global
func QueryRow(intf QueryAble, query string, data ...interface{}) *sql.Row {
	iv := useGlobalIfNil(intf)

	if isQueryAbleNil(iv) {
		return nil
	}

	return iv.QueryRow(query, data...)
}

// Query runs Query on intf or global
func Query(intf QueryAble, query string, data ...interface{}) (*sql.Rows, error) {
	iv := useGlobalIfNil(intf)

	if isQueryAbleNil(iv) {
		return nil, ErrNoDBConnection
	}

	return iv.Query(query, data...)
}

// MustQuery runs Query on intf or global, and will halt on error
func MustQuery(intf QueryAble, query string, data ...interface{}) *sql.Rows {
	iv := useGlobalIfNil(intf)

	if isQueryAbleNil(iv) {
		handleError("MustQuery", ErrNoDBConnection)
		return nil
	}

	rows, err := iv.Query(query, data...)

	handleError("MustQuery", err)

	return rows
}

// ConExec is like Exec but operates on the singleton DB
func ConExec(query string, data ...interface{}) (sql.Result, error) {
	return Exec(GetDatabase(), query, data...)
}

// Exec executes on intf or global if nil
func Exec(intf QueryAble, query string, data ...interface{}) (sql.Result, error) {
	iv := useGlobalIfNil(intf)

	if isQueryAbleNil(iv) {
		return nil, ErrNoDBConnection
	}

	return iv.Exec(query, data...)
}

// ConMustExec runs MustExec on the default connection, and will halt program if error
func ConMustExec(query string, data ...interface{}) sql.Result {
	return MustExec(GetDatabase(), query, data...)
}

// MustExec runs on intf or global if nil, and will halt on execute
func MustExec(intf QueryAble, query string, data ...interface{}) sql.Result {
	iv := useGlobalIfNil(intf)

	if isQueryAbleNil(iv) {
		handleError("MustExec", ErrNoDBConnection)
		return nil
	}

	res, err := iv.Exec(query, data...)

	handleError("MustExec", err)

	return res
}

func handleError(loc string, err error) {
	utils.HaltingError("db "+loc, err)
}
