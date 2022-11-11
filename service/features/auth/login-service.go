package auth

import (
	"errors"
	"github.com/gofrs/uuid"
	"github.com/simonesestito/wasaphoto/service/features/user"
)

type LoginService interface {
	Authenticate(credentials UserLoginCredentials) (authToken string, err error)
	IsAuthenticated(authToken string) (userId string, err error)
}

type UserIdLoginService struct {
	// Dependencies
	UserDao user.Dao
}

func (service UserIdLoginService) Authenticate(credentials UserLoginCredentials) (string, error) {
	foundUser, isPresent := service.UserDao.GetUserByUsername(credentials.Username)
	if !isPresent {
		return "", errors.New("invalid credentials")
	}

	return foundUser.Id, nil
}

func (service UserIdLoginService) IsAuthenticated(authToken string) (string, error) {
	userId, err := uuid.FromString(authToken)
	if err != nil {
		return "", err
	}

	_, isPresent := service.UserDao.GetUserById(userId)
	if isPresent {
		return userId.String(), nil
	} else {
		return "", errors.New("invalid token")
	}
}
