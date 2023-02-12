--
-- Socialapp.
-- Prepared SQL queries for 'CreateUserRequest' definition.
--


--
-- SELECT template for table `create_user_request`
--
SELECT `username`, `password`, `first_name`, `last_name`, `email` FROM `create_user_request` WHERE 1;

--
-- INSERT template for table `create_user_request`
--
INSERT INTO `create_user_request`(`username`, `password`, `first_name`, `last_name`, `email`) VALUES (:username, :password, :first_name, :last_name, :email);

--
-- UPDATE template for table `create_user_request`
--
UPDATE `create_user_request` SET `username` = :username, `password` = :password, `first_name` = :first_name, `last_name` = :last_name, `email` = :email WHERE 1;

--
-- DELETE template for table `create_user_request`
--
DELETE FROM `create_user_request` WHERE 0;

