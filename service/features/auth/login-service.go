package auth

import (
	"errors"
	"github.com/gofrs/uuid"
	"github.com/simonesestito/wasaphoto/service/features/user"
)

type LoginService interface {
	AuthenticateOrSignup(credentials UserLoginCredentials) (authToken string, isNew bool, err error)
	IsAuthenticated(authToken string) (userId string, err error)
}

type UserIdLoginService struct {
	// Dependencies
	UserDao user.Dao
}

var ErrUnknownUser = errors.New("invalid user credentials")

// AuthenticateOrSignup tries to authenticate the user specified with the credentials,
// and returns the user token or an error.
// In case a user with these credentials cannot be found, it creates a new user with that username.
func (service UserIdLoginService) AuthenticateOrSignup(credentials UserLoginCredentials) (string, bool, error) {
	newUuid, err := uuid.NewV4()
	if err != nil {
		return "", false, err
	}

	// Try registering a new user with this ID
	newUser := user.ModelUser{
		Id:       newUuid.Bytes(),
		Name:     credentials.Username,
		Surname:  "",
		Username: credentials.Username,
	}

	foundUserId, isNew, err := service.UserDao.InsertOrGetUserId(newUser)
	if err != nil {
		return "", isNew, err
	}

	return foundUserId.String(), isNew, nil
}

// IsAuthenticated checks if the given authToken can be assigned to a User.
// In case no user is found, it returns ErrUnknownUser
func (service UserIdLoginService) IsAuthenticated(authToken string) (string, error) {
	userId, err := uuid.FromString(authToken)
	if err != nil {
		return "", err
	}

	foundUser, err := service.UserDao.GetUserById(userId)
	if err != nil {
		return "", err
	} else if foundUser == nil {
		return "", ErrUnknownUser
	}

	return foundUser.Uuid().String(), nil
}
