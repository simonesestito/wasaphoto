package ioc

import (
	"github.com/simonesestito/wasaphoto/service/api/route"
	"github.com/simonesestito/wasaphoto/service/features/auth"
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

func (ioc *Container) CreateAuthMiddleware() route.AuthMiddleware {
	return auth.Middleware{LoginService: ioc.CreateAuthService()}
}
