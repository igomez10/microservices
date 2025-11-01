package gandalf

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/igomez10/microservices/socialapp/internal/contexthelper"
	"github.com/igomez10/microservices/socialapp/internal/jwt"
	"github.com/igomez10/microservices/socialapp/internal/tracerhelper"
	"github.com/igomez10/microservices/socialapp/pkg/controller/user"
	db "github.com/igomez10/microservices/socialapp/pkg/dbpgx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.opentelemetry.io/otel/attribute"
)

const (
	// Authentication results for metrics
	authResultAllowlisted = "allowlisted"
	authResultJWT         = "passed_with_jwt"
	authResultBasic       = "passed_with_basic"
	authResultNoAuth      = "noauth"

	// HTTP error messages
	errMsgInvalidBearerToken = `{"code": 401, "message": "Invalid bearer token"}`
	errMsgInvalidBasicAuth   = `{"code": 401, "message": "Invalid basic auth format"}`
	errMsgInvalidCredentials = `{"code": 401, "message": "Invalid username or password"}`
	errMsgErrorGettingUser   = `{"code": 500, "message": "Error while getting user"}`
	errMsgErrorGettingRoles  = `{"code": 500, "message": "Error while getting user roles"}`
	errMsgErrorGettingScopes = `{"code": 500, "message": "Error while getting role scopes"}`
	errMsgScopeNotAllowed    = `{"code": 401, "message": "Scope %q not allowed"}`

	// Scope query limits
	maxScopesPerRole = 10000
)

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
type Middleware struct {
	DB               db.Querier
	DBConn           *pgxpool.Pool
	JWTSecret        string
	AllowlistedPaths map[string]map[string]bool
	AllowBasicAuth   bool
	AuthEndpoint     string
}

// Authenticate is the main authentication middleware
func (m *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracerhelper.GetTracer().Start(r.Context(), "middleware.gandalf")
		defer span.End()

		r = r.WithContext(ctx)

		start := time.Now()
		log := contexthelper.GetLoggerInContext(r.Context())

		// Check if path is allowlisted
		if m.isPathAllowlisted(r) {
			m.handleAllowlistedPath(w, r, next, start, span, log)
			return
		}

		// Try Bearer token authentication
		authHeader := r.Header.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			m.handleBearerAuth(w, r, next, authHeader, start, log)
			return
		}

		// Try Basic authentication
		if m.shouldAllowBasicAuth(r, authHeader) {
			m.handleBasicAuth(w, r, next, start, log)
			return
		}

		// No authentication provided - set noauth scope
		m.handleNoAuth(w, r, next, start)
	})
}

// isPathAllowlisted checks if the requested path and method are allowlisted
func (m *Middleware) isPathAllowlisted(r *http.Request) bool {
	pathMethods, exists := m.AllowlistedPaths[r.URL.Path]
	return exists && pathMethods[r.Method]
}

// shouldAllowBasicAuth determines if basic auth should be allowed for this request
func (m *Middleware) shouldAllowBasicAuth(r *http.Request, authHeader string) bool {
	return m.AllowBasicAuth || (strings.HasPrefix(authHeader, "Basic ") && r.URL.Path == m.AuthEndpoint)
}

// handleAllowlistedPath processes requests to allowlisted paths
func (m *Middleware) handleAllowlistedPath(w http.ResponseWriter, r *http.Request, next http.Handler, start time.Time, span any, log any) {
	if s, ok := span.(interface{ SetAttributes(...attribute.KeyValue) }); ok {
		s.SetAttributes(attribute.Bool("allowlisted", true))
	}

	r = contexthelper.SetRequestedScopesInContext(r, map[string]bool{})

	if l, ok := log.(interface {
		Info() interface {
			Str(string, string) interface {
				Str(string, string) interface {
					Str(string, string) interface{ Msg(string) }
				}
			}
		}
	}); ok {
		l.Info().Str("path", r.URL.Path).Str("method", r.Method).Str("middleware", "gandalf").Msg("allowlisted path")
	}

	authenticationDuration.WithLabelValues(authResultAllowlisted).Observe(time.Since(start).Seconds())
	next.ServeHTTP(w, r)
}

// handleBearerAuth processes Bearer token authentication
func (m *Middleware) handleBearerAuth(w http.ResponseWriter, r *http.Request, next http.Handler, authHeader string, start time.Time, log any) {
	givenToken := strings.TrimPrefix(authHeader, "Bearer ")

	token, err := jwt.NewTokenFromString(givenToken, m.JWTSecret)
	if err != nil {
		if l, ok := log.(interface {
			Error() interface {
				Err(error) interface {
					Str(string, string) interface{ Msg(string) }
				}
			}
		}); ok {
			l.Error().Err(err).Str("token", givenToken).Msg("Failed to parse token")
		}
		m.writeUnauthorized(w, errMsgInvalidBearerToken)
		authenticationDuration.WithLabelValues("failed_jwt").Observe(time.Since(start).Seconds())
		return
	}

	if !m.isTokenValid(token) {
		if l, ok := log.(interface {
			Error() interface {
				Str(string, string) interface{ Msg(string) }
			}
		}); ok {
			l.Error().Str("token", givenToken).Msg("Token invalid or expired")
		}
		m.writeUnauthorized(w, errMsgInvalidBearerToken)
		authenticationDuration.WithLabelValues("failed_jwt_expired").Observe(time.Since(start).Seconds())
		return
	}

	// Add scopes and username to context
	scopes := m.convertTokenScopes(token.Scopes)
	r = contexthelper.SetRequestedScopesInContext(r, scopes)
	r = contexthelper.SetUsernameInContext(r, token.Username)

	authenticationDuration.WithLabelValues(authResultJWT).Observe(time.Since(start).Seconds())
	next.ServeHTTP(w, r)
}

// handleBasicAuth processes Basic authentication
func (m *Middleware) handleBasicAuth(w http.ResponseWriter, r *http.Request, next http.Handler, start time.Time, log any) {
	username, password, ok := r.BasicAuth()
	if !ok {
		if l, ok := log.(interface {
			Error() interface {
				Str(string, string) interface{ Msg(string) }
			}
		}); ok {
			l.Error().Str("username", username).Msg("Basic auth not ok")
		}
		m.writeUnauthorized(w, errMsgInvalidBasicAuth)
		authenticationDuration.WithLabelValues("failed_basic_format").Observe(time.Since(start).Seconds())
		return
	}

	// Get user from database
	usr, err := m.getUserByUsername(r, username, log)
	if err != nil {
		if err == pgx.ErrNoRows {
			if l, ok := log.(interface {
				Debug() interface {
					Err(error) interface {
						Str(string, string) interface{ Msg(string) }
					}
				}
			}); ok {
				l.Debug().Err(err).Str("username", username).Msg("User not found")
			}
			m.writeUnauthorized(w, errMsgInvalidCredentials)
			authenticationDuration.WithLabelValues("failed_basic_user_not_found").Observe(time.Since(start).Seconds())
			return
		}
		if l, ok := log.(interface {
			Error() interface {
				Err(error) interface{ Msg(string) }
			}
		}); ok {
			l.Error().Err(err).Msg("Error while getting user")
		}
		m.writeInternalError(w, errMsgErrorGettingUser)
		authenticationDuration.WithLabelValues("failed_basic_db_error").Observe(time.Since(start).Seconds())
		return
	}

	// Verify password
	if !m.verifyPassword(password, usr) {
		if l, ok := log.(interface {
			Error() interface{ Msg(string) }
		}); ok {
			l.Error().Msg("Invalid password")
		}
		m.writeUnauthorized(w, errMsgInvalidCredentials)
		authenticationDuration.WithLabelValues("failed_basic_invalid_password").Observe(time.Since(start).Seconds())
		return
	}

	// Set username in context
	r = contexthelper.SetUsernameInContext(r, usr.Username)

	// Get and validate scopes
	requestedScopes := m.parseRequestedScopes(r)
	allowedScopes, err := m.getAllowedScopes(r, usr.ID, log)
	if err != nil {
		m.writeInternalError(w, errMsgErrorGettingScopes)
		authenticationDuration.WithLabelValues("failed_basic_scopes_error").Observe(time.Since(start).Seconds())
		return
	}

	// Validate requested scopes
	for _, scopeName := range requestedScopes {
		if _, exists := allowedScopes[scopeName]; !exists {
			if l, ok := log.(interface {
				Error() interface {
					Str(string, string) interface{ Msg(string) }
				}
			}); ok {
				l.Error().Str("scope", scopeName).Msg("Scope not allowed")
			}
			m.writeUnauthorized(w, fmt.Sprintf(errMsgScopeNotAllowed, scopeName))
			authenticationDuration.WithLabelValues("failed_basic_scope_denied").Observe(time.Since(start).Seconds())
			return
		}
	}

	// Convert to map and set in context
	scopeMap := m.convertScopeSliceToMap(requestedScopes)
	r = contexthelper.SetRequestedScopesInContext(r, scopeMap)

	authenticationDuration.WithLabelValues(authResultBasic).Observe(time.Since(start).Seconds())
	next.ServeHTTP(w, r)
}

// handleNoAuth processes requests without authentication
func (m *Middleware) handleNoAuth(w http.ResponseWriter, r *http.Request, next http.Handler, start time.Time) {
	r = contexthelper.SetRequestedScopesInContext(r, map[string]bool{"noauth": true})
	authenticationDuration.WithLabelValues(authResultNoAuth).Observe(time.Since(start).Seconds())
	next.ServeHTTP(w, r)
}

// Helper methods

// isTokenValid checks if a JWT token is valid and not expired
func (m *Middleware) isTokenValid(token *jwt.SocialAPPToken) bool {
	now := time.Now()
	return token.Expires != nil && token.NotBefore != nil &&
		!token.Expires.Time.Before(now) && !token.NotBefore.Time.After(now)
}

// convertTokenScopes converts a slice of scopes to a map
func (m *Middleware) convertTokenScopes(scopes []string) map[string]bool {
	scopeMap := make(map[string]bool, len(scopes))
	for _, scope := range scopes {
		scopeMap[scope] = true
	}
	return scopeMap
}

// getUserByUsername retrieves a user from the database
func (m *Middleware) getUserByUsername(r *http.Request, username string, log any) (db.User, error) {
	ctx, span := tracerhelper.GetTracer().Start(r.Context(), "middleware.gandalf.db.get_user_by_username")
	defer span.End()
	return m.DB.GetUserByUsername(ctx, m.DBConn, username)
}

// verifyPassword checks if the provided password matches the user's hashed password
func (m *Middleware) verifyPassword(password string, usr db.User) bool {
	encryptedPassword := user.EncryptPassword(password, usr.Salt)
	return encryptedPassword == usr.HashedPassword
}

// parseRequestedScopes extracts requested scopes from the request
func (m *Middleware) parseRequestedScopes(r *http.Request) []string {
	scopeStr := r.FormValue("scope")
	if scopeStr == "" {
		return []string{}
	}
	return strings.Split(scopeStr, " ")
}

// getAllowedScopes retrieves all scopes allowed for a user based on their roles
func (m *Middleware) getAllowedScopes(r *http.Request, userID int64, log any) (map[string]db.Scope, error) {
	ctx, rolesSpan := tracerhelper.GetTracer().Start(r.Context(), "middleware.gandalf.db.get_user_roles")
	roles, err := m.DB.GetUserRoles(ctx, m.DBConn, userID)
	rolesSpan.End()

	if err != nil && err != pgx.ErrNoRows {
		if l, ok := log.(interface {
			Error() interface {
				Err(error) interface{ Msg(string) }
			}
		}); ok {
			l.Error().Err(err).Msg("Error while getting user roles")
		}
		return nil, err
	}

	allowedScopes := make(map[string]db.Scope)
	for _, role := range roles {
		scopes, err := m.getRoleScopes(r, role.ID, log)
		if err != nil {
			return nil, err
		}

		for _, scope := range scopes {
			allowedScopes[scope.Name] = scope
		}
	}

	return allowedScopes, nil
}

// getRoleScopes retrieves all scopes for a specific role
func (m *Middleware) getRoleScopes(r *http.Request, roleID int64, log any) ([]db.Scope, error) {
	ctx, span := tracerhelper.GetTracer().Start(r.Context(), "middleware.gandalf.db.get_role_scopes")
	defer span.End()

	scopes, err := m.DB.ListRoleScopes(ctx, m.DBConn, db.ListRoleScopesParams{
		ID:     roleID,
		Limit:  maxScopesPerRole,
		Offset: 0,
	})

	if err != nil && err != pgx.ErrNoRows {
		if l, ok := log.(interface {
			Error() interface {
				Err(error) interface{ Msg(string) }
			}
		}); ok {
			l.Error().Err(err).Msg("Error while getting role scopes")
		}
		return nil, err
	}

	return scopes, nil
}

// convertScopeSliceToMap converts a slice of scope names to a map
func (m *Middleware) convertScopeSliceToMap(scopes []string) map[string]bool {
	scopeMap := make(map[string]bool, len(scopes))
	for _, scope := range scopes {
		scopeMap[scope] = true
	}
	return scopeMap
}

// HTTP response helpers

// writeUnauthorized writes a 401 Unauthorized response
func (m *Middleware) writeUnauthorized(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(message))
}

// writeInternalError writes a 500 Internal Server Error response
func (m *Middleware) writeInternalError(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(message))
}
