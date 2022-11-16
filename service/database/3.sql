-- Count photo likes
DROP VIEW IF EXISTS PhotoLikes;
CREATE VIEW PhotoLikes AS
SELECT Photo.id AS photoId, COALESCE(COUNT(Likes.userId), 0) AS likesCount
FROM Photo
		 LEFT JOIN Likes on Photo.id = Likes.photoId
GROUP BY Photo.id;

-- Count photo comments
DROP VIEW IF EXISTS PhotoComments;
CREATE VIEW PhotoComments AS
SELECT Photo.id AS photoId, COALESCE(COUNT(Comment.id), 0) AS commentsCount
FROM Photo
		 LEFT JOIN Comment on Photo.id = Comment.photoId
GROUP BY Photo.id;

--
-- Fetch aggregate photo data
--
DROP VIEW IF EXISTS PhotoInfo;
CREATE VIEW PhotoInfo AS
SELECT Photo.*,
	   PhotoLikes.likesCount,
	   PhotoComments.commentsCount
FROM Photo
		 LEFT JOIN PhotoLikes ON Photo.id = PhotoLikes.photoId
		 LEFT JOIN PhotoComments ON Photo.id = PhotoComments.photoId;

--
-- Fetch photo data with its author
--
DROP VIEW IF EXISTS PhotoAuthorInfo;
CREATE VIEW PhotoAuthorInfo AS
SELECT PhotoInfo.*,
	   U.name,
	   U.surname,
	   U.username,
	   U.followersCount,
	   U.followingsCount,
	   U.photosCount
FROM PhotoInfo
		 LEFT JOIN UserInfo U on PhotoInfo.authorId = U.id;
