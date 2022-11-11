package database

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/simonesestito/wasaphoto/database"
	"github.com/sirupsen/logrus"
	"io"
	"sort"
)

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sqlx.DB, logger logrus.FieldLogger) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	appDatabase := appSqlDatabase{DB: db}

	currentVersion, err := appDatabase.Version()
	if err != nil {
		return nil, err
	}
	logger.Debugf("Current database schema version: %d", currentVersion)

	newMigrations, err := database.ListMigrationsAfter(currentVersion, logger)
	if err != nil {
		return nil, err
	}

	if len(newMigrations) == 0 {
		logger.Debug("No db schema migrations are needed.")
	}

	// Apply migrations in ascending order
	sort.Slice(newMigrations, func(i, j int) bool {
		return newMigrations[i].Version < newMigrations[j].Version
	})

	for _, file := range newMigrations {
		logger.Infof("Migrating to schema %s", file.Name)

		// Read SQL file
		sqlScript, err := io.ReadAll(file.File)
		if err != nil {
			logger.WithError(err).Errorf("Error reading migration file %s", file.Name)
			return nil, err
		}

		// Execute it in a transaction
		tx, err := db.Begin()

		if err != nil {
			logger.WithError(err).Errorln("Error creating a transaction")
			return nil, err
		}

		_, err = tx.Exec(string(sqlScript))
		if err != nil {
			logger.WithError(err).Errorf("Error running migration file %s", file.Name)
			_ = tx.Rollback()
			return nil, err
		}

		err = appDatabase.SetVersion(file.Version, tx)
		if err != nil {
			logger.WithError(err).Errorf("Error setting schema version %d", file.Version)
			return nil, err
		}

		err = tx.Commit()
		if err != nil {
			logger.WithError(err).Errorf("Error committing schema version %d", file.Version)
			return nil, err
		}

		logger.Debugf("Successfully upgraded to schema %s", file.Name)
	}

	return appDatabase, nil
}

type appSqlDatabase struct {
	DB *sqlx.DB
}

func (db appSqlDatabase) Ping() error {
	return db.DB.Ping()
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
	nextRow := rows.Next()
	switch {
	case err == sql.ErrNoRows || !nextRow:
		return 0, nil
	case err == nil && nextRow:
		var version int
		if err := rows.Scan(&version); err != nil {
			return 0, err
		}
		return version, nil
	default:
		return 0, err
	}
}

// SetVersion sets the database version, without checking it.
// It does not perform any schema upgrade/migration.
func (db appSqlDatabase) SetVersion(version int, tx *sql.Tx) error {
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
