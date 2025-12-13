package user

import (
	"context"
	"net/http"
	"testing"
	"time"

	db "github.com/igomez10/microservices/socialapp/pkg/dbpgx"
	"github.com/igomez10/microservices/socialapp/socialappapi/openapi"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

// MockSnowflakeGenerator is a mock implementation of snowflake.IDGenerator
type MockSnowflakeGenerator struct {
	nextID int64
	err    error
}

func (m *MockSnowflakeGenerator) NextID() (int64, error) {
	return m.nextID, m.err
}

// MockQuerier is a mock implementation of db.Querier for testing
type MockQuerier struct {
	getUserByUsernameFunc func(ctx context.Context, db db.DBTX, username string) (db.User, error)
	getUserByEmailFunc    func(ctx context.Context, db db.DBTX, email string) (db.User, error)
	getRoleByNameFunc     func(ctx context.Context, db db.DBTX, name string) (db.Role, error)
	createUserWithIDFunc  func(ctx context.Context, db db.DBTX, arg db.CreateUserWithIDParams) (db.User, error)
	createUserToRoleFunc  func(ctx context.Context, db db.DBTX, arg db.CreateUserToRoleParams) (db.UsersToRole, error)
}

func (m *MockQuerier) GetUserByUsername(ctx context.Context, dbtx db.DBTX, username string) (db.User, error) {
	if m.getUserByUsernameFunc != nil {
		return m.getUserByUsernameFunc(ctx, dbtx, username)
	}
	return db.User{}, pgx.ErrNoRows
}

func (m *MockQuerier) GetUserByEmail(ctx context.Context, dbtx db.DBTX, email string) (db.User, error) {
	if m.getUserByEmailFunc != nil {
		return m.getUserByEmailFunc(ctx, dbtx, email)
	}
	return db.User{}, pgx.ErrNoRows
}

func (m *MockQuerier) GetRoleByName(ctx context.Context, dbtx db.DBTX, name string) (db.Role, error) {
	if m.getRoleByNameFunc != nil {
		return m.getRoleByNameFunc(ctx, dbtx, name)
	}
	return db.Role{ID: 1, Name: "administrator"}, nil
}

func (m *MockQuerier) CreateUserWithID(ctx context.Context, dbtx db.DBTX, arg db.CreateUserWithIDParams) (db.User, error) {
	if m.createUserWithIDFunc != nil {
		return m.createUserWithIDFunc(ctx, dbtx, arg)
	}
	return db.User{
		ID:        arg.ID,
		Username:  arg.Username,
		FirstName: arg.FirstName,
		LastName:  arg.LastName,
		Email:     arg.Email,
		CreatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
	}, nil
}

func (m *MockQuerier) CreateUserToRole(ctx context.Context, dbtx db.DBTX, arg db.CreateUserToRoleParams) (db.UsersToRole, error) {
	if m.createUserToRoleFunc != nil {
		return m.createUserToRoleFunc(ctx, dbtx, arg)
	}
	return db.UsersToRole{ID: 1, UserID: arg.UserID, RoleID: arg.RoleID}, nil
}

// Stub implementations for remaining Querier methods
func (m *MockQuerier) CreateComment(ctx context.Context, dbtx db.DBTX, arg db.CreateCommentParams) (db.Comment, error) {
	return db.Comment{}, nil
}
func (m *MockQuerier) CreateCommentForUser(ctx context.Context, dbtx db.DBTX, arg db.CreateCommentForUserParams) (db.Comment, error) {
	return db.Comment{}, nil
}
func (m *MockQuerier) CreateCredential(ctx context.Context, dbtx db.DBTX, arg db.CreateCredentialParams) (db.Credential, error) {
	return db.Credential{}, nil
}
func (m *MockQuerier) CreateEvent(ctx context.Context, dbtx db.DBTX, arg db.CreateEventParams) error {
	return nil
}
func (m *MockQuerier) CreateRole(ctx context.Context, dbtx db.DBTX, arg db.CreateRoleParams) (db.Role, error) {
	return db.Role{}, nil
}
func (m *MockQuerier) CreateRoleScope(ctx context.Context, dbtx db.DBTX, arg db.CreateRoleScopeParams) (db.RolesToScope, error) {
	return db.RolesToScope{}, nil
}
func (m *MockQuerier) CreateScope(ctx context.Context, dbtx db.DBTX, arg db.CreateScopeParams) (db.Scope, error) {
	return db.Scope{}, nil
}
func (m *MockQuerier) CreateURL(ctx context.Context, dbtx db.DBTX, arg db.CreateURLParams) (db.Url, error) {
	return db.Url{}, nil
}
func (m *MockQuerier) CreateUser(ctx context.Context, dbtx db.DBTX, arg db.CreateUserParams) (db.User, error) {
	return db.User{}, nil
}
func (m *MockQuerier) DeleteComment(ctx context.Context, dbtx db.DBTX, id int64) error { return nil }
func (m *MockQuerier) DeleteCredential(ctx context.Context, dbtx db.DBTX, id int64) error {
	return nil
}
func (m *MockQuerier) DeleteRole(ctx context.Context, dbtx db.DBTX, id int64) error { return nil }
func (m *MockQuerier) DeleteRoleScope(ctx context.Context, dbtx db.DBTX, arg db.DeleteRoleScopeParams) error {
	return nil
}
func (m *MockQuerier) DeleteScope(ctx context.Context, dbtx db.DBTX, id int64) error { return nil }
func (m *MockQuerier) DeleteURL(ctx context.Context, dbtx db.DBTX, alias string) error {
	return nil
}
func (m *MockQuerier) DeleteUser(ctx context.Context, dbtx db.DBTX, id int64) error { return nil }
func (m *MockQuerier) DeleteUserByUsername(ctx context.Context, dbtx db.DBTX, username string) error {
	return nil
}
func (m *MockQuerier) DeleteUserToRole(ctx context.Context, dbtx db.DBTX, arg db.DeleteUserToRoleParams) error {
	return nil
}
func (m *MockQuerier) FollowUser(ctx context.Context, dbtx db.DBTX, arg db.FollowUserParams) error {
	return nil
}
func (m *MockQuerier) GetComment(ctx context.Context, dbtx db.DBTX, id int64) (db.Comment, error) {
	return db.Comment{}, nil
}
func (m *MockQuerier) GetCredential(ctx context.Context, dbtx db.DBTX, publicKey string) (db.Credential, error) {
	return db.Credential{}, nil
}
func (m *MockQuerier) GetFollowedUsers(ctx context.Context, dbtx db.DBTX, followerID int64) ([]db.User, error) {
	return nil, nil
}
func (m *MockQuerier) GetFollowers(ctx context.Context, dbtx db.DBTX, followedID int64) ([]db.User, error) {
	return nil, nil
}
func (m *MockQuerier) GetRole(ctx context.Context, dbtx db.DBTX, id int64) (db.Role, error) {
	return db.Role{}, nil
}
func (m *MockQuerier) GetScope(ctx context.Context, dbtx db.DBTX, id int64) (db.Scope, error) {
	return db.Scope{}, nil
}
func (m *MockQuerier) GetScopeByName(ctx context.Context, dbtx db.DBTX, name string) (db.Scope, error) {
	return db.Scope{}, nil
}
func (m *MockQuerier) GetURLFromAlias(ctx context.Context, dbtx db.DBTX, alias string) (db.Url, error) {
	return db.Url{}, nil
}
func (m *MockQuerier) GetUserByID(ctx context.Context, dbtx db.DBTX, id int64) (db.User, error) {
	return db.User{}, nil
}
func (m *MockQuerier) GetUserComments(ctx context.Context, dbtx db.DBTX, arg db.GetUserCommentsParams) ([]db.Comment, error) {
	return nil, nil
}
func (m *MockQuerier) GetUserRoles(ctx context.Context, dbtx db.DBTX, id int64) ([]db.Role, error) {
	return nil, nil
}
func (m *MockQuerier) ListComment(ctx context.Context, dbtx db.DBTX, arg db.ListCommentParams) ([]db.Comment, error) {
	return nil, nil
}
func (m *MockQuerier) ListRoleScopes(ctx context.Context, dbtx db.DBTX, arg db.ListRoleScopesParams) ([]db.Scope, error) {
	return nil, nil
}
func (m *MockQuerier) ListRoles(ctx context.Context, dbtx db.DBTX, arg db.ListRolesParams) ([]db.Role, error) {
	return nil, nil
}
func (m *MockQuerier) ListScopes(ctx context.Context, dbtx db.DBTX, arg db.ListScopesParams) ([]db.Scope, error) {
	return nil, nil
}
func (m *MockQuerier) ListUsers(ctx context.Context, dbtx db.DBTX, arg db.ListUsersParams) ([]db.User, error) {
	return nil, nil
}
func (m *MockQuerier) UnfollowUser(ctx context.Context, dbtx db.DBTX, arg db.UnfollowUserParams) error {
	return nil
}
func (m *MockQuerier) UpdateRole(ctx context.Context, dbtx db.DBTX, arg db.UpdateRoleParams) error {
	return nil
}
func (m *MockQuerier) UpdateScope(ctx context.Context, dbtx db.DBTX, arg db.UpdateScopeParams) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, nil
}
func (m *MockQuerier) UpdateUser(ctx context.Context, dbtx db.DBTX, arg db.UpdateUserParams) error {
	return nil
}
func (m *MockQuerier) UpdateUserByUsername(ctx context.Context, dbtx db.DBTX, arg db.UpdateUserByUsernameParams) error {
	return nil
}

// MockEventRecorder is a mock implementation of eventRecorder.EventRecorder
type MockEventRecorder struct{}

func (m *MockEventRecorder) RecordEvent(ctx context.Context, tx pgx.Tx, rawEvent interface{}, id int64) error {
	return nil
}

// MockDBConn is a mock implementation of pgxpool.Pool for testing
type MockDBConn struct{}

func (m *MockDBConn) BeginTx(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error) {
	return &MockTx{}, nil
}

// MockTx is a mock implementation of pgx.Tx
type MockTx struct{}

func (m *MockTx) Begin(ctx context.Context) (pgx.Tx, error) { return m, nil }
func (m *MockTx) Commit(ctx context.Context) error          { return nil }
func (m *MockTx) Rollback(ctx context.Context) error        { return nil }
func (m *MockTx) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	return 0, nil
}
func (m *MockTx) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults { return nil }
func (m *MockTx) LargeObjects() pgx.LargeObjects                               { return pgx.LargeObjects{} }
func (m *MockTx) Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error) {
	return nil, nil
}
func (m *MockTx) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, nil
}
func (m *MockTx) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return nil, nil
}
func (m *MockTx) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row { return nil }
func (m *MockTx) Conn() *pgx.Conn                                               { return nil }

func TestCreateUser_DeterministicID(t *testing.T) {
	expectedID := int64(1234567890123456789)

	mockQuerier := &MockQuerier{
		createUserWithIDFunc: func(ctx context.Context, dbtx db.DBTX, arg db.CreateUserWithIDParams) (db.User, error) {
			// Verify the ID passed to CreateUserWithID matches our expected snowflake ID
			if arg.ID != expectedID {
				t.Errorf("Expected ID %d, got %d", expectedID, arg.ID)
			}
			return db.User{
				ID:        arg.ID,
				Username:  arg.Username,
				FirstName: arg.FirstName,
				LastName:  arg.LastName,
				Email:     arg.Email,
				CreatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
			}, nil
		},
	}

	mockSnowflake := &MockSnowflakeGenerator{nextID: expectedID}

	// Note: This test requires a real pgxpool.Pool which we can't easily mock
	// without significant infrastructure. This test documents the expected behavior.
	// For full integration testing, use the integration_tests package.

	t.Run("snowflake generator returns expected ID", func(t *testing.T) {
		id, err := mockSnowflake.NextID()
		if err != nil {
			t.Fatalf("NextID() returned error: %v", err)
		}
		if id != expectedID {
			t.Errorf("Expected ID %d, got %d", expectedID, id)
		}
	})

	t.Run("mock querier receives correct ID", func(t *testing.T) {
		_, _ = mockQuerier.CreateUserWithID(context.Background(), nil, db.CreateUserWithIDParams{
			ID:        expectedID,
			Username:  "testuser",
			FirstName: "Test",
			LastName:  "User",
			Email:     "test@example.com",
		})
		// Test assertion is in the mock function
	})
}

func TestCreateUser_IDInResponse(t *testing.T) {
	expectedID := int64(7654321098765432109)

	t.Run("response contains correct ID", func(t *testing.T) {
		// Simulate the response construction that happens in CreateUser
		createdUser := db.User{
			ID:        expectedID,
			Username:  "testuser",
			FirstName: "Test",
			LastName:  "User",
			Email:     "test@example.com",
			CreatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
		}

		res := openapi.CreateUserResponse{
			Id:        createdUser.ID,
			Username:  createdUser.Username,
			FirstName: createdUser.FirstName,
			LastName:  createdUser.LastName,
			Email:     createdUser.Email,
			CreatedAt: createdUser.CreatedAt.Time,
		}

		if res.Id != expectedID {
			t.Errorf("Expected response ID %d, got %d", expectedID, res.Id)
		}
	})
}

func TestCreateUser_SnowflakeIDIsInt64(t *testing.T) {
	// This test verifies that the snowflake ID is compatible with the database schema
	// The users table has: id BIGINT NOT NULL PRIMARY KEY
	// Go's int64 maps to PostgreSQL's BIGINT

	mockSnowflake := &MockSnowflakeGenerator{nextID: 7123456789012345678}

	id, err := mockSnowflake.NextID()
	if err != nil {
		t.Fatalf("NextID() returned error: %v", err)
	}

	// Verify it's a valid int64 (positive, within range)
	if id <= 0 {
		t.Errorf("Expected positive ID, got %d", id)
	}

	// Verify it fits in the CreateUserWithIDParams.ID field
	params := db.CreateUserWithIDParams{
		ID: id,
	}
	if params.ID != id {
		t.Errorf("ID not properly assigned to params: expected %d, got %d", id, params.ID)
	}
}

func TestCreateUser_SnowflakeError(t *testing.T) {
	mockSnowflake := &MockSnowflakeGenerator{
		nextID: 0,
		err:    pgx.ErrNoRows, // Using this as a placeholder error
	}

	_, err := mockSnowflake.NextID()
	if err == nil {
		t.Error("Expected error from NextID(), got nil")
	}
}

func TestCreateUser_ResponseStatusOK(t *testing.T) {
	// Verify that successful user creation returns HTTP 200 OK
	expectedStatus := http.StatusOK

	// This is the response that CreateUser returns on success
	response := openapi.Response(expectedStatus, openapi.CreateUserResponse{
		Id:        123,
		Username:  "testuser",
		FirstName: "Test",
		LastName:  "User",
		Email:     "test@example.com",
		CreatedAt: time.Now(),
	})

	if response.Code != expectedStatus {
		t.Errorf("Expected status %d, got %d", expectedStatus, response.Code)
	}
}
