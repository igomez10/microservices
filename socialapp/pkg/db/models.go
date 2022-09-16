// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0

package db

import (
	"database/sql"
	"time"
)

type Comment struct {
	ID        int64        `json:"id"`
	Content   string       `json:"content"`
	LikeCount int32        `json:"like_count"`
	UserID    int64        `json:"user_id"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}

type Credential struct {
	ID          int64        `json:"id"`
	UserID      int64        `json:"user_id"`
	PublicKey   string       `json:"public_key"`
	Description string       `json:"description"`
	Name        string       `json:"name"`
	CreatedAt   time.Time    `json:"created_at"`
	DeletedAt   sql.NullTime `json:"deleted_at"`
}

type Follower struct {
	ID         int64 `json:"id"`
	FollowerID int64 `json:"follower_id"`
	FollowedID int64 `json:"followed_id"`
}

type Token struct {
	ID           int64     `json:"id"`
	CredentialID int64     `json:"credential_id"`
	Token        string    `json:"token"`
	ValidFrom    time.Time `json:"valid_from"`
	ValidUntil   time.Time `json:"valid_until"`
}

type User struct {
	ID                      int64        `json:"id"`
	Username                string       `json:"username"`
	HashedPassword          string       `json:"hashed_password"`
	HashedPasswordExpiresAt time.Time    `json:"hashed_password_expires_at"`
	Salt                    string       `json:"salt"`
	FirstName               string       `json:"first_name"`
	LastName                string       `json:"last_name"`
	Email                   string       `json:"email"`
	CreatedAt               time.Time    `json:"created_at"`
	UpdatedAt               time.Time    `json:"updated_at"`
	DeletedAt               sql.NullTime `json:"deleted_at"`
}
