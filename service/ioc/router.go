package ioc

import (
	"github.com/simonesestito/wasaphoto/service/api"
	"github.com/simonesestito/wasaphoto/service/api/route"
)

func (ioc *Container) CreateControllers() []route.Controller {
	return []route.Controller{
		ioc.createUserController(),
		ioc.createLoginController(),
		ioc.createBanController(),
		ioc.createFollowController(),
		ioc.createPhotoController(),
		ioc.createLikesController(),
		ioc.createCommentsController(),
		ioc.createStreamController(),
	}
}

func (ioc *Container) CreateMiddlewares() []route.Middleware {
	return []route.Middleware{
		api.LimitBodySize(20 * 1024 * 1024),
	}
}
