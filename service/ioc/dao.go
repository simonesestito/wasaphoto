package ioc

import (
	"github.com/simonesestito/wasaphoto/service/features/comments"
	"github.com/simonesestito/wasaphoto/service/features/follow"
	"github.com/simonesestito/wasaphoto/service/features/likes"
	"github.com/simonesestito/wasaphoto/service/features/photo"
	"github.com/simonesestito/wasaphoto/service/features/user"
)

func (ioc *Container) CreateUserDao() user.Dao {
	return user.DbDao{Db: ioc.database}
}

func (ioc *Container) CreateFollowDao() follow.Dao {
	return follow.DbDao{Db: ioc.database}
}

func (ioc *Container) CreatePhotoDao() photo.Dao {
	return photo.DbDao{Db: ioc.database, Time: ioc.CreateTimeProvider()}
}

func (ioc *Container) CreateLikesDao() likes.Dao {
	return likes.DbDao{Db: ioc.database}
}

func (ioc *Container) CreateCommentsDao() comments.Dao {
	return comments.DbDao{Db: ioc.database}
}
