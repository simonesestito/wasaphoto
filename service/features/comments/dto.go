package comments

import (
	"github.com/simonesestito/wasaphoto/service/api"
	"github.com/simonesestito/wasaphoto/service/features/photo"
	"github.com/simonesestito/wasaphoto/service/features/user"
	"time"
)

type IdParams struct {
	photo.IdParam
	CommentId string `json:"commentId" validate:"required,uuid"`
}

type NewComment struct {
	Text string `json:"text" validate:"required,min=1,max=256"`
}

type Comment struct {
	Id          string    `json:"id"`
	PublishDate time.Time `json:"publishDate"`
	Author      user.User `json:"author"`
	NewComment
}

type PhotoCommentsCursor struct {
	photo.IdParam
	api.PaginationInfo
}
