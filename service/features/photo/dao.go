package photo

import (
	"database/sql"
	"github.com/gofrs/uuid"
	"github.com/simonesestito/wasaphoto/service/database"
	"github.com/simonesestito/wasaphoto/service/timeprovider"
)

type Dao interface {
	GetPhotoByIdAs(photoId uuid.UUID, userId uuid.UUID) (*EntityPhotoAuthorInfo, error)
	NewPhotoPerUser(photoId uuid.UUID, userId uuid.UUID, imageUrl string) error
	DeletePhoto(imageUuid uuid.UUID) error
	GetPhotoById(imageUuid uuid.UUID) (*EntityPhotoInfo, error)
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
	if err == sql.ErrNoRows {
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
	if err == sql.ErrNoRows {
		return nil
	} else {
		return err
	}
}

func (db DbDao) GetPhotoById(imageUuid uuid.UUID) (*EntityPhotoInfo, error) {
	photo := EntityPhotoInfo{}
	query := "SELECT * FROM PhotoInfo WHERE id = ?"
	err := db.Db.QueryStructRow(&photo, query, imageUuid.Bytes())
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &photo, err
}
