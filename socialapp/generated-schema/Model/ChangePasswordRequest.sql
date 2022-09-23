--
-- Socialapp.
-- Prepared SQL queries for 'ChangePasswordRequest' definition.
--


--
-- SELECT template for table `ChangePasswordRequest`
--
SELECT `old_password`, `new_password` FROM `ChangePasswordRequest` WHERE 1;

--
-- INSERT template for table `ChangePasswordRequest`
--
INSERT INTO `ChangePasswordRequest`(`old_password`, `new_password`) VALUES (:old_password, :new_password);

--
-- UPDATE template for table `ChangePasswordRequest`
--
UPDATE `ChangePasswordRequest` SET `old_password` = :old_password, `new_password` = :new_password WHERE 1;

--
-- DELETE template for table `ChangePasswordRequest`
--
DELETE FROM `ChangePasswordRequest` WHERE 0;

