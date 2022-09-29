
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
ORDER BY created_at ASC
LIMIT ? OFFSET ?;

-- name: CreateUser :execresult
INSERT INTO users (
    username, hashed_password, salt, first_name, last_name, email, email_token, email_verified_at
) VALUES (
	?, ?, ?, ?, ?, ?, ?, ?
);

-- name: UpdateUser :execresult
UPDATE users 
SET username=?, hashed_password=?, hashed_password_expires_at=?, salt=?, first_name=?, last_name=?, email=?, email_token=?, email_verified_at=?
WHERE id = ? AND deleted_at IS NULL;

-- name: UpdateUserByUsername :execresult
UPDATE users 
SET username = sqlc.arg(new_username), hashed_password=?, hashed_password_expires_at=?, salt=?, first_name=?, last_name=?, email=?, email_token=?, email_verified_at=?
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
	c.created_at DESC
LIMIT ? OFFSET ?;
	

-- name: ListComment :many
SELECT * FROM comments
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

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

-- name: CreateCredential :execresult
INSERT INTO credentials (
  user_id, public_key, description, name
) VALUES (
  ?, ?, ?, ?
);

-- name: DeleteCredential :exec
DELETE FROM credentials
WHERE id = ?;

-- name: GetCredential :one
SELECT * FROM credentials
WHERE public_key = ? AND deleted_at IS NULL LIMIT 1;

-- name: GetToken :one
SELECT * FROM tokens
WHERE token = ? LIMIT 1;

-- name: CreateToken :execresult
INSERT INTO tokens (
	token, user_id, valid_until
) VALUES (
  ?, ?, ?
);

-- name: DeleteToken :exec
UPDATE tokens
SET valid_until = NOW()
WHERE token = ? AND NOW() < valid_until;

-- name: DeleteAllTokensForUser :exec
UPDATE tokens
SET valid_until = NOW()
WHERE user_id = ? AND NOW() < valid_until;


-- name: GetScope :one
SELECT * FROM scopes
WHERE id = ? AND deleted_at IS NULL LIMIT 1;

-- name: GetScopeByName :one
SELECT * FROM scopes
WHERE name = ? AND deleted_at IS NULL LIMIT 1;

-- name: ListScopes :many
SELECT * FROM scopes
WHERE deleted_at IS NULL;

-- name: CreateScope :execresult
INSERT INTO scopes (
	  name, description, deleted_at
) VALUES (
  ?, ?, ?
);

-- name: DeleteScope :exec
UPDATE scopes
SET deleted_at = NOW()
WHERE id = ? AND NOW() < valid_until;

-- name: UpdateScope :execresult
UPDATE scopes
SET name = ?, description = ?
WHERE id = ? AND deleted_at IS NULL;

-- name: GetUserRoles :many
SELECT
	r.*
FROM
	users u
	INNER JOIN users_to_roles ur ON ur.user_id = u.id
	INNER JOIN roles r ON r.id = ur.role_id
WHERE
	u.id = ?
	AND u.deleted_at IS NULL
	AND r.deleted_at IS NULL;

-- name: GetRoleScopes :many
SELECT
	s.*
FROM
	scopes s
	INNER JOIN roles_to_scopes rs ON rs.scope_id = s.id
	INNER JOIN roles r ON r.id = rs.role_id
WHERE
	r.id = ?
	AND r.deleted_at IS NULL
	AND s.deleted_at IS NULL;

-- name: CreateTokenToScope :execresult
INSERT INTO tokens_to_scopes (
	token_id, scope_id
) VALUES (
  ?, ?
);

-- name: GetTokenScopes :many
SELECT
	s.*
FROM
	scopes s
	INNER JOIN tokens_to_scopes ts ON ts.scope_id = s.id
	INNER JOIN tokens t ON t.id = ts.token_id
WHERE
	t.id = ?
	AND t.valid_until > NOW()
	AND s.deleted_at IS NULL;

-- name: CreateUserToRole :execresult
INSERT INTO users_to_roles (
	user_id, role_id
) VALUES (
  ?, ?
);

-- name: GetRoleByName :one
SELECT * FROM roles
WHERE name = ? AND deleted_at IS NULL LIMIT 1;

-- name: GetRole :one
SELECT * FROM roles
WHERE id = ? AND deleted_at IS NULL LIMIT 1;
