package user

import "github.com/simonesestito/wasaphoto/service/api"

type searchParams struct {
	usernameGetParams
	api.PaginationInfo
	ExactMatch bool `json:"exactMatch"`
}

type usernameGetParams struct {
	Username string `json:"username" validate:"required,username"`
}

type banParams struct {
	IdParams
	BannedId string `json:"bannedId" validate:"required,uuid"`
}

type newUser struct {
	Name     string `json:"name" validate:"required,min=2,max=256,singleline"`
	Surname  string `json:"surname" validate:"max=256,singleline"`
	Username string `json:"username" validate:"required,username"`
}

type User struct {
	Id              string `json:"id"`
	FollowersCount  uint   `json:"followersCount"`
	FollowingsCount uint   `json:"followingsCount"`
	PostsCount      uint   `json:"postsCount"`
	Banned          bool   `json:"banned"`
	Following       bool   `json:"following"`
	newUser
}

type IdParams struct {
	UserId string `json:"userId" validate:"required,uuid"`
}

type banResult struct {
	BannedId string `json:"bannedId"`
	BannerId string `json:"bannerId"`
}

type IdUserCursor struct {
	api.PaginationInfo
	IdParams
}
