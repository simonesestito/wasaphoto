package comments

import (
	"github.com/gofrs/uuid"
	"github.com/simonesestito/wasaphoto/service/database"
	"github.com/simonesestito/wasaphoto/service/features/user"
	"github.com/simonesestito/wasaphoto/service/timeprovider"
	"github.com/simonesestito/wasaphoto/service/utils/cursor"
)

type entityComment struct {
	Id          []byte `json:"id"`
	Text        string `json:"text"`
	PublishDate string `json:"publishDate"`
	AuthorId    []byte `json:"authorId"`
	PhotoId     []byte `json:"photoId"`
}

type entityCommentWithAuthor struct {
	user.ModelUserInfo
	entityComment
}

// CommentIdWithAuthorAndPhoto is a simple view with just the IDs
type CommentIdWithAuthorAndPhoto struct {
	CommentId       []byte `json:"commentId"`
	CommentAuthorId []byte `json:"commentAuthorId"`
	PhotoId         []byte `json:"photoId"`
	PhotoAuthorId   []byte `json:"photoAuthorId"`
}

type EntityCommentWithCustom struct {
	user.ModelUserWithCustom
	entityCommentWithAuthor
}

func (entity EntityCommentWithCustom) toDto() Comment {
	publishDate, _ := timeprovider.UTCStringToDate(entity.PublishDate)
	entity.ModelUserWithCustom.ModelUser.Id = entity.entityComment.AuthorId
	return Comment{
		Id:          uuid.FromBytesOrNil(entity.entityComment.Id).String(),
		PublishDate: publishDate,
		Author:      entity.ModelUserWithCustom.ToDto(),
		newComment: newComment{
			Text: entity.entityComment.Text,
		},
	}
}

func dbCommentsListToPage(dbComments []EntityCommentWithCustom) (comments []Comment, pageCursor *string) {
	comments = make([]Comment, len(dbComments))
	for i, dbComment := range dbComments {
		comments[i] = dbComment.toDto()
	}

	// Calculate next cursor
	if len(dbComments) == database.MaxPageItems {
		lastComment := dbComments[len(dbComments)-1]
		nextCursor := cursor.CreateDateIdCursor(lastComment.entityComment.Id, lastComment.PublishDate)
		pageCursor = &nextCursor
	} else {
		pageCursor = nil
	}

	return
}
