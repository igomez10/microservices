-- SHORTLY
-- name: GetURLFromAlias :one
SELECT *
FROM urls
WHERE alias = $1
	AND deleted_at IS NULL
LIMIT 1;
-- name: CreateURL :one
INSERT INTO urls (alias, url)
VALUES ($1, $2)
RETURNING *;
-- name: DeleteURL :exec
UPDATE urls
SET DELETED_AT = NOW()
WHERE alias = $1
	AND deleted_at IS NULL;
