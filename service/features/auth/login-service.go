package auth

import (
	"errors"
	"github.com/gofrs/uuid"
	"github.com/simonesestito/wasaphoto/service/features/user"
	"github.com/sirupsen/logrus"
)

type LoginService interface {
	Authenticate(credentials UserLoginCredentials, logger logrus.FieldLogger) (authToken string, err error)
	IsAuthenticated(authToken string, logger logrus.FieldLogger) (userId string, err error)
}

type UserIdLoginService struct {
	// Dependencies
	UserDao user.Dao
}

func (service UserIdLoginService) Authenticate(credentials UserLoginCredentials, logger logrus.FieldLogger) (string, error) {
	foundUser, err := service.UserDao.GetUserByUsername(credentials.Username)
	if err != nil {
		logger.WithError(err).Errorf("Unexpected error fetching user with username '%s'", credentials.Username)
		return "", errors.New("invalid credentials")
		// FIXME: Sign up new user (name=username, surname="")
	} else if foundUser == nil {
		return "", errors.New("invalid credentials")
	}

	return foundUser.Uuid().String(), nil
}

func (service UserIdLoginService) IsAuthenticated(authToken string, logger logrus.FieldLogger) (string, error) {
	userId, err := uuid.FromString(authToken)
	if err != nil {
		return "", err
	}

	foundUser, err := service.UserDao.GetUserById(userId)
	if err != nil {
		logger.WithError(err).Errorf("Unexpected error fetching user with ID '%s'", userId.String())
		return "", errors.New("invalid token")
	} else if foundUser == nil {
		return "", errors.New("invalid token")
	}

	return foundUser.Uuid().String(), nil
}
