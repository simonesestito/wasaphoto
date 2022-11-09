package auth

import "errors"

type LoginService interface {
	Authenticate(credentials UserLoginCredentials) (authToken string, err error)
	IsAuthenticated(authToken string) (userId string, err error)
}

type UserIdLoginService struct {
	// Dependencies
}

func (service UserIdLoginService) Authenticate(credentials UserLoginCredentials) (string, error) {
	// TODO: use real data source
	if credentials.Username == "mario" {
		return "1111-2222-mario", nil
	} else {
		return "", errors.New("invalid credentials")
	}
}

func (service UserIdLoginService) IsAuthenticated(authToken string) (string, error) {
	// TODO: use real data source
	if authToken == "1111-2222-mario" {
		return authToken, nil // Same as UserID
	} else {
		return "", errors.New("invalid token")
	}
}
