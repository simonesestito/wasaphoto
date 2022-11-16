package database

import (
	"embed"
	"github.com/sirupsen/logrus"
	"io"
	"strconv"
	"strings"
)

//go:embed *
var migrations embed.FS

type Migration struct {
	Version int
	Name    string
	File    io.Reader
}

// ListMigrationsAfter is needed to get all the database schema migrations
// available on disk, and checks if they look like schema migration files.
// The returned migrations list contains only the migrations after the currentVersion
func ListMigrationsAfter(currentVersion int, logger logrus.FieldLogger) ([]Migration, error) {
	allFiles, err := migrations.ReadDir(".")
	if err != nil {
		return nil, err
	}

	migrationFiles := make([]Migration, 0, len(allFiles))

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

				migrationFiles = append(migrationFiles, Migration{
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
