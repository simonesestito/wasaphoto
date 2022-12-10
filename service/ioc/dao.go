package ioc

import (
	"github.com/simonesestito/wasaphoto/service/features/comments"
	"github.com/simonesestito/wasaphoto/service/features/follow"
	"github.com/simonesestito/wasaphoto/service/features/likes"
	"github.com/simonesestito/wasaphoto/service/features/photo"
	"github.com/simonesestito/wasaphoto/service/features/stream"
	"github.com/simonesestito/wasaphoto/service/features/user"
)

func (ioc *Container) createUserDao() user.Dao {
	return user.DbDao{Db: ioc.database}
}

func (ioc *Container) createFollowDao() follow.Dao {
	return follow.DbDao{Db: ioc.database}
}

func (ioc *Container) createPhotoDao() photo.Dao {
	return photo.DbDao{Db: ioc.database, Time: ioc.createTimeProvider()}
}

func (ioc *Container) createLikesDao() likes.Dao {
	return likes.DbDao{Db: ioc.database}
}

func (ioc *Container) createCommentsDao() comments.Dao {
	return comments.DbDao{Db: ioc.database}
}

func (ioc *Container) createStreamDao() stream.Dao {
	return stream.DbDao{Db: ioc.database}
}
