package database

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sqlx.DB, logger logrus.FieldLogger) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	appDatabase := appSqlDatabase{DB: db, PingDB: db}
	if err := appDatabase.runMigrations(logger); err != nil {
		return nil, err
	}

	return appDatabase, nil
}

type appSqlDatabase struct {
	DB     databaseInterface // This may be a *sqlx.DB or *sqlx.Tx
	PingDB pingableDatabase
}

type databaseInterface interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	Queryx(query string, args ...any) (*sqlx.Rows, error)
}

type pingableDatabase interface {
	Ping() error
}

func (db appSqlDatabase) Ping() error {
	return db.PingDB.Ping()
}

func (db appSqlDatabase) Version() (int, error) {
	if err := db.Ping(); err != nil {
		return 0, err
	}

	_, err := db.DB.Exec("CREATE TABLE IF NOT EXISTS SchemaVersion ( version INT NOT NULL PRIMARY KEY )")
	if err != nil {
		return 0, err
	}

	rows, err := db.DB.Query("SELECT version FROM SchemaVersion")
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return 0, nil
	case err == nil:
		defer tryClosingRows(rows)

		if nextRow := rows.Next(); !nextRow {
			return 0, nil
		}

		var version int
		if err := rows.Scan(&version); err != nil {
			return 0, err
		}
		return version, nil
	default:
		return 0, err
	}
}

// setVersion sets the database version, without checking it.
// It does not perform any schema upgrade/migration.
func (db appSqlDatabase) setVersion(version int, tx transaction) error {
	if err := db.Ping(); err != nil {
		return err
	}

	if _, err := tx.Exec("DELETE FROM SchemaVersion WHERE TRUE"); err != nil {
		return err
	}

	if _, err := tx.Exec("INSERT INTO SchemaVersion (version) VALUES (?)", version); err != nil {
		return err
	}

	return nil
}
