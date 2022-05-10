package controller

import (
	"context"
	"database/sql"
	"errors"
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

	// verify regular creation
	// verify long firstname, lastName, email
	// verify invalid email
	// verify empty firstName, lastNmae, email

	// verify duplicate email
	// verify deleted and then created -> integration

	type testCase struct {
		name          string
		firstName     string
		lastName      string
		email         string
		expectedError error
	}

	testCases := []testCase{
		{
			name:          "regularCreation",
			firstName:     "first",
			lastName:      "last",
			email:         "first@last.com",
			expectedError: nil,
		},
		{
			name:          "emptyFirstName",
			firstName:     "",
			lastName:      "last",
			email:         "email",
			expectedError: errors.New(""),
		},
	}

	for i := range testCases {
		currTest := testCases[i]
		t.Run(currTest.name, func(t *testing.T) {

			mockDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer mockDB.Close()

			userToCreate := User{
				FirstName: currTest.firstName,
				LastName:  currTest.lastName,
				Email:     currTest.email,
			}

			created_at := time.Now()
			rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "created_at", "deleted_at"}).
				AddRow(1, userToCreate.FirstName, userToCreate.LastName, userToCreate.Email, created_at, nil)

			mock.ExpectQuery(
				regexp.QuoteMeta(db.CreateUser)).
				WithArgs(userToCreate.FirstName, userToCreate.LastName, userToCreate.Email).
				WillReturnRows(rows)

			ctrlr := UserControllerImpl{}
			ctx := context.Background()
			newUser, err := ctrlr.CreateUser(ctx, mockDB, userToCreate)
			// BUG
			if currTest.expectedError != nil || err != nil {
				if currTest.expectedError != err {
					t.Error("Unexpected error", err)
				} else {
					return
				}
			}

			if userToCreate.Email != newUser.Email {
				t.Fatal("unexpected user was created")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}

		})

	}
}
