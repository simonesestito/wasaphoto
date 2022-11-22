package ioc

import (
	"github.com/simonesestito/wasaphoto/service/api"
	"github.com/simonesestito/wasaphoto/service/api/route"
)

func (ioc *Container) CreateControllers() []route.Controller {
	return []route.Controller{
		ioc.CreateUserController(),
		ioc.CreateLoginController(),
		ioc.CreateBanController(),
		ioc.CreateFollowController(),
		ioc.CreatePhotoController(),
		ioc.CreateLikesController(),
		ioc.CreateCommentsController(),
		ioc.CreateStreamController(),
	}
}

func (ioc *Container) CreateMiddlewares() []route.Middleware {
	return []route.Middleware{
		api.LimitBodySize(20 * 1024 * 1024),
	}
}
