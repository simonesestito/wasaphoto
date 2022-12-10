package likes

import (
	"errors"
	"github.com/gofrs/uuid"
	"github.com/simonesestito/wasaphoto/service/api"
	"github.com/simonesestito/wasaphoto/service/database"
	"github.com/simonesestito/wasaphoto/service/features/photo"
	"github.com/simonesestito/wasaphoto/service/features/user"
)

type Service interface {
	LikePhoto(photoId string, userId string) error
	UnlikePhoto(photoId string, userId string) error
}

type ServiceImpl struct {
	Db           Dao
	BanService   user.BanService
	PhotoService photo.Service
}

func (service ServiceImpl) LikePhoto(photoId string, userId string) error {
	photoUuid := uuid.FromStringOrNil(photoId)
	userUuid := uuid.FromStringOrNil(userId)
	if photoUuid.IsNil() || userUuid.IsNil() {
		return api.ErrWrongUUID
	}

	// Get info about the photo to like
	photoAuthorId, err := service.PhotoService.GetPostAuthorById(photoId)
	if err != nil {
		return err
	}

	// Check if photo author banned me
	iamBanned, err := service.BanService.IsUserBanned(userId, photoAuthorId)
	if err != nil {
		return err
	}
	if iamBanned {
		return api.ErrUserBanned
	}

	// Like photo
	newInsert, err := service.Db.LikePhoto(photoUuid, userUuid)
	if errors.Is(err, database.ErrForeignKey) {
		return api.ErrNotFound
	} else if err != nil {
		return err
	}

	if !newInsert {
		return api.ErrDuplicated
	}

	return nil
}

func (service ServiceImpl) UnlikePhoto(photoId string, userId string) error {
	photoUuid := uuid.FromStringOrNil(photoId)
	userUuid := uuid.FromStringOrNil(userId)
	if photoUuid.IsNil() || userUuid.IsNil() {
		return api.ErrWrongUUID
	}

	_, err := service.Db.UnlikePhoto(photoUuid, userUuid)
	return err
}
