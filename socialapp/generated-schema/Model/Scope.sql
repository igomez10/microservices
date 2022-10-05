--
-- Socialapp.
-- Prepared SQL queries for 'Scope' definition.
--


--
-- SELECT template for table `scope`
--
SELECT `id`, `name`, `description`, `created_at` FROM `scope` WHERE 1;

--
-- INSERT template for table `scope`
--
INSERT INTO `scope`(`id`, `name`, `description`, `created_at`) VALUES (:id, :name, :description, :created_at);

--
-- UPDATE template for table `scope`
--
UPDATE `scope` SET `id` = :id, `name` = :name, `description` = :description, `created_at` = :created_at WHERE 1;

--
-- DELETE template for table `scope`
--
DELETE FROM `scope` WHERE 0;

