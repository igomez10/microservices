--
-- Socialapp.
-- Prepared SQL queries for 'ChangePasswordRequest' definition.
--


--
-- SELECT template for table `change_password_request`
--
SELECT `old_password`, `new_password` FROM `change_password_request` WHERE 1;

--
-- INSERT template for table `change_password_request`
--
INSERT INTO `change_password_request`(`old_password`, `new_password`) VALUES (:old_password, :new_password);

--
-- UPDATE template for table `change_password_request`
--
UPDATE `change_password_request` SET `old_password` = :old_password, `new_password` = :new_password WHERE 1;

--
-- DELETE template for table `change_password_request`
--
DELETE FROM `change_password_request` WHERE 0;

