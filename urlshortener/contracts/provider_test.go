package contracts

import (
	"context"
	"database/sql"
	"fmt"
	"net/http/httptest"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/pact-foundation/pact-go/v2/models"
	"github.com/pact-foundation/pact-go/v2/provider"
	"github.com/stretchr/testify/require"

	"github.com/igomez10/microservices/urlshortener/generated/server"
	"github.com/igomez10/microservices/urlshortener/pkg/controllers/url"
	"github.com/igomez10/microservices/urlshortener/pkg/db"
)

func TestProvider_VerifyUrlShortener(t *testing.T) {
	mock := newMockDB()
	ts := newMockServer(mock)
	defer ts.Close()

	pactFile := filepath.Join("..", "..", "socialapp", "contracts", "pacts", "socialapp-urlshortener.json")

	verifier := provider.NewVerifier()
	err := verifier.VerifyProvider(t, provider.VerifyRequest{
		ProviderBaseURL: ts.URL,
		PactFiles:       []string{pactFile},
		Provider:        "urlshortener",
		StateHandlers: models.StateHandlers{
			"a url alias exists": func(setup bool, state models.ProviderState) (models.ProviderStateResponse, error) {
				mock.Seed("contract-alias", "https://pact.socialapp/golden-url")
				return nil, nil
			},
		},
	})
	require.NoError(t, err)
}

func newMockServer(mock *mockDB) *httptest.Server {
	service := &url.URLApiService{
		DB:     mock,
		DBConn: mockConn{},
	}
	controller := server.NewURLAPIController(service)
	router := server.NewRouter(controller)
	return httptest.NewServer(router)
}

type mockDB struct {
	mu   sync.Mutex
	urls map[string]db.Url
}

func newMockDB() *mockDB {
	return &mockDB{
		urls: make(map[string]db.Url),
	}
}

func (m *mockDB) Seed(alias, shortURL string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := time.Now().UTC()
	m.urls[alias] = db.Url{
		ID:        42,
		Url:       shortURL,
		Alias:     alias,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (m *mockDB) CreateURL(_ context.Context, _ db.DBTX, arg db.CreateURLParams) (db.Url, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := time.Now().UTC()
	entry := db.Url{
		ID:        now.Unix(),
		Url:       arg.Url,
		Alias:     arg.Alias,
		CreatedAt: now,
		UpdatedAt: now,
	}
	m.urls[arg.Alias] = entry
	return entry, nil
}

func (m *mockDB) DeleteURL(_ context.Context, _ db.DBTX, alias string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	entry, ok := m.urls[alias]
	if !ok {
		return sql.ErrNoRows
	}
	entry.DeletedAt = sql.NullTime{Time: time.Now().UTC(), Valid: true}
	m.urls[alias] = entry
	return nil
}

func (m *mockDB) GetURLFromAlias(_ context.Context, _ db.DBTX, alias string) (db.Url, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	entry, ok := m.urls[alias]
	if !ok {
		return db.Url{}, sql.ErrNoRows
	}
	return entry, nil
}

type mockConn struct{}

func (mockConn) ExecContext(context.Context, string, ...interface{}) (sql.Result, error) {
	return mockResult{}, nil
}

func (mockConn) PrepareContext(context.Context, string) (*sql.Stmt, error) {
	return nil, fmt.Errorf("not implemented")
}

func (mockConn) QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error) {
	return nil, fmt.Errorf("not implemented")
}

func (mockConn) QueryRowContext(context.Context, string, ...interface{}) *sql.Row {
	return nil
}

type mockResult struct{}

func (mockResult) LastInsertId() (int64, error) {
	return 0, nil
}

func (mockResult) RowsAffected() (int64, error) {
	return 0, nil
}
