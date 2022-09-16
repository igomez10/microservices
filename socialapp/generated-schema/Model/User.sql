--
-- Socialapp.
-- Prepared SQL queries for 'User' definition.
--


--
-- SELECT template for table `User`
--
SELECT `id`, `username`, `first_name`, `last_name`, `email`, `created_at` FROM `User` WHERE 1;

--
-- INSERT template for table `User`
--
INSERT INTO `User`(`id`, `username`, `first_name`, `last_name`, `email`, `created_at`) VALUES (:id, :username, :first_name, :last_name, :email, :created_at);

--
-- UPDATE template for table `User`
--
UPDATE `User` SET `id` = :id, `username` = :username, `first_name` = :first_name, `last_name` = :last_name, `email` = :email, `created_at` = :created_at WHERE 1;

--
-- DELETE template for table `User`
--
DELETE FROM `User` WHERE 0;

