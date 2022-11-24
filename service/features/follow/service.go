package follow

import (
	"github.com/gofrs/uuid"
	"github.com/simonesestito/wasaphoto/service/api"
	"github.com/simonesestito/wasaphoto/service/database"
	"github.com/simonesestito/wasaphoto/service/features/user"
	"github.com/simonesestito/wasaphoto/service/utils/cursor"
)

type Service interface {
	FollowUser(followerId string, followingId string) error
	UnfollowUser(followerId string, followingId string) error
	ListFollowersAs(userId string, searchAs string, pageCursor string) ([]user.User, *string, error)
}

type ServiceImpl struct {
	Db         Dao
	BanService user.BanService
}

func NewServiceImpl(db Dao, banService user.BanService) ServiceImpl {
	service := ServiceImpl{
		Db:         db,
		BanService: banService,
	}

	// Perform actions when a user is banned
	banService.AddBanListener("unfollowUser", service.UnfollowUser)

	return service
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

	banned, err := service.BanService.IsUserBanned(followerId, followingId)
	if err != nil {
		return err
	}
	if banned {
		return api.ErrUserBanned
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

func (service ServiceImpl) ListFollowersAs(userId string, searchAs string, pageCursor string) ([]user.User, *string, error) {
	userUuid := uuid.FromStringOrNil(userId)
	searchAsUuid := uuid.FromStringOrNil(searchAs)
	if userUuid.IsNil() || searchAsUuid.IsNil() {
		return nil, nil, api.ErrWrongUUID
	}

	// Check if "userId" banned me
	isBanned, err := service.BanService.IsUserBanned(searchAs, userId)
	if err != nil {
		return nil, nil, err
	}

	if isBanned {
		return nil, nil, api.ErrUserBanned
	}

	// Parse cursor
	afterFollowerId, afterUsername, err := cursor.ParseStringIdCursor(pageCursor)
	if err != nil {
		return nil, nil, api.ErrWrongCursor
	}

	// Get followers I can actually see
	dbFollowers, err := service.Db.GetFollowersPageAs(userUuid, searchAsUuid, afterFollowerId, afterUsername)
	if err != nil {
		return nil, nil, err
	}

	// Convert to DTO
	photos, nextCursor := user.DbUsersListToPage(dbFollowers)
	return photos, nextCursor, nil
}
