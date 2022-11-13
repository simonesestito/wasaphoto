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
	AddBanListener(listener BanListener)
}

type BanListener func(bannedId string, bannerId string) error

type BanServiceImpl struct {
	Db        Dao
	listeners []BanListener
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

func (service *BanServiceImpl) AddBanListener(listener BanListener) {
	if service.listeners == nil {
		service.listeners = make([]BanListener, 0, 1)
	}

	service.listeners = append(service.listeners, listener)
}
