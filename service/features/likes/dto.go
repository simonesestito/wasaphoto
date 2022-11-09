package likes

import (
	"github.com/simonesestito/wasaphoto/service/features/photo"
	"github.com/simonesestito/wasaphoto/service/features/user"
)

type LikeParams struct {
	photo.IdParam `json:",squash"`
	user.IdParams `json:",squash"`
}
