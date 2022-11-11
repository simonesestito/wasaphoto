package comments

import (
	"github.com/gofrs/uuid"
	"github.com/simonesestito/wasaphoto/service/features/photo"
	"github.com/simonesestito/wasaphoto/service/features/user"
	"time"
)

type IdParams struct {
	CommentId string `json:"commentId" validate:"required"`
}

type LikeParams struct {
	photo.IdParam `json:",squash"`
	IdParams      `json:",squash"`
}

type NewComment struct {
	Text string `json:"text" validate:"required,min=5,max=500"`
}

type Comment struct {
	Id          uuid.UUID `json:"id"`
	PublishDate time.Time `json:"publishDate"`
	Author      user.User `json:"author"`
	NewComment  `json:",squash"`
}
