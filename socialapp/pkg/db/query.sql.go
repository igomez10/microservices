// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: query.sql

package db

import (
	"context"
)

const CreateComment = `-- name: CreateComment :one
INSERT INTO comments (
  user_id, content
) VALUES (
  $1, $2
)
RETURNING id, content, like_count, created_at, user_id, deleted_at
`

type CreateCommentParams struct {
	UserID  int32  `json:"userID"`
	Content string `json:"content"`
}

func (q *Queries) CreateComment(ctx context.Context, db DBTX, arg CreateCommentParams) (Comment, error) {
	row := db.QueryRowContext(ctx, CreateComment, arg.UserID, arg.Content)
	var i Comment
	err := row.Scan(
		&i.ID,
		&i.Content,
		&i.LikeCount,
		&i.CreatedAt,
		&i.UserID,
		&i.DeletedAt,
	)
	return i, err
}

const CreateCommentForUser = `-- name: CreateCommentForUser :one
INSERT INTO comments (
  user_id, content
) VALUES (
  (SELECT id FROM users WHERE username = $1 AND deleted_at IS NULL), $2
)
RETURNING id, content, like_count, created_at, user_id, deleted_at
`

type CreateCommentForUserParams struct {
	Username string `json:"username"`
	Content  string `json:"content"`
}

func (q *Queries) CreateCommentForUser(ctx context.Context, db DBTX, arg CreateCommentForUserParams) (Comment, error) {
	row := db.QueryRowContext(ctx, CreateCommentForUser, arg.Username, arg.Content)
	var i Comment
	err := row.Scan(
		&i.ID,
		&i.Content,
		&i.LikeCount,
		&i.CreatedAt,
		&i.UserID,
		&i.DeletedAt,
	)
	return i, err
}

const CreateUser = `-- name: CreateUser :one
INSERT INTO users (
  username, first_name, last_name, email
) VALUES (
  $1, $2, $3, $4
)
RETURNING id, username, first_name, last_name, email, created_at, deleted_at
`

type CreateUserParams struct {
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

func (q *Queries) CreateUser(ctx context.Context, db DBTX, arg CreateUserParams) (User, error) {
	row := db.QueryRowContext(ctx, CreateUser,
		arg.Username,
		arg.FirstName,
		arg.LastName,
		arg.Email,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.CreatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const DeleteUser = `-- name: DeleteUser :exec
UPDATE users
SET deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
`

func (q *Queries) DeleteUser(ctx context.Context, db DBTX, id int32) error {
	_, err := db.ExecContext(ctx, DeleteUser, id)
	return err
}

const GetComment = `-- name: GetComment :one
SELECT id, content, like_count, created_at, user_id, deleted_at FROM comments
WHERE id = $1 AND deleted_at IS NULL LIMIT 1
`

func (q *Queries) GetComment(ctx context.Context, db DBTX, id int32) (Comment, error) {
	row := db.QueryRowContext(ctx, GetComment, id)
	var i Comment
	err := row.Scan(
		&i.ID,
		&i.Content,
		&i.LikeCount,
		&i.CreatedAt,
		&i.UserID,
		&i.DeletedAt,
	)
	return i, err
}

const GetUserByEmail = `-- name: GetUserByEmail :one
SELECT id, username, first_name, last_name, email, created_at, deleted_at FROM users
WHERE email = $1 AND deleted_at IS NULL LIMIT 1
`

func (q *Queries) GetUserByEmail(ctx context.Context, db DBTX, email string) (User, error) {
	row := db.QueryRowContext(ctx, GetUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.CreatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const GetUserByID = `-- name: GetUserByID :one
SELECT id, username, first_name, last_name, email, created_at, deleted_at FROM users
WHERE id = $1 AND deleted_at IS NULL LIMIT 1
`

func (q *Queries) GetUserByID(ctx context.Context, db DBTX, id int32) (User, error) {
	row := db.QueryRowContext(ctx, GetUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.CreatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const GetUserByUsername = `-- name: GetUserByUsername :one
SELECT id, username, first_name, last_name, email, created_at, deleted_at FROM users
WHERE username = $1 AND deleted_at IS NULL LIMIT 1
`

func (q *Queries) GetUserByUsername(ctx context.Context, db DBTX, username string) (User, error) {
	row := db.QueryRowContext(ctx, GetUserByUsername, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.CreatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const GetUserComments = `-- name: GetUserComments :many
SELECT
	c.id, c.content, c.like_count, c.created_at, c.user_id, c.deleted_at
FROM
	users u,
	comments c
WHERE
	u.username = $1
	AND u.id = c.user_id
	AND c.deleted_at IS NULL
ORDER BY
	c.created_at DESC
`

func (q *Queries) GetUserComments(ctx context.Context, db DBTX, username string) ([]Comment, error) {
	rows, err := db.QueryContext(ctx, GetUserComments, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Comment
	for rows.Next() {
		var i Comment
		if err := rows.Scan(
			&i.ID,
			&i.Content,
			&i.LikeCount,
			&i.CreatedAt,
			&i.UserID,
			&i.DeletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const ListComment = `-- name: ListComment :many
SELECT id, content, like_count, created_at, user_id, deleted_at FROM comments
WHERE deleted_at IS NULL
ORDER BY name
`

func (q *Queries) ListComment(ctx context.Context, db DBTX) ([]Comment, error) {
	rows, err := db.QueryContext(ctx, ListComment)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Comment
	for rows.Next() {
		var i Comment
		if err := rows.Scan(
			&i.ID,
			&i.Content,
			&i.LikeCount,
			&i.CreatedAt,
			&i.UserID,
			&i.DeletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const ListUsers = `-- name: ListUsers :many
SELECT id, username, first_name, last_name, email, created_at, deleted_at FROM users
WHERE deleted_at IS NULL
ORDER BY first_name
`

func (q *Queries) ListUsers(ctx context.Context, db DBTX) ([]User, error) {
	rows, err := db.QueryContext(ctx, ListUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.CreatedAt,
			&i.DeletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
