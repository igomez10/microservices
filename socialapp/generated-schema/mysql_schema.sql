/* SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO"; */
/* SET AUTOCOMMIT = 0; */
/* START TRANSACTION; */
/* SET time_zone = "+00:00"; */

-- --------------------------------------------------------

--
-- Table structure for table `access_token` generated from model 'AccessToken'
-- OAuth access token response.
--

CREATE TABLE IF NOT EXISTS `access_token` (
  `access_token` TEXT NOT NULL COMMENT 'Access token used for authenticated requests.',
  `token_type` TEXT NOT NULL COMMENT 'Type of token returned by the authorization server.',
  `scopes` JSON DEFAULT NULL COMMENT 'Granted OAuth scopes.',
  `expires_in` INT NOT NULL COMMENT 'Access token lifetime in seconds.'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='OAuth access token response.. Original model name - AccessToken.';

--
-- Table structure for table `change_password_request` generated from model 'ChangePasswordRequest'
-- Request payload for changing a user&#39;s password.
--

CREATE TABLE IF NOT EXISTS `change_password_request` (
  `old_password` TEXT NOT NULL COMMENT 'Current password.',
  `new_password` TEXT NOT NULL COMMENT 'New password.'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Request payload for changing a user&#39;s password.. Original model name - ChangePasswordRequest.';

--
-- Table structure for table `comment` generated from model 'Comment'
-- User-generated comment.
--

CREATE TABLE IF NOT EXISTS `comment` (
  `id` BIGINT DEFAULT NULL COMMENT 'Unique comment identifier.',
  `content` TEXT NOT NULL COMMENT 'Comment text.',
  `like_count` BIGINT DEFAULT NULL COMMENT 'Number of likes for the comment.',
  `created_at` DATETIME DEFAULT NULL COMMENT 'Timestamp when the comment was created.',
  `username` TEXT NOT NULL COMMENT 'Username of the comment author.'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='User-generated comment.. Original model name - Comment.';

--
-- Table structure for table `create_user_request` generated from model 'CreateUserRequest'
-- Request payload for creating a user account.
--

CREATE TABLE IF NOT EXISTS `create_user_request` (
  `username` TEXT NOT NULL COMMENT 'Unique username for the account.',
  `password` TEXT NOT NULL COMMENT 'Plain-text password used during account creation.',
  `first_name` TEXT NOT NULL COMMENT 'User first name.',
  `last_name` TEXT NOT NULL COMMENT 'User last name.',
  `email` TEXT NOT NULL COMMENT 'User email address.'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Request payload for creating a user account.. Original model name - CreateUserRequest.';

--
-- Table structure for table `create_user_response` generated from model 'CreateUserResponse'
-- Response payload returned after user creation.
--

CREATE TABLE IF NOT EXISTS `create_user_response` (
  `id` BIGINT NOT NULL COMMENT 'Unique user identifier.',
  `username` TEXT NOT NULL COMMENT 'Unique username for the account.',
  `first_name` TEXT NOT NULL COMMENT 'User first name.',
  `last_name` TEXT NOT NULL COMMENT 'User last name.',
  `email` TEXT NOT NULL COMMENT 'User email address.',
  `created_at` DATETIME NOT NULL COMMENT 'Timestamp when the user was created.'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Response payload returned after user creation.. Original model name - CreateUserResponse.';

--
-- Table structure for table `error` generated from model 'Error'
-- Standard error payload returned by the API.
--

CREATE TABLE IF NOT EXISTS `error` (
  `code` INT NOT NULL COMMENT 'HTTP-like error code.',
  `message` TEXT NOT NULL COMMENT 'Human-readable error message.'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Standard error payload returned by the API.. Original model name - Error.';

--
-- Table structure for table `reset_password_request` generated from model 'ResetPasswordRequest'
-- Request payload for triggering a password reset.
--

CREATE TABLE IF NOT EXISTS `reset_password_request` (
  `email` TEXT NOT NULL COMMENT 'Email address associated with the user account.'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Request payload for triggering a password reset.. Original model name - ResetPasswordRequest.';

--
-- Table structure for table `role` generated from model 'Role'
-- Role assigned to users for authorization.
--

CREATE TABLE IF NOT EXISTS `role` (
  `id` BIGINT DEFAULT NULL COMMENT 'Unique role identifier.',
  `name` TEXT NOT NULL COMMENT 'Unique role name.',
  `description` TEXT DEFAULT NULL COMMENT 'Description of what the role grants.',
  `created_at` DATETIME DEFAULT NULL COMMENT 'Timestamp when the role was created.'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Role assigned to users for authorization.. Original model name - Role.';

--
-- Table structure for table `scope` generated from model 'Scope'
-- OAuth scope used for access control.
--

CREATE TABLE IF NOT EXISTS `scope` (
  `id` BIGINT DEFAULT NULL COMMENT 'Unique scope identifier.',
  `name` TEXT NOT NULL COMMENT 'Scope name (for example, &#x60;socialapp.roles.read&#x60;).',
  `description` TEXT NOT NULL COMMENT 'Description of the permission granted by this scope.',
  `created_at` DATETIME DEFAULT NULL COMMENT 'Timestamp when the scope was created.'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='OAuth scope used for access control.. Original model name - Scope.';

--
-- Table structure for table `url` generated from model 'URL'
-- URL alias and metadata.
--

CREATE TABLE IF NOT EXISTS `url` (
  `url` TEXT NOT NULL COMMENT 'Destination URL.',
  `alias` TEXT NOT NULL COMMENT 'Short alias used to reference the URL.',
  `created_at` DATETIME DEFAULT NULL COMMENT 'Timestamp when the alias was created.',
  `updated_at` DATETIME DEFAULT NULL COMMENT 'Timestamp when the alias was last updated.',
  `deleted_at` DATETIME DEFAULT NULL COMMENT 'Timestamp when the alias was deleted.'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='URL alias and metadata.. Original model name - URL.';

--
-- Table structure for table `user` generated from model 'User'
-- User profile information.
--

CREATE TABLE IF NOT EXISTS `user` (
  `id` BIGINT DEFAULT NULL COMMENT 'Unique user identifier.',
  `username` TEXT NOT NULL COMMENT 'Unique username for the account.',
  `first_name` TEXT NOT NULL COMMENT 'User first name.',
  `last_name` TEXT NOT NULL COMMENT 'User last name.',
  `email` TEXT NOT NULL COMMENT 'User email address.',
  `created_at` DATETIME DEFAULT NULL COMMENT 'Timestamp when the user was created.'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='User profile information.. Original model name - User.';


--
-- OAuth2 framework tables
-- Thanks to https://github.com/dsquier/oauth2-server-php-mysql repo
--

--
-- Table structure for table `oauth_clients`
--
CREATE TABLE IF NOT EXISTS `oauth_clients` (
  `client_id`            VARCHAR(80)    NOT NULL,
  `client_secret`        VARCHAR(80)    DEFAULT NULL,
  `redirect_uri`         VARCHAR(2000)  DEFAULT NULL,
  `grant_types`          VARCHAR(80)    DEFAULT NULL,
  `scope`                VARCHAR(4000)  DEFAULT NULL,
  `user_id`              VARCHAR(80)    DEFAULT NULL,
  PRIMARY KEY (`client_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Table structure for table `oauth_access_tokens`
--
CREATE TABLE IF NOT EXISTS `oauth_access_tokens` (
  `access_token`         VARCHAR(40)    NOT NULL,
  `client_id`            VARCHAR(80)    DEFAULT NULL,
  `user_id`              VARCHAR(80)    DEFAULT NULL,
  `expires`              TIMESTAMP      NOT NULL,
  `scope`                VARCHAR(4000)  DEFAULT NULL,
  PRIMARY KEY (`access_token`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Table structure for table `oauth_authorization_codes`
--
CREATE TABLE IF NOT EXISTS `oauth_authorization_codes` (
  `authorization_code`  VARCHAR(40)    NOT NULL,
  `client_id`           VARCHAR(80)    DEFAULT NULL,
  `user_id`             VARCHAR(80)    DEFAULT NULL,
  `redirect_uri`        VARCHAR(2000)  NOT NULL,
  `expires`             TIMESTAMP      NOT NULL,
  `scope`               VARCHAR(4000)  DEFAULT NULL,
  `id_token`            VARCHAR(1000)  DEFAULT NULL,
  PRIMARY KEY (`authorization_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Table structure for table `oauth_refresh_tokens`
--
CREATE TABLE IF NOT EXISTS `oauth_refresh_tokens` (
  `refresh_token`       VARCHAR(40)    NOT NULL,
  `client_id`           VARCHAR(80)    DEFAULT NULL,
  `user_id`             VARCHAR(80)    DEFAULT NULL,
  `expires`             TIMESTAMP      on update CURRENT_TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `scope`               VARCHAR(4000)  DEFAULT NULL,
  PRIMARY KEY (`refresh_token`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Table structure for table `oauth_users`
--
CREATE TABLE IF NOT EXISTS `oauth_users` (
  `username`            VARCHAR(80)    DEFAULT NULL,
  `password`            VARCHAR(255)   DEFAULT NULL,
  `first_name`          VARCHAR(80)    DEFAULT NULL,
  `last_name`           VARCHAR(80)    DEFAULT NULL,
  `email`               VARCHAR(2000)  DEFAULT NULL,
  `email_verified`      TINYINT(1)     DEFAULT NULL,
  `scope`               VARCHAR(4000)  DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Table structure for table `oauth_scopes`
--
CREATE TABLE IF NOT EXISTS `oauth_scopes` (
  `scope`               VARCHAR(80)  NOT NULL,
  `is_default`          TINYINT(1)   DEFAULT NULL,
  PRIMARY KEY (`scope`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Table structure for table `oauth_jwt`
--
CREATE TABLE IF NOT EXISTS `oauth_jwt` (
  `client_id`           VARCHAR(80)    NOT NULL,
  `subject`             VARCHAR(80)    DEFAULT NULL,
  `public_key`          VARCHAR(2000)  NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Table structure for table `oauth_jti`
--
CREATE TABLE IF NOT EXISTS `oauth_jti` (
  `issuer`              VARCHAR(80)    NOT NULL,
  `subject`             VARCHAR(80)    DEFAULT NULL,
  `audience`            VARCHAR(80)    DEFAULT NULL,
  `expires`             TIMESTAMP      NOT NULL,
  `jti`                 VARCHAR(2000)  NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Table structure for table `oauth_public_keys`
--
CREATE TABLE IF NOT EXISTS `oauth_public_keys` (
  `client_id`            VARCHAR(80)    DEFAULT NULL,
  `public_key`           VARCHAR(2000)  DEFAULT NULL,
  `private_key`          VARCHAR(2000)  DEFAULT NULL,
  `encryption_algorithm` VARCHAR(100)   DEFAULT 'RS256'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
