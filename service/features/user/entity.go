package user

import "github.com/gofrs/uuid"

// ModelUser is the entity for the User database table
type ModelUser struct {
	Id       []byte `json:"id"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Username string `json:"username"`
}

func (user ModelUser) Uuid() uuid.UUID {
	return uuid.FromBytesOrNil(user.Id)
}

// ModelUserInfo represents the actual database entity UserInfo (Data Layer in our architecture)
type ModelUserInfo struct {
	ModelUser
	FollowersCount  uint `json:"followersCount"`
	FollowingsCount uint `json:"followingsCount"`
	PostsCount      uint `json:"photosCount"`
}

type ModelUserWithBan struct {
	ModelUserInfo
	Banned int64 `json:"banned"`
}

func (user ModelUserWithBan) ToDto() User {
	return User{
		Id:              user.Uuid().String(),
		FollowersCount:  user.FollowersCount,
		FollowingsCount: user.FollowingsCount,
		PostsCount:      user.PostsCount,
		Banned:          user.Banned > 0,
		NewUser: NewUser{
			Name:     user.Name,
			Surname:  user.Surname,
			Username: user.Username,
		},
	}
}
