--
-- Socialapp.
-- Prepared SQL queries for 'CreateUserRequest' definition.
--


--
-- SELECT template for table `create_user_request`
--
SELECT `id`, `username`, `password`, `first_name`, `last_name`, `email`, `created_at` FROM `create_user_request` WHERE 1;

--
-- INSERT template for table `create_user_request`
--
INSERT INTO `create_user_request`(`id`, `username`, `password`, `first_name`, `last_name`, `email`, `created_at`) VALUES (:id, :username, :password, :first_name, :last_name, :email, :created_at);

--
-- UPDATE template for table `create_user_request`
--
UPDATE `create_user_request` SET `id` = :id, `username` = :username, `password` = :password, `first_name` = :first_name, `last_name` = :last_name, `email` = :email, `created_at` = :created_at WHERE 1;

--
-- DELETE template for table `create_user_request`
--
DELETE FROM `create_user_request` WHERE 0;

