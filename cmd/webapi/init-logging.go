package main

import (
	"github.com/sirupsen/logrus"
	"log"
	"os"
)

func initLogging(cfg webAPIConfiguration) *logrus.Logger {
	logger := logrus.New()

	// Set output
	if cfg.Log.FileName == "-" {
		logger.SetOutput(os.Stdout)
	} else {
		logFile, err := os.Create(cfg.Log.FileName)
		if err != nil {
			log.Fatalf("Error creating log file: %s", err.Error())
		}
		logger.SetOutput(logFile)
	}

	// Set level
	if cfg.Log.Debug {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}

	return logger
}
