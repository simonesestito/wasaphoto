package photo

import (
	"github.com/simonesestito/wasaphoto/service/features/user"
	"github.com/simonesestito/wasaphoto/service/timeprovider"
	"time"
)

type EntityPhoto struct {
	Id          []byte `json:"id"`
	ImageUrl    string `json:"imageUrl"`
	AuthorId    []byte `json:"authorId"`
	PublishDate string `json:"publishDate"`
}

type EntityPhotoInfo struct {
	EntityPhoto
	LikesCount    uint `json:"likesCount"`
	CommentsCount uint `json:"commentsCount"`
}

type EntityPhotoInfoWithCustom struct {
	EntityPhotoInfo
	Liked int64 `json:"liked"`
}

type EntityPhotoAuthorInfo struct {
	user.ModelUserWithCustom
	EntityPhotoInfoWithCustom
}

func (photo EntityPhotoAuthorInfo) ToDto() Photo {
	publishDate, _ := time.Parse(timeprovider.UTCFormat, photo.PublishDate)
	return Photo{
		Id:            photo.Uuid().String(),
		Author:        photo.ModelUserWithCustom.ToDto(),
		PublishDate:   publishDate,
		LikesCount:    photo.LikesCount,
		CommentsCount: photo.CommentsCount,
		Liked:         photo.Liked > 0,
		ImageUrl:      photo.ImageUrl,
	}
}
