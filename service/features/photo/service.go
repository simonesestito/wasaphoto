package photo

import (
	"bytes"
	"github.com/gofrs/uuid"
	"github.com/simonesestito/wasaphoto/service/api"
	"github.com/simonesestito/wasaphoto/service/storage"
)

type Service interface {
	CreatePost(userId string, imageData []byte) (Photo, error)
	DeletePostAs(imageId string, userId string) error
	GetPostAuthorById(imageId string) (string, error)
}

type ServiceImpl struct {
	Db             Dao
	Storage        storage.Storage
	ImageProcessor ImageProcessor
}

func (service ServiceImpl) CreatePost(userId string, imageData []byte) (Photo, error) {
	userUuid := uuid.FromStringOrNil(userId)
	if userUuid == uuid.Nil {
		return Photo{}, api.ErrWrongUUID
	}

	// Process image
	imageData, err := service.ImageProcessor.CompressPhotoToWebp(imageData)
	if err != nil {
		return Photo{}, err
	}

	// Generate new UUID
	photoUuid, err := uuid.NewV4()
	if err != nil {
		return Photo{}, err
	}

	// Save processed photo
	photoPath := service.pathForPhotoFile(photoUuid)
	savedFilePath, err := service.Storage.SaveFile(photoPath, imageData)
	if err != nil {
		return Photo{}, err
	}

	// Handle errors in inserting the image in the DB, preparing a rollback
	isCommitted := false
	defer func() {
		if !isCommitted {
			// Rollback!
			_ = service.Storage.DeleteFile(photoPath)
			_ = service.Db.DeletePhoto(photoUuid)
		}
	}()

	// Create new photo struct
	err = service.Db.NewPhotoPerUser(photoUuid, userUuid, savedFilePath)
	if err != nil {
		return Photo{}, err
	}

	// Get just created photo
	photo, err := service.Db.GetPhotoByIdAs(photoUuid, userUuid)
	if err != nil {
		return Photo{}, err
	}

	// Commit!
	isCommitted = true

	// Get current server URL
	return photo.ToDto(), nil
}

func (ServiceImpl) pathForPhotoFile(photoUuid uuid.UUID) string {
	return "photos/" + photoUuid.String() + ".webp"
}

func (service ServiceImpl) DeletePostAs(imageId string, userId string) error {
	imageUuid := uuid.FromStringOrNil(imageId)
	userUuid := uuid.FromStringOrNil(userId)
	if imageUuid.IsNil() || userUuid.IsNil() {
		return api.ErrWrongUUID
	}

	// Get photo
	imageToDelete, err := service.Db.GetPhotoByIdAs(imageUuid, userUuid)
	if err != nil {
		return err
	} else if imageToDelete == nil {
		return api.ErrNotFound
	}

	// Check authorization
	if !bytes.Equal(imageToDelete.AuthorId, userUuid.Bytes()) {
		return api.ErrOthersData
	}

	// Delete photo from database
	err = service.Db.DeletePhoto(imageUuid)
	if err != nil {
		return err
	}

	// Delete photo file from storage
	return service.Storage.DeleteFile(service.pathForPhotoFile(imageUuid))
}

func (service ServiceImpl) GetPostAuthorById(imageId string) (string, error) {
	imageUuid := uuid.FromStringOrNil(imageId)
	if imageUuid.IsNil() {
		return "", api.ErrWrongUUID
	}

	// Get photo
	imageEntity, err := service.Db.GetPhotoById(imageUuid)
	if err != nil {
		return "", err
	} else if imageEntity == nil {
		return "", api.ErrNotFound
	}

	// Return author ID
	authorUuid := uuid.FromBytesOrNil(imageEntity.AuthorId)
	return authorUuid.String(), nil
}
