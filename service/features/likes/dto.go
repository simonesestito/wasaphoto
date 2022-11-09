package likes

import (
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/features/photo"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/features/user"
)

type LikeParams struct {
	photo.IdParam `json:",squash"`
	user.IdParams `json:",squash"`
}
