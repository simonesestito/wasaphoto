package database

import (
	"database/sql"
	"errors"
)

// AppDatabase is an abstraction over platform *sql.DB
type AppDatabase interface {
	Ping() error
	Version() (int, error)
	QueryStructRow(destPointer any, query string, args ...any) error
	Exec(query string, args ...any) error
	ExecRows(query string, args ...any) (int64, error)
}

var ErrNoResult = sql.ErrNoRows
var ErrDuplicated = errors.New("operation duplicated")
