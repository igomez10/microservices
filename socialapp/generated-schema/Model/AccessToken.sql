--
-- Socialapp.
-- Prepared SQL queries for 'AccessToken' definition.
--


--
-- SELECT template for table `access_token`
--
SELECT `access_token`, `token_type`, `scopes`, `expires_in` FROM `access_token` WHERE 1;

--
-- INSERT template for table `access_token`
--
INSERT INTO `access_token`(`access_token`, `token_type`, `scopes`, `expires_in`) VALUES (:access_token, :token_type, :scopes, :expires_in);

--
-- UPDATE template for table `access_token`
--
UPDATE `access_token` SET `access_token` = :access_token, `token_type` = :token_type, `scopes` = :scopes, `expires_in` = :expires_in WHERE 1;

--
-- DELETE template for table `access_token`
--
DELETE FROM `access_token` WHERE 0;

