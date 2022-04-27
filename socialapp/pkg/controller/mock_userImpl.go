package controller

import (
	"context"
	"database/sql"
	"strconv"
)

const GETUSERBYID = "GetUserByID"
const GETUSERBYEMAIL = "GetUserByEmail"
const CREATEUSER = "CreateUser"
const LISTUSERS = "ListUsers"
const DELETEUSERBYID = "DeleteUserByID"

type userControllerImplMock struct {
	callRegister map[string][]interface{}
}

func (u *userControllerImplMock) GetUserByID(ctx context.Context, dbConn *sql.DB, id int) (User, error) {
	callName := GETUSERBYID
	u.callRegister[callName] = append(u.callRegister[callName], strconv.Itoa(id))
	return User{}, nil
}

func (u *userControllerImplMock) GetUserByEmail(ctx context.Context, dbConn *sql.DB, email string) (User, error) {
	callName := GETUSERBYEMAIL
	u.callRegister[callName] = append(u.callRegister[callName], email)
	return User{}, nil
}

func (u *userControllerImplMock) CreateUser(ctx context.Context, dbConn *sql.DB, user User) error {
	callName := CREATEUSER
	u.callRegister[callName] = append(u.callRegister[callName], user)
	return nil
}

func (u *userControllerImplMock) DeleteUserByID(id int, dbConn *sql.DB) error {
	callName := DELETEUSERBYID
	u.callRegister[callName] = append(u.callRegister[callName], id)
	return nil
}

func (u *userControllerImplMock) ListUsers(ctx context.Context, dbConn *sql.DB) ([]User, error) {
	callName := LISTUSERS
	u.callRegister[callName] = append(u.callRegister[callName], true)
	return []User{}, nil
}
