package database

import (
	"database/sql"
	"errors"
)

// transactionsDb is a database interface with the ability to create new transactions
type transactionsDb interface {
	Begin() (*sql.Tx, error)
}

// A transaction is a database interface which can be committed or rollback.
type transaction interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	Commit() error
	Rollback() error
}

func (db appSqlDatabase) BeginTx() (transaction, error) {
	transactDb, ok := db.DB.(transactionsDb)
	if !ok {
		return nil, errors.New("cannot create a transaction on this database implementation")
	}

	return transactDb.Begin()
}
