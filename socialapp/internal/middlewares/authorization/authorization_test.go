package authorization

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/igomez10/microservices/socialapp/internal/contexthelper"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestAuthorize(t *testing.T) {
	testCases := []struct {
		name            string
		requiredScopes  map[string]bool
		tokenScopes     map[string]bool
		scopesInContext bool
		expectStatus    int
		expectCalled    bool
		expectBodyPart  string
	}{
		{
			name: "Success - All required scopes present",
			requiredScopes: map[string]bool{
				"read:users":  true,
				"write:users": true,
			},
			tokenScopes: map[string]bool{
				"read:users":  true,
				"write:users": true,
				"admin":       true,
			},
			scopesInContext: true,
			expectStatus:    http.StatusOK,
			expectCalled:    true,
			expectBodyPart:  "success",
		},
		{
			name: "Failure - Missing scopes in context",
			requiredScopes: map[string]bool{
				"read:users": true,
			},
			tokenScopes:     nil,
			scopesInContext: false,
			expectStatus:    http.StatusForbidden,
			expectCalled:    false,
			expectBodyPart:  "No scopes in context",
		},
		{
			name: "Failure - Missing required scope",
			requiredScopes: map[string]bool{
				"read:users":  true,
				"write:users": true,
				"admin":       true,
			},
			tokenScopes: map[string]bool{
				"read:users": true,
			},
			scopesInContext: true,
			expectStatus:    http.StatusForbidden,
			expectCalled:    false,
			expectBodyPart:  "missing from token",
		},
		{
			name:            "Success - Empty required scopes",
			requiredScopes:  map[string]bool{},
			tokenScopes:     map[string]bool{"some:scope": true},
			scopesInContext: true,
			expectStatus:    http.StatusOK,
			expectCalled:    true,
			expectBodyPart:  "success",
		},
		{
			name: "Success - Exact scope match",
			requiredScopes: map[string]bool{
				"read:users":   true,
				"write:users":  true,
				"delete:users": true,
			},
			tokenScopes: map[string]bool{
				"read:users":   true,
				"write:users":  true,
				"delete:users": true,
			},
			scopesInContext: true,
			expectStatus:    http.StatusOK,
			expectCalled:    true,
			expectBodyPart:  "success",
		},
		{
			name: "Success - Token has more scopes than required",
			requiredScopes: map[string]bool{
				"read:users": true,
			},
			tokenScopes: map[string]bool{
				"read:users":   true,
				"write:users":  true,
				"delete:users": true,
				"admin":        true,
				"superadmin":   true,
			},
			scopesInContext: true,
			expectStatus:    http.StatusOK,
			expectCalled:    true,
			expectBodyPart:  "success",
		},
		{
			name: "Failure - Single required scope missing",
			requiredScopes: map[string]bool{
				"admin": true,
			},
			tokenScopes: map[string]bool{
				"read:users":  true,
				"write:users": true,
			},
			scopesInContext: true,
			expectStatus:    http.StatusForbidden,
			expectCalled:    false,
			expectBodyPart:  "admin",
		},
		{
			name: "Failure - Multiple required scopes, one missing",
			requiredScopes: map[string]bool{
				"read:users":   true,
				"write:users":  true,
				"delete:users": true,
			},
			tokenScopes: map[string]bool{
				"read:users":  true,
				"write:users": true,
			},
			scopesInContext: true,
			expectStatus:    http.StatusForbidden,
			expectCalled:    false,
			expectBodyPart:  "delete:users",
		},
		{
			name: "Failure - Token has no scopes at all",
			requiredScopes: map[string]bool{
				"read:users": true,
			},
			tokenScopes:     map[string]bool{},
			scopesInContext: true,
			expectStatus:    http.StatusUnauthorized,
			expectCalled:    false,
			expectBodyPart:  "No scopes in context and required scopes are not empty",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			middleware := &Middleware{
				RequiredScopes: tc.requiredScopes,
			}

			// Create test handler
			handlerCalled := false
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handlerCalled = true
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"message":"success"}`))
			})

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			ctx := req.Context()

			// Add logger to context
			logger := zerolog.Nop()
			ctx = context.WithValue(ctx, "logger", logger)

			// Add scopes to context if needed
			if tc.scopesInContext {
				req = contexthelper.SetRequestedScopesInContext(req.WithContext(ctx), tc.tokenScopes)
			} else {
				req = req.WithContext(ctx)
			}

			rec := httptest.NewRecorder()

			// Execute
			handler := middleware.Authorize(testHandler)
			handler.ServeHTTP(rec, req)

			// Assertions
			assert.Equal(t, tc.expectStatus, rec.Code, "Status code mismatch")
			assert.Equal(t, tc.expectCalled, handlerCalled, "Handler call expectation mismatch")
			assert.Contains(t, rec.Body.String(), tc.expectBodyPart, "Response body should contain expected text")
		})
	}
}

func TestAuthorize_ContextPropagation(t *testing.T) {
	testCases := []struct {
		name           string
		requiredScopes map[string]bool
		tokenScopes    map[string]bool
		contextKey     string
		contextValue   string
		expectValue    bool
	}{
		{
			name: "Context preserved on success",
			requiredScopes: map[string]bool{
				"read:users": true,
			},
			tokenScopes: map[string]bool{
				"read:users": true,
			},
			contextKey:   "test-key",
			contextValue: "test-value",
			expectValue:  true,
		},
		{
			name: "Context available even on failure",
			requiredScopes: map[string]bool{
				"admin": true,
			},
			tokenScopes: map[string]bool{
				"read:users": true,
			},
			contextKey:   "custom-key",
			contextValue: "custom-value",
			expectValue:  false, // Handler not called, so we can't check context in handler
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			middleware := &Middleware{
				RequiredScopes: tc.requiredScopes,
			}

			var receivedContext context.Context
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				receivedContext = r.Context()
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			ctx := req.Context()

			// Add logger to context
			logger := zerolog.Nop()
			ctx = context.WithValue(ctx, "logger", logger)

			// Add custom value to verify context propagation
			ctx = context.WithValue(ctx, tc.contextKey, tc.contextValue)

			req = contexthelper.SetRequestedScopesInContext(req.WithContext(ctx), tc.tokenScopes)
			rec := httptest.NewRecorder()

			handler := middleware.Authorize(testHandler)
			handler.ServeHTTP(rec, req)

			if tc.expectValue {
				assert.NotNil(t, receivedContext, "Context should be propagated to handler")
				assert.Equal(t, tc.contextValue, receivedContext.Value(tc.contextKey), "Custom context values should be preserved")
			}
		})
	}
}

func TestAuthorize_HTTPMethods(t *testing.T) {
	methods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodPatch,
	}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			middleware := &Middleware{
				RequiredScopes: map[string]bool{
					"read:users": true,
				},
			}

			handlerCalled := false
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handlerCalled = true
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest(method, "/test", nil)
			ctx := req.Context()

			logger := zerolog.Nop()
			ctx = context.WithValue(ctx, "logger", logger)

			tokenScopes := map[string]bool{
				"read:users": true,
			}
			req = contexthelper.SetRequestedScopesInContext(req.WithContext(ctx), tokenScopes)
			rec := httptest.NewRecorder()

			handler := middleware.Authorize(testHandler)
			handler.ServeHTTP(rec, req)

			assert.True(t, handlerCalled, "Handler should be called for all HTTP methods with valid scopes")
			assert.Equal(t, http.StatusOK, rec.Code)
		})
	}
}

func TestAuthorize_ErrorResponses(t *testing.T) {
	testCases := []struct {
		name           string
		requiredScopes map[string]bool
		tokenScopes    map[string]bool
		scopesInCtx    bool
		expectStatus   int
		expectCode     int
		expectMessage  string
	}{
		{
			name: "No scopes in context - returns 403 with appropriate message",
			requiredScopes: map[string]bool{
				"read:users": true,
			},
			tokenScopes:   nil,
			scopesInCtx:   false,
			expectStatus:  http.StatusForbidden,
			expectCode:    403,
			expectMessage: "No scopes in context",
		},
		{
			name: "Missing specific scope - returns 403 with scope name",
			requiredScopes: map[string]bool{
				"delete:users": true,
			},
			tokenScopes: map[string]bool{
				"read:users": true,
			},
			scopesInCtx:   true,
			expectStatus:  http.StatusForbidden,
			expectCode:    403,
			expectMessage: "delete:users",
		},
		{
			name: "Missing admin scope",
			requiredScopes: map[string]bool{
				"admin": true,
			},
			tokenScopes: map[string]bool{
				"read:users":  true,
				"write:users": true,
			},
			scopesInCtx:   true,
			expectStatus:  http.StatusForbidden,
			expectCode:    403,
			expectMessage: "admin",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			middleware := &Middleware{
				RequiredScopes: tc.requiredScopes,
			}

			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				t.Error("Handler should not be called")
			})

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			ctx := req.Context()

			logger := zerolog.Nop()
			ctx = context.WithValue(ctx, "logger", logger)

			if tc.scopesInCtx {
				req = contexthelper.SetRequestedScopesInContext(req.WithContext(ctx), tc.tokenScopes)
			} else {
				req = req.WithContext(ctx)
			}

			rec := httptest.NewRecorder()

			handler := middleware.Authorize(testHandler)
			handler.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectStatus, rec.Code)
			assert.Contains(t, rec.Body.String(), tc.expectMessage)

			// Verify the response contains the code field
			body := rec.Body.String()
			assert.Contains(t, body, "code")
			assert.Contains(t, body, "message")
		})
	}
}

func TestAuthorize_ComplexScopeScenarios(t *testing.T) {
	testCases := []struct {
		name           string
		requiredScopes map[string]bool
		tokenScopes    map[string]bool
		expectPass     bool
		description    string
	}{
		{
			name: "Wildcard-like pattern - many required, all present",
			requiredScopes: map[string]bool{
				"read:users":   true,
				"write:users":  true,
				"delete:users": true,
				"read:posts":   true,
				"write:posts":  true,
			},
			tokenScopes: map[string]bool{
				"read:users":   true,
				"write:users":  true,
				"delete:users": true,
				"read:posts":   true,
				"write:posts":  true,
				"admin":        true,
			},
			expectPass:  true,
			description: "All required scopes present with extras",
		},
		{
			name: "Wildcard-like pattern - many required, one missing",
			requiredScopes: map[string]bool{
				"read:users":   true,
				"write:users":  true,
				"delete:users": true,
				"read:posts":   true,
				"write:posts":  true,
			},
			tokenScopes: map[string]bool{
				"read:users":   true,
				"write:users":  true,
				"delete:users": true,
				"read:posts":   true,
				// Missing write:posts
			},
			expectPass:  false,
			description: "Missing one of many required scopes",
		},
		{
			name:           "Empty required scopes allows everything",
			requiredScopes: map[string]bool{},
			tokenScopes: map[string]bool{
				// Any scopes or empty - doesn't matter when no scopes required
			},
			expectPass:  true,
			description: "Public endpoint - no scopes required",
		},
		{
			name: "Single scope requirement",
			requiredScopes: map[string]bool{
				"service:access": true,
			},
			tokenScopes: map[string]bool{
				"service:access": true,
			},
			expectPass:  true,
			description: "Simplest case - single scope",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			middleware := &Middleware{
				RequiredScopes: tc.requiredScopes,
			}

			handlerCalled := false
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handlerCalled = true
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			ctx := req.Context()

			logger := zerolog.Nop()
			ctx = context.WithValue(ctx, "logger", logger)

			req = contexthelper.SetRequestedScopesInContext(req.WithContext(ctx), tc.tokenScopes)
			rec := httptest.NewRecorder()

			handler := middleware.Authorize(testHandler)
			handler.ServeHTTP(rec, req)

			if tc.expectPass {
				assert.True(t, handlerCalled, tc.description)
				assert.Equal(t, http.StatusOK, rec.Code)
			} else {
				assert.False(t, handlerCalled, tc.description)
				assert.Equal(t, http.StatusForbidden, rec.Code)
			}
		})
	}
}
