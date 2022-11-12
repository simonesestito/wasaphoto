package user

import (
	"github.com/gofrs/uuid"
	"github.com/simonesestito/wasaphoto/service/api"
)

type Service interface {
	GetUserAs(searchedId string, searchAsId string) (*User, error)
}

type ServiceImpl struct {
	Db Dao
}

func (service ServiceImpl) GetUserAs(searchedId string, searchAsId string) (*User, error) {
	searchedUuid := uuid.FromStringOrNil(searchedId)
	searchAsUuid := uuid.FromStringOrNil(searchAsId)
	if searchedUuid == uuid.Nil || searchAsUuid == uuid.Nil {
		return nil, api.ErrWrongUUID
	}

	// Check access permission
	ban, err := service.Db.IsUserBannedBy(searchAsUuid, searchedUuid)
	if err != nil {
		return nil, err
	} else if ban {
		return nil, api.ErrUserBanned
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
