package photo

import (
	"fmt"
	"github.com/simonesestito/wasaphoto/service/api"
	"github.com/simonesestito/wasaphoto/service/features/user"
	"net/http"
	"strings"
	"time"
)

type Photo struct {
	Id            string    `json:"id"`
	Author        user.User `json:"author"`
	PublishDate   time.Time `json:"publishDate"`
	LikesCount    uint      `json:"likesCount"`
	CommentsCount uint      `json:"commentsCount"`
	Liked         bool      `json:"liked"`
	ImageUrl      string    `json:"imageUrl"`
}

func (photo *Photo) AddImageHost(r *http.Request) {
	if strings.HasPrefix(photo.ImageUrl, "/") {
		var schema string
		if r.TLS == nil {
			schema = "http"
		} else {
			schema = "https"
		}

		photo.ImageUrl = fmt.Sprintf("%s://%s%s", schema, r.Host, photo.ImageUrl)
	}
}

type IdParam struct {
	PhotoId string `json:"photoId" validate:"required,uuid"`
}

type UserPhotosCursor struct {
	api.PaginationInfo
	user.IdParams
}
