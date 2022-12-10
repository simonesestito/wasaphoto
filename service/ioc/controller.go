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
	return auth.Middleware{LoginService: ioc.createAuthService()}
}

func (ioc *Container) createLoginController() auth.LoginController {
	return auth.LoginController{
		AuthService: ioc.createAuthService(),
	}
}

func (ioc *Container) createUserController() user.Controller {
	return user.Controller{
		Service: ioc.createUserService(),
	}
}

func (ioc *Container) createBanController() user.BanController {
	return user.BanController{
		Service: ioc.createBanService(),
	}
}

func (ioc *Container) createFollowController() follow.Controller {
	return follow.Controller{Service: ioc.createFollowService()}
}

func (ioc *Container) createPhotoController() photo.Controller {
	return photo.Controller{Service: ioc.createPhotoService()}
}

func (ioc *Container) createLikesController() likes.Controller {
	return likes.Controller{
		Service: ioc.createLikesService(),
	}
}

func (ioc *Container) createStreamController() stream.Controller {
	return stream.Controller{
		Service: ioc.createStreamService(),
	}
}

func (ioc *Container) createCommentsController() comments.Controller {
	return comments.Controller{Service: ioc.createCommentsService()}
}
