package user

import "github.com/simonesestito/wasaphoto/service/api"

type SearchParams struct {
	UsernameGetParams
	api.PaginationInfo
	ExactMatch bool `json:"exactMatch"`
}

type UsernameGetParams struct {
	Username string `json:"username" validate:"required,username"`
}

type BanParams struct {
	IdParams
	BannedId string `json:"bannedId" validate:"required,uuid"`
}

type NewUser struct {
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
	NewUser
}

type IdParams struct {
	UserId string `json:"userId" validate:"required,uuid"`
}

type BanResult struct {
	BannedId string `json:"bannedId"`
	BannerId string `json:"bannerId"`
}

type IdUserCursor struct {
	api.PaginationInfo
	IdParams
}
