package cache

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

// mockRedisClient implements a simple in-memory cache for testing
type mockRedisClient struct {
	data map[string][]byte
}

func newMockRedisClient() *mockRedisClient {
	return &mockRedisClient{
		data: make(map[string][]byte),
	}
}

func (m *mockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	cmd := redis.NewStringCmd(ctx)
	if val, ok := m.data[key]; ok {
		cmd.SetVal(string(val))
	} else {
		cmd.SetErr(redis.Nil)
	}
	return cmd
}

func (m *mockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		bytes = []byte(fmt.Sprint(v))
	}
	m.data[key] = bytes
	return redis.NewStatusCmd(ctx)
}

func TestCacheMiddleware_CacheHit(t *testing.T) {
	// Setup
	mockClient := newMockRedisClient()
	cache := &Cache{Client: mockClient}

	// Pre-populate cache
	expectedBody := []byte(`{"message":"cached response"}`)
	cacheKey := "GET+/test/path"
	mockClient.data[cacheKey] = expectedBody

	// Create test handler that should not be called
	handlerCalled := false
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"fresh response"}`))
	})

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/test/path", nil)
	rec := httptest.NewRecorder()

	// Execute
	middleware := cache.Middleware(testHandler)
	middleware.ServeHTTP(rec, req)

	// Assertions
	assert.False(t, handlerCalled, "Handler should not be called on cache hit")
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, expectedBody, rec.Body.Bytes())
	assert.Equal(t, "HIT", rec.Header().Get("X-Cache"))
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.Equal(t, "public, max-age=3600", rec.Header().Get("Cache-Control"))
}

func TestCacheMiddleware_CacheMiss(t *testing.T) {
	// Setup
	mockClient := newMockRedisClient()
	cache := &Cache{Client: mockClient}

	// Create test handler
	expectedBody := []byte(`{"message":"fresh response"}`)
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(expectedBody)
	})

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/test/path", nil)
	rec := httptest.NewRecorder()

	// Execute
	middleware := cache.Middleware(testHandler)
	middleware.ServeHTTP(rec, req)

	// Assertions
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, expectedBody, rec.Body.Bytes())
	assert.NotEqual(t, "HIT", rec.Header().Get("X-Cache"))

	// Verify response was cached
	cacheKey := "GET+/test/path"
	cachedValue, exists := mockClient.data[cacheKey]
	assert.True(t, exists, "Response should be cached")
	assert.Equal(t, expectedBody, cachedValue)
}

func TestCacheMiddleware_CacheSkipWithNoStore(t *testing.T) {
	// Setup
	mockClient := newMockRedisClient()
	cache := &Cache{Client: mockClient}

	// Pre-populate cache
	cacheKey := "GET+/test/path"
	mockClient.data[cacheKey] = []byte(`{"message":"cached response"}`)

	// Create test handler
	expectedBody := []byte(`{"message":"fresh response"}`)
	handlerCalled := false
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
		w.Write(expectedBody)
	})

	// Create request with Cache-Control: no-store
	req := httptest.NewRequest(http.MethodGet, "/test/path", nil)
	req.Header.Set("Cache-Control", "no-store")
	rec := httptest.NewRecorder()

	// Execute
	middleware := cache.Middleware(testHandler)
	middleware.ServeHTTP(rec, req)

	// Assertions
	assert.True(t, handlerCalled, "Handler should be called when cache is skipped")
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, expectedBody, rec.Body.Bytes())
	assert.NotEqual(t, "HIT", rec.Header().Get("X-Cache"))
}

func TestCacheMiddleware_NonGetRequest(t *testing.T) {
	// Setup
	mockClient := newMockRedisClient()
	cache := &Cache{Client: mockClient}

	// Create test handler
	handlerCalled := false
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"message":"created"}`))
	})

	// Test POST request
	req := httptest.NewRequest(http.MethodPost, "/test/path", nil)
	rec := httptest.NewRecorder()

	// Execute
	middleware := cache.Middleware(testHandler)
	middleware.ServeHTTP(rec, req)

	// Assertions
	assert.True(t, handlerCalled, "Handler should be called for non-GET requests")
	assert.Equal(t, http.StatusCreated, rec.Code)

	// Verify nothing was cached
	assert.Equal(t, 0, len(mockClient.data), "POST responses should not be cached")
}

func TestCacheMiddleware_DoNotCacheNon200Responses(t *testing.T) {
	testCases := []struct {
		name       string
		statusCode int
	}{
		{"404 Not Found", http.StatusNotFound},
		{"500 Internal Server Error", http.StatusInternalServerError},
		{"400 Bad Request", http.StatusBadRequest},
		{"201 Created", http.StatusCreated},
		{"409 Conflict", http.StatusConflict},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			mockClient := newMockRedisClient()
			cache := &Cache{Client: mockClient}

			// Create test handler
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.statusCode)
				w.Write([]byte(`{"error":"something went wrong"}`))
			})

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/test/path", nil)
			rec := httptest.NewRecorder()

			// Execute
			middleware := cache.Middleware(testHandler)
			middleware.ServeHTTP(rec, req)

			// Assertions
			assert.Equal(t, tc.statusCode, rec.Code)

			// Verify nothing was cached
			assert.Equal(t, 0, len(mockClient.data), "Non-200 responses should not be cached")
		})
	}
}

func TestCacheMiddleware_CacheWithQueryParams(t *testing.T) {
	// Setup
	mockClient := newMockRedisClient()
	cache := &Cache{Client: mockClient}

	// Create test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"response with params"}`))
	})

	// Create request with query params
	req := httptest.NewRequest(http.MethodGet, "/test/path?foo=bar&baz=qux", nil)
	rec := httptest.NewRecorder()

	// Execute
	middleware := cache.Middleware(testHandler)
	middleware.ServeHTTP(rec, req)

	// Verify response was cached with query params in key
	cacheKey := "GET+/test/path?foo=bar&baz=qux"
	_, exists := mockClient.data[cacheKey]
	assert.True(t, exists, "Response should be cached with query params in key")
}

func TestCacheMiddleware_DifferentPathsDifferentCacheEntries(t *testing.T) {
	// Setup
	mockClient := newMockRedisClient()
	cache := &Cache{Client: mockClient}

	paths := []string{"/path1", "/path2", "/path3"}

	for _, path := range paths {
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf(`{"path":"%s"}`, r.URL.Path)))
		})

		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()

		middleware := cache.Middleware(testHandler)
		middleware.ServeHTTP(rec, req)
	}

	// Verify each path has its own cache entry
	assert.Equal(t, len(paths), len(mockClient.data), "Each path should have its own cache entry")

	for _, path := range paths {
		cacheKey := "GET+" + path
		_, exists := mockClient.data[cacheKey]
		assert.True(t, exists, fmt.Sprintf("Cache entry should exist for %s", path))
	}
}

func TestCacheMiddleware_CacheKeyFormat(t *testing.T) {
	// Setup
	mockClient := newMockRedisClient()
	cache := &Cache{Client: mockClient}

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"test"}`))
	})

	// Test various URLs
	testCases := []struct {
		url         string
		expectedKey string
	}{
		{"/simple", "GET+/simple"},
		{"/with/path", "GET+/with/path"},
		{"/with?query=params", "GET+/with?query=params"},
		{"/with?multiple=params&another=one", "GET+/with?multiple=params&another=one"},
	}

	for _, tc := range testCases {
		t.Run(tc.url, func(t *testing.T) {
			// Clear cache between tests
			mockClient.data = make(map[string][]byte)

			req := httptest.NewRequest(http.MethodGet, tc.url, nil)
			rec := httptest.NewRecorder()

			middleware := cache.Middleware(testHandler)
			middleware.ServeHTTP(rec, req)

			_, exists := mockClient.data[tc.expectedKey]
			assert.True(t, exists, fmt.Sprintf("Cache key should be '%s'", tc.expectedKey))
		})
	}
}

func TestCacheMiddleware_SubsequentRequestsUseCachedResponse(t *testing.T) {
	// Setup
	mockClient := newMockRedisClient()
	cache := &Cache{Client: mockClient}

	callCount := 0
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"call":%d}`, callCount)))
	})

	middleware := cache.Middleware(testHandler)

	// First request - should miss cache and call handler
	req1 := httptest.NewRequest(http.MethodGet, "/test/path", nil)
	rec1 := httptest.NewRecorder()
	middleware.ServeHTTP(rec1, req1)

	assert.Equal(t, 1, callCount, "Handler should be called once")
	assert.Equal(t, `{"call":1}`, rec1.Body.String())
	assert.NotEqual(t, "HIT", rec1.Header().Get("X-Cache"))

	// Second request - should hit cache and NOT call handler
	req2 := httptest.NewRequest(http.MethodGet, "/test/path", nil)
	rec2 := httptest.NewRecorder()
	middleware.ServeHTTP(rec2, req2)

	assert.Equal(t, 1, callCount, "Handler should still only be called once (cache hit)")
	assert.Equal(t, `{"call":1}`, rec2.Body.String(), "Should return cached response")
	assert.Equal(t, "HIT", rec2.Header().Get("X-Cache"))
}

func TestNewCache(t *testing.T) {
	// This is an integration-like test that would require a real Redis connection
	// Skip in unit tests unless Redis is available
	t.Skip("Requires Redis connection - implement as integration test")

	// Example of how to test with real Redis:
	// config := CacheConfig{
	// 	RedisOpts: &redis.Options{
	// 		Addr: "localhost:6379",
	// 	},
	// }
	// cache := NewCache(config)
	// require.NotNil(t, cache)
	// require.NotNil(t, cache.Client)
}

func TestCacheMiddleware_ContextPropagation(t *testing.T) {
	// Setup
	mockClient := newMockRedisClient()
	cache := &Cache{Client: mockClient}

	// Create test handler that checks context
	contextChecked := false
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify context is still valid
		if r.Context() != nil {
			contextChecked = true
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"test"}`))
	})

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/test/path", nil)
	rec := httptest.NewRecorder()

	// Execute
	middleware := cache.Middleware(testHandler)
	middleware.ServeHTTP(rec, req)

	// Assertions
	assert.True(t, contextChecked, "Context should be propagated to handler")
}
