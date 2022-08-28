--
-- Socialapp.
-- Prepared SQL queries for 'Error' definition.
--


--
-- SELECT template for table `Error`
--
SELECT `code`, `message` FROM `Error` WHERE 1;

--
-- INSERT template for table `Error`
--
INSERT INTO `Error`(`code`, `message`) VALUES (:code, :message);

--
-- UPDATE template for table `Error`
--
UPDATE `Error` SET `code` = :code, `message` = :message WHERE 1;

--
-- DELETE template for table `Error`
--
DELETE FROM `Error` WHERE 0;

