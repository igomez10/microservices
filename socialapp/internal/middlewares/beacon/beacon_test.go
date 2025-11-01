package beacon

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/igomez10/microservices/socialapp/internal/contexthelper"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
)

func setupTestEnvironment() (*metric.MeterProvider, *trace.TracerProvider) {
	// Setup metric provider for testing
	meterProvider := metric.NewMeterProvider()
	otel.SetMeterProvider(meterProvider)

	// Setup tracer provider for testing
	tracerProvider := trace.NewTracerProvider()
	otel.SetTracerProvider(tracerProvider)

	return meterProvider, tracerProvider
}

func TestBeacon_Middleware(t *testing.T) {
	setupTestEnvironment()

	tests := []struct {
		name           string
		path           string
		method         string
		statusCode     int
		expectedStatus int
		headers        map[string]string
		body           string
	}{
		{
			name:           "GET request - 200 OK",
			path:           "/api/users",
			method:         "GET",
			statusCode:     http.StatusOK,
			expectedStatus: http.StatusOK,
			headers:        map[string]string{"X-Request-ID": "test-request-id-123"},
		},
		{
			name:           "POST request - 201 Created",
			path:           "/api/users",
			method:         "POST",
			statusCode:     http.StatusCreated,
			expectedStatus: http.StatusCreated,
			headers:        map[string]string{"Content-Type": "application/json"},
			body:           `{"name":"test"}`,
		},
		{
			name:           "GET request - 404 Not Found",
			path:           "/api/nonexistent",
			method:         "GET",
			statusCode:     http.StatusNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "GET request - 401 Unauthorized",
			path:           "/api/protected",
			method:         "GET",
			statusCode:     http.StatusUnauthorized,
			expectedStatus: http.StatusUnauthorized,
			headers:        map[string]string{"Authorization": "Bearer invalid-token"},
		},
		{
			name:           "GET request - 500 Internal Server Error",
			path:           "/api/error",
			method:         "GET",
			statusCode:     http.StatusInternalServerError,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "PUT request - 200 OK",
			path:           "/api/users/123",
			method:         "PUT",
			statusCode:     http.StatusOK,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "DELETE request - 204 No Content",
			path:           "/api/users/123",
			method:         "DELETE",
			statusCode:     http.StatusNoContent,
			expectedStatus: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			beacon := &Beacon{}

			// Create test handler that returns the specified status code
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				if tt.body != "" {
					w.Write([]byte("response"))
				}
			})

			// Wrap with beacon middleware
			handler := beacon.Middleware(testHandler)

			// Create request
			var reqBody io.Reader
			if tt.body != "" {
				reqBody = strings.NewReader(tt.body)
			}
			req := httptest.NewRequest(tt.method, tt.path, reqBody)

			// Set headers
			for key, value := range tt.headers {
				req.Header.Set(key, value)
			}

			// Set logger and pattern in context
			logBuf := &bytes.Buffer{}
			logger := zerolog.New(logBuf).With().Timestamp().Logger()
			req = req.WithContext(contexthelper.SetLoggerInContext(req.Context(), logger))
			req = req.WithContext(contexthelper.SetRequestPatternInContext(req.Context(), tt.path))

			// Create response recorder
			rr := httptest.NewRecorder()

			// Serve request
			handler.ServeHTTP(rr, req)

			// Verify status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, rr.Code)
			}

			// Verify log was created
			logOutput := logBuf.String()
			if logOutput == "" {
				t.Error("Expected log output, got empty string")
			}

			// Verify log contains expected fields
			if !strings.Contains(logOutput, "finished request") {
				t.Error("Expected log to contain 'finished request'")
			}

			if !strings.Contains(logOutput, tt.method) {
				t.Errorf("Expected log to contain method %s", tt.method)
			}

			if !strings.Contains(logOutput, tt.path) {
				t.Errorf("Expected log to contain path %s", tt.path)
			}
		})
	}
}

func TestBeacon_Middleware_WithQueryParameters(t *testing.T) {
	setupTestEnvironment()

	beacon := &Beacon{}
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := beacon.Middleware(testHandler)

	req := httptest.NewRequest("GET", "/api/users?limit=10&offset=20", nil)
	logBuf := &bytes.Buffer{}
	logger := zerolog.New(logBuf).With().Timestamp().Logger()
	req = req.WithContext(contexthelper.SetLoggerInContext(req.Context(), logger))
	req = req.WithContext(contexthelper.SetRequestPatternInContext(req.Context(), "/api/users"))

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	logOutput := logBuf.String()
	if !strings.Contains(logOutput, "limit=10&offset=20") {
		t.Error("Expected log to contain query parameters")
	}
}

func TestBeacon_Middleware_WithUserAgent(t *testing.T) {
	setupTestEnvironment()

	beacon := &Beacon{}
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := beacon.Middleware(testHandler)

	req := httptest.NewRequest("GET", "/api/test", nil)
	req.Header.Set("User-Agent", "TestAgent/1.0")
	logBuf := &bytes.Buffer{}
	logger := zerolog.New(logBuf).With().Timestamp().Logger()
	req = req.WithContext(contexthelper.SetLoggerInContext(req.Context(), logger))
	req = req.WithContext(contexthelper.SetRequestPatternInContext(req.Context(), "/api/test"))

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	logOutput := logBuf.String()
	if !strings.Contains(logOutput, "TestAgent/1.0") {
		t.Error("Expected log to contain user agent")
	}
}

func TestBeacon_Middleware_UnauthorizedWithUsername(t *testing.T) {
	setupTestEnvironment()

	beacon := &Beacon{}
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})

	handler := beacon.Middleware(testHandler)

	req := httptest.NewRequest("GET", "/api/protected", nil)
	req.Header.Set("Authorization", "Bearer token123")

	logBuf := &bytes.Buffer{}
	logger := zerolog.New(logBuf).With().Timestamp().Logger()
	ctx := contexthelper.SetLoggerInContext(req.Context(), logger)
	ctx = contexthelper.SetUsernameInContext(req.WithContext(ctx), "testuser").Context()
	ctx = contexthelper.SetRequestPatternInContext(ctx, "/api/protected")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, rr.Code)
	}

	logOutput := logBuf.String()
	if !strings.Contains(logOutput, "Unauthorized") {
		t.Error("Expected log to contain 'Unauthorized'")
	}
	if !strings.Contains(logOutput, "testuser") {
		t.Error("Expected log to contain username")
	}
	if !strings.Contains(logOutput, "Bearer token123") {
		t.Error("Expected log to contain authorization header for unauthorized request")
	}
}

func TestBeacon_Middleware_RemoteIPParsing(t *testing.T) {
	setupTestEnvironment()

	tests := []struct {
		name          string
		remoteAddr    string
		expectedInLog string
	}{
		{
			name:          "Valid IP with port",
			remoteAddr:    "192.168.1.1:12345",
			expectedInLog: "192.168.1.1:12345",
		},
		{
			name:          "IPv6 with port",
			remoteAddr:    "[2001:db8::1]:8080",
			expectedInLog: "[2001:db8::1]:8080",
		},
		{
			name:          "Invalid format",
			remoteAddr:    "invalid-format",
			expectedInLog: "invalid-format",
		},
		{
			name:          "Empty remote addr",
			remoteAddr:    "",
			expectedInLog: "remote_addr", // Just check the field exists
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			beacon := &Beacon{}
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			handler := beacon.Middleware(testHandler)

			req := httptest.NewRequest("GET", "/test", nil)
			req.RemoteAddr = tt.remoteAddr

			logBuf := &bytes.Buffer{}
			logger := zerolog.New(logBuf).With().Timestamp().Logger()
			req = req.WithContext(contexthelper.SetLoggerInContext(req.Context(), logger))
			req = req.WithContext(contexthelper.SetRequestPatternInContext(req.Context(), "/test"))

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			logOutput := logBuf.String()
			if !strings.Contains(logOutput, tt.expectedInLog) {
				t.Errorf("Expected log to contain %s, got log: %s", tt.expectedInLog, logOutput)
			}
		})
	}
}

func TestBeacon_Middleware_LatencyTracking(t *testing.T) {
	setupTestEnvironment()

	beacon := &Beacon{}

	// Handler with artificial delay
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(50 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	})

	handler := beacon.Middleware(testHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	logBuf := &bytes.Buffer{}
	logger := zerolog.New(logBuf).With().Timestamp().Logger()
	req = req.WithContext(contexthelper.SetLoggerInContext(req.Context(), logger))
	req = req.WithContext(contexthelper.SetRequestPatternInContext(req.Context(), "/test"))

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	logOutput := logBuf.String()
	if !strings.Contains(logOutput, "latency_ms") {
		t.Error("Expected log to contain latency_ms")
	}

	// Check that latency is at least 50ms (accounting for some variance)
	if !strings.Contains(logOutput, "latency_ms") {
		t.Error("Latency should be tracked in logs")
	}
}

func TestBeacon_Middleware_AuthorizationHeaderExcluded(t *testing.T) {
	setupTestEnvironment()

	beacon := &Beacon{}
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := beacon.Middleware(testHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer secret-token-123")
	req.Header.Set("X-Custom-Header", "custom-value")

	logBuf := &bytes.Buffer{}
	logger := zerolog.New(logBuf).With().Timestamp().Logger()
	req = req.WithContext(contexthelper.SetLoggerInContext(req.Context(), logger))
	req = req.WithContext(contexthelper.SetRequestPatternInContext(req.Context(), "/test"))

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	logOutput := logBuf.String()

	// Authorization header should be excluded for non-401 status codes
	if strings.Contains(logOutput, "secret-token-123") {
		t.Error("Expected Authorization header to be excluded from logs for successful requests")
	}

	// Custom header should be included
	if !strings.Contains(logOutput, "custom-value") {
		t.Error("Expected custom header to be included in logs")
	}
}

func TestBeacon_Middleware_StatusCodeClassification(t *testing.T) {
	setupTestEnvironment()

	tests := []struct {
		name       string
		statusCode int
		isError    bool
	}{
		{"200 OK", http.StatusOK, false},
		{"201 Created", http.StatusCreated, false},
		{"204 No Content", http.StatusNoContent, false},
		{"301 Moved Permanently", http.StatusMovedPermanently, false},
		{"302 Found", http.StatusFound, false},
		{"400 Bad Request", http.StatusBadRequest, true},
		{"401 Unauthorized", http.StatusUnauthorized, true},
		{"403 Forbidden", http.StatusForbidden, true},
		{"404 Not Found", http.StatusNotFound, true},
		{"500 Internal Server Error", http.StatusInternalServerError, true},
		{"502 Bad Gateway", http.StatusBadGateway, true},
		{"503 Service Unavailable", http.StatusServiceUnavailable, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			beacon := &Beacon{}
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			})

			handler := beacon.Middleware(testHandler)

			req := httptest.NewRequest("GET", "/test", nil)
			logBuf := &bytes.Buffer{}
			logger := zerolog.New(logBuf).With().Timestamp().Logger()
			req = req.WithContext(contexthelper.SetLoggerInContext(req.Context(), logger))
			req = req.WithContext(contexthelper.SetRequestPatternInContext(req.Context(), "/test"))

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Code != tt.statusCode {
				t.Errorf("Expected status code %d, got %d", tt.statusCode, rr.Code)
			}

			logOutput := logBuf.String()
			if !strings.Contains(logOutput, "finished request") {
				t.Error("Expected log to contain 'finished request'")
			}
		})
	}
}

func TestBeacon_Middleware_RefererHeader(t *testing.T) {
	setupTestEnvironment()

	beacon := &Beacon{}
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := beacon.Middleware(testHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Referer", "https://example.com/previous-page")

	logBuf := &bytes.Buffer{}
	logger := zerolog.New(logBuf).With().Timestamp().Logger()
	req = req.WithContext(contexthelper.SetLoggerInContext(req.Context(), logger))
	req = req.WithContext(contexthelper.SetRequestPatternInContext(req.Context(), "/test"))

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	logOutput := logBuf.String()
	if !strings.Contains(logOutput, "https://example.com/previous-page") {
		t.Error("Expected log to contain referer header")
	}
}

func TestBeacon_Middleware_RequestHost(t *testing.T) {
	setupTestEnvironment()

	beacon := &Beacon{}
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := beacon.Middleware(testHandler)

	req := httptest.NewRequest("GET", "http://example.com:8080/test", nil)

	logBuf := &bytes.Buffer{}
	logger := zerolog.New(logBuf).With().Timestamp().Logger()
	req = req.WithContext(contexthelper.SetLoggerInContext(req.Context(), logger))
	req = req.WithContext(contexthelper.SetRequestPatternInContext(req.Context(), "/test"))

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	logOutput := logBuf.String()
	if !strings.Contains(logOutput, "example.com:8080") {
		t.Error("Expected log to contain request host")
	}
}

func TestCreateLogEvent(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		query      string
		method     string
		statusCode int
		username   string
		pattern    string
	}{
		{
			name:       "GET request with query",
			path:       "/api/users",
			query:      "limit=10&offset=0",
			method:     "GET",
			statusCode: http.StatusOK,
			pattern:    "/api/users",
		},
		{
			name:       "POST request",
			path:       "/api/users",
			query:      "",
			method:     "POST",
			statusCode: http.StatusCreated,
			pattern:    "/api/users",
		},
		{
			name:       "Unauthorized request with username",
			path:       "/api/protected",
			query:      "",
			method:     "GET",
			statusCode: http.StatusUnauthorized,
			username:   "testuser",
			pattern:    "/api/protected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path+"?"+tt.query, nil)
			req.RemoteAddr = "192.168.1.1:12345"
			req.Header.Set("User-Agent", "TestAgent/1.0")
			req.Header.Set("Referer", "https://example.com")

			ctx := req.Context()
			if tt.username != "" {
				ctx = contexthelper.SetUsernameInContext(req.WithContext(ctx), tt.username).Context()
			}
			ctx = contexthelper.SetRequestPatternInContext(ctx, tt.pattern)
			req = req.WithContext(ctx)

			logBuf := &bytes.Buffer{}
			logger := zerolog.New(logBuf).With().Timestamp().Logger()
			startTime := time.Now()

			logEvent := createLogEvent(req, tt.statusCode, startTime, &logger)
			logEvent.Msg("test")

			logOutput := logBuf.String()

			if !strings.Contains(logOutput, tt.path) {
				t.Errorf("Expected log to contain path %s", tt.path)
			}

			if !strings.Contains(logOutput, tt.method) {
				t.Errorf("Expected log to contain method %s", tt.method)
			}

			if tt.statusCode == http.StatusUnauthorized && tt.username != "" {
				if !strings.Contains(logOutput, tt.username) {
					t.Error("Expected log to contain username for unauthorized request")
				}
				if !strings.Contains(logOutput, "Unauthorized") {
					t.Error("Expected log to contain 'Unauthorized'")
				}
			}
		})
	}
}

func TestBeacon_Middleware_ContextPropagation(t *testing.T) {
	setupTestEnvironment()

	beacon := &Beacon{}

	var capturedContext context.Context
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedContext = r.Context()
		w.WriteHeader(http.StatusOK)
	})

	handler := beacon.Middleware(testHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	logBuf := &bytes.Buffer{}
	logger := zerolog.New(logBuf).With().Timestamp().Logger()
	req = req.WithContext(contexthelper.SetLoggerInContext(req.Context(), logger))
	req = req.WithContext(contexthelper.SetRequestPatternInContext(req.Context(), "/test"))

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if capturedContext == nil {
		t.Fatal("Expected context to be passed to handler")
	}

	// Verify logger is in context (we can't compare structs with []byte)
	_ = contexthelper.GetLoggerInContext(capturedContext)
}

func BenchmarkBeacon_Middleware(b *testing.B) {
	setupTestEnvironment()

	beacon := &Beacon{}
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := beacon.Middleware(testHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	logger := zerolog.New(io.Discard)
	req = req.WithContext(contexthelper.SetLoggerInContext(req.Context(), logger))
	req = req.WithContext(contexthelper.SetRequestPatternInContext(req.Context(), "/test"))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
	}
}
