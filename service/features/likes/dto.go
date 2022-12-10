package likes

import (
	"github.com/simonesestito/wasaphoto/service/features/photo"
	"github.com/simonesestito/wasaphoto/service/features/user"
)

type likeParams struct {
	photo.IdParam
	user.IdParams
}

type photoLike struct {
	PhotoId string `json:"photoId"`
	UserId  string `json:"userId"`
}
