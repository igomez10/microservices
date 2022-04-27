
-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 AND deleted_at IS NULL LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 AND deleted_at IS NULL LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
WHERE deleted_at IS NULL
ORDER BY first_name;

-- name: CreateUser :one
INSERT INTO users (
  first_name, last_name, email
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: DeleteUser :exec
UPDATE users
SET deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetComment :one
SELECT * FROM comments

WHERE id = $1 AND deleted_at IS NULL LIMIT 1;

-- name: ListComment :many
SELECT * FROM comments
WHERE deleted_at IS NULL
ORDER BY name;

-- name: CreateComment :one
INSERT INTO comments (
  user_id, content, like_count
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- -- name: DeleteComment :exec
UPDATE comments
SET deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;
