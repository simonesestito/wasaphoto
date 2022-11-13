package user

import (
	"github.com/gofrs/uuid"
	"github.com/simonesestito/wasaphoto/service/api"
	"github.com/simonesestito/wasaphoto/service/database"
)

type Service interface {
	GetUserAs(searchedId string, searchAsId string) (*User, error)
	UpdateUserDetails(id string, newUser NewUser) (User, error)
	UpdateUsername(id string, username string) (User, error)
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

func (service ServiceImpl) UpdateUserDetails(id string, newUser NewUser) (User, error) {
	userUuid := uuid.FromStringOrNil(id)
	if userUuid == uuid.Nil {
		return User{}, api.ErrWrongUUID
	}

	err := service.Db.EditUser(userUuid, ModelUser{
		Name:     newUser.Name,
		Surname:  newUser.Surname,
		Username: newUser.Username,
	})

	if err == database.ErrDuplicated {
		return User{}, api.ErrAlreadyTaken
	} else if err != nil {
		return User{}, err
	}

	updatedUser, err := service.Db.GetUserByIdAs(userUuid, userUuid)
	if err != nil {
		return User{}, err
	}
	return updatedUser.ToDto(), nil
}

func (service ServiceImpl) UpdateUsername(id string, username string) (User, error) {
	userUuid := uuid.FromStringOrNil(id)
	if userUuid == uuid.Nil {
		return User{}, api.ErrWrongUUID
	}

	err := service.Db.EditUsername(userUuid, username)

	if err == database.ErrDuplicated {
		return User{}, api.ErrAlreadyTaken
	} else if err != nil {
		return User{}, err
	}

	updatedUser, err := service.Db.GetUserByIdAs(userUuid, userUuid)
	if err != nil {
		return User{}, err
	}
	return updatedUser.ToDto(), nil
}
