package user

import (
	"errors"
	"github.com/gofrs/uuid"
)

type Service interface {
	GetUserAs(searchedId string, searchAsId string) (*User, error)
}

// ErrUserBanned is used in case the current user has no permission
// to read the requested information because he is banned
// by the owner of that data.
var ErrUserBanned = errors.New("forbidden because of user ban")

// ErrWrongUUID is used to indicate the given ID cannot be interpreted as a UUID
var ErrWrongUUID = errors.New("wrong UUID supplied")

type ServiceImpl struct {
	Db Dao
}

func (service ServiceImpl) GetUserAs(searchedId string, searchAsId string) (*User, error) {
	searchedUuid := uuid.FromStringOrNil(searchedId)
	searchAsUuid := uuid.FromStringOrNil(searchAsId)
	if searchedUuid == uuid.Nil || searchAsUuid == uuid.Nil {
		return nil, ErrWrongUUID
	}

	// Check access permission
	ban, err := service.Db.IsUserBannedBy(searchAsUuid, searchedUuid)
	if err != nil {
		return nil, err
	} else if ban {
		return nil, ErrUserBanned
	}

	// Access granted
	user, err := service.Db.GetUserByIdAs(searchedUuid, searchAsUuid)
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, nil
	}

	result := user.ToDto()
	return &result, nil
}
