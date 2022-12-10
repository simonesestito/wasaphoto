package follow

import (
	"github.com/simonesestito/wasaphoto/service/features/user"
)

type followerParams struct {
	user.IdParams
	FollowedId string `json:"followedId" validate:"required,uuid"`
}

type userFollow struct {
	FollowingId string `json:"followingId"`
	FollowerId  string `json:"followerId"`
}
