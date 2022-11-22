package photo

import (
	"github.com/gofrs/uuid"
	"github.com/simonesestito/wasaphoto/service/features/user"
	"github.com/simonesestito/wasaphoto/service/timeprovider"
	"github.com/simonesestito/wasaphoto/service/utils/cursor"
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
	photo.ModelUserWithCustom.Id = photo.AuthorId

	return Photo{
		Id:            uuid.FromBytesOrNil(photo.EntityPhoto.Id).String(),
		Author:        photo.ModelUserWithCustom.ToDto(),
		PublishDate:   publishDate,
		LikesCount:    photo.LikesCount,
		CommentsCount: photo.CommentsCount,
		Liked:         photo.Liked > 0,
		ImageUrl:      photo.ImageUrl,
	}
}

func DbPhotosListToPage(dbPhotos []EntityPhotoAuthorInfo) (photos []Photo, pageCursor *string) {
	photos = make([]Photo, len(dbPhotos))
	for i, dbPhoto := range dbPhotos {
		photos[i] = dbPhoto.ToDto()
	}

	// Calculate next cursor
	if len(dbPhotos) > 0 {
		lastPhoto := dbPhotos[len(dbPhotos)-1]
		nextCursor := cursor.CreateDateIdCursor(lastPhoto.EntityPhoto.Id, lastPhoto.PublishDate)
		pageCursor = &nextCursor
	} else {
		pageCursor = nil
	}

	return
}
