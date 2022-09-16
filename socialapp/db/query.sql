
-- name: GetUserByID :one
SELECT * FROM users
WHERE id = ? AND deleted_at IS NULL LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = ? AND deleted_at IS NULL LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = ? AND deleted_at IS NULL LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
WHERE deleted_at IS NULL
ORDER BY first_name;

-- name: CreateUser :execresult
INSERT INTO users (
    username, hashed_password, hashed_password_expires_at, salt, first_name, last_name, email
) VALUES (
	?, ?, ?, ?, ?, ?, ?
);

-- name: UpdateUser :execresult
UPDATE users 
SET username = ?, first_name = ?, last_name=?, email=?
WHERE id = ? AND deleted_at IS NULL;

-- name: UpdateUserByUsername :execresult
UPDATE users 
SET username = sqlc.arg(new_username), first_name=?, last_name=?, email=?
WHERE username = sqlc.arg(old_username) AND deleted_at IS NULL;

-- name: DeleteUser :exec
UPDATE users
SET deleted_at = NOW()
WHERE id = ? AND deleted_at IS NULL;

-- name: DeleteUserByUsername :exec
UPDATE users
SET deleted_at = NOW()
WHERE username = ? AND deleted_at IS NULL;

-- name: GetComment :one
SELECT * FROM comments
WHERE id = ? AND deleted_at IS NULL LIMIT 1;

-- name: GetUserComments :many
SELECT
	c.*
FROM
	comments c JOIN users u
	ON c.user_id = u.id
WHERE
	u.username = ?
	AND c.deleted_at IS NULL
	AND u.deleted_at IS NULL
ORDER BY
	c.created_at DESC;

-- name: ListComment :many
SELECT * FROM comments
WHERE deleted_at IS NULL
ORDER BY created_at DESC;

-- name: CreateComment :execresult
INSERT INTO comments (
  user_id, content
) VALUES (
  ?, ?
);

-- name: CreateCommentForUser :execresult
INSERT INTO comments (
  user_id, content
) VALUES (
  (SELECT id FROM users WHERE username = ? AND deleted_at IS NULL), ?
);

-- name: DeleteComment :exec
UPDATE comments
SET deleted_at = NOW()
WHERE id = ? AND deleted_at IS NULL;

-- name: FollowUser :exec
INSERT INTO followers (
  follower_id, followed_id
) VALUES (
  ?, ?
);

-- name: UnfollowUser :exec
DELETE FROM followers
WHERE follower_id = ? AND followed_id = ?;

-- name: GetFollowers :many
SELECT
	u.*
FROM
	users u,
	followers f
WHERE
	f.followed_id = ?
	AND f.follower_id = u.id
	AND u.deleted_at IS NULL
ORDER BY
	u.first_name;

-- name: GetFollowedUsers :many
SELECT
	u.*
FROM
	users u,
	followers f
WHERE
	f.follower_id = ?
	AND f.followed_id = u.id
	AND u.deleted_at IS NULL
ORDER BY
	u.first_name;

-- name: GetToken :one
SELECT * FROM tokens
WHERE token = ? AND deleted_at IS NULL LIMIT 1;

-- name: CreateToken :execresult
INSERT INTO tokens (
  credential_id, token
) VALUES (
  ?, ?
);

-- name: DeleteToken :exec
UPDATE tokens
SET deleted_at = NOW()
WHERE token = ? AND deleted_at IS NULL;
