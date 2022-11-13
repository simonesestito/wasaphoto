package ioc

import (
	"github.com/simonesestito/wasaphoto/service/api"
	"github.com/simonesestito/wasaphoto/service/api/route"
)

func (ioc *Container) CreateRouter() (api.Router, error) {
	// List endpoints to register
	controllers := []route.Controller{
		ioc.CreateUserController(),
		ioc.CreateLoginController(),
		ioc.CreateBanController(),
		ioc.CreateFollowController(),
		ioc.CreatePhotoController(),
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
