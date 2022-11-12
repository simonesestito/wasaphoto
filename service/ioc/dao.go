package ioc

import "github.com/simonesestito/wasaphoto/service/features/user"

func (ioc *Container) CreateUserDao() user.Dao {
	return user.DbDao{DB: ioc.database}
}
