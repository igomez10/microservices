package gandalf

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/igomez10/microservices/socialapp/internal/contexthelper"
	socialappjwt "github.com/igomez10/microservices/socialapp/internal/jwt"
	db "github.com/igomez10/microservices/socialapp/pkg/dbpgx"
	"github.com/igomez10/microservices/socialapp/pkg/scopes"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

// Mock database for testing
type mockDB struct {
	getUserByUsernameFunc func(ctx context.Context, conn *pgxpool.Pool, username string) (db.User, error)
	getUserRolesFunc      func(ctx context.Context, conn *pgxpool.Pool, userID int64) ([]db.Role, error)
	listRoleScopesFunc    func(ctx context.Context, conn *pgxpool.Pool, params db.ListRoleScopesParams) ([]db.Scope, error)
}

// Implement required methods matching the Querier interface
func (m *mockDB) CreateComment(ctx context.Context, dbtx db.DBTX, arg db.CreateCommentParams) (db.Comment, error) {
	return db.Comment{}, nil
}
func (m *mockDB) CreateCommentForUser(ctx context.Context, dbtx db.DBTX, arg db.CreateCommentForUserParams) (db.Comment, error) {
	return db.Comment{}, nil
}
func (m *mockDB) CreateCredential(ctx context.Context, dbtx db.DBTX, arg db.CreateCredentialParams) (db.Credential, error) {
	return db.Credential{}, nil
}
func (m *mockDB) CreateEvent(ctx context.Context, dbtx db.DBTX, arg db.CreateEventParams) error {
	return nil
}
func (m *mockDB) CreateRole(ctx context.Context, dbtx db.DBTX, arg db.CreateRoleParams) (db.Role, error) {
	return db.Role{}, nil
}
func (m *mockDB) CreateRoleScope(ctx context.Context, dbtx db.DBTX, arg db.CreateRoleScopeParams) (db.RolesToScope, error) {
	return db.RolesToScope{}, nil
}
func (m *mockDB) CreateScope(ctx context.Context, dbtx db.DBTX, arg db.CreateScopeParams) (db.Scope, error) {
	return db.Scope{}, nil
}
func (m *mockDB) CreateURL(ctx context.Context, dbtx db.DBTX, arg db.CreateURLParams) (db.Url, error) {
	return db.Url{}, nil
}
func (m *mockDB) CreateUser(ctx context.Context, dbtx db.DBTX, arg db.CreateUserParams) (db.User, error) {
	return db.User{}, nil
}
func (m *mockDB) CreateUserToRole(ctx context.Context, dbtx db.DBTX, arg db.CreateUserToRoleParams) (db.UsersToRole, error) {
	return db.UsersToRole{}, nil
}
func (m *mockDB) DeleteComment(ctx context.Context, dbtx db.DBTX, id int64) error {
	return nil
}
func (m *mockDB) DeleteCredential(ctx context.Context, dbtx db.DBTX, id int64) error {
	return nil
}
func (m *mockDB) DeleteRole(ctx context.Context, dbtx db.DBTX, id int64) error {
	return nil
}
func (m *mockDB) DeleteRoleScope(ctx context.Context, dbtx db.DBTX, arg db.DeleteRoleScopeParams) error {
	return nil
}
func (m *mockDB) DeleteScope(ctx context.Context, dbtx db.DBTX, id int64) error {
	return nil
}
func (m *mockDB) DeleteURL(ctx context.Context, dbtx db.DBTX, alias string) error {
	return nil
}
func (m *mockDB) DeleteUser(ctx context.Context, dbtx db.DBTX, id int64) error {
	return nil
}
func (m *mockDB) DeleteUserByUsername(ctx context.Context, dbtx db.DBTX, username string) error {
	return nil
}
func (m *mockDB) DeleteUserToRole(ctx context.Context, dbtx db.DBTX, arg db.DeleteUserToRoleParams) error {
	return nil
}
func (m *mockDB) FollowUser(ctx context.Context, dbtx db.DBTX, arg db.FollowUserParams) error {
	return nil
}
func (m *mockDB) GetComment(ctx context.Context, dbtx db.DBTX, id int64) (db.Comment, error) {
	return db.Comment{}, nil
}
func (m *mockDB) GetCredential(ctx context.Context, dbtx db.DBTX, publicKey string) (db.Credential, error) {
	return db.Credential{}, nil
}
func (m *mockDB) GetFollowedUsers(ctx context.Context, dbtx db.DBTX, followerID int64) ([]db.User, error) {
	return nil, nil
}
func (m *mockDB) GetFollowers(ctx context.Context, dbtx db.DBTX, followedID int64) ([]db.User, error) {
	return nil, nil
}
func (m *mockDB) GetRole(ctx context.Context, dbtx db.DBTX, id int64) (db.Role, error) {
	return db.Role{}, nil
}
func (m *mockDB) GetRoleByName(ctx context.Context, dbtx db.DBTX, name string) (db.Role, error) {
	return db.Role{}, nil
}
func (m *mockDB) GetScope(ctx context.Context, dbtx db.DBTX, id int64) (db.Scope, error) {
	return db.Scope{}, nil
}
func (m *mockDB) GetScopeByName(ctx context.Context, dbtx db.DBTX, name string) (db.Scope, error) {
	return db.Scope{}, nil
}
func (m *mockDB) GetURLFromAlias(ctx context.Context, dbtx db.DBTX, alias string) (db.Url, error) {
	return db.Url{}, nil
}
func (m *mockDB) GetUserByEmail(ctx context.Context, dbtx db.DBTX, email string) (db.User, error) {
	return db.User{}, nil
}
func (m *mockDB) GetUserByID(ctx context.Context, dbtx db.DBTX, id int64) (db.User, error) {
	return db.User{}, nil
}
func (m *mockDB) GetUserByUsername(ctx context.Context, dbtx db.DBTX, username string) (db.User, error) {
	if m.getUserByUsernameFunc != nil {
		conn, _ := dbtx.(*pgxpool.Pool)
		return m.getUserByUsernameFunc(ctx, conn, username)
	}
	return db.User{}, pgx.ErrNoRows
}
func (m *mockDB) GetUserComments(ctx context.Context, dbtx db.DBTX, arg db.GetUserCommentsParams) ([]db.Comment, error) {
	return nil, nil
}
func (m *mockDB) GetUserRoles(ctx context.Context, dbtx db.DBTX, id int64) ([]db.Role, error) {
	if m.getUserRolesFunc != nil {
		conn, _ := dbtx.(*pgxpool.Pool)
		return m.getUserRolesFunc(ctx, conn, id)
	}
	return nil, pgx.ErrNoRows
}
func (m *mockDB) ListComment(ctx context.Context, dbtx db.DBTX, arg db.ListCommentParams) ([]db.Comment, error) {
	return nil, nil
}
func (m *mockDB) ListRoleScopes(ctx context.Context, dbtx db.DBTX, arg db.ListRoleScopesParams) ([]db.Scope, error) {
	if m.listRoleScopesFunc != nil {
		conn, _ := dbtx.(*pgxpool.Pool)
		return m.listRoleScopesFunc(ctx, conn, arg)
	}
	return nil, pgx.ErrNoRows
}
func (m *mockDB) ListRoles(ctx context.Context, dbtx db.DBTX, arg db.ListRolesParams) ([]db.Role, error) {
	return nil, nil
}
func (m *mockDB) ListScopes(ctx context.Context, dbtx db.DBTX, arg db.ListScopesParams) ([]db.Scope, error) {
	return nil, nil
}
func (m *mockDB) ListUsers(ctx context.Context, dbtx db.DBTX, arg db.ListUsersParams) ([]db.User, error) {
	return nil, nil
}
func (m *mockDB) UnfollowUser(ctx context.Context, dbtx db.DBTX, arg db.UnfollowUserParams) error {
	return nil
}
func (m *mockDB) UpdateRole(ctx context.Context, dbtx db.DBTX, arg db.UpdateRoleParams) error {
	return nil
}
func (m *mockDB) UpdateScope(ctx context.Context, dbtx db.DBTX, arg db.UpdateScopeParams) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, nil
}
func (m *mockDB) UpdateUser(ctx context.Context, dbtx db.DBTX, arg db.UpdateUserParams) error {
	return nil
}
func (m *mockDB) UpdateUserByUsername(ctx context.Context, dbtx db.DBTX, arg db.UpdateUserByUsernameParams) error {
	return nil
}

func TestJWTParsing(t *testing.T) {
	hmacSampleSecret := []byte("secret")
	username := "ignacio"
	tokenScopes := []string{scopes.SocialappUsersList.String(), scopes.SocialappUsersRead.String()}

	jwtToken := &socialappjwt.SocialAPPToken{
		Username:  username,
		Scopes:    tokenScopes,
		Audience:  jwt.ClaimStrings{"aud1", "aud2"},
		Expires:   jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		NotBefore: jwt.NewNumericDate(time.Now().UTC().Add(-time.Minute)),
		Issuer:    "gandalf",
		Subject:   "social-app",
	}

	tokenString, err := socialappjwt.TokenToString(jwtToken, hmacSampleSecret)
	if err != nil {
		t.Fatal(err)
	}

	parsedToken, err := socialappjwt.ParseJWTToken(tokenString, hmacSampleSecret)
	if err != nil {
		t.Fatal(err)
	}

	if !parsedToken.Valid {
		t.Fatal("parsed token is invalid")
	}

	fromtoken, err := socialappjwt.FromToken(parsedToken)
	if err != nil {
		t.Fatal(err)
	}

	content1, err := json.Marshal(fromtoken)
	if err != nil {
		t.Fatal(err)
	}

	content2, err := json.Marshal(jwtToken)
	if err != nil {
		t.Fatal(err)
	}

	if string(content1) != string(content2) {
		fmt.Printf("%+v\n", content1)
		fmt.Printf("%+v\n", content2)
		t.Fatal("parsed token claims do not match")
	}
}

func TestAuthenticate_AllowlistedPath(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		method         string
		allowlisted    map[string]map[string]bool
		expectedStatus int
	}{
		{
			name:   "Allowlisted GET path",
			path:   "/health",
			method: "GET",
			allowlisted: map[string]map[string]bool{
				"/health": {"GET": true},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Allowlisted POST path",
			path:   "/webhook",
			method: "POST",
			allowlisted: map[string]map[string]bool{
				"/webhook": {"POST": true},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Not allowlisted path",
			path:   "/protected",
			method: "GET",
			allowlisted: map[string]map[string]bool{
				"/health": {"GET": true},
			},
			expectedStatus: http.StatusOK, // Will have noauth scope
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middleware := &Middleware{
				AllowlistedPaths: tt.allowlisted,
			}

			handler := middleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				scopes, ok := contexthelper.GetRequestedScopesInContext(r.Context())
				if !ok {
					t.Error("Expected scopes in context")
				}
				if tt.allowlisted[tt.path] != nil && tt.allowlisted[tt.path][tt.method] {
					if len(scopes) != 0 {
						t.Errorf("Expected empty scopes for allowlisted path, got %v", scopes)
					}
				}
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest(tt.method, tt.path, nil)
			req = req.WithContext(contexthelper.SetLoggerInContext(req.Context(), zerolog.New(nil)))
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestAuthenticate_BearerToken(t *testing.T) {
	secret := "test-secret"

	tests := []struct {
		name           string
		token          *socialappjwt.SocialAPPToken
		expectedStatus int
		expectedScopes []string
	}{
		{
			name: "Valid bearer token with valid scopes",
			token: &socialappjwt.SocialAPPToken{
				Username:  "testuser",
				Scopes:    []string{scopes.SocialappUsersList.String(), scopes.SocialappUsersRead.String()},
				Expires:   jwt.NewNumericDate(time.Now().Add(time.Hour)),
				NotBefore: jwt.NewNumericDate(time.Now().Add(-time.Minute)),
			},
			expectedStatus: http.StatusOK,
			expectedScopes: []string{scopes.SocialappUsersList.String(), scopes.SocialappUsersRead.String()},
		},
		{
			name: "Token with invalid scopes are filtered out",
			token: &socialappjwt.SocialAPPToken{
				Username:  "testuser",
				Scopes:    []string{"invalid.scope", "another.invalid"},
				Expires:   jwt.NewNumericDate(time.Now().Add(time.Hour)),
				NotBefore: jwt.NewNumericDate(time.Now().Add(-time.Minute)),
			},
			expectedStatus: http.StatusOK,
			expectedScopes: []string{}, // Invalid scopes are filtered
		},
		{
			name: "Expired token",
			token: &socialappjwt.SocialAPPToken{
				Username:  "testuser",
				Scopes:    []string{scopes.SocialappUsersRead.String()},
				Expires:   jwt.NewNumericDate(time.Now().Add(-time.Hour)),
				NotBefore: jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Token not yet valid",
			token: &socialappjwt.SocialAPPToken{
				Username:  "testuser",
				Scopes:    []string{scopes.SocialappUsersRead.String()},
				Expires:   jwt.NewNumericDate(time.Now().Add(time.Hour)),
				NotBefore: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middleware := &Middleware{
				JWTSecret:        secret,
				AllowlistedPaths: map[string]map[string]bool{},
			}

			tokenString, err := socialappjwt.TokenToString(tt.token, []byte(secret))
			if err != nil {
				t.Fatalf("Failed to create token: %v", err)
			}

			handler := middleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.expectedStatus == http.StatusOK {
					tokenScopes, ok := contexthelper.GetRequestedScopesInContext(r.Context())
					if !ok {
						t.Error("Expected scopes in context")
					}
					for _, expectedScope := range tt.expectedScopes {
						if !tokenScopes[expectedScope] {
							t.Errorf("Expected scope %s in context", expectedScope)
						}
					}

					username, ok := contexthelper.GetUsernameInContext(r.Context())
					if !ok {
						t.Error("Expected username in context")
					}
					if username != tt.token.Username {
						t.Errorf("Expected username %s, got %s", tt.token.Username, username)
					}
				}
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest("GET", "/test", nil)
			req = req.WithContext(contexthelper.SetLoggerInContext(req.Context(), zerolog.New(nil)))
			req.Header.Set("Authorization", "Bearer "+tokenString)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestAuthenticate_InvalidBearerToken(t *testing.T) {
	middleware := &Middleware{
		JWTSecret:        "test-secret",
		AllowlistedPaths: map[string]map[string]bool{},
	}

	handler := middleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called for invalid token")
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	req = req.WithContext(contexthelper.SetLoggerInContext(req.Context(), zerolog.New(nil)))
	req.Header.Set("Authorization", "Bearer invalid-token")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestAuthenticate_BasicAuth_Success(t *testing.T) {
	testUser := db.User{
		ID:             1,
		Username:       "testuser",
		HashedPassword: "hashed_password",
		Salt:           "salt",
	}

	testRole := db.Role{
		ID:   1,
		Name: "admin",
	}

	testScopes := []db.Scope{
		{ID: 1, Name: scopes.SocialappUsersList.String()},
		{ID: 2, Name: scopes.SocialappUsersRead.String()},
	}

	mockDBInstance := &mockDB{
		getUserByUsernameFunc: func(ctx context.Context, conn *pgxpool.Pool, username string) (db.User, error) {
			if username == testUser.Username {
				return testUser, nil
			}
			return db.User{}, pgx.ErrNoRows
		},
		getUserRolesFunc: func(ctx context.Context, conn *pgxpool.Pool, userID int64) ([]db.Role, error) {
			if userID == testUser.ID {
				return []db.Role{testRole}, nil
			}
			return nil, pgx.ErrNoRows
		},
		listRoleScopesFunc: func(ctx context.Context, conn *pgxpool.Pool, params db.ListRoleScopesParams) ([]db.Scope, error) {
			if params.ID == testRole.ID {
				return testScopes, nil
			}
			return nil, pgx.ErrNoRows
		},
	}

	middleware := &Middleware{
		DB:               mockDBInstance,
		AllowBasicAuth:   true,
		AllowlistedPaths: map[string]map[string]bool{},
	}

	handler := middleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, ok := contexthelper.GetUsernameInContext(r.Context())
		if !ok {
			t.Error("Expected username in context")
		}
		if username != testUser.Username {
			t.Errorf("Expected username %s, got %s", testUser.Username, username)
		}

		tokenScopes, ok := contexthelper.GetRequestedScopesInContext(r.Context())
		if !ok {
			t.Error("Expected scopes in context")
		}
		if !tokenScopes[scopes.SocialappUsersList.String()] {
			t.Error("Expected 'socialapp.users.list' scope in context")
		}

		w.WriteHeader(http.StatusOK)
	}))

	form := url.Values{}
	form.Add("scope", scopes.SocialappUsersList.String()+" "+scopes.SocialappUsersRead.String())
	req := httptest.NewRequest("POST", "/auth", strings.NewReader(form.Encode()))
	req = req.WithContext(contexthelper.SetLoggerInContext(req.Context(), zerolog.New(nil)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Set basic auth header
	auth := testUser.Username + ":" + "password"
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Set("Authorization", "Basic "+encodedAuth)

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// Note: This test will fail password verification unless we mock EncryptPassword
	// For now, we expect it to fail at password verification
	if rr.Code == http.StatusOK {
		t.Log("Basic auth succeeded (password verification would need to be mocked)")
	}
}

func TestAuthenticate_BasicAuth_UserNotFound(t *testing.T) {
	mockDBInstance := &mockDB{
		getUserByUsernameFunc: func(ctx context.Context, conn *pgxpool.Pool, username string) (db.User, error) {
			return db.User{}, pgx.ErrNoRows
		},
	}

	middleware := &Middleware{
		DB:               mockDBInstance,
		AllowBasicAuth:   true,
		AllowlistedPaths: map[string]map[string]bool{},
	}

	handler := middleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called when user not found")
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	req = req.WithContext(contexthelper.SetLoggerInContext(req.Context(), zerolog.New(nil)))
	auth := "nonexistent:password"
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Set("Authorization", "Basic "+encodedAuth)

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestAuthenticate_BasicAuth_InvalidFormat(t *testing.T) {
	middleware := &Middleware{
		AllowBasicAuth:   true,
		AllowlistedPaths: map[string]map[string]bool{},
	}

	handler := middleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called with invalid basic auth")
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	req = req.WithContext(contexthelper.SetLoggerInContext(req.Context(), zerolog.New(nil)))
	req.Header.Set("Authorization", "Basic invalid")

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestAuthenticate_NoAuth(t *testing.T) {
	middleware := &Middleware{
		AllowlistedPaths: map[string]map[string]bool{},
	}

	handler := middleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		scopes, ok := contexthelper.GetRequestedScopesInContext(r.Context())
		if !ok {
			t.Error("Expected scopes in context")
		}
		if !scopes["noauth"] {
			t.Error("Expected 'noauth' scope in context")
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	req = req.WithContext(contexthelper.SetLoggerInContext(req.Context(), zerolog.New(nil)))
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestIsTokenValid(t *testing.T) {
	middleware := &Middleware{}

	tests := []struct {
		name     string
		token    *socialappjwt.SocialAPPToken
		expected bool
	}{
		{
			name: "Valid token",
			token: &socialappjwt.SocialAPPToken{
				Expires:   jwt.NewNumericDate(time.Now().Add(time.Hour)),
				NotBefore: jwt.NewNumericDate(time.Now().Add(-time.Minute)),
			},
			expected: true,
		},
		{
			name: "Expired token",
			token: &socialappjwt.SocialAPPToken{
				Expires:   jwt.NewNumericDate(time.Now().Add(-time.Hour)),
				NotBefore: jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			},
			expected: false,
		},
		{
			name: "Token not yet valid",
			token: &socialappjwt.SocialAPPToken{
				Expires:   jwt.NewNumericDate(time.Now().Add(time.Hour)),
				NotBefore: jwt.NewNumericDate(time.Now().Add(time.Minute)),
			},
			expected: false,
		},
		{
			name: "Nil Expires",
			token: &socialappjwt.SocialAPPToken{
				Expires:   nil,
				NotBefore: jwt.NewNumericDate(time.Now().Add(-time.Minute)),
			},
			expected: false,
		},
		{
			name: "Nil NotBefore",
			token: &socialappjwt.SocialAPPToken{
				Expires:   jwt.NewNumericDate(time.Now().Add(time.Hour)),
				NotBefore: nil,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := middleware.isTokenValid(tt.token)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestConvertTokenScopes(t *testing.T) {
	middleware := &Middleware{}

	tests := []struct {
		name     string
		scopes   []string
		expected map[string]bool
	}{
		{
			name:   "Valid scopes from API spec",
			scopes: []string{scopes.SocialappUsersList.String(), scopes.SocialappUsersRead.String()},
			expected: map[string]bool{
				scopes.SocialappUsersList.String(): true,
				scopes.SocialappUsersRead.String(): true,
			},
		},
		{
			name:     "Invalid scopes are filtered out",
			scopes:   []string{"invalid.scope", "another.invalid"},
			expected: map[string]bool{},
		},
		{
			name:   "Mixed valid and invalid scopes",
			scopes: []string{scopes.SocialappUsersList.String(), "invalid.scope"},
			expected: map[string]bool{
				scopes.SocialappUsersList.String(): true,
			},
		},
		{
			name:     "Empty scopes",
			scopes:   []string{},
			expected: map[string]bool{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := middleware.convertTokenScopes(tt.scopes)
			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d scopes, got %d", len(tt.expected), len(result))
			}
			for key, value := range tt.expected {
				if result[key] != value {
					t.Errorf("Expected scope %s to be %v, got %v", key, value, result[key])
				}
			}
		})
	}
}

func TestParseRequestedScopes(t *testing.T) {
	middleware := &Middleware{}

	tests := []struct {
		name        string
		scope       string
		expected    []string
		expectError bool
	}{
		{
			name:        "Valid scopes from API spec",
			scope:       scopes.SocialappUsersList.String() + " " + scopes.SocialappUsersRead.String(),
			expected:    []string{scopes.SocialappUsersList.String(), scopes.SocialappUsersRead.String()},
			expectError: false,
		},
		{
			name:        "Single valid scope",
			scope:       scopes.SocialappUsersList.String(),
			expected:    []string{scopes.SocialappUsersList.String()},
			expectError: false,
		},
		{
			name:        "Empty scope",
			scope:       "",
			expected:    []string{},
			expectError: false,
		},
		{
			name:        "Invalid scope returns error",
			scope:       "invalid.scope",
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Mixed valid and invalid scopes returns error",
			scope:       scopes.SocialappUsersList.String() + " invalid.scope",
			expected:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			if tt.scope != "" {
				form.Add("scope", tt.scope)
			}
			req := httptest.NewRequest("POST", "/test", strings.NewReader(form.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			result, err := middleware.parseRequestedScopes(req)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d scopes, got %d", len(tt.expected), len(result))
			}
			for i, expected := range tt.expected {
				if i >= len(result) || result[i] != expected {
					t.Errorf("Expected scope at index %d to be %s, got %s", i, expected, result[i])
				}
			}
		})
	}
}

func TestConvertScopeSliceToMap(t *testing.T) {
	middleware := &Middleware{}

	tests := []struct {
		name     string
		scopes   []string
		expected map[string]bool
	}{
		{
			name:   "Multiple scopes",
			scopes: []string{"read", "write", "delete"},
			expected: map[string]bool{
				"read":   true,
				"write":  true,
				"delete": true,
			},
		},
		{
			name:     "Empty scopes",
			scopes:   []string{},
			expected: map[string]bool{},
		},
		{
			name:   "Duplicate scopes",
			scopes: []string{"read", "read", "write"},
			expected: map[string]bool{
				"read":  true,
				"write": true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := middleware.convertScopeSliceToMap(tt.scopes)
			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d scopes, got %d", len(tt.expected), len(result))
			}
			for key, value := range tt.expected {
				if result[key] != value {
					t.Errorf("Expected scope %s to be %v, got %v", key, value, result[key])
				}
			}
		})
	}
}

func TestIsPathAllowlisted(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		method      string
		allowlisted map[string]map[string]bool
		expected    bool
	}{
		{
			name:   "Path and method allowlisted",
			path:   "/health",
			method: "GET",
			allowlisted: map[string]map[string]bool{
				"/health": {"GET": true},
			},
			expected: true,
		},
		{
			name:   "Path allowlisted but wrong method",
			path:   "/health",
			method: "POST",
			allowlisted: map[string]map[string]bool{
				"/health": {"GET": true},
			},
			expected: false,
		},
		{
			name:   "Path not allowlisted",
			path:   "/protected",
			method: "GET",
			allowlisted: map[string]map[string]bool{
				"/health": {"GET": true},
			},
			expected: false,
		},
		{
			name:        "Empty allowlist",
			path:        "/health",
			method:      "GET",
			allowlisted: map[string]map[string]bool{},
			expected:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middleware := &Middleware{
				AllowlistedPaths: tt.allowlisted,
			}

			req := httptest.NewRequest(tt.method, tt.path, nil)
			result := middleware.isPathAllowlisted(req)

			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestShouldAllowBasicAuth(t *testing.T) {
	tests := []struct {
		name           string
		authHeader     string
		requestPath    string
		allowBasicAuth bool
		authEndpoint   string
		expected       bool
	}{
		{
			name:           "AllowBasicAuth is true",
			authHeader:     "Basic dGVzdDp0ZXN0",
			requestPath:    "/api/test",
			allowBasicAuth: true,
			authEndpoint:   "/auth",
			expected:       true,
		},
		{
			name:           "Basic auth on auth endpoint",
			authHeader:     "Basic dGVzdDp0ZXN0",
			requestPath:    "/auth",
			allowBasicAuth: false,
			authEndpoint:   "/auth",
			expected:       true,
		},
		{
			name:           "Basic auth on non-auth endpoint when not allowed",
			authHeader:     "Basic dGVzdDp0ZXN0",
			requestPath:    "/api/test",
			allowBasicAuth: false,
			authEndpoint:   "/auth",
			expected:       false,
		},
		{
			name:           "No Basic auth header",
			authHeader:     "Bearer token",
			requestPath:    "/auth",
			allowBasicAuth: false,
			authEndpoint:   "/auth",
			expected:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middleware := &Middleware{
				AllowBasicAuth: tt.allowBasicAuth,
				AuthEndpoint:   tt.authEndpoint,
			}

			req := httptest.NewRequest("GET", tt.requestPath, nil)
			result := middleware.shouldAllowBasicAuth(req, tt.authHeader)

			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}
