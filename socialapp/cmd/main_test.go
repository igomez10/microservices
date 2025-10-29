package main

import (
	"bytes"
	"context"
	"net"
	"testing"
	"time"

	db "github.com/igomez10/microservices/socialapp/pkg/dbpgx"
	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)

// TestFetchURLIntegration tests the integration of the fetchURL function
func TestFetchURLIntegration(t *testing.T) {
	ctx := context.Background()
	dbConn, err := pgx.Connect(ctx, "postgres://postgres:password@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	defer dbConn.Close(ctx)

	if dbConn == nil {
		t.Fatal("db is nil")
	}
	defer dbConn.Close(ctx)

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

// func TestCreateUsers(t *testing.T) {
// 	dbConn, err := sql.Open("postgres", "postgres://postgres:password@localhost:5432/postgres?sslmode=disable")
// 	if err != nil {
// 		log.Fatal().Err(err)
// 	}
// 	defer dbConn.Close()

// 	if dbConn == nil {
// 		log.Fatal().Msg("db is nil")
// 	}
// 	defer dbConn.Close()

// 	queries := db.New()

// 	ctx := context.Background()
// 	UserApiService := &user.UserApiService{DB: queries, DBConn: dbConn}
// 	counter := 0
// 	for {
// 		for i := 0; i < 10; i++ {
// 			UserApiService.CreateUser(ctx, openapi.User{
// 				Username:  fmt.Sprintf("Test-%d-%d", time.Now().UnixNano(), i),
// 				FirstName: "first",
// 				LastName:  "last",
// 				Email:     fmt.Sprintf("Test-%d-%d@test.com", time.Now().UnixNano(), i),
// 			})

// 			time.Sleep(100 * time.Millisecond)
// 		}
// 		counter++
// 		fmt.Printf("new users: %d\n", counter*10)
// 		time.Sleep(200 * time.Millisecond)
// 	}
// }

// func TestListUsers(t *testing.T) {
// 	// call get localhost:8080/users
// 	for {
// 		http.Get("http://localhost:8080/users")
// 		time.Sleep(200 * time.Millisecond)
// 	}
// }

func TestMultiLevelLog(t *testing.T) {
	memoryFile := bytes.NewBuffer([]byte{})
	secondFile := bytes.NewBuffer([]byte{})
	server, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 0})
	if err != nil {
		t.Fatal(err)
	}
	streamBytes := bytes.NewBuffer([]byte{})
	go func() {
		conn, err := server.Accept()
		if err != nil {
			t.Log(err)
			return
		}
		defer conn.Close()
		buf := make([]byte, 1024)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				t.Log(err)
				return
			}
			streamBytes.Write(buf[:n])
		}
	}()
	defer server.Close()
	stream, err := net.DialTCP("tcp", nil, server.Addr().(*net.TCPAddr))
	if err != nil {
		t.Fatal(err)
	}
	defer stream.Close()

	ml := zerolog.MultiLevelWriter(memoryFile, secondFile, stream)
	logger := zerolog.New(ml).With().Timestamp().Logger()

	logger.Info().Msg("info message")

	// Give the goroutine time to read from the stream
	time.Sleep(100 * time.Millisecond)

	if !bytes.Contains(memoryFile.Bytes(), []byte("info message")) {
		t.Error("info message not found in memoryFile")
	}

	if !bytes.Contains(secondFile.Bytes(), []byte("info message")) {
		t.Error("info message not found in secondFile")
	}

	if !bytes.Contains(streamBytes.Bytes(), []byte("info message")) {
		t.Error("info message not found in stream")
	}
}
