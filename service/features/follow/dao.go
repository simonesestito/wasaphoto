package follow

import (
	"database/sql"
	"github.com/gofrs/uuid"
	"github.com/simonesestito/wasaphoto/service/database"
	"github.com/simonesestito/wasaphoto/service/features/user"
)

type Dao interface {
	FollowUser(followerUuid uuid.UUID, followingUuid uuid.UUID) (bool, error)
	UnfollowUser(followerUuid uuid.UUID, followingUuid uuid.UUID) (bool, error)
	GetFollowersPageAs(userUuid uuid.UUID, searchAsUuid uuid.UUID, afterFollowerId uuid.UUID, afterUsername string) ([]user.ModelUserWithCustom, error)
	GetFollowingsPageAs(userUuid uuid.UUID, searchAsUuid uuid.UUID, afterFollowerId uuid.UUID, afterUsername string) ([]user.ModelUserWithCustom, error)
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

func (db DbDao) GetFollowersPageAs(userUuid uuid.UUID, searchAsUuid uuid.UUID, afterFollowerId uuid.UUID, afterUsername string) ([]user.ModelUserWithCustom, error) {
	query := `
		SELECT UserInfo.*,
		       EXISTS(SELECT * FROM Ban WHERE bannedId = UserInfo.id AND bannerId = ?) AS banned,
		       EXISTS(SELECT * FROM Follow WHERE followedId = UserInfo.id AND followerId = ?) AS following
		FROM UserInfo
		LEFT JOIN Follow on UserInfo.id = Follow.followerId
		WHERE Follow.followedId = ?
		 	  -- Cursor pagination
			  AND (username, id) > (?, ?)
		ORDER BY username, id
		LIMIT ?`

	rows, err := db.Db.QueryStructRows(
		user.ModelUserWithCustom{},
		query,
		searchAsUuid.Bytes(),
		searchAsUuid.Bytes(),
		userUuid.Bytes(),
		afterUsername,
		afterFollowerId.Bytes(),
		database.MaxPageItems,
	)

	var photos []user.ModelUserWithCustom
	var entity any
	for entity, err = rows.Next(); err == nil; entity, err = rows.Next() {
		photos = append(photos, entity.(user.ModelUserWithCustom))
	}
	if err != sql.ErrNoRows {
		return nil, err
	}

	return photos, nil
}

func (db DbDao) GetFollowingsPageAs(userUuid uuid.UUID, searchAsUuid uuid.UUID, afterFollowerId uuid.UUID, afterUsername string) ([]user.ModelUserWithCustom, error) {
	query := `
		SELECT UserInfo.*,
		       EXISTS(SELECT * FROM Ban WHERE bannedId = UserInfo.id AND bannerId = ?) AS banned,
		       EXISTS(SELECT * FROM Follow WHERE followedId = UserInfo.id AND followerId = ?) AS following
		FROM UserInfo
		LEFT JOIN Follow on UserInfo.id = Follow.followedId
		WHERE Follow.followerId = ?
		 	  -- Cursor pagination
			  AND (username, id) > (?, ?)
		ORDER BY username, id
		LIMIT ?`

	rows, err := db.Db.QueryStructRows(
		user.ModelUserWithCustom{},
		query,
		searchAsUuid.Bytes(),
		searchAsUuid.Bytes(),
		userUuid.Bytes(),
		afterUsername,
		afterFollowerId.Bytes(),
		database.MaxPageItems,
	)

	var photos []user.ModelUserWithCustom
	var entity any
	for entity, err = rows.Next(); err == nil; entity, err = rows.Next() {
		photos = append(photos, entity.(user.ModelUserWithCustom))
	}
	if err != sql.ErrNoRows {
		return nil, err
	}

	return photos, nil
}
