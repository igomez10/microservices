--
-- Socialapp.
-- Prepared SQL queries for 'CreateCommentRequest' definition.
--


--
-- SELECT template for table `create_comment_request`
--
SELECT `content`, `username` FROM `create_comment_request` WHERE 1;

--
-- INSERT template for table `create_comment_request`
--
INSERT INTO `create_comment_request`(`content`, `username`) VALUES (:content, :username);

--
-- UPDATE template for table `create_comment_request`
--
UPDATE `create_comment_request` SET `content` = :content, `username` = :username WHERE 1;

--
-- DELETE template for table `create_comment_request`
--
DELETE FROM `create_comment_request` WHERE 0;

