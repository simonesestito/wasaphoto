package comments

import (
	"github.com/gofrs/uuid"
	"github.com/simonesestito/wasaphoto/service/database"
	"github.com/simonesestito/wasaphoto/service/features/user"
	"github.com/simonesestito/wasaphoto/service/timeprovider"
	"github.com/simonesestito/wasaphoto/service/utils/cursor"
)

type EntityComment struct {
	Id          []byte `json:"id"`
	Text        string `json:"text"`
	PublishDate string `json:"publishDate"`
	AuthorId    []byte `json:"authorId"`
	PhotoId     []byte `json:"photoId"`
}

type EntityCommentWithAuthor struct {
	user.ModelUserInfo
	EntityComment
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
	EntityCommentWithAuthor
}

func (entity EntityCommentWithCustom) ToDto() Comment {
	publishDate, _ := timeprovider.UTCStringToDate(entity.PublishDate)
	entity.ModelUserWithCustom.ModelUser.Id = entity.EntityComment.AuthorId
	return Comment{
		Id:          uuid.FromBytesOrNil(entity.EntityComment.Id).String(),
		PublishDate: publishDate,
		Author:      entity.ModelUserWithCustom.ToDto(),
		NewComment: NewComment{
			Text: entity.EntityComment.Text,
		},
	}
}

func DbCommentsListToPage(dbComments []EntityCommentWithCustom) (comments []Comment, pageCursor *string) {
	comments = make([]Comment, len(dbComments))
	for i, dbComment := range dbComments {
		comments[i] = dbComment.ToDto()
	}

	// Calculate next cursor
	if len(dbComments) == database.MaxPageItems {
		lastComment := dbComments[len(dbComments)-1]
		nextCursor := cursor.CreateDateIdCursor(lastComment.EntityComment.Id, lastComment.PublishDate)
		pageCursor = &nextCursor
	} else {
		pageCursor = nil
	}

	return
}
