--
-- Socialapp.
-- Prepared SQL queries for 'Comment' definition.
--


--
-- SELECT template for table `Comment`
--
SELECT `id`, `content`, `like_count`, `created_at`, `username` FROM `Comment` WHERE 1;

--
-- INSERT template for table `Comment`
--
INSERT INTO `Comment`(`id`, `content`, `like_count`, `created_at`, `username`) VALUES (:id, :content, :like_count, :created_at, :username);

--
-- UPDATE template for table `Comment`
--
UPDATE `Comment` SET `id` = :id, `content` = :content, `like_count` = :like_count, `created_at` = :created_at, `username` = :username WHERE 1;

--
-- DELETE template for table `Comment`
--
DELETE FROM `Comment` WHERE 0;

