package photo

import (
	"github.com/gofrs/uuid"
	"github.com/simonesestito/wasaphoto/service/database"
	"github.com/simonesestito/wasaphoto/service/timeprovider"
)

type Dao interface {
	GetPhotoById(photoId uuid.UUID) (*EntityPhoto, error)
	NewPhotoPerUser(photoId uuid.UUID, userId uuid.UUID, imageUrl string) error
}

type DbDao struct {
	Time timeprovider.TimeProvider
	Db   database.AppDatabase
}

func (db DbDao) GetPhotoById(photoId uuid.UUID) (*EntityPhoto, error) {
	photo := EntityPhoto{}
	err := db.Db.QueryStructRow(&photo, "SELECT * FROM Photo WHERE id = ?", photoId.Bytes())
	return &photo, err
}

func (db DbDao) NewPhotoPerUser(photoId uuid.UUID, userId uuid.UUID, imageUrl string) error {
	currentTime := db.Time.UTCString()
	return db.Db.Exec("INSERT INTO Photo (id, imageUrl, authorId, publishDate) VALUES (?, ?, ?, ?)", photoId.Bytes(), imageUrl, userId.Bytes(), currentTime)
}
