--
-- URL Shortener.
-- Prepared SQL queries for 'Error' definition.
--


--
-- SELECT template for table `error`
--
SELECT `code`, `message` FROM `error` WHERE 1;

--
-- INSERT template for table `error`
--
INSERT INTO `error`(`code`, `message`) VALUES (:code, :message);

--
-- UPDATE template for table `error`
--
UPDATE `error` SET `code` = :code, `message` = :message WHERE 1;

--
-- DELETE template for table `error`
--
DELETE FROM `error` WHERE 0;

