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
	credential_id BIGINT NOT NULL REFERENCES credentials (id),
	token VARCHAR(512) NOT NULL UNIQUE,
	valid_from TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	valid_until TIMESTAMP NOT NULL DEFAULT '2030-01-01 00:00:00'
);

--  SEEDING DATA
INSERT INTO users (
    id , username , hashed_password , hashed_password_expires_at , salt , first_name , last_name , email , created_at , updated_at , deleted_at ) 
VALUES (1, 'admin', '3c4a79782143337be4492be072abcfe979dd703c00541a8c39a0f3df4bab2029c050cf46fddc47090b5b04ac537b3e78189e3de16e601e859f95c51ac9f6dafb', '2030-01-01 00:00:00', 'salt', 'first_name', 'last_name', 'email@email.com', '2020-01-01 00:00:00', '2020-01-01 00:00:00', NULL);

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

INSERT INTO credentials (id, user_id, public_key, description, name, created_at, deleted_at) VALUES
(1, 1, 'public_key', 'key description', 'mykey', '2022-08-20 11:50:28.522646', NULL);

INSERT INTO tokens (id, credential_id, token, valid_from, valid_until) VALUES
(1, 1, 'token', '2022-08-20 11:50:28.522646', '2030-08-20 11:50:28.522646');
