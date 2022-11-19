package comments

import (
	"github.com/simonesestito/wasaphoto/service/features/photo"
	"github.com/simonesestito/wasaphoto/service/features/user"
	"time"
)

type IdParams struct {
	photo.IdParam `json:",squash"`
	CommentId     string `json:"commentId" validate:"required"`
}

type NewComment struct {
	Text string `json:"text" validate:"required,min=5,max=500"`
}

type Comment struct {
	Id          string    `json:"id"`
	PublishDate time.Time `json:"publishDate"`
	Author      user.User `json:"author"`
	NewComment  `json:",squash"`
}
