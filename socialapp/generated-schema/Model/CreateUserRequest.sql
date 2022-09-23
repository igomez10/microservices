--
-- Socialapp.
-- Prepared SQL queries for 'CreateUserRequest' definition.
--


--
-- SELECT template for table `CreateUserRequest`
--
SELECT `id`, `username`, `password`, `first_name`, `last_name`, `email`, `created_at` FROM `CreateUserRequest` WHERE 1;

--
-- INSERT template for table `CreateUserRequest`
--
INSERT INTO `CreateUserRequest`(`id`, `username`, `password`, `first_name`, `last_name`, `email`, `created_at`) VALUES (:id, :username, :password, :first_name, :last_name, :email, :created_at);

--
-- UPDATE template for table `CreateUserRequest`
--
UPDATE `CreateUserRequest` SET `id` = :id, `username` = :username, `password` = :password, `first_name` = :first_name, `last_name` = :last_name, `email` = :email, `created_at` = :created_at WHERE 1;

--
-- DELETE template for table `CreateUserRequest`
--
DELETE FROM `CreateUserRequest` WHERE 0;

