package photo

import (
	"github.com/simonesestito/wasaphoto/service/api"
	"github.com/simonesestito/wasaphoto/service/features/user"
	"time"
)

type Photo struct {
	Id            string    `json:"id"`
	Author        user.User `json:"author"`
	PublishDate   time.Time `json:"publishDate"`
	LikesCount    uint      `json:"likesCount"`
	CommentsCount uint      `json:"commentsCount"`
	Liked         bool      `json:"liked"`
	ImageUrl      string    `json:"imageUrl" validate:"required,datauri"`
}

type IdParam struct {
	PhotoId string `json:"photoId" validate:"required"`
}

type UserPhotosCursor struct {
	api.PaginationInfo
	user.IdParams
}
