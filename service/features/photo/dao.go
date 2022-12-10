package photo

import (
	"database/sql"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/simonesestito/wasaphoto/service/database"
	"github.com/simonesestito/wasaphoto/service/timeprovider"
)

type Dao interface {
	GetPhotoByIdAs(photoId uuid.UUID, userId uuid.UUID) (*EntityPhotoAuthorInfo, error)
	NewPhotoPerUser(photoId uuid.UUID, userId uuid.UUID, imageUrl string) error
	DeletePhoto(imageUuid uuid.UUID) error
	GetPhotoById(imageUuid uuid.UUID) (*EntityPhotoInfo, error)
	ListUsersPhotoAfter(authorUuid uuid.UUID, searchAsUuid uuid.UUID, afterPhotoId uuid.UUID, beforeDate string) ([]EntityPhotoAuthorInfo, error)
}

type DbDao struct {
	Time timeprovider.TimeProvider
	Db   database.AppDatabase
}

func (db DbDao) GetPhotoByIdAs(photoId uuid.UUID, userId uuid.UUID) (*EntityPhotoAuthorInfo, error) {
	photo := EntityPhotoAuthorInfo{}
	query := `
SELECT P.*,
EXISTS(SELECT * FROM Ban B WHERE B.bannedId = P.authorId AND B.bannerId = ?) AS banned,
EXISTS(SELECT * FROM Follow F WHERE F.followedId = P.authorId AND F.followerId = ?) AS following,
EXISTS(SELECT * FROM Likes L WHERE L.photoId = P.authorId AND L.userId = ?) AS liked
FROM PhotoAuthorInfo P
WHERE P.id = ?
`
	err := db.Db.QueryStructRow(&photo, query, userId.Bytes(), userId.Bytes(), userId.Bytes(), photoId.Bytes())

	// Fix shadowed properties
	photo.ModelUser.Id = photo.entityPhoto.AuthorId

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return &photo, err
}

func (db DbDao) NewPhotoPerUser(photoId uuid.UUID, userId uuid.UUID, imageUrl string) error {
	currentTime := db.Time.UTCString()
	return db.Db.Exec("INSERT INTO Photo (id, imageUrl, authorId, publishDate) VALUES (?, ?, ?, ?)", photoId.Bytes(), imageUrl, userId.Bytes(), currentTime)
}

func (db DbDao) DeletePhoto(imageUuid uuid.UUID) error {
	err := db.Db.Exec("DELETE FROM Photo WHERE id = ?", imageUuid.Bytes())
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	} else {
		return err
	}
}

func (db DbDao) GetPhotoById(imageUuid uuid.UUID) (*EntityPhotoInfo, error) {
	photo := EntityPhotoInfo{}
	query := "SELECT * FROM PhotoInfo WHERE id = ?"
	err := db.Db.QueryStructRow(&photo, query, imageUuid.Bytes())
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &photo, err
}

func (db DbDao) ListUsersPhotoAfter(authorUuid uuid.UUID, searchAsUuid uuid.UUID, afterPhotoId uuid.UUID, beforeDate string) ([]EntityPhotoAuthorInfo, error) {
	query := `
		SELECT PhotoAuthorInfo.*,
		       EXISTS(SELECT * FROM Likes WHERE Likes.photoId = PhotoAuthorInfo.id AND Likes.userId = ?) AS liked,
		       EXISTS(SELECT * FROM Ban WHERE bannedId = PhotoAuthorInfo.authorId AND bannerId = ?) AS banned,
		       EXISTS(SELECT * FROM Follow WHERE followedId = PhotoAuthorInfo.authorId AND followerId = ?) AS following
		FROM PhotoAuthorInfo
		WHERE PhotoAuthorInfo.authorId = ?
		 	  -- Cursor pagination
			  AND (publishDate, id) < (?, ?)
		ORDER BY publishDate DESC, id DESC
		LIMIT ?`

	rows, err := db.Db.QueryStructRows(
		EntityPhotoAuthorInfo{},
		query,
		searchAsUuid.Bytes(),
		searchAsUuid.Bytes(),
		searchAsUuid.Bytes(),
		authorUuid.Bytes(),
		beforeDate,
		afterPhotoId.Bytes(),
		database.MaxPageItems,
	)

	if err != nil {
		return nil, err
	}

	return ParsePhotoEntity(rows)
}

func ParsePhotoEntity(rows database.StructRows) ([]EntityPhotoAuthorInfo, error) {
	var (
		photos []EntityPhotoAuthorInfo
		entity any
		err    error
	)

	for entity, err = rows.Next(); err == nil; entity, err = rows.Next() {
		newPhoto, ok := entity.(EntityPhotoAuthorInfo)
		if ok {
			photos = append(photos, newPhoto)
		} else {
			return nil, errors.New("invalid cast from db map to application entity")
		}
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	return photos, nil
}
