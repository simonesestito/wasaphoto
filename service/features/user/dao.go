package user

import (
	"github.com/gofrs/uuid"
	"github.com/simonesestito/wasaphoto/service/database"
)

type Dao interface {
	GetUserById(id uuid.UUID) (*ModelUser, error)
	GetUserByUsername(username string) (*ModelUser, error)
	InsertUser(user ModelUser) error
}

type DbDao struct {
	DB database.AppDatabase
}

func (dao DbDao) GetUserById(id uuid.UUID) (*ModelUser, error) {
	user := &ModelUser{}
	err := dao.DB.QueryStructRow(user, "SELECT * FROM User WHERE id = ?", id.Bytes())
	switch {
	case err == database.ErrNoResult:
		return nil, nil
	case err != nil:
		return nil, err
	default:
		return user, nil
	}
}

func (dao DbDao) GetUserByUsername(username string) (*ModelUser, error) {
	user := &ModelUser{}
	err := dao.DB.QueryStructRow(user, "SELECT * FROM User WHERE username = ?", username)
	switch {
	case err == database.ErrNoResult:
		return nil, nil
	case err != nil:
		return nil, err
	default:
		return user, nil
	}
}

func (dao DbDao) InsertUser(user ModelUser) error {
	return dao.DB.Exec("INSERT INTO User (id, name, surname, username) VALUES (?, ?, ?, ?)",
		user.Id, user.Name, user.Surname, user.Username)
}
