package photo

import (
	"github.com/gofrs/uuid"
	"github.com/simonesestito/wasaphoto/service/api"
	"github.com/simonesestito/wasaphoto/service/features/user"
	"github.com/simonesestito/wasaphoto/service/storage"
	"strings"
	"time"
)

type Service interface {
	CreatePost(userId string, imageData []byte) (Photo, error)
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
	filePath := "static/user_content/photos/" + strings.ReplaceAll(photoUuid.String(), "-", "") + ".webp"
	err = service.Storage.SaveFile(filePath, imageData)
	if err != nil {
		return Photo{}, err
	}

	// Create new photo struct
	err = service.Db.NewPhotoPerUser(photoUuid, userUuid, "/"+filePath)
	if err != nil {
		return Photo{}, err
	}

	// Get just created photo
	photo, err := service.Db.GetPhotoById(photoUuid)
	if err != nil {
		return Photo{}, err
	}

	// TODO: Use nested struct from DB
	return Photo{
		Id:            photoUuid.String(),
		Author:        user.User{},
		PublishDate:   time.Now(),
		LikesCount:    0,
		CommentsCount: 0,
		Liked:         false,
		ImageUrl:      photo.ImageUrl,
	}, nil
}
