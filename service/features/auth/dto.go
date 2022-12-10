package auth

type userLoginCredentials struct {
	Username string `json:"username" validate:"required,username"`
}

type userLoginResult struct {
	UserId string `json:"userId"`
}
