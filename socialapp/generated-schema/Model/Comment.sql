--
-- Socialapp.
-- Prepared SQL queries for 'Comment' definition.
--


--
-- SELECT template for table `comment`
--
SELECT `id`, `content`, `like_count`, `created_at`, `username` FROM `comment` WHERE 1;

--
-- INSERT template for table `comment`
--
INSERT INTO `comment`(`id`, `content`, `like_count`, `created_at`, `username`) VALUES (:id, :content, :like_count, :created_at, :username);

--
-- UPDATE template for table `comment`
--
UPDATE `comment` SET `id` = :id, `content` = :content, `like_count` = :like_count, `created_at` = :created_at, `username` = :username WHERE 1;

--
-- DELETE template for table `comment`
--
DELETE FROM `comment` WHERE 0;

