package user

import "github.com/gofrs/uuid"

func (dao DbDao) IsUserBannedBy(bannedId uuid.UUID, bannerId uuid.UUID) (bool, error) {
	result := struct {
		Banned int64 `json:"banned"`
	}{}

	err := dao.DB.QueryStructRow(&result,
		"SELECT EXISTS(SELECT * FROM Ban WHERE bannedId = ? AND bannerId) AS banned",
		bannedId.Bytes(),
		bannerId.Bytes(),
	)

	if err != nil {
		return false, err
	}

	return result.Banned > 0, nil
}
