/* SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO"; */
/* SET AUTOCOMMIT = 0; */
/* START TRANSACTION; */
/* SET time_zone = "+00:00"; */

-- --------------------------------------------------------

--
-- Table structure for table `Comment` generated from model 'Comment'
--

CREATE TABLE IF NOT EXISTS `Comment` (
  `id` BIGINT DEFAULT NULL,
  `content` TEXT NOT NULL,
  `like_count` BIGINT DEFAULT NULL,
  `created_at` DATETIME DEFAULT NULL,
  `username` TEXT NOT NULL,
  `deleted_at` DATETIME DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Table structure for table `Error` generated from model 'Error'
--

CREATE TABLE IF NOT EXISTS `Error` (
  `code` INT NOT NULL,
  `message` TEXT NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Table structure for table `User` generated from model 'User'
--

CREATE TABLE IF NOT EXISTS `User` (
  `id` BIGINT DEFAULT NULL,
  `username` TEXT NOT NULL,
  `first_name` TEXT NOT NULL,
  `last_name` TEXT NOT NULL,
  `email` TEXT NOT NULL,
  `created_at` DATETIME DEFAULT NULL,
  `deleted_at` DATETIME DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


