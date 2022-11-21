package likes

import (
	"github.com/simonesestito/wasaphoto/service/features/photo"
	"github.com/simonesestito/wasaphoto/service/features/user"
)

type LikeParams struct {
	photo.IdParam
	user.IdParams
}

type PhotoLike struct {
	PhotoId string `json:"photoId"`
	UserId  string `json:"userId"`
}
