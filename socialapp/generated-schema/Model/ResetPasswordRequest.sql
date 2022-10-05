--
-- Socialapp.
-- Prepared SQL queries for 'ResetPasswordRequest' definition.
--


--
-- SELECT template for table `reset_password_request`
--
SELECT `email` FROM `reset_password_request` WHERE 1;

--
-- INSERT template for table `reset_password_request`
--
INSERT INTO `reset_password_request`(`email`) VALUES (:email);

--
-- UPDATE template for table `reset_password_request`
--
UPDATE `reset_password_request` SET `email` = :email WHERE 1;

--
-- DELETE template for table `reset_password_request`
--
DELETE FROM `reset_password_request` WHERE 0;

