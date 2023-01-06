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
	// Check if the actual URL is relative to this host, or it's already an absolute URL
	if strings.HasPrefix(photo.ImageUrl, "/") {
		// Detect if HTTP or HTTPS
		var schema string
		if r.TLS == nil && strings.HasPrefix(r.Host, "localhost:") {
			schema = "http"
		} else {
			// Force HTTPS on domains different from localhost
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
