--
-- Socialapp.
-- Prepared SQL queries for 'LikeRequest' definition.
--


--
-- SELECT template for table `like_request`
--
SELECT `comment_id` FROM `like_request` WHERE 1;

--
-- INSERT template for table `like_request`
--
INSERT INTO `like_request`(`comment_id`) VALUES (:comment_id);

--
-- UPDATE template for table `like_request`
--
UPDATE `like_request` SET `comment_id` = :comment_id WHERE 1;

--
-- DELETE template for table `like_request`
--
DELETE FROM `like_request` WHERE 0;

