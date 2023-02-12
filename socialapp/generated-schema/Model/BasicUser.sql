--
-- Socialapp.
-- Prepared SQL queries for 'BasicUser' definition.
--


--
-- SELECT template for table `basic_user`
--
SELECT `username`, `first_name`, `last_name`, `email` FROM `basic_user` WHERE 1;

--
-- INSERT template for table `basic_user`
--
INSERT INTO `basic_user`(`username`, `first_name`, `last_name`, `email`) VALUES (:username, :first_name, :last_name, :email);

--
-- UPDATE template for table `basic_user`
--
UPDATE `basic_user` SET `username` = :username, `first_name` = :first_name, `last_name` = :last_name, `email` = :email WHERE 1;

--
-- DELETE template for table `basic_user`
--
DELETE FROM `basic_user` WHERE 0;

