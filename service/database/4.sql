DROP VIEW IF EXISTS CommentWithAuthor;
CREATE VIEW CommentWithAuthor AS
SELECT Comment.*,
	   UserInfo.name,
	   UserInfo.surname,
	   UserInfo.username,
	   UserInfo.followersCount,
	   UserInfo.followingsCount,
	   UserInfo.photosCount
FROM Comment
		 LEFT JOIN UserInfo ON UserInfo.id = Comment.authorId;

DROP VIEW IF EXISTS CommentIdWithAuthorAndPhoto;
CREATE VIEW CommentIdWithAuthorAndPhoto AS
SELECT Comment.id       AS commentId,
	   Comment.authorId AS commentAuthorId,
	   Photo.id         AS photoId,
	   Photo.authorId   AS photoAuthorId
FROM Comment
		 LEFT JOIN Photo on Comment.photoId = Photo.id;

