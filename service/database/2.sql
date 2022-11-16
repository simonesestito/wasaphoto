--
-- SQL Schema version 1
--

-- Count of the followers of each user
DROP VIEW IF EXISTS Followers;
CREATE VIEW Followers AS
SELECT User.id AS followedId, COALESCE(COUNT(Follow.followerId), 0) AS followersCount
FROM User
		 LEFT JOIN Follow ON User.id = Follow.followedId
GROUP BY User.id;

-- Count of the followings of each user
DROP VIEW IF EXISTS Followings;
CREATE VIEW Followings AS
SELECT User.id AS followerId, COALESCE(COUNT(Follow.followedId), 0) AS followingsCount
FROM User
		 LEFT JOIN Follow on User.id = Follow.followerId
GROUP BY User.id;

-- Count the posts of each user
DROP VIEW IF EXISTS UserPhotosCount;
CREATE VIEW UserPhotosCount AS
SELECT User.id AS authorId, COALESCE(COUNT(Photo.id), 0) AS photosCount
FROM User
		 LEFT JOIN Photo on User.id = Photo.authorId
GROUP BY User.id;

-- Aggregate all useful information about a user
DROP VIEW IF EXISTS UserInfo;
CREATE VIEW UserInfo AS
SELECT User.*, followersCount, followingsCount, photosCount
FROM User
		 LEFT JOIN Followers ON Followers.followedId = User.id
		 LEFT JOIN Followings ON Followings.followerId = User.id
		 LEFT JOIN UserPhotosCount ON UserPhotosCount.authorId = User.id;
