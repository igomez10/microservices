--
-- Socialapp.
-- Prepared SQL queries for 'AccessToken' definition.
--


--
-- SELECT template for table `AccessToken`
--
SELECT `access_token`, `scopes`, `expires_at` FROM `AccessToken` WHERE 1;

--
-- INSERT template for table `AccessToken`
--
INSERT INTO `AccessToken`(`access_token`, `scopes`, `expires_at`) VALUES (:access_token, :scopes, :expires_at);

--
-- UPDATE template for table `AccessToken`
--
UPDATE `AccessToken` SET `access_token` = :access_token, `scopes` = :scopes, `expires_at` = :expires_at WHERE 1;

--
-- DELETE template for table `AccessToken`
--
DELETE FROM `AccessToken` WHERE 0;

