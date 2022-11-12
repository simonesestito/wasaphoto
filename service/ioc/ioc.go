package ioc

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/simonesestito/wasaphoto/service/database"
	"github.com/simonesestito/wasaphoto/service/timeprovider"
	"github.com/sirupsen/logrus"
)

type Container struct {
	// External dependencies here
	forcedTime timeprovider.TimeProvider
	logger     *logrus.Logger
	database   database.AppDatabase
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
	}, nil
}

func (ioc *Container) CreateTimeProvider() timeprovider.TimeProvider {
	if ioc.forcedTime != nil {
		return ioc.forcedTime
	}

	return timeprovider.RealTimeProvider{}
}
