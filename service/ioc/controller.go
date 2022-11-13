package ioc

import (
	"github.com/simonesestito/wasaphoto/service/api/route"
	"github.com/simonesestito/wasaphoto/service/features/auth"
	"github.com/simonesestito/wasaphoto/service/features/follow"
	"github.com/simonesestito/wasaphoto/service/features/photo"
	"github.com/simonesestito/wasaphoto/service/features/user"
)

func (ioc *Container) CreateLoginController() auth.LoginController {
	return auth.LoginController{
		AuthService: ioc.CreateAuthService(),
	}
}

func (ioc *Container) CreateUserController() user.Controller {
	return user.Controller{
		Service: ioc.CreateUserService(),
	}
}

func (ioc *Container) CreateBanController() user.BanController {
	return user.BanController{
		Service: ioc.CreateBanService(),
	}
}

func (ioc *Container) CreateFollowController() follow.Controller {
	return follow.Controller{Service: ioc.CreateFollowService()}
}

func (ioc *Container) CreatePhotoController() photo.Controller {
	return photo.Controller{Service: ioc.CreatePhotoService()}
}

func (ioc *Container) CreateAuthMiddleware() route.AuthMiddleware {
	return auth.Middleware{LoginService: ioc.CreateAuthService()}
}
