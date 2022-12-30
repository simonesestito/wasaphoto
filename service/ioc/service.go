package ioc

import (
	"github.com/simonesestito/wasaphoto/service/features/auth"
	"github.com/simonesestito/wasaphoto/service/features/comments"
	"github.com/simonesestito/wasaphoto/service/features/follow"
	"github.com/simonesestito/wasaphoto/service/features/likes"
	"github.com/simonesestito/wasaphoto/service/features/photo"
	"github.com/simonesestito/wasaphoto/service/features/stream"
	"github.com/simonesestito/wasaphoto/service/features/user"
)

func (ioc *Container) createAuthService() auth.LoginService {
	return auth.UserIdLoginService{
		UserDao: ioc.createUserDao(),
	}
}

func (ioc *Container) createUserService() user.Service {
	return user.ServiceImpl{
		Db: ioc.createUserDao(),
	}
}

// createBanService creates a Singleton instance of the BanService,
// because it's required for this service, on the contrary of others.
func (ioc *Container) createBanService() user.BanService {
	const key = "user.BanService"
	if previousInstance, ok := ioc.instances[key]; ok {
		castedInstance, ok := previousInstance.(user.BanService)
		if ok {
			return castedInstance
		} else {
			ioc.logger.Fatalf("Unable to recycle old storage instance in ioc.CreateBanService")
		}
	}

	// Create a new BanService
	newInstance := user.BanServiceImpl{
		Db: ioc.createUserDao(),
	}
	ioc.instances[key] = &newInstance
	return &newInstance
}

func (ioc *Container) createFollowService() follow.Service {
	return follow.NewServiceImpl(
		ioc.createFollowDao(),
		ioc.createBanService(),
		ioc.createUserService(),
	)
}

func (ioc *Container) createPhotoService() photo.Service {
	return photo.ServiceImpl{
		Db:          ioc.createPhotoDao(),
		Storage:     ioc.CreateStorage(),
		UserService: ioc.createUserService(),
		BanService:  ioc.createBanService(),
	}
}

func (ioc *Container) createLikesService() likes.Service {
	return likes.ServiceImpl{
		Db:           ioc.createLikesDao(),
		BanService:   ioc.createBanService(),
		PhotoService: ioc.createPhotoService(),
	}
}

func (ioc *Container) createCommentsService() comments.Service {
	return comments.ServiceImpl{
		Db:           ioc.createCommentsDao(),
		BanService:   ioc.createBanService(),
		PhotoService: ioc.createPhotoService(),
		TimeProvider: ioc.createTimeProvider(),
	}
}

func (ioc *Container) createStreamService() stream.Service {
	return stream.ServiceImpl{Db: ioc.createStreamDao()}
}
