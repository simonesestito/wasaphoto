package follow

import (
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/features/user"
)

type FollowerParams struct {
	user.IdParams `json:",squash"`
	FollowedId    string `json:"followedId" validate:"required"`
}
