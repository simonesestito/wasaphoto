package ioc

import (
	"github.com/simonesestito/wasaphoto/service/api"
	"github.com/simonesestito/wasaphoto/service/api/route"
	"github.com/simonesestito/wasaphoto/service/features/auth"
	"github.com/simonesestito/wasaphoto/service/features/user"
	"github.com/simonesestito/wasaphoto/service/timeprovider"
	"github.com/sirupsen/logrus"
)

type Container struct {
	// External dependencies here
	ForcedTime timeprovider.TimeProvider
	Logger     *logrus.Logger
}

func (ioc *Container) CreateTimeProvider() timeprovider.TimeProvider {
	if ioc.ForcedTime != nil {
		return ioc.ForcedTime
	}

	return timeprovider.RealTimeProvider{}
}

func (ioc *Container) CreateRouter() (api.Router, error) {
	// List endpoints to register
	controllers := []route.Controller{
		//ioc.CreateUserController(),
		ioc.CreateLoginController(),
	}

	// List middlewares
	middlewares := []route.Middleware{
		api.LimitBodySize(1024 * 1024),
	}

	// Create router
	router := api.NewRouter(
		ioc.CreateAuthMiddleware(),
		middlewares,
		ioc.Logger,
	)

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
	return auth.UserIdLoginService{}
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
	return auth.Middleware{ioc.CreateAuthService()}
}
