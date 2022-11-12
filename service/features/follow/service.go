package follow

import (
	"github.com/gofrs/uuid"
	"github.com/simonesestito/wasaphoto/service/api"
	"github.com/simonesestito/wasaphoto/service/database"
)

type Service interface {
	FollowUser(followerId string, followingId string) error
	UnfollowUser(followerId string, followingId string) error
}

type ServiceImpl struct {
	Db Dao
}

func (service ServiceImpl) FollowUser(followerId string, followingId string) error {
	followerUuid := uuid.FromStringOrNil(followerId)
	followingUuid := uuid.FromStringOrNil(followingId)
	if followerUuid == uuid.Nil || followingUuid == uuid.Nil {
		return api.ErrWrongUUID
	}

	if followerId == followingId {
		return api.ErrSelfOperation
	}

	newInsert, err := service.Db.FollowUser(followerUuid, followingUuid)
	if err == database.ErrForeignKey {
		return api.ErrNotFound
	} else if err != nil {
		return err
	}

	if !newInsert {
		return api.ErrDuplicated
	}

	return nil
}

func (service ServiceImpl) UnfollowUser(followerId string, followingId string) error {
	followerUuid := uuid.FromStringOrNil(followerId)
	followingUuid := uuid.FromStringOrNil(followingId)
	if followerUuid == uuid.Nil || followingUuid == uuid.Nil {
		return api.ErrWrongUUID
	}

	_, err := service.Db.UnfollowUser(followerUuid, followingUuid)
	return err
}
