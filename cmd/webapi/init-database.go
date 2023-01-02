package main

import (
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

func initDatabase(cfg webAPIConfiguration, logger *logrus.Logger) (*sqlx.DB, func()) {
	var logDbFilename string
	if strings.HasPrefix(cfg.DB.Filename, "/") {
		// Absolute path
		logDbFilename = cfg.DB.Filename
	} else {
		// Relative path
		wd, _ := os.Getwd()
		logDbFilename = wd + "/" + cfg.DB.Filename
	}

	logger.Infoln("initializing database support", logDbFilename)

	dbConn, err := sqlx.Open("sqlite3", cfg.DB.Filename+"?_foreign_keys=on")
	if err != nil {
		logger.WithError(err).Fatalln("error opening SQLite DB")
	}

	logger.Debugln("trying to ping the database")
	if err = dbConn.Ping(); err != nil {
		_ = dbConn.Close()
		logger.WithError(err).Fatalln("error pinging the DB, make sure the file's destination folder is writable")
	}

	mustWriteDatabase(dbConn, logger)

	return dbConn, func() {
		logger.Debug("database stopping")
		_ = dbConn.Close()
	}
}

func mustWriteDatabase(db *sqlx.DB, logger logrus.FieldLogger) {
	logger.Debugln("trying to write on the database")

	var (
		tx  *sqlx.Tx
		err error
	)

	// Function to log the help message and close the database connection
	showFatalMessage := func(err error) {
		if tx != nil {
			_ = tx.Rollback()
		}

		logger.Errorln("Make sure to bind the entire directory that contains the database file, instead of just binding the DB file itself")
		_ = db.Close()
		logger.WithError(err).Fatal()
	}

	// In a transaction, create a new test table, write to it, and finally drop it
	tx, err = db.Beginx()
	if err != nil {
		showFatalMessage(err)
	}

	_, err = tx.Exec("CREATE TABLE WriteTestTable (id INT NOT NULL PRIMARY KEY)")
	if err != nil {
		showFatalMessage(err)
	}

	_, err = tx.Exec("INSERT INTO WriteTestTable (id) VALUES(1)")
	if err != nil {
		showFatalMessage(err)
	}

	_, err = tx.Exec("DROP TABLE WriteTestTable")
	if err != nil {
		showFatalMessage(err)
	}

	err = tx.Commit()
	if err != nil {
		showFatalMessage(err)
	}
}
