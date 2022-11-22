package stream

import (
	"github.com/gofrs/uuid"
	"github.com/simonesestito/wasaphoto/service/api"
	"github.com/simonesestito/wasaphoto/service/features/photo"
	"github.com/simonesestito/wasaphoto/service/timeprovider"
	"github.com/simonesestito/wasaphoto/service/utils/cursor"
)

type Service interface {
	GetStreamPage(userId string, cursor string) ([]photo.Photo, *string, error)
}

type ServiceImpl struct {
	Db Dao
}

func (service ServiceImpl) GetStreamPage(userId string, pageCursor string) ([]photo.Photo, *string, error) {
	userUuid := uuid.FromStringOrNil(userId)
	photoUuid, photoDate, err := cursor.ParseDateIdCursor(pageCursor)
	if err != nil {
		return nil, nil, api.ErrWrongCursor
	}

	// Get photos
	dbPhotos, err := service.Db.GetMyFollowingsPhotosSortedByDate(userUuid, photoUuid, timeprovider.DateToUTCString(photoDate))
	if err != nil {
		return nil, nil, err
	}

	photos, nextCursor := photo.DbPhotosListToPage(dbPhotos)
	return photos, nextCursor, nil
}
