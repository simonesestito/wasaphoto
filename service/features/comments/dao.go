package comments

import (
	"database/sql"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/simonesestito/wasaphoto/service/database"
)

type Dao interface {
	CreateComment(newComment entityComment) error
	GetCommentByIdAs(commentId uuid.UUID, userId uuid.UUID) (*EntityCommentWithCustom, error)
	DeleteByIdPhotoAndAuthor(commentUuid uuid.UUID, photoUuid uuid.UUID, userUuid uuid.UUID) (bool, error)
	GetCommentInfoIds(commentUuid uuid.UUID) (*CommentIdWithAuthorAndPhoto, error)
	GetCommentsAfter(photoUuid uuid.UUID, userUuid uuid.UUID, afterComment uuid.UUID, beforeDate string) ([]EntityCommentWithCustom, error)
}

type DbDao struct {
	Db database.AppDatabase
}

func (db DbDao) CreateComment(newComment entityComment) error {
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

	// Fix shadowed properties
	entity.ModelUserWithCustom.ModelUser.Id = entity.entityComment.AuthorId

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return entity, nil
	}
}

func (db DbDao) DeleteByIdPhotoAndAuthor(commentUuid uuid.UUID, photoUuid uuid.UUID, userUuid uuid.UUID) (bool, error) {
	rows, err := db.Db.ExecRows(
		"DELETE FROM Comment WHERE id = ? AND photoId = ? AND authorId = ?",
		commentUuid.Bytes(),
		photoUuid.Bytes(),
		userUuid.Bytes(),
	)
	return rows > 0, err
}

func (db DbDao) GetCommentInfoIds(commentUuid uuid.UUID) (*CommentIdWithAuthorAndPhoto, error) {
	entity := &CommentIdWithAuthorAndPhoto{}
	err := db.Db.QueryStructRow(entity, "SELECT * FROM CommentIdWithAuthorAndPhoto WHERE commentId = ?", commentUuid.Bytes())
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return entity, nil
	}
}

func (db DbDao) GetCommentsAfter(photoUuid uuid.UUID, userUuid uuid.UUID, afterComment uuid.UUID, beforeDate string) ([]EntityCommentWithCustom, error) {
	query := `
		SELECT CommentWithAuthor.*,
			   EXISTS(SELECT * FROM Ban WHERE bannedId = CommentWithAuthor.authorId AND bannerId = ?) AS banned,
			   EXISTS(SELECT * FROM Follow WHERE followedId = CommentWithAuthor.authorId AND followerId = ?) AS following
		FROM CommentWithAuthor
		WHERE CommentWithAuthor.photoId = ?
		 	  -- Cursor pagination
			  AND (publishDate, id) < (?, ?)
			  -- Hide comments from users who banned me
			  AND NOT EXISTS(SELECT * FROM Ban WHERE bannedId = ? AND bannerId = CommentWithAuthor.authorId)
		ORDER BY publishDate DESC, id DESC
		LIMIT ?`

	rows, err := db.Db.QueryStructRows(
		EntityCommentWithCustom{},
		query,
		userUuid.Bytes(),
		userUuid.Bytes(),
		photoUuid.Bytes(),
		beforeDate,
		afterComment.Bytes(),
		userUuid.Bytes(),
		database.MaxPageItems,
	)

	if err != nil {
		return nil, err
	}

	var comments []EntityCommentWithCustom
	var entity any
	for entity, err = rows.Next(); err == nil; entity, err = rows.Next() {
		newComment, ok := entity.(EntityCommentWithCustom)
		if ok {
			comments = append(comments, newComment)
		} else {
			return nil, errors.New("invalid cast from db map to application entity")
		}
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	return comments, nil
}
