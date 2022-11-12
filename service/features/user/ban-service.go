package user

import (
	"github.com/gofrs/uuid"
	"github.com/simonesestito/wasaphoto/service/api"
	"github.com/simonesestito/wasaphoto/service/database"
)

type BanService interface {
	BanUser(bannedId string, bannerId string) error
	UnbanUser(bannedId string, bannerId string) error
}

type BanServiceImpl struct {
	Db Dao
}

func (service BanServiceImpl) BanUser(bannedId string, bannerId string) error {
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

	return nil
}

func (service BanServiceImpl) UnbanUser(bannedId string, bannerId string) error {
	bannedUuid := uuid.FromStringOrNil(bannedId)
	bannerUuid := uuid.FromStringOrNil(bannerId)
	if bannedUuid == uuid.Nil || bannerUuid == uuid.Nil {
		return api.ErrWrongUUID
	}

	_, err := service.Db.UnbanUser(bannedUuid, bannerUuid)
	return err
}
