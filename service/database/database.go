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
	QueryStructRows(entityStruct any, query string, args ...any) (StructRows, error)
	Exec(query string, args ...any) error
	ExecRows(query string, args ...any) (int64, error)
	BeginTx() (transaction, error)
}

var ErrNoResult = sql.ErrNoRows
var ErrDuplicated = errors.New("operation duplicated")
var ErrForeignKey = errors.New("foreign key failed")

const MaxPageItems = 2 // According to the docs
