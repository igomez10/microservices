package controller

import (
	"context"
	"database/sql"
	"socialapp/pkg/db"
	"time"
)

type User struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	CreatedAt *time.Time
	DeletedAt *time.Time
}

type UserController interface {
	GetUserByID(ctx context.Context, dbConn *sql.DB, id int) (User, error)
	GetUserByEmail(ctx context.Context, dbConn *sql.DB, email string) (User, error)
	CreateUser(ctx context.Context, dbConn *sql.DB, user User) (User, error)
	DeleteUserByID(ctx context.Context, dbConn *sql.DB, id int) error
	ListUsers(ctx context.Context, dbConn *sql.DB) ([]User, error)
}

func FromDBUser(dbUser db.User) User {
	retUser := User{}
	retUser.ID = int(dbUser.ID)
	retUser.FirstName = dbUser.FirstName
	retUser.LastName = dbUser.LastName
	retUser.Email = dbUser.Email
	retUser.CreatedAt = &dbUser.CreatedAt

	if dbUser.DeletedAt.Valid {
		// user was marked as deleted in db
		retUser.DeletedAt = &dbUser.DeletedAt.Time
	}
	return retUser
}

func ToDBUser(retUser User) db.User {
	dbUser := db.User{}
	dbUser.ID = int32(retUser.ID)
	dbUser.FirstName = retUser.FirstName
	dbUser.LastName = retUser.LastName
	dbUser.Email = retUser.Email
	dbUser.CreatedAt = *retUser.CreatedAt

	if retUser.DeletedAt != nil {
		// user was marked as deleted in db
		dbUser.DeletedAt.Time = *retUser.DeletedAt
	}
	return dbUser
}
