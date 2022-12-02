package ioc

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/simonesestito/wasaphoto/service/database"
	"github.com/simonesestito/wasaphoto/service/storage"
	"github.com/simonesestito/wasaphoto/service/timeprovider"
	"github.com/sirupsen/logrus"
)

//
// This file (and this package in general)
// implements the Inversion Of Control principle
// using a Dependency Container,
// in order to create other app components still
// using Dependency Injection.
//
// More info in service/README.md file
//

type Container struct {
	// External dependencies here
	forcedTime timeprovider.TimeProvider
	logger     *logrus.Logger
	database   database.AppDatabase

	// instances collects singleton instances for those
	// dependencies which need to be a shared instance.
	// Not meant to be goroutine safe.
	instances map[string]any
}

func New(timeProvider timeprovider.TimeProvider, logger *logrus.Logger, rawDatabase *sqlx.DB) (Container, error) {
	if logger == nil {
		return Container{}, errors.New("logger is required")
	}

	if rawDatabase == nil {
		return Container{}, errors.New("rawDatabase is required")
	}

	appDatabase, err := database.New(rawDatabase, logger)
	if err != nil {
		return Container{}, errors.New(fmt.Sprintf("error wrapping database: %s", err.Error()))
	}

	return Container{
		forcedTime: timeProvider,
		logger:     logger,
		database:   appDatabase,
		instances:  make(map[string]any),
	}, nil
}

func (ioc *Container) CreateTimeProvider() timeprovider.TimeProvider {
	if ioc.forcedTime != nil {
		return ioc.forcedTime
	}

	return timeprovider.RealTimeProvider{}
}

func (ioc *Container) CreateStorage() storage.Storage {
	const key = "storage.Storage"
	if previousInstance, ok := ioc.instances[key]; ok {
		return previousInstance.(storage.Storage)
	}

	// Create a new storage.Storage
	newInstance := storage.FilesystemStorage{}
	ioc.instances[key] = &newInstance
	return &newInstance
}
