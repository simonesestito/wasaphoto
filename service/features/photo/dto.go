package photo

import (
	"github.com/simonesestito/wasaphoto/service/api"
	"github.com/simonesestito/wasaphoto/service/features/user"
	"github.com/simonesestito/wasaphoto/service/utils"
	"github.com/sirupsen/logrus"
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

func (photo *Photo) AddImageHost(r *http.Request, logger logrus.FieldLogger) {
	// Check if the actual URL is relative to this host, or it's already an absolute URL
	if strings.HasPrefix(photo.ImageUrl, "/") {
		photo.ImageUrl = utils.GetUrlPrefix(r, logger) + photo.ImageUrl
	}
}

type IdParam struct {
	PhotoId string `json:"photoId" validate:"required,uuid"`
}

type UserPhotosCursor struct {
	api.PaginationInfo
	user.IdParams
}
