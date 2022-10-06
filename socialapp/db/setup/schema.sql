CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY NOT NULL AUTO_INCREMENT,
    username VARCHAR(100) NOT NULL,
    hashed_password VARCHAR(256) NOT NULL DEFAULT '',
    hashed_password_expires_at TIMESTAMP NOT NULL DEFAULT '2000-01-01 00:00:00',
    salt VARCHAR(100) NOT NULL DEFAULT '',
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL,
    email_token VARCHAR(100) NOT NULL DEFAULT '',
    email_verified_at TIMESTAMP NULL DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL 
);

CREATE TABLE IF NOT EXISTS comments (
    id BIGINT PRIMARY KEY NOT NULL AUTO_INCREMENT,
    content VARCHAR(8192) NOT NULL,
    like_count INTEGER NOT NULL DEFAULT 0,
    user_id BIGINT NOT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS followers (
    id BIGINT PRIMARY KEY NOT NULL AUTO_INCREMENT,
    follower_id BIGINT NOT NULL REFERENCES users(id),
    followed_id BIGINT NOT NULL REFERENCES users(id),
    UNIQUE (follower_id, followed_id)
);

CREATE TABLE IF NOT EXISTS credentials ( -- long term api keys
    id BIGINT PRIMARY KEY NOT NULL AUTO_INCREMENT,
    user_id BIGINT NOT NULL REFERENCES users(id),
    public_key VARCHAR(512) NOT NULL,
    description VARCHAR(100) NOT NULL,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    UNIQUE (user_id, public_key)
);

CREATE TABLE IF NOT EXISTS tokens ( -- short term tokens
	id BIGINT PRIMARY KEY NOT NULL AUTO_INCREMENT,
	user_id BIGINT NOT NULL REFERENCES users(id),
	token VARCHAR(512) NOT NULL UNIQUE,
	valid_from TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	valid_until TIMESTAMP NOT NULL DEFAULT '2030-01-01 00:00:00'
);


CREATE TABLE IF NOT EXISTS roles (
    id BIGINT  PRIMARY KEY NOT NULL AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);


CREATE TABLE IF NOT EXISTS roles_to_scopes (
    id  BIGINT  PRIMARY KEY NOT NULL AUTO_INCREMENT,
    role_id BIGINT NOT NULL REFERENCES roles(id),
    scope_id INTEGER NOT NULL REFERENCES scopes(id),
    UNIQUE (role_id, scope_id)
);

CREATE TABLE IF NOT EXISTS scopes (
    id BIGINT  PRIMARY KEY NOT NULL AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS users_to_roles (
    id BIGINT  PRIMARY KEY NOT NULL AUTO_INCREMENT,
    role_id BIGINT NOT NULL REFERENCES roles(id),
    user_id BIGINT NOT NULL REFERENCES users(id),
    UNIQUE (role_id, user_id)
);

CREATE TABLE IF NOT EXISTS  tokens_to_scopes (
    id BIGINT  PRIMARY KEY NOT NULL AUTO_INCREMENT,
    token_id BIGINT NOT NULL REFERENCES tokens(id),
    scope_id BIGINT NOT NULL REFERENCES scopes(id),
    INDEX (token_id)
);

--  SEEDING DATA
INSERT INTO users (
    id , username , hashed_password , hashed_password_expires_at , salt , first_name , last_name , email , created_at , updated_at , deleted_at ) 
VALUES (1, 'admin', 'LPWJbW+u2por79jDW+uClI+VvxrX7HpT5eX53kdNd4U=', '2030-01-01 00:00:00', 'MTUyNzJkYjgtZjVjNi00YjIxLTk3ZDktZDJkMTEzODM5NjQ1', 'first_name', 'last_name', 'email@email.com', '2020-01-01 00:00:00', '2020-01-01 00:00:00', NULL);

INSERT INTO comments (id, content, like_count, created_at, user_id, deleted_at) VALUES
(1, 'something', 0, '2022-08-20 11:53:21.218349', 1, NULL),
(1, 'something', 0, '2022-08-20 11:53:21.218349', 1, NULL),
(2, 'something', 0, '2022-08-20 11:53:21.218349', 1, NULL),
(3, 'something', 0, '2022-08-20 11:53:21.218349', 1, NULL),
(4, 'something', 0, '2022-08-20 11:53:21.218349', 2, NULL),
(5, 'something', 0, '2022-08-20 11:53:21.218349', 2, NULL),
(6, 'something', 0, '2022-08-20 11:53:21.218349', 2, NULL),
(7, 'something', 0, '2022-08-20 11:53:21.218349', 2, NULL);

INSERT INTO followers (follower_id, followed_id) VALUES
(1, 2);

INSERT INTO roles (id, name, description, created_at, deleted_at) VALUES
(1, 'administrator', 'administrator', '2020-01-01 00:00:00', NULL),
(2, 'user', 'socialapp user', '2020-01-01 00:00:00', NULL);

INSERT INTO `scopes` (`id`, `name`, `description`, `created_at`, `deleted_at`) VALUES
(1, 'socialapp.users.list', 'socialapp.users.list', '2020-01-01 00:00:00', NULL),
(2, 'socialapp.users.create', 'socialapp.users.create', '2020-01-01 00:00:00', NULL),
(3, 'socialapp.users.update', 'socialapp.users.update', '2020-01-01 00:00:00', NULL),
(4, 'socialapp.users.delete', 'socialapp.users.delete', '2020-01-01 00:00:00', NULL),
(5, 'socialapp.comments.list', 'socialapp.comments.list', '2020-01-01 00:00:00', NULL),
(6, 'socialapp.comments.create', 'socialapp.comments.create', '2020-01-01 00:00:00', NULL),
(7, 'socialapp.comments.update', 'socialapp.comments.update', '2020-01-01 00:00:00', NULL),
(8, 'socialapp.comments.delete', 'socialapp.comments.delete', '2020-01-01 00:00:00', NULL),
(9, 'socialapp.followers.list', 'socialapp.followers.list', '2020-01-01 00:00:00', NULL),
(10, 'socialapp.following.list', 'socialapp.following.list', '2020-01-01 00:00:00', NULL),
(11, 'socialapp.users.read', 'socialapp.users.read', '2022-09-28 20:07:31', NULL),
(12, 'socialapp.follower.create', 'socialapp.follower.create', '2022-09-28 20:13:04', NULL),
(13, 'socialapp.follower.read', 'socialapp.follower.read', '2022-09-28 20:14:03', NULL),
(14, 'socialapp.follower.delete', 'socialapp.follower.delete', '2022-09-28 20:14:55', NULL),
(15, 'socialapp.feed.read', 'socialapp.feed.read', '2022-09-28 20:21:07', NULL),
(16, 'socialapp.roles.create', 'socialapp.roles.create', '2022-09-28 20:21:07', NULL),
(17, 'socialapp.roles.read', 'socialapp.roles.read', '2022-09-28 20:21:07', NULL),
(18, 'socialapp.roles.update', 'socialapp.roles.update', '2022-09-28 20:21:07', NULL),
(19, 'socialapp.roles.delete', 'socialapp.roles.delete', '2022-09-28 20:21:07', NULL),
(20, 'socialapp.roles.list', 'socialapp.roles.list', '2022-09-28 20:21:07', NULL),
(21, 'socialapp.scopes.list', 'socialapp.scopes.list', '2022-09-28 20:21:07', NULL),
(22, 'socialapp.scopes.read', 'socialapp.scopes.read', '2022-09-28 20:21:07', NULL),
(23, 'socialapp.scopes.create', 'socialapp.scopes.create', '2022-09-28 20:21:07', NULL),
(24, 'socialapp.scopes.delete', 'socialapp.scopes.delete', '2022-09-28 20:21:07', NULL),
(25, 'socialapp.scopes.update', 'socialapp.scopes.update', '2022-09-28 20:21:07', NULL),
(26, 'socialapp.roles.listscopes', 'socialapp.roles.listscopes','2022-09-28 20:21:07', NULL);

INSERT INTO roles_to_scopes (id, role_id, scope_id) VALUES
(1, 1, 1),
(2, 1, 2),
(3, 1, 3),
(4, 1, 4),
(5, 1, 5),
(6, 1, 6),
(7, 1, 7),
(8, 1, 8),
(9, 1, 9),
(10, 1, 10),
(11, 2, 1),
(12, 2, 5),
(13, 2, 6),
(14, 2, 7),
(15, 2, 8),
(16, 2, 9),
(17, 2, 10),
(18, 1, 11),
(19, 1, 12),
(20, 1, 13),
(21, 1, 14),
(22, 1, 15),
(23, 1, 16),
(24, 1, 17),
(25, 1, 18),
(26, 1, 19),
(27, 1, 20),
(28, 2, 11),
(29, 2, 12),
(30, 2, 13),
(31, 2, 14),
(32, 2, 15),
(33, 2, 16),
(34, 2, 17),
(35, 2, 18),
(36, 2, 19),
(37, 2, 20),
(38, 1, 21),
(39, 1, 22),
(40, 1, 23),
(41, 1, 24),
(42, 1, 25),
(43, 2, 21),
(44, 2, 22),
(45, 2, 23),
(46, 2, 24),
(47, 2, 25),
(48, 1, 26),
(49, 2, 26);


INSERT INTO users_to_roles (id, role_id, user_id) VALUES
(1, 1, 1);


