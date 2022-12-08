package user

import (
	"database/sql"
	"errors"
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
	GetUserByUsernameAs(username string, searchAsId uuid.UUID) (*ModelUserWithCustom, error)
	ListUsersByUsernameAs(username string, searchAsId uuid.UUID, afterUsername string, afterId uuid.UUID) ([]ModelUserWithCustom, error)
}

type DbDao struct {
	Db database.AppDatabase
}

func (dao DbDao) GetUserById(id uuid.UUID) (*ModelUser, error) {
	user := &ModelUser{}
	err := dao.Db.QueryStructRow(user, "SELECT * FROM User WHERE id = ?", id.Bytes())
	switch {
	case errors.Is(err, database.ErrNoResult):
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
	defer func() {
		_ = tx.Rollback()
	}()

	// Get user if exists
	existingUserId := make([]byte, 16)
	result, err := tx.Query("SELECT id FROM User WHERE username = ?", user.Username)
	if err != nil && result != nil {
		defer result.Close()
	}

	if !errors.Is(err, sql.ErrNoRows) && err != nil {
		// Unexpected error
		return uuid.Nil, false, err
	}

	if errors.Is(err, sql.ErrNoRows) || !result.Next() {
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
	case errors.Is(err, database.ErrNoResult):
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

func (dao DbDao) GetUserByUsernameAs(username string, searchAsId uuid.UUID) (*ModelUserWithCustom, error) {
	query := `
		SELECT UserInfo.*,
		       EXISTS(SELECT * FROM Ban WHERE bannedId = UserInfo.id AND bannerId = ?) AS banned,
		       EXISTS(SELECT * FROM Follow WHERE followedId = UserInfo.id AND followerId = ?) AS following
		FROM UserInfo
		WHERE UserInfo.username = ?`

	row := &ModelUserWithCustom{}
	err := dao.Db.QueryStructRow(
		row,
		query,
		searchAsId.Bytes(),
		searchAsId.Bytes(),
		username,
	)

	if errors.Is(err, database.ErrNoResult) {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return row, nil
	}
}

func (dao DbDao) ListUsersByUsernameAs(username string, searchAsId uuid.UUID, afterUsername string, afterId uuid.UUID) ([]ModelUserWithCustom, error) {
	query := `
		SELECT UserInfo.*,
		       EXISTS(SELECT * FROM Ban WHERE bannedId = UserInfo.id AND bannerId = ?) AS banned,
		       EXISTS(SELECT * FROM Follow WHERE followedId = UserInfo.id AND followerId = ?) AS following
		FROM UserInfo
		WHERE username LIKE ?
		 	  -- Cursor pagination
			  AND (username, id) > (?, ?)
			  AND NOT EXISTS(SELECT * FROM Ban WHERE Ban.bannerId = UserInfo.id AND Ban.bannedId = ?)
		ORDER BY username, id
		LIMIT ?`

	rows, err := dao.Db.QueryStructRows(
		ModelUserWithCustom{},
		query,
		searchAsId.Bytes(),
		searchAsId.Bytes(),
		"%"+username+"%",
		afterUsername,
		afterId.Bytes(),
		searchAsId.Bytes(),
		database.MaxPageItems,
	)

	if err != nil {
		return nil, err
	}

	return ParseUserEntities(rows)
}

func ParseUserEntities(rows database.StructRows) ([]ModelUserWithCustom, error) {
	var (
		users  []ModelUserWithCustom
		entity any
		err    error
	)

	for entity, err = rows.Next(); err == nil; entity, err = rows.Next() {
		newUser, ok := entity.(ModelUserWithCustom)
		if ok {
			users = append(users, newUser)
		} else {
			return nil, errors.New("invalid cast from db map to application entity")
		}
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	return users, nil
}
