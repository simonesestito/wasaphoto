package user

import (
	"github.com/gofrs/uuid"
	"github.com/simonesestito/wasaphoto/service/database"
	"github.com/simonesestito/wasaphoto/service/utils/cursor"
)

// ModelUser is the entity for the User database table
type ModelUser struct {
	Id       []byte `json:"id"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Username string `json:"username"`
}

// ModelUserInfo represents the actual database entity UserInfo (Data Layer in our architecture)
type ModelUserInfo struct {
	ModelUser
	FollowersCount  uint `json:"followersCount"`
	FollowingsCount uint `json:"followingsCount"`
	PostsCount      uint `json:"photosCount"`
}

// ModelUserWithCustom is ModelUserInfo with all fields which
// depend on the actual user performing the query.
type ModelUserWithCustom struct {
	ModelUserInfo
	Banned    int64 `json:"banned"`
	Following int64 `json:"following"`
}

func (user ModelUserWithCustom) ToDto() User {
	return User{
		Id:              uuid.FromBytesOrNil(user.Id).String(),
		FollowersCount:  user.FollowersCount,
		FollowingsCount: user.FollowingsCount,
		PostsCount:      user.PostsCount,
		Banned:          user.Banned > 0,
		Following:       user.Following > 0,
		newUser: newUser{
			Name:     user.Name,
			Surname:  user.Surname,
			Username: user.Username,
		},
	}
}

func DbUsersListToPage(dbUsers []ModelUserWithCustom) (users []User, pageCursor *string) {
	users = make([]User, len(dbUsers))
	for i, dbUser := range dbUsers {
		users[i] = dbUser.ToDto()
	}

	// Calculate next cursor
	if len(dbUsers) == database.MaxPageItems {
		lastUser := dbUsers[len(dbUsers)-1]
		nextCursor := cursor.CreateStringIdCursor(lastUser.Id, lastUser.Username)
		pageCursor = &nextCursor
	} else {
		pageCursor = nil
	}

	return
}
