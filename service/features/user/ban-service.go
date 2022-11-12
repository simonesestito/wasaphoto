package user

import (
	"errors"
	"github.com/gofrs/uuid"
	"github.com/simonesestito/wasaphoto/service/database"
)

type BanService interface {
	BanUser(bannedId string, bannerId string) error
}

type BanServiceImpl struct {
	Db Dao
}

// ErrSelfOperation is used to indicate a user is performing
// an operation both as subject and object,
// and that is not possible in this circumstance.
var ErrSelfOperation = errors.New("operation not allowed on yourself")

// ErrNotFound is used if the object of an operation cannot be found
var ErrNotFound = errors.New("subject resource not found")

// ErrDuplicated is used if an item was already added in a set
var ErrDuplicated = database.ErrDuplicated

func (service BanServiceImpl) BanUser(bannedId string, bannerId string) error {
	bannedUuid := uuid.FromStringOrNil(bannedId)
	bannerUuid := uuid.FromStringOrNil(bannerId)
	if bannedUuid == uuid.Nil || bannerUuid == uuid.Nil {
		return ErrWrongUUID
	}

	if bannedUuid == bannerUuid {
		return ErrSelfOperation
	}

	newBan, err := service.Db.BanUser(bannedUuid, bannerUuid)
	if err != nil {
		return err
	}

	if !newBan {
		return ErrDuplicated
	}

	return nil
}
