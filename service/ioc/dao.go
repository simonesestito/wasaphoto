package ioc

import (
	"github.com/simonesestito/wasaphoto/service/features/follow"
	"github.com/simonesestito/wasaphoto/service/features/user"
)

func (ioc *Container) CreateUserDao() user.Dao {
	return user.DbDao{DB: ioc.database}
}

func (ioc *Container) CreateFollowDao() follow.Dao {
	return follow.DbDao{Db: ioc.database}
}
