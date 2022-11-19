package comments

import (
	"bytes"
	"github.com/gofrs/uuid"
	"github.com/simonesestito/wasaphoto/service/api"
	"github.com/simonesestito/wasaphoto/service/database"
	"github.com/simonesestito/wasaphoto/service/features/photo"
	"github.com/simonesestito/wasaphoto/service/features/user"
	"github.com/simonesestito/wasaphoto/service/timeprovider"
)

type Service interface {
	CommentPhoto(photoId string, userId string, comment NewComment) (Comment, error)
	DeleteCommentOnPhotoIfAuthor(commentId string, photoId string, userId string) error
}

type ServiceImpl struct {
	Db           Dao
	BanService   user.BanService
	PhotoService photo.Service
	TimeProvider timeprovider.TimeProvider
}

func (service ServiceImpl) CommentPhoto(photoId string, userId string, comment NewComment) (Comment, error) {
	photoUuid := uuid.FromStringOrNil(photoId)
	userUuid := uuid.FromStringOrNil(userId)
	if photoUuid.IsNil() || userUuid.IsNil() {
		return Comment{}, api.ErrWrongUUID
	}

	// Get info about the photo to like
	photoAuthorId, err := service.PhotoService.GetPostAuthorById(photoId)
	if err == database.ErrNoResult {
		return Comment{}, api.ErrNotFound
	}

	// Check if photo author banned me
	iamBanned, err := service.BanService.IsUserBanned(userId, photoAuthorId)
	if err != nil {
		return Comment{}, err
	}
	if iamBanned {
		return Comment{}, api.ErrUserBanned
	}

	newCommentUuid, err := uuid.NewV4()
	if err != nil {
		return Comment{}, err
	}

	// Publish the comment
	err = service.Db.CreateComment(EntityComment{
		Id:          newCommentUuid.Bytes(),
		Text:        comment.Text,
		PublishDate: service.TimeProvider.UTCString(),
		AuthorId:    userUuid.Bytes(),
		PhotoId:     photoUuid.Bytes(),
	})
	if err == database.ErrForeignKey {
		return Comment{}, api.ErrNotFound
	} else if err != nil {
		return Comment{}, err
	}

	newComment, err := service.Db.GetCommentByIdAs(newCommentUuid, userUuid)
	if err != nil {
		return Comment{}, err
	}

	return newComment.ToDto(), nil
}

func (service ServiceImpl) DeleteCommentOnPhotoIfAuthor(commentId string, photoId string, userId string) error {
	photoUuid := uuid.FromStringOrNil(photoId)
	commentUuid := uuid.FromStringOrNil(commentId)
	userUuid := uuid.FromStringOrNil(userId)
	if photoUuid.IsNil() || commentUuid.IsNil() || userUuid.IsNil() {
		return api.ErrWrongUUID
	}

	// Get comment to delete with all necessary info
	commentInfoIds, err := service.Db.GetCommentInfoIds(commentUuid)
	if err != nil {
		return err
	} else if commentInfoIds == nil {
		return nil // Comment doesn't exist
	}

	// Check provided info against real ones from the database source of truth
	if !bytes.Equal(photoUuid.Bytes(), commentInfoIds.PhotoId) {
		// Wrong Photo ID provided
		return api.ErrNotFound
	} else if !bytes.Equal(userUuid.Bytes(), commentInfoIds.CommentAuthorId) {
		// Wrong user is trying to delete this comment
		return api.ErrOthersData
	}

	// Delete comment
	_, err = service.Db.DeleteByIdPhotoAndAuthor(commentUuid, photoUuid, userUuid)
	return err
}
