package ioc

import (
	"github.com/simonesestito/wasaphoto/service/features/auth"
	"github.com/simonesestito/wasaphoto/service/features/follow"
	"github.com/simonesestito/wasaphoto/service/features/user"
)

func (ioc *Container) CreateAuthService() auth.LoginService {
	return auth.UserIdLoginService{
		UserDao: ioc.CreateUserDao(),
	}
}

func (ioc *Container) CreateUserService() user.Service {
	return user.ServiceImpl{
		Db: ioc.CreateUserDao(),
	}
}

// CreateBanService creates a Singleton instance of the BanService,
// because it's required for this service, on the contrary of others.
func (ioc *Container) CreateBanService() user.BanService {
	const key = "user.BanService"
	if previousInstance, ok := ioc.instances[key]; ok {
		return previousInstance.(user.BanService)
	}

	// Create a new BanService
	newInstance := user.BanServiceImpl{
		Db: ioc.CreateUserDao(),
	}
	ioc.instances[key] = &newInstance
	return &newInstance
}

func (ioc *Container) CreateFollowService() follow.Service {
	return follow.NewServiceImpl(
		ioc.CreateFollowDao(),
		ioc.CreateBanService(),
	)
}
