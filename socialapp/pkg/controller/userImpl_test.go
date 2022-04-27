package controller

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"socialapp/pkg/db"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/lib/pq"
)

func TestExampleUserControllerImpl_ListUsers(t *testing.T) {

	dbConn, err := sql.Open("postgres", "postgres://postgres:password@localhost:5433/postgres?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer dbConn.Close()

	ctrlr := UserControllerImpl{}

	ctx := context.Background()
	users, err := ctrlr.ListUsers(ctx, dbConn)
	if err != nil {
		panic(err)
	}

	for i := range users {
		fmt.Println(users[i])
	}
}

func TestCreateUser(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	dummyUser := User{
		FirstName: "first",
		LastName:  "last",
		Email:     "first@last.com",
	}

	created_at := time.Now()
	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "create_at", "deleted_at"}).
		AddRow(1, dummyUser.FirstName, dummyUser.LastName, dummyUser.Email, created_at, nil)

	mock.ExpectQuery(
		regexp.QuoteMeta(db.CreateUser)).
		WithArgs(dummyUser.FirstName, dummyUser.LastName, dummyUser.Email).
		WillReturnRows(rows)

	ctrlr := UserControllerImpl{}
	ctx := context.Background()
	newUser, err := ctrlr.CreateUser(ctx, mockDB, dummyUser)
	if err != nil {
		t.Fatal(err)
	}

	if dummyUser.Email != newUser.Email {
		t.Fatal("unexpected user was created")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
