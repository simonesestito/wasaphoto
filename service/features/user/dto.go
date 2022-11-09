package user

import "git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api"

type SearchParams struct {
	UsernameGetParams  `json:",squash"`
	api.PaginationInfo `json:",squash"`
}

type UsernameGetParams struct {
	Username string `json:"username" validate:"required,min=3"`
}

type BanParams struct {
	IdParams `json:",squash"`
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
	NewUser
}

type IdParams struct {
	UserId string `json:"userId" validate:"required"`
}
