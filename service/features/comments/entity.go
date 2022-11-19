package comments

import (
	"github.com/gofrs/uuid"
	"github.com/simonesestito/wasaphoto/service/features/user"
	"github.com/simonesestito/wasaphoto/service/utils"
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
	publishDate, _ := utils.UTCStringToDate(entity.PublishDate)
	return Comment{
		Id:          uuid.FromBytesOrNil(entity.EntityComment.Id).String(),
		PublishDate: publishDate,
		Author:      entity.ModelUserWithCustom.ToDto(),
		NewComment: NewComment{
			Text: entity.EntityComment.Text,
		},
	}
}
