package follow

import (
	"github.com/gofrs/uuid"
	"github.com/simonesestito/wasaphoto/service/database"
)

type Dao interface {
	FollowUser(followerUuid uuid.UUID, followingUuid uuid.UUID) (bool, error)
	UnfollowUser(followerUuid uuid.UUID, followingUuid uuid.UUID) (bool, error)
}

type DbDao struct {
	Db database.AppDatabase
}

func (db DbDao) FollowUser(followerUuid uuid.UUID, followingUuid uuid.UUID) (bool, error) {
	rows, err := db.Db.ExecRows("INSERT INTO Follow (followerId, followedId) VALUES (?, ?)",
		followerUuid.Bytes(), followingUuid.Bytes())
	return rows > 0, err
}

func (db DbDao) UnfollowUser(followerUuid uuid.UUID, followingUuid uuid.UUID) (bool, error) {
	rows, err := db.Db.ExecRows("DELETE FROM Follow WHERE followerId = ? AND followedId = ?",
		followerUuid.Bytes(), followingUuid.Bytes())
	return rows > 0, err
}
