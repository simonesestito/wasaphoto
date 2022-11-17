package database

import (
	"errors"
	"github.com/jmoiron/sqlx"
)

// transactionsDb is a database interface with the ability to create new transactions
type transactionsDb interface {
	databaseInterface
	Beginx() (*sqlx.Tx, error)
}

// A Transaction is a database interface which can be committed or rollback.
type Transaction interface {
	databaseInterface
	Commit() error
	Rollback() error
}

func (db appSqlDatabase) BeginTx() (Transaction, error) {
	transactDb, ok := db.DB.(transactionsDb)
	if !ok {
		return nil, errors.New("cannot create a transaction on this database implementation")
	}

	return transactDb.Beginx()
}
