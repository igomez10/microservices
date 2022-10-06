--
-- Socialapp.
-- Prepared SQL queries for 'User' definition.
--


--
-- SELECT template for table `user`
--
SELECT `id`, `username`, `first_name`, `last_name`, `email`, `created_at`, `roles` FROM `user` WHERE 1;

--
-- INSERT template for table `user`
--
INSERT INTO `user`(`id`, `username`, `first_name`, `last_name`, `email`, `created_at`, `roles`) VALUES (:id, :username, :first_name, :last_name, :email, :created_at, :roles);

--
-- UPDATE template for table `user`
--
UPDATE `user` SET `id` = :id, `username` = :username, `first_name` = :first_name, `last_name` = :last_name, `email` = :email, `created_at` = :created_at, `roles` = :roles WHERE 1;

--
-- DELETE template for table `user`
--
DELETE FROM `user` WHERE 0;

