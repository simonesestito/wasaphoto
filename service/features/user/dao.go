package user

import (
	"github.com/gofrs/uuid"
	"github.com/simonesestito/wasaphoto/service/database"
)

type Dao interface {
	GetUserById(id uuid.UUID) (User, bool)
	GetUserByUsername(username string) (User, bool)
}

type DbDao struct {
	DB database.AppDatabase
}

func (dao DbDao) GetUserById(id uuid.UUID) (User, bool) {
	// TODO: Do query
	return User{}, false
}

func (dao DbDao) GetUserByUsername(username string) (User, bool) {
	// TODO: Do query
	return User{}, false
}
