package user

import (
	"database/sql"
	"github.com/gofrs/uuid"
	"github.com/simonesestito/wasaphoto/service/database"
)

type Dao interface {
	GetUserById(id uuid.UUID) (*ModelUser, error)
	InsertOrGetUserId(user ModelUser) (id uuid.UUID, isNew bool, err error)
	IsUserBannedBy(bannedId uuid.UUID, bannerId uuid.UUID) (bool, error)
	GetUserByIdAs(id uuid.UUID, searchAsId uuid.UUID) (*ModelUserWithCustom, error)
	BanUser(bannedId uuid.UUID, bannerId uuid.UUID) (bool, error)
	UnbanUser(bannedUuid uuid.UUID, bannerUuid uuid.UUID) (bool, error)
	EditUser(userUuid uuid.UUID, user ModelUser) error
	EditUsername(userUuid uuid.UUID, username string) error
	GetBannedUsersAs(userUuid uuid.UUID, searchAsId uuid.UUID) ([]ModelUserWithCustom, error)
}

type DbDao struct {
	Db database.AppDatabase
}

func (dao DbDao) GetUserById(id uuid.UUID) (*ModelUser, error) {
	user := &ModelUser{}
	err := dao.Db.QueryStructRow(user, "SELECT * FROM User WHERE id = ?", id.Bytes())
	switch {
	case err == database.ErrNoResult:
		return nil, nil
	case err != nil:
		return nil, err
	default:
		return user, nil
	}
}

// InsertOrGetUserId tries to insert the new given user,
// or returns the ID of the existing one in case of conflicts (e.g.: on username).
//
// Since concurrency issues may happen, a transaction is used.
// Boolean return value indicates whether the user is new
func (dao DbDao) InsertOrGetUserId(user ModelUser) (uuid.UUID, bool, error) {
	tx, err := dao.Db.BeginTx()
	if err != nil {
		return uuid.Nil, false, err
	}
	defer tx.Rollback()

	// Get user if exists
	existingUserId := make([]byte, 16)
	result, err := tx.Query("SELECT id FROM User WHERE username = ?", user.Username)
	defer result.Close()

	if err != sql.ErrNoRows && err != nil {
		// Unexpected error
		return uuid.Nil, false, err
	}

	if err == sql.ErrNoRows || !result.Next() {
		// Username not found, create it!
		_, err = tx.Exec("INSERT INTO User (id, name, surname, username) VALUES (?, ?, ?, ?)", user.Id, user.Name, user.Surname, user.Username)
		if err != nil {
			return uuid.Nil, true, err
		}

		return uuid.FromBytesOrNil(user.Id), true, tx.Commit()
	} else {
		// User found!
		if err := result.Scan(&existingUserId); err != nil {
			return uuid.Nil, false, err
		}

		userUuid, err := uuid.FromBytes(existingUserId)
		return userUuid, false, err
	}
}

// GetUserByIdAs also adds personal fields such as "banned" which are
// relative to the actual user looking for this data
func (dao DbDao) GetUserByIdAs(id uuid.UUID, searchAsId uuid.UUID) (*ModelUserWithCustom, error) {
	user := &ModelUserWithCustom{}
	query := "SELECT UserInfo.*, " +
		"EXISTS(SELECT * FROM Ban WHERE bannedId = ? AND bannerId = ?) AS banned, " +
		"EXISTS(SELECT * FROM Follow WHERE followedId = ? AND followerId = ?) AS following " +
		"FROM UserInfo WHERE id = ?"
	err := dao.Db.QueryStructRow(user, query, id.Bytes(), searchAsId.Bytes(), id.Bytes(), searchAsId.Bytes(), id.Bytes())
	switch {
	case err == database.ErrNoResult:
		return nil, nil
	case err != nil:
		return nil, err
	default:
		return user, nil
	}
}

func (dao DbDao) EditUser(userUuid uuid.UUID, user ModelUser) error {
	query := "UPDATE User SET name = ?, surname = ?, username = ? WHERE id = ?"
	return dao.Db.Exec(query, user.Name, user.Surname, user.Username, userUuid.Bytes())
}

func (dao DbDao) EditUsername(userUuid uuid.UUID, username string) error {
	query := "UPDATE User SET username = ? WHERE id = ?"
	return dao.Db.Exec(query, username, userUuid.Bytes())
}
