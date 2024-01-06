--
-- URL Shortener.
-- Prepared SQL queries for 'URL' definition.
--


--
-- SELECT template for table `url`
--
SELECT `url`, `alias`, `created_at`, `updated_at`, `deleted_at` FROM `url` WHERE 1;

--
-- INSERT template for table `url`
--
INSERT INTO `url`(`url`, `alias`, `created_at`, `updated_at`, `deleted_at`) VALUES (:url, :alias, :created_at, :updated_at, :deleted_at);

--
-- UPDATE template for table `url`
--
UPDATE `url` SET `url` = :url, `alias` = :alias, `created_at` = :created_at, `updated_at` = :updated_at, `deleted_at` = :deleted_at WHERE 1;

--
-- DELETE template for table `url`
--
DELETE FROM `url` WHERE 0;

