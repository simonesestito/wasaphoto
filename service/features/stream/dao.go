package stream

import (
	"github.com/gofrs/uuid"
	"github.com/simonesestito/wasaphoto/service/database"
	"github.com/simonesestito/wasaphoto/service/features/photo"
)

type Dao interface {
	GetMyFollowingsPhotosSortedByDate(userId uuid.UUID, afterId uuid.UUID, beforeDate string) ([]photo.EntityPhotoAuthorInfo, error)
}

type DbDao struct {
	Db database.AppDatabase
}

func (db DbDao) GetMyFollowingsPhotosSortedByDate(userId uuid.UUID, afterId uuid.UUID, beforeDate string) ([]photo.EntityPhotoAuthorInfo, error) {
	query := `
		SELECT PhotoAuthorInfo.*,
		       EXISTS(SELECT * FROM Likes WHERE Likes.photoId = PhotoAuthorInfo.id AND Likes.userId = ?) AS liked,
		       EXISTS(SELECT * FROM Ban WHERE bannedId = PhotoAuthorInfo.authorId AND bannerId = ?) AS banned,
		       EXISTS(SELECT * FROM Follow WHERE followedId = PhotoAuthorInfo.authorId AND followerId = ?) AS following
		FROM PhotoAuthorInfo
		LEFT JOIN Follow ON Follow.followedId = PhotoAuthorInfo.authorId
		WHERE Follow.followerId = ?
		 	  -- Cursor pagination
			  AND (publishDate, id) < (?, ?)
		ORDER BY publishDate DESC, id DESC
		LIMIT ?`

	rows, err := db.Db.QueryStructRows(
		photo.EntityPhotoAuthorInfo{},
		query,
		userId.Bytes(),
		userId.Bytes(),
		userId.Bytes(),
		userId.Bytes(),
		beforeDate,
		afterId.Bytes(),
		database.MaxPageItems,
	)

	if err != nil {
		return nil, err
	}

	return photo.ParsePhotoEntity(rows)
}
