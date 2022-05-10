package controller

import (
	"context"
	"database/sql"
	"fmt"
	"socialapp/pkg/db"
)

type UserControllerImpl struct {
}

func NewUserController(opts ...func(*UserControllerImpl)) *UserControllerImpl {
	newCtrlr := &UserControllerImpl{}

	for i := range opts {
		opts[i](newCtrlr)
	}

	return newCtrlr
}

func (u *UserControllerImpl) GetUserByID(ctx context.Context, dbConn *sql.DB, id int) (User, error) {
	queries := db.New()
	user, err := queries.GetUserByID(ctx, dbConn, int32(id))
	if err != nil {
		return User{}, fmt.Errorf("Failed to get user: %+v", err)
	}

	return FromDBUser(user), nil
}

func (u *UserControllerImpl) GetUserByEmail(ctx context.Context, dbConn *sql.DB, email string) (User, error) {
	queries := db.New()
	user, err := queries.GetUserByEmail(ctx, dbConn, email)
	if err != nil {
		return User{}, fmt.Errorf("Failed to get user: %+v", err)
	}

	return FromDBUser(user), nil
}

func (u *UserControllerImpl) CreateUser(ctx context.Context, dbConn *sql.DB, user User) (User, error) {
	// validate input
	if len(user.FirstName) == 0 ||
		len(user.LastName) == 0 ||
		len(user.Email) == 0 {
		return User{}, fmt.Errorf("Invalid input")
	}

	queries := db.New()
	createUserParams := db.CreateUserParams{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}

	createdUser, err := queries.CreateUser(ctx, dbConn, createUserParams)
	if err != nil {
		return User{}, fmt.Errorf("Failed to create user: %+v", err)
	}

	return FromDBUser(createdUser), nil
}

func (u *UserControllerImpl) DeleteUserByID(ctx context.Context, dbConn *sql.DB, id int) error {
	queries := db.New()
	if err := queries.DeleteUser(ctx, dbConn, int32(id)); err != nil {
		return err
	}
	return nil
}

func (u *UserControllerImpl) ListUsers(ctx context.Context, dbConn *sql.DB) ([]User, error) {
	queries := db.New()
	users, err := queries.ListUsers(ctx, dbConn)
	if err != nil {
		return []User{}, err
	}

	res := make([]User, len(users))
	for i := range res {
		res[i] = FromDBUser(users[i])
	}
	return res, nil
}
