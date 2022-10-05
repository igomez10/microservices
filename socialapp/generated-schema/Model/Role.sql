--
-- Socialapp.
-- Prepared SQL queries for 'Role' definition.
--


--
-- SELECT template for table `role`
--
SELECT `id`, `name`, `description`, `created_at` FROM `role` WHERE 1;

--
-- INSERT template for table `role`
--
INSERT INTO `role`(`id`, `name`, `description`, `created_at`) VALUES (:id, :name, :description, :created_at);

--
-- UPDATE template for table `role`
--
UPDATE `role` SET `id` = :id, `name` = :name, `description` = :description, `created_at` = :created_at WHERE 1;

--
-- DELETE template for table `role`
--
DELETE FROM `role` WHERE 0;

