package gandalf

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/igomez10/microservices/socialapp/internal/contexthelper"
	"github.com/igomez10/microservices/socialapp/internal/jwt"
	"github.com/igomez10/microservices/socialapp/internal/tracerhelper"
	"github.com/igomez10/microservices/socialapp/pkg/controller/user"
	db "github.com/igomez10/microservices/socialapp/pkg/dbpgx"
	"github.com/igomez10/microservices/socialapp/pkg/scopes"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.opentelemetry.io/otel/attribute"
)

// AuthenticationResult represents the result of an authentication attempt
type AuthenticationResult string

const (
	// Authentication results for metrics
	authResultAllowlisted = "allowlisted"
	authResultJWT         = "passed_with_jwt"
	authResultBasic       = "passed_with_basic"
	authResultNoAuth      = "noauth"
	authResultFailedJWT   = "failed_jwt"
	authResultFailedBasic = "failed_basic"

	// HTTP error messages
	errMsgInvalidBearerToken = `{"code": 401, "message": "Invalid bearer token"}`
	errMsgInvalidBasicAuth   = `{"code": 401, "message": "Invalid basic auth format"}`
	errMsgInvalidCredentials = `{"code": 401, "message": "Invalid username or password"}`
	errMsgErrorGettingUser   = `{"code": 500, "message": "Error while getting user"}`
	errMsgErrorGettingRoles  = `{"code": 500, "message": "Error while getting user roles"}`
	errMsgErrorGettingScopes = `{"code": 500, "message": "Error while getting role scopes"}`
	errMsgScopeNotAllowed    = `{"code": 401, "message": "Scope %q not allowed"}`
	errMsgInvalidScope       = `{"code": 400, "message": "Invalid scope %q: scope is not defined in API specification"}`

	// Scope query limits
	maxScopesPerRole = 10000
)

// Prometheus metrics for Gandalf middleware
var gandalfTokenCache = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "gandalf_token_cache",
	Help: "The total number of gandalf token cache operations",
}, []string{"cache", "status"})

// authenticationDuration tracks the duration of authentication operations
var authenticationDuration = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "authentication_duration_seconds",
		Help:    "Duration of authentication operations in seconds",
		Buckets: prometheus.DefBuckets,
	},
	[]string{"auth_result"},
)

// Middleware handles authentication for incoming requests
// It supports Bearer token authentication (JWT) and Basic authentication
type Middleware struct {
	// DB is the database querier used to fetch user information and roles
	DB db.Querier

	// DBConn is the database connection pool used for database operations
	DBConn *pgxpool.Pool

	// JWTSecret is the secret key used to sign and verify JWT tokens
	JWTSecret string

	// AllowlistedPaths defines which paths and methods are allowed without authentication
	AllowlistedPaths map[string]map[string]bool

	// AllowBasicAuth controls whether basic authentication is allowed for non-auth endpoints
	AllowBasicAuth bool

	// AuthEndpoint defines the endpoint that accepts basic authentication (e.g., "/auth")
	AuthEndpoint string
}

// Authenticate is the main authentication middleware handler.
// It processes incoming HTTP requests and applies authentication based on:
// 1. Allowlisted paths (no auth required)
// 2. Bearer token authentication (JWT)
// 3. Basic authentication (username/password)
// 4. No authentication (if no auth header provided)
func (m *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracerhelper.GetTracer().Start(r.Context(), "middleware.gandalf")
		defer span.End()

		r = r.WithContext(ctx)

		start := time.Now()
		logger := contexthelper.GetLoggerInContext(r.Context())

		// Log the incoming request
		logger.Info("Processing authentication request",
			slog.String("path", r.URL.Path),
			slog.String("method", r.Method),
			slog.String("remote_addr", r.RemoteAddr),
		)

		// Check if path is allowlisted
		if m.isPathAllowlisted(r) {
			m.handleAllowlistedPath(w, r, next, start, span, logger)
			return
		}

		// Try Bearer token authentication
		authHeader := r.Header.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			m.handleBearerAuth(w, r, next, authHeader, start, logger)
			return
		}

		// Try Basic authentication
		if m.shouldAllowBasicAuth(r, authHeader) {
			m.handleBasicAuth(w, r, next, start, logger)
			return
		}

		// No authentication provided - set noauth scope and continue
		m.handleNoAuth(w, r, next, start)
	})
}

// isPathAllowlisted checks if the requested path and method are allowlisted
// This allows certain paths to be accessed without authentication, useful for health checks,
// public endpoints, or other unauthenticated access points.
func (m *Middleware) isPathAllowlisted(r *http.Request) bool {
	pathMethods, exists := m.AllowlistedPaths[r.URL.Path]
	return exists && pathMethods[r.Method]
}

// shouldAllowBasicAuth determines if basic auth should be allowed for this request.
// Basic auth is allowed either when:
// 1. AllowBasicAuth is true (global setting)
// 2. The request is to the auth endpoint and contains a Basic auth header
func (m *Middleware) shouldAllowBasicAuth(r *http.Request, authHeader string) bool {
	return m.AllowBasicAuth || (strings.HasPrefix(authHeader, "Basic ") && r.URL.Path == m.AuthEndpoint)
}

// handleAllowlistedPath processes requests to allowlisted paths.
// These paths are exempt from authentication and proceed directly to the next handler.
func (m *Middleware) handleAllowlistedPath(w http.ResponseWriter, r *http.Request, next http.Handler, start time.Time, span any, logger *slog.Logger) {
	if s, ok := span.(interface{ SetAttributes(...attribute.KeyValue) }); ok {
		s.SetAttributes(attribute.Bool("allowlisted", true))
	}

	logger.Info("Allowlisted path accessed",
		slog.String("path", r.URL.Path),
		slog.String("method", r.Method),
	)

	r = contexthelper.SetRequestedScopesInContext(r, map[string]bool{})

	authenticationDuration.WithLabelValues(authResultAllowlisted).Observe(time.Since(start).Seconds())
	next.ServeHTTP(w, r)
}

// handleBearerAuth processes Bearer token authentication (JWT).
// It validates the JWT token and extracts user information and scopes.
func (m *Middleware) handleBearerAuth(w http.ResponseWriter, r *http.Request, next http.Handler, authHeader string, start time.Time, logger *slog.Logger) {
	givenToken := strings.TrimPrefix(authHeader, "Bearer ")

	// Parse the JWT token
	token, err := jwt.NewTokenFromString(givenToken, m.JWTSecret)
	if err != nil {
		logger.Error("Failed to parse JWT token",
			slog.String("error", err.Error()),
			slog.String("token_length", fmt.Sprintf("%d", len(givenToken))),
		)
		m.writeUnauthorized(w, errMsgInvalidBearerToken)
		authenticationDuration.WithLabelValues(authResultFailedJWT).Observe(time.Since(start).Seconds())
		return
	}

	// Validate token expiration and validity
	if !m.isTokenValid(token) {
		logger.Error("JWT token invalid or expired",
			slog.String("username", token.Username),
			slog.Time("expires_at", token.Expires.Time),
			slog.Time("not_before", token.NotBefore.Time),
		)
		m.writeUnauthorized(w, errMsgInvalidBearerToken)
		authenticationDuration.WithLabelValues(authResultFailedJWT).Observe(time.Since(start).Seconds())
		return
	}

	// Add scopes and username to context for downstream handlers
	scopes := m.convertTokenScopes(token.Scopes)
	r = contexthelper.SetRequestedScopesInContext(r, scopes)
	r = contexthelper.SetUsernameInContext(r, token.Username)

	logger.Info("Bearer token authentication successful",
		slog.String("username", token.Username),
		slog.Int("scope_count", len(scopes)),
	)

	authenticationDuration.WithLabelValues(authResultJWT).Observe(time.Since(start).Seconds())
	next.ServeHTTP(w, r)
}

// handleBasicAuth processes Basic authentication (username/password).
// It validates credentials against the database and checks requested scopes.
func (m *Middleware) handleBasicAuth(w http.ResponseWriter, r *http.Request, next http.Handler, start time.Time, logger *slog.Logger) {
	username, password, ok := r.BasicAuth()
	if !ok {
		logger.Error("Basic auth header invalid",
			slog.String("username", username),
		)
		m.writeUnauthorized(w, errMsgInvalidBasicAuth)
		authenticationDuration.WithLabelValues(authResultFailedBasic).Observe(time.Since(start).Seconds())
		return
	}

	logger.Debug("Basic auth attempt",
		slog.String("username", username),
	)

	// Get user from database
	usr, err := m.getUserByUsername(r, username)
	if err != nil {
		if err == pgx.ErrNoRows {
			logger.Debug("User not found in database",
				slog.String("username", username),
				slog.String("error", err.Error()),
			)
			m.writeUnauthorized(w, errMsgInvalidCredentials)
			authenticationDuration.WithLabelValues(authResultFailedBasic).Observe(time.Since(start).Seconds())
			return
		}
		logger.Error("Database error while getting user",
			slog.String("username", username),
			slog.String("error", err.Error()),
		)
		m.writeInternalError(w, errMsgErrorGettingUser)
		authenticationDuration.WithLabelValues(authResultFailedBasic).Observe(time.Since(start).Seconds())
		return
	}

	// Verify password
	if !m.verifyPassword(password, usr) {
		logger.Error("Invalid password provided",
			slog.String("username", username),
		)
		m.writeUnauthorized(w, errMsgInvalidCredentials)
		authenticationDuration.WithLabelValues(authResultFailedBasic).Observe(time.Since(start).Seconds())
		return
	}

	logger.Debug("Password verification successful",
		slog.String("username", usr.Username),
	)

	// Set username in context
	r = contexthelper.SetUsernameInContext(r, usr.Username)

	// Parse and validate requested scopes
	requestedScopes, err := m.parseRequestedScopes(r)
	if err != nil {
		// Invalid scope requested - return 400 Bad Request
		logger.Error("Invalid scope requested",
			slog.String("error", err.Error()),
			slog.String("scope", r.FormValue("scope")),
		)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(errMsgInvalidScope, err.Error())))
		authenticationDuration.WithLabelValues(authResultFailedBasic).Observe(time.Since(start).Seconds())
		return
	}

	logger.Debug("Valid scopes requested",
		slog.Int("scope_count", len(requestedScopes)),
	)

	// Get allowed scopes for this user
	allowedScopes, err := m.getAllowedScopes(r, usr.ID, logger)
	if err != nil {
		logger.Error("Database error while getting user scopes",
			slog.String("username", usr.Username),
			slog.String("error", err.Error()),
		)
		m.writeInternalError(w, errMsgErrorGettingScopes)
		authenticationDuration.WithLabelValues(authResultFailedBasic).Observe(time.Since(start).Seconds())
		return
	}

	logger.Debug("Retrieved allowed scopes for user",
		slog.String("username", usr.Username),
		slog.Int("allowed_scope_count", len(allowedScopes)),
	)

	// Validate requested scopes
	for _, scopeName := range requestedScopes {
		if _, exists := allowedScopes[scopeName]; !exists {
			logger.Error("Requested scope not allowed for user",
				slog.String("username", usr.Username),
				slog.String("scope", scopeName),
			)
			m.writeUnauthorized(w, fmt.Sprintf(errMsgScopeNotAllowed, scopeName))
			authenticationDuration.WithLabelValues(authResultFailedBasic).Observe(time.Since(start).Seconds())
			return
		}
	}

	logger.Debug("All requested scopes are allowed",
		slog.String("username", usr.Username),
	)

	// Convert to map and set in context
	scopeMap := m.convertScopeSliceToMap(requestedScopes)
	r = contexthelper.SetRequestedScopesInContext(r, scopeMap)

	logger.Info("Basic authentication successful",
		slog.String("username", usr.Username),
		slog.Int("scope_count", len(requestedScopes)),
	)

	authenticationDuration.WithLabelValues(authResultBasic).Observe(time.Since(start).Seconds())
	next.ServeHTTP(w, r)
}

// handleNoAuth processes requests without authentication.
// These requests proceed to the next handler with empty scopes.
func (m *Middleware) handleNoAuth(w http.ResponseWriter, r *http.Request, next http.Handler, start time.Time) {
	logger := contexthelper.GetLoggerInContext(r.Context())

	logger.Debug("No authentication provided",
		slog.String("path", r.URL.Path),
		slog.String("method", r.Method),
	)

	r = contexthelper.SetRequestedScopesInContext(r, map[string]bool{})
	authenticationDuration.WithLabelValues(authResultNoAuth).Observe(time.Since(start).Seconds())
	next.ServeHTTP(w, r)
}

// Helper methods

// isTokenValid checks if a JWT token is valid and not expired.
// A token is considered valid if:
// 1. It has an expiration time that is not in the past
// 2. It has a "not before" time that is not in the future
func (m *Middleware) isTokenValid(token *jwt.SocialAPPToken) bool {
	now := time.Now()
	return token.Expires != nil && token.NotBefore != nil &&
		!token.Expires.Time.Before(now) && !token.NotBefore.Time.After(now)
}

// convertTokenScopes converts and validates a slice of scopes to a map.
// Only scopes that are defined in the API specification are included.
// Invalid scopes in JWT tokens are silently ignored to maintain backwards compatibility.
func (m *Middleware) convertTokenScopes(tokenScopes []string) map[string]bool {
	scopeMap := make(map[string]bool, len(tokenScopes))
	for _, scopeName := range tokenScopes {
		// Validate that the scope is defined in the API specification
		if _, valid := scopes.FromString(scopeName); valid {
			scopeMap[scopeName] = true
		}
		// Silently ignore undefined scopes in JWT tokens to maintain backwards compatibility
		// They simply won't be granted
	}
	return scopeMap
}

// getUserByUsername retrieves a user from the database by username.
// This method is used during basic authentication to verify credentials.
func (m *Middleware) getUserByUsername(r *http.Request, username string) (db.User, error) {
	ctx, span := tracerhelper.GetTracer().Start(r.Context(), "middleware.gandalf.db.get_user_by_username")
	defer span.End()
	return m.DB.GetUserByUsername(ctx, m.DBConn, username)
}

// verifyPassword checks if the provided password matches the user's hashed password.
// It uses the same hashing algorithm and salt that was used when the password was originally stored.
func (m *Middleware) verifyPassword(password string, usr db.User) bool {
	encryptedPassword := user.EncryptPassword(password, usr.Salt)
	return encryptedPassword == usr.HashedPassword
}

// parseRequestedScopes extracts and validates requested scopes from the request.
// Returns error if any requested scope is not defined in the API specification.
func (m *Middleware) parseRequestedScopes(r *http.Request) ([]string, error) {
	scopeStr := r.FormValue("scope")
	if scopeStr == "" {
		return []string{}, nil
	}

	requestedScopes := strings.Split(scopeStr, " ")
	validScopes := make([]string, 0, len(requestedScopes))

	// Validate that all requested scopes are defined in the API specification
	for _, scopeName := range requestedScopes {
		if _, valid := scopes.FromString(scopeName); !valid {
			// Return error for undefined scope
			return nil, fmt.Errorf("undefined scope: %s", scopeName)
		}
		validScopes = append(validScopes, scopeName)
	}

	return validScopes, nil
}

// getAllowedScopes retrieves all scopes allowed for a user based on their roles.
// It fetches the user's roles and then gets all scopes for each role.
func (m *Middleware) getAllowedScopes(r *http.Request, userID int64, logger *slog.Logger) (map[string]db.Scope, error) {
	ctx, rolesSpan := tracerhelper.GetTracer().Start(r.Context(), "middleware.gandalf.db.get_user_roles")
	roles, err := m.DB.GetUserRoles(ctx, m.DBConn, userID)
	rolesSpan.End()

	if err != nil && err != pgx.ErrNoRows {
		logger.Error("Database error while getting user roles",
			slog.Int64("user_id", userID),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	allowedScopes := make(map[string]db.Scope)
	for _, role := range roles {
		scopes, err := m.getRoleScopes(r, role.ID, logger)
		if err != nil {
			return nil, err
		}

		for _, scope := range scopes {
			allowedScopes[scope.Name] = scope
		}
	}

	return allowedScopes, nil
}

// getRoleScopes retrieves all scopes for a specific role.
// This method is used to build the complete set of scopes available to a user.
func (m *Middleware) getRoleScopes(r *http.Request, roleID int64, logger *slog.Logger) ([]db.Scope, error) {
	ctx, span := tracerhelper.GetTracer().Start(r.Context(), "middleware.gandalf.db.get_role_scopes")
	defer span.End()

	scopes, err := m.DB.ListRoleScopes(ctx, m.DBConn, db.ListRoleScopesParams{
		ID:     roleID,
		Limit:  maxScopesPerRole,
		Offset: 0,
	})

	if err != nil && err != pgx.ErrNoRows {
		logger.Error("Database error while getting role scopes",
			slog.Int64("role_id", roleID),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	return scopes, nil
}

// convertScopeSliceToMap converts a slice of scope names to a map.
// This is used to efficiently check if a scope is present in the set of allowed scopes.
func (m *Middleware) convertScopeSliceToMap(scopes []string) map[string]bool {
	scopeMap := make(map[string]bool, len(scopes))
	for _, scope := range scopes {
		scopeMap[scope] = true
	}
	return scopeMap
}

// HTTP response helpers

// writeUnauthorized writes a 401 Unauthorized response with the provided message.
func (m *Middleware) writeUnauthorized(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(message))
}

// writeInternalError writes a 500 Internal Server Error response with the provided message.
func (m *Middleware) writeInternalError(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(message))
}
