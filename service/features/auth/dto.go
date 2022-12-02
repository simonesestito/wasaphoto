package auth

type UserLoginCredentials struct {
	Username string `json:"username" validate:"required,username"`
}

type UserLoginResult struct {
	UserId string `json:"userId"`
}
