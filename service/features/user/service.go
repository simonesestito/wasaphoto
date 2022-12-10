package user

import (
	"errors"
	"github.com/gofrs/uuid"
	"github.com/simonesestito/wasaphoto/service/api"
	"github.com/simonesestito/wasaphoto/service/database"
	"github.com/simonesestito/wasaphoto/service/utils/cursor"
)

type Service interface {
	GetUserAs(searchedId string, searchAsId string) (*User, error)
	UpdateUserDetails(id string, newUser newUser) (User, error)
	UpdateUsername(id string, username string) (User, error)
	GetUserByUsernameAs(username string, searchAsId string) (*User, error)
	ListUsersByUsernameAs(username string, searchAsId string, pageCursor string) ([]User, *string, error)
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

func (service ServiceImpl) UpdateUserDetails(id string, newUser newUser) (User, error) {
	userUuid := uuid.FromStringOrNil(id)
	if userUuid == uuid.Nil {
		return User{}, api.ErrWrongUUID
	}

	err := service.Db.EditUser(userUuid, ModelUser{
		Name:     newUser.Name,
		Surname:  newUser.Surname,
		Username: newUser.Username,
	})

	if errors.Is(err, database.ErrDuplicated) {
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

	if errors.Is(err, database.ErrDuplicated) {
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

func (service ServiceImpl) GetUserByUsernameAs(username string, searchAsId string) (*User, error) {
	searchAsUuid := uuid.FromStringOrNil(searchAsId)
	dbUser, err := service.Db.GetUserByUsernameAs(username, searchAsUuid)
	if err != nil {
		return nil, err
	} else if dbUser == nil {
		return nil, nil
	}

	user := dbUser.ToDto()
	return &user, nil
}

func (service ServiceImpl) ListUsersByUsernameAs(username string, searchAsId string, pageCursor string) ([]User, *string, error) {
	searchAsUuid := uuid.FromStringOrNil(searchAsId)
	afterId, afterUsername, err := cursor.ParseStringIdCursor(pageCursor)
	if err != nil {
		return nil, nil, api.ErrWrongCursor
	}

	dbUsers, err := service.Db.ListUsersByUsernameAs(username, searchAsUuid, afterUsername, afterId)
	if err != nil {
		return nil, nil, err
	}

	// Convert to DTO
	users, nextCursor := DbUsersListToPage(dbUsers)
	return users, nextCursor, nil
}
