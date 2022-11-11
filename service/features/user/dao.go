package user

import (
	"github.com/gofrs/uuid"
	"github.com/simonesestito/wasaphoto/service/database"
)

type Dao interface {
	GetUserById(id uuid.UUID) (*ModelUser, error)
	GetUserByUsername(username string) (*ModelUser, error)
}

type DbDao struct {
	DB database.AppDatabase
}

func (dao DbDao) GetUserById(id uuid.UUID) (*ModelUser, error) {
	user := &ModelUser{}
	if err := dao.DB.QueryStructRow(user, "SELECT * FROM UserInfo WHERE id = ?", id.Bytes()); err != nil {
		return nil, err
	}

	return user, nil
}

func (dao DbDao) GetUserByUsername(username string) (*ModelUser, error) {
	// TODO: Do query
	return &ModelUser{}, nil
}
