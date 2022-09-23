--
-- Socialapp.
-- Prepared SQL queries for 'ResetPasswordRequest' definition.
--


--
-- SELECT template for table `ResetPasswordRequest`
--
SELECT `email` FROM `ResetPasswordRequest` WHERE 1;

--
-- INSERT template for table `ResetPasswordRequest`
--
INSERT INTO `ResetPasswordRequest`(`email`) VALUES (:email);

--
-- UPDATE template for table `ResetPasswordRequest`
--
UPDATE `ResetPasswordRequest` SET `email` = :email WHERE 1;

--
-- DELETE template for table `ResetPasswordRequest`
--
DELETE FROM `ResetPasswordRequest` WHERE 0;

