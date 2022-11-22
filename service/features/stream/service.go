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

	photos := make([]photo.Photo, len(dbPhotos))
	for i, dbPhoto := range dbPhotos {
		photos[i] = dbPhoto.ToDto()
	}

	// Calculate next cursor
	if len(dbPhotos) > 0 {
		lastPhoto := dbPhotos[len(dbPhotos)-1]
		nextCursor := cursor.CreateDateIdCursor(lastPhoto.EntityPhoto.Id, lastPhoto.PublishDate)
		return photos, &nextCursor, nil
	} else {
		return photos, nil, nil
	}
}
