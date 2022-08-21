
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

-- name: DeleteUser :exec
UPDATE users
SET deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetComment :one
SELECT * FROM comments
WHERE id = $1 AND deleted_at IS NULL LIMIT 1;

-- name: GetUserComments :many
SELECT
	c.*
FROM
	users u,
	comments c
WHERE
	u.username = $1
	AND u.id = c.user_id
	AND c.deleted_at IS NULL
ORDER BY
	c.created_at DESC;

-- name: ListComment :many
SELECT * FROM comments
WHERE deleted_at IS NULL
ORDER BY name;

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
