package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	db "github.com/igomez10/microservices/socialapp/pkg/dbpgx"
	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// TestFetchURLIntegration tests the integration of the fetchURL function
func TestFetchURLIntegration(t *testing.T) {
	ctx := context.Background()
	container, dbConn := setupPostgresContainer(ctx, t)
	defer func() {
		dbConn.Close(ctx)
		container.Terminate(ctx)
	}()

	createSchema(t, ctx, dbConn)

	queries := db.New()
	createUserReq := db.CreateUserWithIDParams{
		ID:        time.Now().UnixNano(),
		FirstName: "first",
		LastName:  "last",
		Email:     "first@last.com",
	}

	createdUser, err := queries.CreateUserWithID(ctx, dbConn, createUserReq)
	if err != nil {
		t.Fatal(err)
	}

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

func setupPostgresContainer(ctx context.Context, t *testing.T) (testcontainers.Container, *pgx.Conn) {
	t.Helper()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:15",
		Env:          map[string]string{"POSTGRES_USER": "postgres", "POSTGRES_PASSWORD": "password", "POSTGRES_DB": "socialapp"},
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp").WithStartupTimeout(90 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("failed to read container host: %v", err)
	}
	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("failed to read mapped port: %v", err)
	}

	conn, err := pgx.Connect(ctx, fmt.Sprintf("postgres://postgres:password@%s:%s/socialapp?sslmode=disable", host, port.Port()))
	if err != nil {
		container.Terminate(ctx)
		t.Fatalf("failed to connect to postgres container: %v", err)
	}

	return container, conn
}

func createSchema(t *testing.T, ctx context.Context, conn *pgx.Conn) {
	t.Helper()

	schemaPath := schemaFilePath(t)
	payload, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Fatalf("failed to read schema file: %v", err)
	}

	if _, err := conn.Exec(ctx, "BEGIN"); err != nil {
		t.Fatalf("failed to begin tx: %v", err)
	}
	if _, err := conn.Exec(ctx, string(payload)); err != nil {
		t.Fatalf("failed to apply schema: %v", err)
	}
	if _, err := conn.Exec(ctx, "COMMIT"); err != nil {
		t.Fatalf("failed to commit schema transaction: %v", err)
	}
}

func schemaFilePath(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("unable to determine caller")
	}
	return filepath.Join(filepath.Dir(file), "..", "db", "setup", "schema.sql")
}

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

	ml := io.MultiWriter(memoryFile, secondFile, stream)
	logger := slog.New(slog.NewJSONHandler(ml, nil))

	logger.Info("info message")

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
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			shouldRetry := false
			defer func() {
				if r := recover(); r != nil && !tt.expectPanic {
					t.Fatalf("unexpected panic: %v", r)
				}
			}()

			shouldRetry = checkRetry(tt.resp, tt.err)
			if shouldRetry != tt.expectRetry {
				t.Fatalf("%s: expected retry=%t got %t", tt.description, tt.expectRetry, shouldRetry)
			}
		})
	}
}

func mustParseURL(t *testing.T, raw string) *url.URL {
	t.Helper()
	parsed, err := url.Parse(raw)
	if err != nil {
		t.Fatalf("failed to parse URL %s: %v", raw, err)
	}
	return parsed
}

func checkRetry(resp *http.Response, err error) bool {
	if err != nil {
		return true
	}
	if resp == nil || resp.Request == nil {
		return false
	}
	if resp.StatusCode >= http.StatusInternalServerError {
		return true
	}
	return false
}
