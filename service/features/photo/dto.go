package photo

import (
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/features/user"
	"time"
)

type NewPhoto struct {
	ImageUrl string `json:"imageUrl" validate:"required,datauri"`
}

type Photo struct {
	Id            string    `json:"id"`
	Author        user.User `json:"author"`
	PublishDate   time.Time `json:"publishDate"`
	LikesCount    uint      `json:"likesCount"`
	CommentsCount uint      `json:"commentsCount"`
	Liked         bool      `json:"liked"`
	NewPhoto      `json:",squash"`
}

type IdParam struct {
	PhotoId string `json:"photoId" validate:"required"`
}
