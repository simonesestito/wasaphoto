package user

import (
	"database/sql"
	"github.com/gofrs/uuid"
)

func (dao DbDao) IsUserBannedBy(bannedId uuid.UUID, bannerId uuid.UUID) (bool, error) {
	result := struct {
		Banned int64 `json:"banned"`
	}{}

	err := dao.Db.QueryStructRow(&result, "SELECT EXISTS(SELECT * FROM Ban WHERE bannedId = ? AND bannerId = ?) AS banned", bannedId.Bytes(), bannerId.Bytes())

	if err != nil {
		return false, err
	}

	return result.Banned > 0, nil
}

func (dao DbDao) BanUser(bannedId uuid.UUID, bannerId uuid.UUID) (bool, error) {
	rows, err := dao.Db.ExecRows("INSERT INTO Ban (bannedId, bannerId) VALUES (?, ?)", bannedId.Bytes(), bannerId.Bytes())
	return rows > 0, err
}

func (dao DbDao) UnbanUser(bannedUuid uuid.UUID, bannerUuid uuid.UUID) (bool, error) {
	rows, err := dao.Db.ExecRows("DELETE FROM Ban WHERE bannedId = ? AND bannerId = ?", bannedUuid.Bytes(), bannerUuid.Bytes())
	return rows > 0, err
}

func (dao DbDao) GetBannedUsersAs(id uuid.UUID, searchAsId uuid.UUID) ([]ModelUserWithCustom, error) {
	query := `SELECT UserInfo.*,
		EXISTS(SELECT * FROM Ban WHERE bannedId = UserInfo.id AND bannerId = ?) AS banned,
		EXISTS(SELECT * FROM Follow WHERE followedId = UserInfo.id AND followerId = ?) AS following
		FROM Ban
		LEFT JOIN UserInfo ON Ban.bannedId = UserInfo.id
		WHERE Ban.bannerId = ?`

	rows, err := dao.Db.QueryStructRows(ModelUserWithCustom{}, query, searchAsId.Bytes(), searchAsId.Bytes(), id.Bytes())
	defer rows.Close()

	var users []ModelUserWithCustom
	var entity any
	for entity, err = rows.Next(); err == nil; entity, err = rows.Next() {
		users = append(users, entity.(ModelUserWithCustom))
	}
	if err != sql.ErrNoRows {
		return nil, err
	}

	return users, nil
}
