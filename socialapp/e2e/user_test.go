package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/igomez10/microservices/socialapp/internal/middlewares/cache"
	"github.com/igomez10/microservices/socialapp/internal/server"
	"github.com/igomez10/microservices/socialapp/pkg/dbpgx"
	"github.com/igomez10/microservices/socialapp/socialappapi/openapi"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	tcRedis "github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
	"log/slog"
)

// TestEnv holds all the resources needed for e2e tests
type TestEnv struct {
	PostgresContainer *postgres.PostgresContainer
	RedisContainer    *tcRedis.RedisContainer
	DBPool            *pgxpool.Pool
	RedisClient       *redis.Client
	Server            *httptest.Server
	BaseURL           string
}

// setupTestEnv creates a test environment with PostgreSQL and Redis containers
func setupTestEnv(t *testing.T) *TestEnv {
	ctx := context.Background()

	// Start PostgreSQL container
	pgContainer, err := postgres.Run(ctx,
		"postgres:15-bookworm",
		postgres.WithDatabase("socialapp"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	require.NoError(t, err, "Failed to start PostgreSQL container")

	// Get PostgreSQL connection string
	pgConnStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err, "Failed to get PostgreSQL connection string")

	// Connect to PostgreSQL
	dbPool, err := pgxpool.New(ctx, pgConnStr)
	require.NoError(t, err, "Failed to connect to PostgreSQL")

	// Initialize database schema
	err = initializeDatabase(ctx, dbPool)
	require.NoError(t, err, "Failed to initialize database schema")

	// Start Redis container
	redisContainer, err := tcRedis.Run(ctx,
		"redis:7-alpine",
		testcontainers.WithWaitStrategy(
			wait.ForLog("Ready to accept connections").
				WithStartupTimeout(30*time.Second),
		),
	)
	require.NoError(t, err, "Failed to start Redis container")

	// Get Redis connection string
	redisConnStr, err := redisContainer.ConnectionString(ctx)
	require.NoError(t, err, "Failed to get Redis connection string")

	// Connect to Redis
	redisOpts, err := redis.ParseURL(redisConnStr)
	require.NoError(t, err, "Failed to parse Redis URL")
	redisClient := redis.NewClient(redisOpts)

	// Create test server using the server package
	testServer := createTestServer(t, ctx, dbPool, redisOpts)

	return &TestEnv{
		PostgresContainer: pgContainer,
		RedisContainer:    redisContainer,
		DBPool:            dbPool,
		RedisClient:       redisClient,
		Server:            testServer,
		BaseURL:           testServer.URL,
	}
}

// teardownTestEnv cleans up all test resources
func teardownTestEnv(t *testing.T, env *TestEnv) {
	ctx := context.Background()

	if env.Server != nil {
		env.Server.Close()
	}

	if env.RedisClient != nil {
		env.RedisClient.Close()
	}

	if env.DBPool != nil {
		env.DBPool.Close()
	}

	if env.RedisContainer != nil {
		if err := env.RedisContainer.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate Redis container: %v", err)
		}
	}

	if env.PostgresContainer != nil {
		if err := env.PostgresContainer.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate PostgreSQL container: %v", err)
		}
	}
}

// initializeDatabase runs the schema.sql file to set up the database
func initializeDatabase(ctx context.Context, pool *pgxpool.Pool) error {
	// Find schema.sql - try multiple paths
	schemaPaths := []string{
		"../db/setup/schema.sql",
		"db/setup/schema.sql",
		filepath.Join("..", "db", "setup", "schema.sql"),
	}

	var schemaSQL []byte
	var err error
	for _, path := range schemaPaths {
		schemaSQL, err = os.ReadFile(path)
		if err == nil {
			break
		}
	}

	if err != nil {
		return fmt.Errorf("failed to read schema.sql: %w", err)
	}

	// Execute the schema
	_, err = pool.Exec(ctx, string(schemaSQL))
	if err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	return nil
}

// newConfigForTest creates a minimal server.Config suitable for testing
func newConfigForTest(dbPool *pgxpool.Pool, redisOpts *redis.Options) server.Config {
	queries := dbpgx.New()

	var cacheMiddleware *cache.Cache
	if redisOpts != nil {
		cacheMiddleware = cache.NewCache(cache.CacheConfig{
			RedisOpts: redisOpts,
		})
	}

	openAPIContent, err := os.ReadFile("../openapi.yaml")
	if err != nil {
		panic(err)
	}
	return server.Config{
		InstanceID:      "test-instance",
		AppName:         "socialapp-test",
		AppPort:         0, // Will be assigned by httptest
		LogLevel:        slog.LevelInfo,
		LogDestinations: []io.Writer{io.Discard},
		DBPool:          dbPool,
		Queries:         queries,
		Cache:           cacheMiddleware,
		DefaultTimeout:  20 * time.Second,
		JWTSecret:       "test-jwt-secret",
		URLShortenerURL: &url.URL{Scheme: "http", Host: "localhost:8089"},
		KibanaURL:       "http://localhost:5601",
		OpenAPIContent:  openAPIContent,
	}
}

// createTestServer creates an httptest.Server using the server.Config and server.NewRouter
func createTestServer(t *testing.T, ctx context.Context, dbPool *pgxpool.Pool, redisOpts *redis.Options) *httptest.Server {
	// Create configuration for testing
	config := newConfigForTest(dbPool, redisOpts)

	// Create the router using the server package
	router, err := server.NewRouter(ctx, config)
	require.NoError(t, err, "Failed to create router")

	// Create httptest server
	return httptest.NewServer(router)
}

// TestUser represents a user created for testing purposes
type TestUser struct {
	Username string
	Password string
	Email    string
	UserID   int64
}

// createTestUser creates a user via the API with known credentials.
// The CreateUser API automatically assigns the "administrator" role to new users,
// so created users have all scopes available.
func createTestUser(t *testing.T, env *TestEnv) TestUser {
	t.Helper()

	// Generate unique user credentials
	timestamp := time.Now().UnixNano()
	username := fmt.Sprintf("testuser_%d", timestamp)
	password := "TestPassword123!"
	email := fmt.Sprintf("testuser_%d@example.com", timestamp)

	// Create user via API
	createUserReq := openapi.CreateUserRequest{
		Username:  username,
		Password:  password,
		FirstName: "Test",
		LastName:  "User",
		Email:     email,
	}

	body, err := json.Marshal(createUserReq)
	require.NoError(t, err)

	resp, err := http.Post(
		env.BaseURL+"/v1/users",
		"application/json",
		bytes.NewReader(body),
	)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode,
		"Failed to create test user: %s", string(respBody))

	var createdUser openapi.CreateUserResponse
	err = json.Unmarshal(respBody, &createdUser)
	require.NoError(t, err)

	return TestUser{
		Username: username,
		Password: password,
		Email:    email,
		UserID:   createdUser.Id,
	}
}

// createAdminTestUser creates a user with the administrator role (has all scopes)
// This is an alias for createTestUser since the API auto-assigns the admin role
func createAdminTestUser(t *testing.T, env *TestEnv) TestUser {
	return createTestUser(t, env)
}

// TestCreateUser_Success tests successful user creation
func TestCreateUser_Success(t *testing.T) {
	env := setupTestEnv(t)
	defer teardownTestEnv(t, env)

	// Create user request
	createUserReq := openapi.CreateUserRequest{
		Username:  "testuser",
		Password:  "TestPassword123!",
		FirstName: "Test",
		LastName:  "User",
		Email:     "testuser@example.com",
	}

	body, err := json.Marshal(createUserReq)
	require.NoError(t, err)

	// Make HTTP request
	resp, err := http.Post(
		env.BaseURL+"/v1/users",
		"application/json",
		bytes.NewReader(body),
	)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	// Assert status code
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected 200 OK, got %d: %s", resp.StatusCode, string(respBody))

	// Parse response
	var createUserResp openapi.CreateUserResponse
	err = json.Unmarshal(respBody, &createUserResp)
	require.NoError(t, err, "Failed to unmarshal response: %s", string(respBody))

	// Assert response fields
	assert.NotZero(t, createUserResp.Id, "Expected non-zero ID")
	assert.Equal(t, createUserReq.Username, createUserResp.Username)
	assert.Equal(t, createUserReq.FirstName, createUserResp.FirstName)
	assert.Equal(t, createUserReq.LastName, createUserResp.LastName)
	assert.Equal(t, createUserReq.Email, createUserResp.Email)
	assert.False(t, createUserResp.CreatedAt.IsZero(), "Expected non-zero CreatedAt")
}

// TestCreateUser_DuplicateUsername tests that creating a user with duplicate username returns 409
func TestCreateUser_DuplicateUsername(t *testing.T) {
	env := setupTestEnv(t)
	defer teardownTestEnv(t, env)

	// Create first user
	createUserReq := openapi.CreateUserRequest{
		Username:  "duplicateuser",
		Password:  "TestPassword123!",
		FirstName: "First",
		LastName:  "User",
		Email:     "first@example.com",
	}

	body, err := json.Marshal(createUserReq)
	require.NoError(t, err)

	resp, err := http.Post(
		env.BaseURL+"/v1/users",
		"application/json",
		bytes.NewReader(body),
	)
	require.NoError(t, err)
	resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode, "First user creation should succeed")

	// Try to create second user with same username but different email
	createUserReq2 := openapi.CreateUserRequest{
		Username:  "duplicateuser", // Same username
		Password:  "TestPassword123!",
		FirstName: "Second",
		LastName:  "User",
		Email:     "second@example.com", // Different email
	}

	body2, err := json.Marshal(createUserReq2)
	require.NoError(t, err)

	resp2, err := http.Post(
		env.BaseURL+"/v1/users",
		"application/json",
		bytes.NewReader(body2),
	)
	require.NoError(t, err)
	defer resp2.Body.Close()

	// Should return 409 Conflict
	assert.Equal(t, http.StatusConflict, resp2.StatusCode, "Expected 409 Conflict for duplicate username")
}

// TestCreateUser_DuplicateEmail tests that creating a user with duplicate email returns 409
func TestCreateUser_DuplicateEmail(t *testing.T) {
	env := setupTestEnv(t)
	defer teardownTestEnv(t, env)

	// Create first user
	createUserReq := openapi.CreateUserRequest{
		Username:  "emailuser1",
		Password:  "TestPassword123!",
		FirstName: "First",
		LastName:  "User",
		Email:     "duplicate@example.com",
	}

	body, err := json.Marshal(createUserReq)
	require.NoError(t, err)

	resp, err := http.Post(
		env.BaseURL+"/v1/users",
		"application/json",
		bytes.NewReader(body),
	)
	require.NoError(t, err)
	resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode, "First user creation should succeed")

	// Try to create second user with different username but same email
	createUserReq2 := openapi.CreateUserRequest{
		Username:  "emailuser2", // Different username
		Password:  "TestPassword123!",
		FirstName: "Second",
		LastName:  "User",
		Email:     "duplicate@example.com", // Same email
	}

	body2, err := json.Marshal(createUserReq2)
	require.NoError(t, err)

	resp2, err := http.Post(
		env.BaseURL+"/v1/users",
		"application/json",
		bytes.NewReader(body2),
	)
	require.NoError(t, err)
	defer resp2.Body.Close()

	// Should return 409 Conflict
	assert.Equal(t, http.StatusConflict, resp2.StatusCode, "Expected 409 Conflict for duplicate email")
}

// TestCreateUser_InvalidRequest tests that invalid requests are properly rejected
func TestCreateUser_InvalidRequest(t *testing.T) {
	env := setupTestEnv(t)
	defer teardownTestEnv(t, env)

	testCases := []struct {
		name           string
		request        map[string]interface{}
		expectedStatus int
	}{
		{
			name: "missing username",
			request: map[string]interface{}{
				"password":   "TestPassword123!",
				"first_name": "Test",
				"last_name":  "User",
				"email":      "test@example.com",
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "missing email",
			request: map[string]interface{}{
				"username":   "testuser",
				"password":   "TestPassword123!",
				"first_name": "Test",
				"last_name":  "User",
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "missing password",
			request: map[string]interface{}{
				"username":   "testuser",
				"first_name": "Test",
				"last_name":  "User",
				"email":      "test@example.com",
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "empty body",
			request:        map[string]interface{}{},
			expectedStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body, err := json.Marshal(tc.request)
			require.NoError(t, err)

			resp, err := http.Post(
				env.BaseURL+"/v1/users",
				"application/json",
				bytes.NewReader(body),
			)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tc.expectedStatus, resp.StatusCode, "Test case: %s", tc.name)
		})
	}
}

// TestCreateUser_MultipleUsers tests creating multiple users successfully
func TestCreateUser_MultipleUsers(t *testing.T) {
	env := setupTestEnv(t)
	defer teardownTestEnv(t, env)

	users := []openapi.CreateUserRequest{
		{
			Username:  "user1",
			Password:  "TestPassword123!",
			FirstName: "User",
			LastName:  "One",
			Email:     "user1@example.com",
		},
		{
			Username:  "user2",
			Password:  "TestPassword123!",
			FirstName: "User",
			LastName:  "Two",
			Email:     "user2@example.com",
		},
		{
			Username:  "user3",
			Password:  "TestPassword123!",
			FirstName: "User",
			LastName:  "Three",
			Email:     "user3@example.com",
		},
	}

	createdIDs := make(map[int64]bool)

	for _, userReq := range users {
		body, err := json.Marshal(userReq)
		require.NoError(t, err)

		resp, err := http.Post(
			env.BaseURL+"/v1/users",
			"application/json",
			bytes.NewReader(body),
		)
		require.NoError(t, err)

		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode, "User %s creation failed: %s", userReq.Username, string(respBody))

		var createUserResp openapi.CreateUserResponse
		err = json.Unmarshal(respBody, &createUserResp)
		require.NoError(t, err)

		// Verify ID is unique
		assert.False(t, createdIDs[createUserResp.Id], "Duplicate ID found: %d", createUserResp.Id)
		createdIDs[createUserResp.Id] = true
	}

	// All users should have unique IDs
	assert.Len(t, createdIDs, len(users), "Expected %d unique IDs", len(users))
}

// TestCreateUser_SnowflakeIDReturned verifies that created users have valid snowflake IDs
func TestCreateUser_SnowflakeIDReturned(t *testing.T) {
	env := setupTestEnv(t)
	defer teardownTestEnv(t, env)

	createUserReq := openapi.CreateUserRequest{
		Username:  "snowflakeuser",
		Password:  "TestPassword123!",
		FirstName: "Snowflake",
		LastName:  "User",
		Email:     "snowflake@example.com",
	}

	body, err := json.Marshal(createUserReq)
	require.NoError(t, err)

	resp, err := http.Post(
		env.BaseURL+"/v1/users",
		"application/json",
		bytes.NewReader(body),
	)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode, "Failed: %s", string(respBody))

	var createUserResp openapi.CreateUserResponse
	err = json.Unmarshal(respBody, &createUserResp)
	require.NoError(t, err)

	// Verify the ID is a valid positive int64 (snowflake IDs are positive)
	assert.True(t, createUserResp.Id > 0, "Snowflake ID should be positive, got %d", createUserResp.Id)

	// Snowflake IDs are typically large numbers (> 1 million due to timestamp component)
	// The epoch is 2020-01-01, so any ID generated after that will be substantial
	assert.True(t, createUserResp.Id > 1000000, "Snowflake ID should be a large number, got %d", createUserResp.Id)
}

// TestCreateUser_IDsMonotonicallyIncreasing verifies that IDs are monotonically increasing
func TestCreateUser_IDsMonotonicallyIncreasing(t *testing.T) {
	env := setupTestEnv(t)
	defer teardownTestEnv(t, env)

	var previousID int64 = 0

	for i := 0; i < 5; i++ {
		createUserReq := openapi.CreateUserRequest{
			Username:  fmt.Sprintf("monotonicuser%d", i),
			Password:  "TestPassword123!",
			FirstName: "Monotonic",
			LastName:  fmt.Sprintf("User%d", i),
			Email:     fmt.Sprintf("monotonic%d@example.com", i),
		}

		body, err := json.Marshal(createUserReq)
		require.NoError(t, err)

		resp, err := http.Post(
			env.BaseURL+"/v1/users",
			"application/json",
			bytes.NewReader(body),
		)
		require.NoError(t, err)

		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, resp.StatusCode, "Failed: %s", string(respBody))

		var createUserResp openapi.CreateUserResponse
		err = json.Unmarshal(respBody, &createUserResp)
		require.NoError(t, err)

		// Verify IDs are monotonically increasing
		assert.True(t, createUserResp.Id > previousID,
			"ID should be greater than previous: got %d, previous was %d", createUserResp.Id, previousID)
		previousID = createUserResp.Id
	}
}
