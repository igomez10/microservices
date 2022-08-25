package main

import (
	"context"
	"database/sql"
	"fmt"
	"socialapp/pkg/controller/user"
	"socialapp/pkg/db"
	"socialapp/socialappapi/openapi"
	"sync"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
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

func TestCreateUsers(t *testing.T) {
	dbConn, err := sql.Open("postgres", "postgres://postgres:password@localhost:5433/postgres?sslmode=disable")
	if err != nil {
		log.Fatal().Err(err)
	}
	defer dbConn.Close()

	if dbConn == nil {
		log.Fatal().Msg("db is nil")
	}
	defer dbConn.Close()

	queries := db.New()

	ctx := context.Background()
	UserApiService := &user.UserApiService{DB: queries, DBConn: dbConn}
	for {
		wg := sync.WaitGroup{}
		numUsers := 500
		wg.Add(numUsers)
		for i := 0; i < numUsers; i++ {
			go func(i int) {
				UserApiService.CreateUser(ctx, openapi.User{
					Username:  fmt.Sprintf("Test-%d-%d", time.Now().UnixNano(), i),
					FirstName: "first",
					LastName:  "last",
					Email:     fmt.Sprintf("Test-%d-%d@test.com", time.Now().UnixNano(), i),
				})
				wg.Done()
				// time.Sleep(10 * time.Millisecond)
			}(i)
		}
		wg.Wait()
	}
}
