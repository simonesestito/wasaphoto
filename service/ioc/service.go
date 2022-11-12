package ioc

import (
	"github.com/simonesestito/wasaphoto/service/features/auth"
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
