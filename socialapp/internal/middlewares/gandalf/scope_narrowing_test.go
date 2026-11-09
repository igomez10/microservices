package gandalf

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/igomez10/microservices/socialapp/internal/contexthelper"
	"github.com/igomez10/microservices/socialapp/pkg/controller/user"
	db "github.com/igomez10/microservices/socialapp/pkg/dbpgx"
	"github.com/igomez10/microservices/socialapp/pkg/scopes"
	"github.com/jackc/pgx/v5/pgxpool"
)

// A user whose role grants only users.read asks for users.read plus admin scopes, the way the
// frontend requests every scope at login.
func basicAuthScopeRequest(t *testing.T, requested ...string) (*httptest.ResponseRecorder, map[string]bool, bool) {
	t.Helper()
	testUser := db.User{ID: 7, Username: "regular", Salt: "salt", HashedPassword: user.EncryptPassword("password", "salt")}
	m := &Middleware{
		DB: &mockDB{
			getUserByUsernameFunc: func(context.Context, *pgxpool.Pool, string) (db.User, error) { return testUser, nil },
			getUserRolesFunc: func(context.Context, *pgxpool.Pool, int64) ([]db.Role, error) {
				return []db.Role{{ID: 2, Name: "user"}}, nil
			},
			listRoleScopesFunc: func(context.Context, *pgxpool.Pool, db.ListRoleScopesParams) ([]db.Scope, error) {
				return []db.Scope{{ID: 11, Name: scopes.SocialappUsersRead.String()}}, nil
			},
		},
		AllowBasicAuth:   true,
		AllowlistedPaths: map[string]map[string]bool{},
	}

	var granted map[string]bool
	called := false
	handler := m.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		granted, _ = contexthelper.GetRequestedScopesInContext(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	form := url.Values{"scope": {strings.Join(requested, " ")}}
	req := httptest.NewRequest(http.MethodPost, "/auth", strings.NewReader(form.Encode()))
	req = req.WithContext(contexthelper.SetLoggerInContext(req.Context(), newTestLogger()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("regular:password")))
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	return rr, granted, called
}

func TestBasicAuth_DropsScopesTheUserIsNotAllowed(t *testing.T) {
	rr, granted, called := basicAuthScopeRequest(t,
		scopes.SocialappUsersRead.String(), scopes.SocialappRolesCreate.String(), scopes.SocialappUsersRolesUpdate.String())
	if !called || rr.Code != http.StatusOK {
		t.Fatalf("expected the request to pass with the allowed subset, got status %d", rr.Code)
	}
	if !granted[scopes.SocialappUsersRead.String()] {
		t.Errorf("allowed scope %s missing from granted scopes %v", scopes.SocialappUsersRead, granted)
	}
	if granted[scopes.SocialappRolesCreate.String()] || granted[scopes.SocialappUsersRolesUpdate.String()] {
		t.Errorf("disallowed scopes must be dropped, got %v", granted)
	}
}

func TestBasicAuth_RejectsWhenNoRequestedScopeIsAllowed(t *testing.T) {
	rr, _, called := basicAuthScopeRequest(t, scopes.SocialappRolesCreate.String())
	if called || rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 when nothing requested is allowed, got status %d (handler called: %v)", rr.Code, called)
	}
}
