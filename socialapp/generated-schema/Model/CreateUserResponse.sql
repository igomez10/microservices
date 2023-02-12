--
-- Socialapp.
-- Prepared SQL queries for 'CreateUserResponse' definition.
--


--
-- SELECT template for table `create_user_response`
--
SELECT `id`, `username`, `first_name`, `last_name`, `email`, `created_at` FROM `create_user_response` WHERE 1;

--
-- INSERT template for table `create_user_response`
--
INSERT INTO `create_user_response`(`id`, `username`, `first_name`, `last_name`, `email`, `created_at`) VALUES (:id, :username, :first_name, :last_name, :email, :created_at);

--
-- UPDATE template for table `create_user_response`
--
UPDATE `create_user_response` SET `id` = :id, `username` = :username, `first_name` = :first_name, `last_name` = :last_name, `email` = :email, `created_at` = :created_at WHERE 1;

--
-- DELETE template for table `create_user_response`
--
DELETE FROM `create_user_response` WHERE 0;

