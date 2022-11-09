package follow

import (
	"github.com/simonesestito/wasaphoto/service/features/user"
)

type FollowerParams struct {
	user.IdParams `json:",squash"`
	FollowedId    string `json:"followedId" validate:"required"`
}
