package auth

type UserLoginCredentials struct {
	Username string `json:"username" validate:"required,min=3,max=16"`
}

type UserLoginResult struct {
	UserId string `json:"userId"`
}
