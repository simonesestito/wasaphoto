package ioc

import (
	"github.com/simonesestito/wasaphoto/service/api/route"
	"github.com/simonesestito/wasaphoto/service/features/auth"
	"github.com/simonesestito/wasaphoto/service/features/comments"
	"github.com/simonesestito/wasaphoto/service/features/follow"
	"github.com/simonesestito/wasaphoto/service/features/likes"
	"github.com/simonesestito/wasaphoto/service/features/photo"
	"github.com/simonesestito/wasaphoto/service/features/stream"
	"github.com/simonesestito/wasaphoto/service/features/user"
)

func (ioc *Container) CreateAuthMiddleware() route.AuthMiddleware {
	return auth.Middleware{LoginService: ioc.CreateAuthService()}
}

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

func (ioc *Container) CreateLikesController() likes.Controller {
	return likes.Controller{
		Service: ioc.CreateLikesService(),
	}
}

func (ioc *Container) CreateStreamController() stream.Controller {
	return stream.Controller{
		Service: ioc.CreateStreamService(),
	}
}

func (ioc *Container) CreateCommentsController() comments.Controller {
	return comments.Controller{Service: ioc.CreateCommentsService()}
}
