package gandalf

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/igomez10/microservices/socialapp/internal/contexthelper"
	"github.com/igomez10/microservices/socialapp/pkg/controller/user"
	db "github.com/igomez10/microservices/socialapp/pkg/dbpgx"
	"github.com/igomez10/microservices/socialapp/pkg/scopes"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Go's oauth2 clientcredentials sends Basic base64(QueryEscape(id):QueryEscape(secret)) (RFC 6749 §2.3.1);
// browsers and curl send the raw secret. Both must log in, and a wrong secret must not.
func TestBasicAuth_AcceptsFormEncodedSecret(t *testing.T) {
	const password = "a+b/c=d%e f"
	testUser := db.User{ID: 7, Username: "admin-ish", Salt: "salt", HashedPassword: user.EncryptPassword(password, "salt")}
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
	handler := m.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }))

	for name, tc := range map[string]struct {
		secret string
		want   int
	}{
		"raw (browser, curl)":      {password, http.StatusOK},
		"form-encoded (Go oauth2)": {url.QueryEscape(password), http.StatusOK},
		"wrong":                    {"a+b/c=d%e g", http.StatusUnauthorized},
		"wrong, form-encoded":      {url.QueryEscape("a+b/c=d%e g"), http.StatusUnauthorized},
	} {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/auth?scope="+scopes.SocialappUsersRead.String(), nil)
			req = req.WithContext(contexthelper.SetLoggerInContext(req.Context(), newTestLogger()))
			req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(url.QueryEscape(testUser.Username)+":"+tc.secret)))
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			if rr.Code != tc.want {
				t.Fatalf("got %d, want %d: %s", rr.Code, tc.want, rr.Body.String())
			}
		})
	}
}
