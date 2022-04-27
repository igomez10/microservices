package main

import (
	"context"
	"database/sql"
	"socialapp/pkg/db"
	"testing"

	_ "github.com/lib/pq"
)

// integration
func TestFetchURLIntegration(t *testing.T) {
	ctx := context.Background()
	dbConn, err := sql.Open("postgres", "postgres://postgres:password@localhost:5433/postgres?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	defer dbConn.Close()

	if dbConn == nil {
		t.Fatal("db is nil")
	}
	defer dbConn.Close()

	queries := db.New()
	createUserReq := db.CreateUserParams{
		FirstName: "first",
		LastName:  "last",
		Email:     "first@last.com",
	}

	createdUser, err := queries.CreateUser(ctx, dbConn, createUserReq)
	if err != nil {
		t.Fatal(err)
	}

	// -------------
	actualUser, err := queries.GetUserByID(ctx, dbConn, createdUser.ID)
	if err != nil {
		t.Fatal(err)
	}

	if actualUser.Email != createUserReq.Email ||
		actualUser.FirstName != createUserReq.FirstName ||
		actualUser.LastName != createUserReq.LastName {
		t.Error(actualUser, createUserReq)
	}

}
