package follow

import (
	"github.com/simonesestito/wasaphoto/service/features/user"
)

type FollowerParams struct {
	user.IdParams
	FollowedId string `json:"followedId" validate:"required"`
}

type UserFollow struct {
	FollowingId string `json:"followingId"`
	FollowerId  string `json:"followerId"`
}
