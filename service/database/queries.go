package database

import (
	"database/sql"
)

// Exec executes a database query with the given arguments.
// In case the query returns an error different from sql.ErrNoResult,
// it'll be returned.
func (db appSqlDatabase) Exec(query string, args ...any) error {
	_, err := db.DB.Exec(query, args...)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	return nil
}

type Closable interface {
	Close() error
}

// tryClosingRows without error handling (it tries).
func tryClosingRows(rows Closable) {
	_ = rows.Close()
}
