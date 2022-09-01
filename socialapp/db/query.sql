
-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 AND deleted_at IS NULL LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 AND deleted_at IS NULL LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1 AND deleted_at IS NULL LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
WHERE deleted_at IS NULL
ORDER BY first_name;

-- name: CreateUser :one
INSERT INTO users (
  username, first_name, last_name, email
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateUser :one
UPDATE users 
SET username = $2, first_name = $3, last_name=$4, email=$5
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: UpdateUserByUsername :one
UPDATE users 
SET username = @new_username::text, first_name = $1, last_name=$2, email=$3
WHERE username = @old_username::text AND deleted_at IS NULL
RETURNING *;

-- name: DeleteUser :exec
UPDATE users
SET deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;

-- name: DeleteUserByUsername :exec
UPDATE users
SET deleted_at = NOW()
WHERE username = $1 AND deleted_at IS NULL;

-- name: GetComment :one
SELECT * FROM comments
WHERE id = $1 AND deleted_at IS NULL LIMIT 1;

-- name: GetUserComments :many
SELECT
	c.*
FROM
	comments c JOIN users u
	ON c.user_id = u.id
WHERE
	u.username = $1
	AND c.deleted_at IS NULL
	AND u.deleted_at IS NULL
ORDER BY
	c.created_at DESC;

-- name: ListComment :many
SELECT * FROM comments
WHERE deleted_at IS NULL
ORDER BY created_at DESC;

-- name: CreateComment :one
INSERT INTO comments (
  user_id, content
) VALUES (
  $1, $2
)
RETURNING *;

-- name: CreateCommentForUser :one
INSERT INTO comments (
  user_id, content
) VALUES (
  (SELECT id FROM users WHERE username = $1 AND deleted_at IS NULL), $2
)
RETURNING *;

-- -- name: DeleteComment :exec
UPDATE comments
SET deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;

-- name: FollowUser :exec
INSERT INTO followers (
  follower_id, followed_id
) VALUES (
  $1, $2
);

-- name: UnfollowUser :exec
DELETE FROM followers
WHERE follower_id = $1 AND followed_id = $2;

-- name: GetFollowers :many
SELECT
	u.*
FROM
	users u,
	followers f
WHERE
	f.followed_id = $1
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
	f.follower_id = $1
	AND f.followed_id = u.id
	AND u.deleted_at IS NULL
ORDER BY
	u.first_name;
