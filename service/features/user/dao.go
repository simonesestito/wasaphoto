package user

import (
	"github.com/gofrs/uuid"
	"github.com/simonesestito/wasaphoto/service/database"
)

type Dao interface {
	GetUserById(id uuid.UUID) (*ModelUser, error)
	GetUserByUsername(username string) (*ModelUser, error)
	InsertUser(user ModelUser) error
	IsUserBannedBy(bannedId uuid.UUID, bannerId uuid.UUID) (bool, error)
	GetUserByIdAs(id uuid.UUID, searchAsId uuid.UUID) (*ModelUserWithBan, error)
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

// GetUserByIdAs also adds personal fields such as "banned" which are
// relative to the actual user looking for this data
func (dao DbDao) GetUserByIdAs(id uuid.UUID, searchAsId uuid.UUID) (*ModelUserWithBan, error) {
	user := &ModelUserWithBan{}
	query := "SELECT UserInfo.*, EXISTS(SELECT * FROM Ban WHERE bannedId = ? AND bannerId = ?) AS banned FROM UserInfo WHERE id = ?"
	err := dao.DB.QueryStructRow(user, query, id.Bytes(), searchAsId.Bytes(), id.Bytes())
	switch {
	case err == database.ErrNoResult:
		return nil, nil
	case err != nil:
		return nil, err
	default:
		return user, nil
	}
}
