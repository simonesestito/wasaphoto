package database

import (
	"database/sql"
	"strings"
)

// Exec executes a database query with the given arguments.
// In case the query returns an error different from sql.ErrNoResult,
// it'll be returned.
func (db appSqlDatabase) Exec(query string, args ...any) error {
	_, err := db.ExecRows(query, args...)
	return err
}

// ExecRows executes a query as Exec, but it also reports
// how many rows have been affected.
func (db appSqlDatabase) ExecRows(query string, args ...any) (int64, error) {
	result, err := db.DB.Exec(query, args...)

	if err == sql.ErrNoRows {
		return 0, nil
	} else if err != nil && strings.HasPrefix(err.Error(), "UNIQUE") {
		return 0, ErrDuplicated
	} else if err != nil {
		return 0, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rows, nil
}

type Closable interface {
	Close() error
}

// tryClosingRows without error handling (it tries).
func tryClosingRows(rows Closable) {
	_ = rows.Close()
}
