package photo

import (
	"bytes"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/simonesestito/wasaphoto/service/api"
	"github.com/simonesestito/wasaphoto/service/features/user"
	"github.com/simonesestito/wasaphoto/service/storage"
	"github.com/simonesestito/wasaphoto/service/timeprovider"
	"github.com/simonesestito/wasaphoto/service/utils/cursor"
	"github.com/sirupsen/logrus"
)

type Service interface {
	CreatePost(userId string, imageData []byte, logger logrus.FieldLogger) (Photo, error)
	DeletePostAs(imageId string, userId string) error
	GetPostAuthorById(imageId string) (string, error)
	GetUsersPhotosPage(id string, searchAs string, cursor string) ([]Photo, *string, error)
	GetPhotoByIdAs(photoId string, searchAs string) (*Photo, error)
}

type ServiceImpl struct {
	Db             Dao
	Storage        storage.Storage
	ImageProcessor imageProcessor
	UserService    user.Service
	BanService     user.BanService
}

func (service ServiceImpl) CreatePost(userId string, imageData []byte, logger logrus.FieldLogger) (Photo, error) {
	userUuid := uuid.FromStringOrNil(userId)
	if userUuid == uuid.Nil {
		return Photo{}, api.ErrWrongUUID
	}

	// Process image
	imageData, err := service.ImageProcessor.compressPhotoToWebp(imageData, logger)
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
	return photo.toDto(), nil
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

func (service ServiceImpl) GetUsersPhotosPage(id string, searchAs string, pageCursor string) ([]Photo, *string, error) {
	authorUuid := uuid.FromStringOrNil(id)
	searchAsUuid := uuid.FromStringOrNil(searchAs)
	if authorUuid.IsNil() || searchAsUuid.IsNil() {
		return nil, nil, api.ErrWrongUUID
	}

	nextPhotoId, nextDate, err := cursor.ParseDateIdCursor(pageCursor)
	if err != nil {
		return nil, nil, api.ErrWrongCursor
	}

	// Check if the searched user exists and ban status
	searchedUser, err := service.UserService.GetUserAs(id, searchAs)
	switch {
	case errors.Is(err, api.ErrUserBanned):
		return nil, nil, api.ErrUserBanned
	case err != nil:
		return nil, nil, err
	case searchedUser == nil:
		return nil, nil, api.ErrNotFound
	}

	// Get photos
	dbPhotos, err := service.Db.ListUsersPhotoAfter(authorUuid, searchAsUuid, nextPhotoId, timeprovider.DateToUTCString(nextDate))
	if err != nil {
		return nil, nil, err
	}

	photos, nextCursor := DbPhotosListToPage(dbPhotos)
	return photos, nextCursor, nil
}

func (service ServiceImpl) GetPhotoByIdAs(photoId string, searchAs string) (*Photo, error) {
	photoUuid := uuid.FromStringOrNil(photoId)
	searchAsUuid := uuid.FromStringOrNil(searchAs)
	if photoUuid.IsNil() || searchAsUuid.IsNil() {
		return nil, api.ErrWrongUUID
	}

	// Get the required photo
	dbPhoto, err := service.Db.GetPhotoByIdAs(photoUuid, searchAsUuid)
	if err != nil {
		return nil, err
	} else if dbPhoto == nil {
		return nil, nil
	}

	photo := dbPhoto.toDto()

	// Check author ban
	isBanned, err := service.BanService.IsUserBanned(searchAs, photo.Author.Id)
	if err != nil {
		return nil, err
	} else if isBanned {
		return nil, api.ErrUserBanned
	}

	// Success! Return photo
	return &photo, nil
}
