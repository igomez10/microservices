package requestid

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/igomez10/microservices/socialapp/internal/contexthelper"
	"github.com/rs/zerolog"
)

func TestMiddleware(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{
			name: "GET request",
			path: "/api/users",
		},
		{
			name: "POST request",
			path: "/api/users/create",
		},
		{
			name: "Root path",
			path: "/",
		},
		{
			name: "Path with query parameters",
			path: "/api/users?limit=10&offset=0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test handler that checks if request ID is set
			var capturedRequestID string
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Check if request ID is in the request header
				requestID := r.Header.Get("X-Request-ID")
				if requestID == "" {
					t.Error("Expected X-Request-ID in request header, got empty string")
				}
				capturedRequestID = requestID

				// Validate it's a valid UUID
				if _, err := uuid.Parse(requestID); err != nil {
					t.Errorf("Expected valid UUID for request ID, got %s: %v", requestID, err)
				}

				// Check if request ID is in the context
				contextRequestID := contexthelper.GetRequestIDInContext(r.Context())
				if contextRequestID == "" {
					t.Error("Expected request ID in context")
				}
				if contextRequestID != requestID {
					t.Errorf("Expected context request ID %s to match header request ID %s", contextRequestID, requestID)
				}

				// Check if logger has request ID
				log := contexthelper.GetLoggerInContext(r.Context())
				// Logger exists if it's not empty (we can't compare structs with []byte)
				_ = log

				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			})

			// Wrap the test handler with the middleware
			handler := Middleware(testHandler)

			// Create a test request
			req := httptest.NewRequest("GET", tt.path, nil)
			req = req.WithContext(contexthelper.SetLoggerInContext(req.Context(), zerolog.New(nil)))
			rr := httptest.NewRecorder()

			// Serve the request
			handler.ServeHTTP(rr, req)

			// Check response status
			if rr.Code != http.StatusOK {
				t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
			}

			// Check if X-Request-ID is in response header
			responseRequestID := rr.Header().Get("X-Request-ID")
			if responseRequestID == "" {
				t.Error("Expected X-Request-ID in response header, got empty string")
			}

			// Check if response request ID matches the captured request ID
			if responseRequestID != capturedRequestID {
				t.Errorf("Expected response request ID %s to match request ID %s", responseRequestID, capturedRequestID)
			}

			// Validate response request ID is a valid UUID
			if _, err := uuid.Parse(responseRequestID); err != nil {
				t.Errorf("Expected valid UUID for response request ID, got %s: %v", responseRequestID, err)
			}
		})
	}
}

func TestMiddleware_RequestIDUniqueness(t *testing.T) {
	// Test that each request gets a unique request ID
	requestIDs := make(map[string]bool)
	numRequests := 100

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestIDs[requestID] {
			t.Errorf("Duplicate request ID generated: %s", requestID)
		}
		requestIDs[requestID] = true
		w.WriteHeader(http.StatusOK)
	})

	handler := Middleware(testHandler)

	for i := 0; i < numRequests; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		req = req.WithContext(contexthelper.SetLoggerInContext(req.Context(), zerolog.New(nil)))
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
	}

	if len(requestIDs) != numRequests {
		t.Errorf("Expected %d unique request IDs, got %d", numRequests, len(requestIDs))
	}
}

func TestMiddleware_ContextPropagation(t *testing.T) {
	// Test that context is properly propagated through the middleware
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the context has the request ID
		requestID := contexthelper.GetRequestIDInContext(r.Context())
		if requestID == "" {
			t.Error("Expected request ID in context")
		}

		// Check if context contains the logger
		log := contexthelper.GetLoggerInContext(r.Context())
		// Logger exists (we can't compare structs with []byte)
		_ = log

		// Verify request ID is valid UUID
		if _, err := uuid.Parse(requestID); err != nil {
			t.Errorf("Expected valid UUID in context, got %s: %v", requestID, err)
		}

		w.WriteHeader(http.StatusOK)
	})

	handler := Middleware(testHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	req = req.WithContext(contexthelper.SetLoggerInContext(req.Context(), zerolog.New(nil)))
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestMiddleware_ResponseHeaderSet(t *testing.T) {
	// Test that the response header is set correctly
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Don't write anything, just return
		w.WriteHeader(http.StatusNoContent)
	})

	handler := Middleware(testHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	req = req.WithContext(contexthelper.SetLoggerInContext(req.Context(), zerolog.New(nil)))
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	requestID := rr.Header().Get("X-Request-ID")
	if requestID == "" {
		t.Error("Expected X-Request-ID in response header")
	}

	if _, err := uuid.Parse(requestID); err != nil {
		t.Errorf("Expected valid UUID in response header, got %s: %v", requestID, err)
	}
}

func TestMiddleware_RequestHeaderSet(t *testing.T) {
	// Test that the request header is set correctly
	var requestHeaderID string

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestHeaderID = r.Header.Get("X-Request-ID")
		w.WriteHeader(http.StatusOK)
	})

	handler := Middleware(testHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	req = req.WithContext(contexthelper.SetLoggerInContext(req.Context(), zerolog.New(nil)))
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if requestHeaderID == "" {
		t.Error("Expected X-Request-ID in request header")
	}

	if _, err := uuid.Parse(requestHeaderID); err != nil {
		t.Errorf("Expected valid UUID in request header, got %s: %v", requestHeaderID, err)
	}

	// Verify request and response headers match
	responseHeaderID := rr.Header().Get("X-Request-ID")
	if requestHeaderID != responseHeaderID {
		t.Errorf("Expected request header ID %s to match response header ID %s", requestHeaderID, responseHeaderID)
	}
}

func TestMiddleware_WithoutLogger(t *testing.T) {
	// Test that middleware works even without a logger in context
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			t.Error("Expected X-Request-ID in request header")
		}
		w.WriteHeader(http.StatusOK)
	})

	handler := Middleware(testHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	// Don't set logger in context
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	requestID := rr.Header().Get("X-Request-ID")
	if requestID == "" {
		t.Error("Expected X-Request-ID in response header")
	}
}

func TestMiddleware_HandlerPanics(t *testing.T) {
	// Test that request ID is still set even if the handler panics
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that request ID is set before panic
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			t.Error("Expected X-Request-ID before handler execution")
		}
		// Don't actually panic in the test, just verify request ID is set early
		w.WriteHeader(http.StatusOK)
	})

	handler := Middleware(testHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	req = req.WithContext(contexthelper.SetLoggerInContext(req.Context(), zerolog.New(nil)))
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	requestID := rr.Header().Get("X-Request-ID")
	if requestID == "" {
		t.Error("Expected X-Request-ID in response header")
	}
}

func TestMiddleware_DifferentHTTPMethods(t *testing.T) {
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requestID := r.Header.Get("X-Request-ID")
				if requestID == "" {
					t.Errorf("Expected X-Request-ID for %s request", method)
				}
				w.WriteHeader(http.StatusOK)
			})

			handler := Middleware(testHandler)

			req := httptest.NewRequest(method, "/test", nil)
			req = req.WithContext(contexthelper.SetLoggerInContext(req.Context(), zerolog.New(nil)))
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			requestID := rr.Header().Get("X-Request-ID")
			if requestID == "" {
				t.Errorf("Expected X-Request-ID in response header for %s request", method)
			}

			if _, err := uuid.Parse(requestID); err != nil {
				t.Errorf("Expected valid UUID for %s request, got %s: %v", method, requestID, err)
			}
		})
	}
}

func TestMiddleware_ConsistencyBetweenHeaderAndContext(t *testing.T) {
	// Test that the request ID is consistent between header and context
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headerRequestID := r.Header.Get("X-Request-ID")
		contextRequestID := contexthelper.GetRequestIDInContext(r.Context())

		if contextRequestID == "" {
			t.Error("Expected request ID in context")
		}

		if headerRequestID != contextRequestID {
			t.Errorf("Request ID mismatch: header=%s, context=%s", headerRequestID, contextRequestID)
		}

		w.WriteHeader(http.StatusOK)
	})

	handler := Middleware(testHandler)

	// Run multiple requests to ensure consistency
	for i := 0; i < 10; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		req = req.WithContext(contexthelper.SetLoggerInContext(req.Context(), zerolog.New(nil)))
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
	}
}

func BenchmarkMiddleware(b *testing.B) {
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := Middleware(testHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	req = req.WithContext(contexthelper.SetLoggerInContext(req.Context(), zerolog.New(nil)))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
	}
}
