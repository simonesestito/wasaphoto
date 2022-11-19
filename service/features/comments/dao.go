package comments

import (
	"database/sql"
	"github.com/gofrs/uuid"
	"github.com/simonesestito/wasaphoto/service/database"
)

type Dao interface {
	CreateComment(newComment EntityComment) error
	GetCommentByIdAs(commentId uuid.UUID, userId uuid.UUID) (*EntityCommentWithCustom, error)
}

type DbDao struct {
	Db database.AppDatabase
}

func (db DbDao) CreateComment(newComment EntityComment) error {
	return db.Db.Exec("INSERT INTO Comment (id, `text`, publishDate, authorId, photoId) VALUES (?, ?, ?, ?, ?)",
		newComment.Id,
		newComment.Text,
		newComment.PublishDate,
		newComment.AuthorId,
		newComment.PhotoId,
	)
}

func (db DbDao) GetCommentByIdAs(commentId uuid.UUID, userId uuid.UUID) (*EntityCommentWithCustom, error) {
	entity := &EntityCommentWithCustom{}

	query := `
SELECT CommentWithAuthor.*,
       EXISTS(SELECT * FROM Ban WHERE bannedId = CommentWithAuthor.authorId AND bannerId = ?) AS banned,
       EXISTS(SELECT * FROM Follow WHERE followedId = CommentWithAuthor.authorId AND followerId = ?) AS following
FROM CommentWithAuthor
WHERE CommentWithAuthor.id = ?`

	err := db.Db.QueryStructRow(entity, query, userId.Bytes(), userId.Bytes(), commentId.Bytes())
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return entity, nil
	}
}
