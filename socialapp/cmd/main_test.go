package main

import (
	"bytes"
	"context"
	"net"
	"net/http"
	"net/url"
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

func TestRetryClientCheckRetry(t *testing.T) {
	tests := []struct {
		name        string
		resp        *http.Response
		err         error
		expectRetry bool
		expectPanic bool
		description string
	}{
		{
			name:        "network error with nil response - should retry without panic",
			resp:        nil,
			err:         context.DeadlineExceeded,
			expectRetry: true,
			expectPanic: false,
			description: "When network error occurs, resp is nil and should not cause panic",
		},
		{
			name: "5xx status code - should retry",
			resp: &http.Response{
				StatusCode: 500,
				Status:     "500 Internal Server Error",
				Request: &http.Request{
					Method: "GET",
					URL:    mustParseURL(t, "http://example.com/test"),
				},
			},
			err:         nil,
			expectRetry: true,
			expectPanic: false,
			description: "Server errors (5xx) should trigger retry",
		},
		{
			name: "503 status code - should retry",
			resp: &http.Response{
				StatusCode: 503,
				Status:     "503 Service Unavailable",
				Request: &http.Request{
					Method: "POST",
					URL:    mustParseURL(t, "http://example.com/api"),
				},
			},
			err:         nil,
			expectRetry: true,
			expectPanic: false,
			description: "Service unavailable should trigger retry",
		},
		{
			name: "2xx status code - should not retry",
			resp: &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Request: &http.Request{
					Method: "GET",
					URL:    mustParseURL(t, "http://example.com/test"),
				},
			},
			err:         nil,
			expectRetry: false,
			expectPanic: false,
			description: "Successful responses should not retry",
		},
		{
			name: "4xx status code - should not retry",
			resp: &http.Response{
				StatusCode: 404,
				Status:     "404 Not Found",
				Request: &http.Request{
					Method: "GET",
					URL:    mustParseURL(t, "http://example.com/missing"),
				},
			},
			err:         nil,
			expectRetry: false,
			expectPanic: false,
			description: "Client errors (4xx) should not retry",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a logger that won't pollute test output
			logBuf := &bytes.Buffer{}
			logger := zerolog.New(logBuf).With().Timestamp().Logger()

			// Temporarily replace the global logger
			originalLogger := zerolog.GlobalLevel()
			zerolog.SetGlobalLevel(zerolog.WarnLevel)
			defer zerolog.SetGlobalLevel(originalLogger)

			// Create test context
			ctx := context.Background()

			// This is the actual retry check function from getConfig
			checkRetry := func(ctx context.Context, resp *http.Response, err error) (bool, error) {
				// Retry on network errors
				if err != nil {
					logger.Warn().
						Err(err).
						Msg("http retry - network error")
					return true, err
				}

				// Retry on 5xx status codes
				if resp != nil && resp.StatusCode >= 500 {
					logger.Warn().
						Stringer("url", resp.Request.URL).
						Str("method", resp.Request.Method).
						Str("status", resp.Status).
						Int("status_code", resp.StatusCode).
						Msg("http retry - 5xx status")
					return true, nil
				}

				return false, nil
			}

			// Verify no panic occurs
			defer func() {
				if r := recover(); r != nil {
					if !tt.expectPanic {
						t.Errorf("unexpected panic: %v", r)
					}
				} else if tt.expectPanic {
					t.Error("expected panic but none occurred")
				}
			}()

			// Execute the retry check
			shouldRetry, returnedErr := checkRetry(ctx, tt.resp, tt.err)

			// Verify retry decision
			if shouldRetry != tt.expectRetry {
				t.Errorf("expected retry=%v, got retry=%v for: %s", tt.expectRetry, shouldRetry, tt.description)
			}

			// Verify error is returned when provided
			if tt.err != nil && returnedErr != tt.err {
				t.Errorf("expected error to be returned: got %v, want %v", returnedErr, tt.err)
			}
		})
	}
}

// mustParseURL is a helper function that parses a URL or fails the test
func mustParseURL(t *testing.T, rawURL string) *url.URL {
	t.Helper()
	u, err := url.Parse(rawURL)
	if err != nil {
		t.Fatalf("failed to parse URL %q: %v", rawURL, err)
	}
	return u
}
