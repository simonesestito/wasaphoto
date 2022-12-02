package likes

import (
	"github.com/gofrs/uuid"
	"github.com/simonesestito/wasaphoto/service/database"
)

type Dao interface {
	LikePhoto(photoUuid uuid.UUID, userUuid uuid.UUID) (bool, error)
	UnlikePhoto(photoUuid uuid.UUID, userUuid uuid.UUID) (bool, error)
}

type DbDao struct {
	Db database.AppDatabase
}

func (db DbDao) LikePhoto(photoUuid uuid.UUID, userUuid uuid.UUID) (bool, error) {
	rows, err := db.Db.ExecRows("INSERT INTO Likes (userId, photoId) VALUES (?, ?)",
		userUuid.Bytes(), photoUuid.Bytes())
	return rows > 0, err
}

func (db DbDao) UnlikePhoto(photoUuid uuid.UUID, userUuid uuid.UUID) (bool, error) {
	rows, err := db.Db.ExecRows("DELETE FROM Likes WHERE photoId = ? AND userId = ?",
		photoUuid.Bytes(), userUuid.Bytes())
	return rows > 0, err
}
