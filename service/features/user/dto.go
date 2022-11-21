package user

import "github.com/simonesestito/wasaphoto/service/api"

type SearchParams struct {
	UsernameGetParams
	api.PaginationInfo
	ExactMatch bool `json:"exactMatch"`
}

type UsernameGetParams struct {
	Username string `json:"username" validate:"required,min=3"`
}

type BanParams struct {
	IdParams
	BannedId string `json:"bannedId" validate:"required"`
}

type NewUser struct {
	Name     string `json:"name" validate:"required,min=2"`
	Surname  string `json:"surname"`
	Username string `json:"username" validate:"required,min=3,max=16"`
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
	UserId string `json:"userId" validate:"required"`
}

type BanResult struct {
	BannedId string `json:"bannedId"`
	BannerId string `json:"bannerId"`
}
