package user

import (
	"github.com/gofrs/uuid"
	"github.com/simonesestito/wasaphoto/service/api"
	"github.com/simonesestito/wasaphoto/service/database"
)

type BanService interface {
	BanUser(bannedId string, bannerId string) error
	UnbanUser(bannedId string, bannerId string) error
	IsUserBanned(bannedId string, bannerId string) (bool, error)
	AddBanListener(tag string, listener BanListener)
	GetBannedUsers(id string) ([]User, error)
}

type BanListener func(bannedId string, bannerId string) error

type BanServiceImpl struct {
	Db        Dao
	listeners map[string]BanListener
}

func (service *BanServiceImpl) BanUser(bannedId string, bannerId string) error {
	bannedUuid := uuid.FromStringOrNil(bannedId)
	bannerUuid := uuid.FromStringOrNil(bannerId)
	if bannedUuid == uuid.Nil || bannerUuid == uuid.Nil {
		return api.ErrWrongUUID
	}

	if bannedUuid == bannerUuid {
		return api.ErrSelfOperation
	}

	newBan, err := service.Db.BanUser(bannedUuid, bannerUuid)
	if err == database.ErrForeignKey {
		return api.ErrNotFound
	} else if err != nil {
		return err
	}

	if !newBan {
		return api.ErrDuplicated
	}

	// Execute listeners to perform actions based on other services needs on user ban
	if service.listeners != nil {
		for _, listener := range service.listeners {
			if err := listener(bannedId, bannerId); err != nil {
				return err
			}
		}
	}

	return nil
}

func (service *BanServiceImpl) UnbanUser(bannedId string, bannerId string) error {
	bannedUuid := uuid.FromStringOrNil(bannedId)
	bannerUuid := uuid.FromStringOrNil(bannerId)
	if bannedUuid == uuid.Nil || bannerUuid == uuid.Nil {
		return api.ErrWrongUUID
	}

	_, err := service.Db.UnbanUser(bannedUuid, bannerUuid)
	return err
}

func (service *BanServiceImpl) IsUserBanned(bannedId string, bannerId string) (bool, error) {
	bannedUuid := uuid.FromStringOrNil(bannedId)
	bannerUuid := uuid.FromStringOrNil(bannerId)
	if bannedUuid == uuid.Nil || bannerUuid == uuid.Nil {
		return false, api.ErrWrongUUID
	}

	return service.Db.IsUserBannedBy(bannedUuid, bannerUuid)
}

func (service *BanServiceImpl) AddBanListener(tag string, listener BanListener) {
	if service.listeners == nil {
		service.listeners = make(map[string]BanListener)
	}

	// Use a tag to avoid adding the same listener twice
	service.listeners[tag] = listener
}

func (service *BanServiceImpl) GetBannedUsers(id string) ([]User, error) {
	userUuid := uuid.FromStringOrNil(id)
	if userUuid.IsNil() {
		return nil, api.ErrWrongUUID
	}

	userEntities, err := service.Db.GetBannedUsersAs(userUuid, userUuid)
	if err != nil {
		return nil, err
	}

	users := make([]User, len(userEntities))
	for i, entity := range userEntities {
		users[i] = entity.ToDto()
	}

	return users, nil
}
