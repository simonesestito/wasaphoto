package database

import (
	"embed"
	"github.com/sirupsen/logrus"
	"io"
	"sort"
	"strconv"
	"strings"
)

//go:embed *
var migrations embed.FS

type migration struct {
	Version int
	Name    string
	File    io.Reader
}

func (db appSqlDatabase) runMigrations(logger logrus.FieldLogger) error {
	currentVersion, err := db.Version()
	if err != nil {
		return err
	}
	logger.Debugf("Current database schema version: %d", currentVersion)

	newMigrations, err := listMigrationsAfter(currentVersion, logger)
	if err != nil {
		return err
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
			return err
		}

		// Execute it in a transaction
		tx, err := db.BeginTx()

		if err != nil {
			logger.WithError(err).Errorln("Error creating a transaction")
			return err
		}

		_, err = tx.Exec(string(sqlScript))
		if err != nil {
			logger.WithError(err).Errorf("Error running migration file %s", file.Name)
			_ = tx.Rollback()
			return err
		}

		err = db.setVersion(file.Version, tx)
		if err != nil {
			logger.WithError(err).Errorf("Error setting schema version %d", file.Version)
			return err
		}

		err = tx.Commit()
		if err != nil {
			logger.WithError(err).Errorf("Error committing schema version %d", file.Version)
			return err
		}

		logger.Debugf("Successfully upgraded to schema %s", file.Name)
	}

	return nil
}

// listMigrationsAfter is needed to get all the database schema migrations
// available on disk, and checks if they look like schema migration files.
// The returned migrations list contains only the migrations after the currentVersion
func listMigrationsAfter(currentVersion int, logger logrus.FieldLogger) ([]migration, error) {
	allFiles, err := migrations.ReadDir(".")
	if err != nil {
		return nil, err
	}

	migrationFiles := make([]migration, 0, len(allFiles))

	// Find files that are named like a schema migration
	for _, entry := range allFiles {
		if !entry.Type().IsRegular() {
			logger.Debugf("Skipping non regular file in database migrations directory: %s", entry.Name())
			continue
		}

		if strings.HasSuffix(entry.Name(), ".sql") {
			sqlName := strings.TrimSuffix(entry.Name(), ".sql")
			schemaVersion, err := strconv.Atoi(sqlName)
			if err != nil {
				logger.WithError(err).Warn("SQL file found not named with a version number:", entry.Name())
			} else if schemaVersion > currentVersion {
				file, err := migrations.Open(entry.Name())
				if err != nil {
					// Fatal unexpected error.
					return nil, err
				}

				migrationFiles = append(migrationFiles, migration{
					Version: schemaVersion,
					Name:    entry.Name(),
					File:    file,
				})
			}
		} else if !strings.HasSuffix(entry.Name(), ".go") {
			logger.Warnf("Ignored file '%s' in database migrations directory", entry.Name())
		}
	}

	return migrationFiles, nil
}
