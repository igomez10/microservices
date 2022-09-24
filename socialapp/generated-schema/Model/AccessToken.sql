--
-- Socialapp.
-- Prepared SQL queries for 'AccessToken' definition.
--


--
-- SELECT template for table `AccessToken`
--
SELECT `access_token`, `token_type`, `scopes`, `expires_in` FROM `AccessToken` WHERE 1;

--
-- INSERT template for table `AccessToken`
--
INSERT INTO `AccessToken`(`access_token`, `token_type`, `scopes`, `expires_in`) VALUES (:access_token, :token_type, :scopes, :expires_in);

--
-- UPDATE template for table `AccessToken`
--
UPDATE `AccessToken` SET `access_token` = :access_token, `token_type` = :token_type, `scopes` = :scopes, `expires_in` = :expires_in WHERE 1;

--
-- DELETE template for table `AccessToken`
--
DELETE FROM `AccessToken` WHERE 0;

