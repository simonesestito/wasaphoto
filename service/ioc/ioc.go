package ioc

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/simonesestito/wasaphoto/service/api"
	"github.com/simonesestito/wasaphoto/service/api/route"
	"github.com/simonesestito/wasaphoto/service/database"
	"github.com/simonesestito/wasaphoto/service/features/auth"
	"github.com/simonesestito/wasaphoto/service/features/user"
	"github.com/simonesestito/wasaphoto/service/timeprovider"
	"github.com/sirupsen/logrus"
)

type Container struct {
	// External dependencies here
	forcedTime timeprovider.TimeProvider
	logger     *logrus.Logger
	database   database.AppDatabase
}

func New(timeProvider timeprovider.TimeProvider, logger *logrus.Logger, rawDatabase *sql.DB) (Container, error) {
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

func (ioc *Container) CreateRouter() (api.Router, error) {
	// List endpoints to register
	controllers := []route.Controller{
		ioc.CreateUserController(),
		ioc.CreateLoginController(),
	}

	// List middlewares
	middlewares := []route.Middleware{
		api.LimitBodySize(1024 * 1024),
	}

	// Create router
	router := api.NewRouter(ioc.CreateAuthMiddleware(), middlewares, ioc.logger)

	// Register routes
	for _, controller := range controllers {
		for _, routeInfo := range controller.ListRoutes() {
			if err := router.Register(routeInfo); err != nil {
				return router, err
			}
		}
	}

	return router, nil
}

func (ioc *Container) CreateAuthService() auth.LoginService {
	return auth.UserIdLoginService{
		UserDao: ioc.CreateUserDao(),
	}
}

func (ioc *Container) CreateLoginController() auth.LoginController {
	return auth.LoginController{
		AuthService: ioc.CreateAuthService(),
	}
}

func (ioc *Container) CreateUserController() user.Controller {
	return user.Controller{}
}

func (ioc *Container) CreateAuthMiddleware() route.AuthMiddleware {
	return auth.Middleware{LoginService: ioc.CreateAuthService()}
}

func (ioc *Container) CreateUserDao() user.Dao {
	return user.DbDao{DB: ioc.database}
}
