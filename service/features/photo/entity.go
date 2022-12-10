package photo

import (
	"github.com/gofrs/uuid"
	"github.com/simonesestito/wasaphoto/service/database"
	"github.com/simonesestito/wasaphoto/service/features/user"
	"github.com/simonesestito/wasaphoto/service/timeprovider"
	"github.com/simonesestito/wasaphoto/service/utils/cursor"
	"time"
)

type entityPhoto struct {
	Id          []byte `json:"id"`
	ImageUrl    string `json:"imageUrl"`
	AuthorId    []byte `json:"authorId"`
	PublishDate string `json:"publishDate"`
}

type EntityPhotoInfo struct {
	entityPhoto
	LikesCount    uint `json:"likesCount"`
	CommentsCount uint `json:"commentsCount"`
}

type entityPhotoInfoWithCustom struct {
	EntityPhotoInfo
	Liked int64 `json:"liked"`
}

type EntityPhotoAuthorInfo struct {
	user.ModelUserWithCustom
	entityPhotoInfoWithCustom
}

func (photo EntityPhotoAuthorInfo) toDto() Photo {
	publishDate, _ := time.Parse(timeprovider.UTCFormat, photo.PublishDate)
	photo.ModelUserWithCustom.Id = photo.AuthorId

	return Photo{
		Id:            uuid.FromBytesOrNil(photo.entityPhoto.Id).String(),
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
		photos[i] = dbPhoto.toDto()
	}

	// Calculate next cursor
	if len(dbPhotos) == database.MaxPageItems {
		lastPhoto := dbPhotos[len(dbPhotos)-1]
		nextCursor := cursor.CreateDateIdCursor(lastPhoto.entityPhoto.Id, lastPhoto.PublishDate)
		pageCursor = &nextCursor
	} else {
		pageCursor = nil
	}

	return
}
