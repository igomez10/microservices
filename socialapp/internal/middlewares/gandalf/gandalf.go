package gandalf

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/igomez10/microservices/socialapp/internal/contexthelper"
	"github.com/igomez10/microservices/socialapp/internal/jwt"
	"github.com/igomez10/microservices/socialapp/internal/middlewares/cache"
	"github.com/igomez10/microservices/socialapp/internal/tracerhelper"
	"github.com/igomez10/microservices/socialapp/pkg/controller/user"
	db "github.com/igomez10/microservices/socialapp/pkg/dbpgx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.opentelemetry.io/otel/attribute"
)

var gandalf_token_cache = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "gandalf_token_cache",
	Help: "The total number of gandalf token cache",
}, []string{"cache", "status"})

// gandalf_duration_microseconds_quantile is a histogram to track the duration of the gandalf.
var gandalf_duration_microseconds = promauto.NewSummaryVec(prometheus.SummaryOpts{
	Name:       "gandalf_duration_microseconds",
	Help:       "Summary for the runtime of the gandalf.",
	Objectives: map[float64]float64{0.25: 0.05, 0.50: 0.05, 0.75: 0.05, 1: 0.01}, // These are the default settings
}, []string{"auth_result"})

type Middleware struct {
	DB               db.Querier
	DBConn           *pgxpool.Pool
	Cache            *cache.Cache
	JWTSecret        string
	AllowlistedPaths map[string]map[string]bool
	AllowBasicAuth   bool
	AuthEndpoint     string
}

func (m *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracerhelper.GetTracer().Start(r.Context(), "middleware.gandalf")
		defer span.End()

		r = r.WithContext(ctx)

		start := time.Now()
		var authResult string
		log := contexthelper.GetLoggerInContext(r.Context())
		// get token from header
		isAllowlisted := m.AllowlistedPaths[r.URL.Path] != nil && m.AllowlistedPaths[r.URL.Path][r.Method]
		if isAllowlisted {
			span.SetAttributes(attribute.KeyValue{
				Key:   attribute.Key("allowlisted"),
				Value: attribute.StringValue(fmt.Sprintf("true")),
			})
			r = contexthelper.SetRequestedScopesInContext(r, map[string]bool{})
			log.Info().
				Str("path", r.URL.Path).
				Str("method", r.Method).
				Str("middleware", "gandalf").
				Msg("allowlisted path")
			authResult = "allowlisted"
		} else {
			authHeader := r.Header.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				givenToken := strings.TrimPrefix(authHeader, "Bearer ")
				token, err := jwt.NewTokenFromString(givenToken, m.JWTSecret)
				if err != nil {
					log.Error().
						Err(err).
						Str("token", givenToken).
						Msg("Failed to parse token")
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte(`{"code": 401, "message": "Invalid bearer token"}`))
					return
				}

				if token.Expires.Before(time.Now()) || token.NotBefore.After(time.Now()) {
					log.Error().
						Str("token", givenToken).
						Msg("Token invalid")
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte(`{"code": 401, "message": "Invalid bearer token"}`))
					return
				}

				// add scopes to request context
				scopes := map[string]bool{}
				for i := range token.Scopes {
					scopes[token.Scopes[i]] = true
				}
				r = contexthelper.SetRequestedScopesInContext(r, scopes)

				// add username to request context
				r = contexthelper.SetUsernameInContext(r, token.Username)

				authResult = "passed_with_jwt"
			} else if m.AllowBasicAuth || (strings.HasPrefix(authHeader, "Basic ") && r.URL.Path == m.AuthEndpoint) {
				// check grant type is client_credentials
				username, password, ok := r.BasicAuth()
				if !ok {
					log.Error().
						Str("username", username).
						Msg("Basic auth not ok")
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte(`{"code": 401, "message": "Invalid basic auth format"}`))
					return
				}

				ctx, dbUsernameSpan := tracerhelper.GetTracer().Start(r.Context(), "middleware.gandalf.db.get_user_by_username")
				usr, err := m.DB.GetUserByUsername(ctx, m.DBConn, username)
				dbUsernameSpan.End()
				switch err {
				case nil:
					// exit switch
				case pgx.ErrNoRows:
					log.Error().
						Err(err).
						Str("username", username).
						Msg("User not found")
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte(`{"code": 401, "message": "Invalid username or password"}`))
					return
				default:
					log.Error().
						Err(err).
						Msg("Error while getting user")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"code": 500, "message": "Error while getting user"}`))
					return
				}

				encryptedPassword := user.EncryptPassword(password, usr.Salt)
				if encryptedPassword != usr.HashedPassword {
					log.Error().
						Msg("Invalid password")
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte(`{"code": 401, "message": "Invalid username or password"}`))
					return
				}

				// passed authentication
				r = contexthelper.SetUsernameInContext(r, usr.Username)

				// get requested scopes
				requestedScopes := []string{}
				if len(r.FormValue("scope")) > 0 {
					requestedScopes = strings.Split(r.FormValue("scope"), " ")
				}

				// validate every requested scope exists in the DB

				// get user roles from DB
				ctx, dbRolesSpan := tracerhelper.GetTracer().Start(r.Context(), "middleware.gandalf.db.get_user_roles")
				roles, err := m.DB.GetUserRoles(ctx, m.DBConn, usr.ID)
				dbRolesSpan.End()
				switch err {
				case nil:
					// exit switch
				case pgx.ErrNoRows:
					// no roles found
				default:
					log.Error().
						Err(err).
						Msg("Error while getting user roles")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"code": 500, "message": "Error while getting user roles"}`))
					return
				}

				allowedScopes := map[string]db.Scope{}
				for i := range roles {
					// get role scopes from DB
					ctx, dbRoleScopesSpan := tracerhelper.GetTracer().Start(r.Context(), "middleware.gandalf.db.get_role_scopes")
					scopes, err := m.DB.ListRoleScopes(ctx, m.DBConn, db.ListRoleScopesParams{
						ID:     roles[i].ID,
						Limit:  10000,
						Offset: 0,
					})
					dbRoleScopesSpan.End()
					switch err {
					case nil:
						// exit switch
					case pgx.ErrNoRows:
						// no scopes found
					default:
						log.Error().
							Err(err).
							Msg("Error while getting role scopes")
						w.WriteHeader(http.StatusInternalServerError)
						w.Write([]byte(`{"code": 500, "message": "Error while getting role scopes"}`))
						return
					}

					for j := range scopes {
						allowedScopes[scopes[j].Name] = scopes[j]
					}
				}

				// remove duplicated scopes
				mapReqScopes := map[string]bool{}
				for _, scopeName := range requestedScopes {
					mapReqScopes[scopeName] = true
				}

				// verify requested scopes are allowed
				for i := range requestedScopes {
					if _, exist := allowedScopes[requestedScopes[i]]; !exist {
						log.Error().
							Str("scope", requestedScopes[i]).
							Msg("Scope not allowed")
						w.WriteHeader(http.StatusUnauthorized)
						w.Write([]byte(fmt.Sprintf(`{"code": 401, "message": "Scope %q not allowed"}`, requestedScopes[i])))
						return
					}
				}

				r = contexthelper.SetRequestedScopesInContext(r, mapReqScopes)
				authResult = "passed_with_basic"
			} else {
				mapReqScopes := map[string]bool{
					"noauth": true,
				}
				r = contexthelper.SetRequestedScopesInContext(r, mapReqScopes)
				authResult = "noauth"
			}
		}

		gandalf_duration_microseconds.WithLabelValues(authResult).Observe(float64(time.Since(start).Microseconds()))

		next.ServeHTTP(w, r)
	})

}
