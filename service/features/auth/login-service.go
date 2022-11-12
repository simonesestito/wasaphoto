package auth

import (
	"errors"
	"github.com/gofrs/uuid"
	"github.com/simonesestito/wasaphoto/service/features/user"
)

type LoginService interface {
	Authenticate(credentials UserLoginCredentials) (authToken string, err error)
	IsAuthenticated(authToken string) (userId string, err error)
	SignUp(credentials UserLoginCredentials) (authToken string, err error)
}

type UserIdLoginService struct {
	// Dependencies
	UserDao user.Dao
}

var ErrUnknownUser = errors.New("invalid user credentials")

// Authenticate tries to authenticate the user specified with the credentials,
// and returns the user token or an error.
// In case a user with these credentials cannot be found, it returns ErrUnknownUser
func (service UserIdLoginService) Authenticate(credentials UserLoginCredentials) (string, error) {
	foundUser, err := service.UserDao.GetUserByUsername(credentials.Username)
	switch {
	case foundUser == nil:
		return "", ErrUnknownUser
	case err != nil:
		return "", err
	default:
		return foundUser.Uuid().String(), nil
	}
}

// SignUp creates a new user based on the given credentials.
// The name will be the username, and the surname will be empty.
func (service UserIdLoginService) SignUp(credentials UserLoginCredentials) (string, error) {
	newUuid, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	newUser := user.ModelUser{
		Id:       newUuid.Bytes(),
		Name:     credentials.Username,
		Surname:  "",
		Username: credentials.Username,
	}

	if err := service.UserDao.InsertUser(newUser); err != nil {
		return "", err
	}

	return newUuid.String(), nil
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
