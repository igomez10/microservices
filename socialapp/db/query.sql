
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
ORDER BY created_at ASC
LIMIT $1 OFFSET $2;

-- name: CreateUser :one
INSERT INTO users (
    username, hashed_password, salt, first_name, last_name, email, email_token, email_verified_at
) VALUES (
	$1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: UpdateUser :exec
UPDATE users 
SET username=$1, hashed_password=$2, hashed_password_expires_at=$3, salt=$4, first_name=$5, last_name=$6, email=$7, email_token=$8, email_verified_at=$9, updated_at=$10
WHERE id=$11 AND deleted_at IS NULL;

-- name: UpdateUserByUsername :exec
UPDATE users 
SET username = sqlc.arg(new_username), hashed_password=$1, hashed_password_expires_at=$2, salt=$3, first_name=$4, last_name=$5, email=$6, email_token=$7, email_verified_at=$8, updated_at=$9
WHERE username = sqlc.arg(old_username) AND deleted_at IS NULL;

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
	c.created_at DESC
LIMIT $2 OFFSET $3;
	

-- name: ListComment :many
SELECT * FROM comments
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

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

-- name: DeleteComment :exec
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

-- name: CreateCredential :one
INSERT INTO credentials (
  user_id, public_key, description, name
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: DeleteCredential :exec
DELETE FROM credentials
WHERE id = $1;

-- name: GetCredential :one
SELECT * FROM credentials
WHERE public_key = $1 AND deleted_at IS NULL LIMIT 1;

-- name: GetToken :one
SELECT * FROM tokens
WHERE token = $1 LIMIT 1;

-- name: CreateToken :one
INSERT INTO tokens (
	token, user_id, valid_until
) VALUES (
	$1, $2, $3
)
RETURNING *;

-- name: DeleteToken :exec
UPDATE tokens
SET valid_until = NOW()
WHERE token = $1 AND NOW() < valid_until;

-- name: DeleteAllTokensForUser :exec
UPDATE tokens
SET valid_until = NOW()
WHERE user_id = $1 AND NOW() < valid_until;


-- name: GetScope :one
SELECT * FROM scopes
WHERE id = $1 AND deleted_at IS NULL LIMIT 1;

-- name: GetScopeByName :one
SELECT * FROM scopes
WHERE name = $1 AND deleted_at IS NULL LIMIT 1;

-- name: ListScopes :many
SELECT * FROM scopes
WHERE deleted_at IS NULL
LIMIT $1 OFFSET $2;

-- name: CreateScope :one
INSERT INTO scopes (
	  name, description
) VALUES (
	$1, $2
)
RETURNING *;

-- name: DeleteScope :exec
UPDATE scopes
SET deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;

-- name: UpdateScope :execresult
UPDATE scopes
SET name = $1, description = $2
WHERE id = $3 AND deleted_at IS NULL;

-- name: GetUserRoles :many
SELECT
	r.*
FROM
	users u
	INNER JOIN users_to_roles ur ON ur.user_id = u.id
	INNER JOIN roles r ON r.id = ur.role_id
WHERE
	u.id = $1
	AND u.deleted_at IS NULL
	AND r.deleted_at IS NULL;

-- name: ListRoleScopes :many
SELECT
	s.*
FROM
	scopes s
	INNER JOIN roles_to_scopes rs ON rs.scope_id = s.id
	INNER JOIN roles r ON r.id = rs.role_id
WHERE
	r.id = $1
	AND r.deleted_at IS NULL
	AND s.deleted_at IS NULL
ORDER BY
	s.name
LIMIT $2 OFFSET $3;

-- name: CreateRoleScope :one
INSERT INTO roles_to_scopes (
	role_id, scope_id
) VALUES (
	$1, $2
)
RETURNING *;

-- name: DeleteRoleScope :exec
DELETE FROM roles_to_scopes
WHERE role_id = $1 AND scope_id = $2;

-- name: CreateTokenToScope :one
INSERT INTO tokens_to_scopes (
	token_id, scope_id
) VALUES (
  $1, $2
)
RETURNING *;

-- name: GetTokenScopes :many
SELECT
	s.*
FROM
	scopes s
	INNER JOIN tokens_to_scopes ts ON ts.scope_id = s.id
	INNER JOIN tokens t ON t.id = ts.token_id
WHERE
	t.id = $1
	AND t.valid_until > NOW()
	AND s.deleted_at IS NULL;

-- name: CreateUserToRole :one
INSERT INTO users_to_roles (
	user_id, role_id
) VALUES (
	$1, $2
)
RETURNING *;

-- name: DeleteUserToRole :exec
DELETE FROM users_to_roles
WHERE user_id = $1 AND role_id = $2;

-- name: GetRoleByName :one
SELECT * FROM roles
WHERE name = $1 AND deleted_at IS NULL LIMIT 1;

-- name: GetRole :one
SELECT * FROM roles
WHERE id = $1 AND deleted_at IS NULL LIMIT 1;

-- name: ListRoles :many
SELECT * FROM roles
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CreateRole :one
INSERT INTO roles (name, description) 
VALUES ($1, $2)
RETURNING *;

-- name: UpdateRole :exec
UPDATE roles 
SET name = $1, description = $2
WHERE id = $3 AND deleted_at IS NULL;

-- name: DeleteRole :exec
UPDATE roles 
SET deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;



-- SHORTLY

-- name: GetURLFromAlias :one
SELECT * FROM urls
WHERE alias = $1 AND deleted_at IS NULL LIMIT 1;

-- name: CreateURL :one
INSERT INTO urls (
  alias, url
) VALUES (
  $1, $2
)
RETURNING *;

-- name: DeleteURL :exec
DELETE FROM urls
WHERE alias = $1;
