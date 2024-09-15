package gandalf

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/igomez10/microservices/socialapp/internal/contexthelper"
	"github.com/igomez10/microservices/socialapp/internal/middlewares/cache"
	"github.com/igomez10/microservices/socialapp/internal/tracerhelper"
	"github.com/igomez10/microservices/socialapp/pkg/controller/user"
	"github.com/igomez10/microservices/socialapp/pkg/db"
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
	DBConn           *sql.DB
	Cache            *cache.Cache
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
		if m.AllowlistedPaths[r.URL.Path] != nil && m.AllowlistedPaths[r.URL.Path][r.Method] {
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
				var token db.Token
				// check token in cache
				cachedTokenBytes, err := m.Cache.Client.Get(r.Context(), "token_"+givenToken).Result()
				if err != nil {
					gandalf_token_cache.WithLabelValues("token", "miss").Inc()
					// token not in cache, get from db
					dbtoken, err := m.DB.GetToken(r.Context(), m.DBConn, givenToken)
					switch err {
					case nil:
						token = dbtoken
						// save token in cache
						buf := &bytes.Buffer{}
						encoder := gob.NewEncoder(buf)
						if err := encoder.Encode(token); err != nil {
							log.Error().
								Err(err).
								Str("middleware", "gandalf").
								Msg("Failed to marshal token")
						} else {
							m.Cache.Client.Set(r.Context(), "token_"+givenToken, buf.Bytes(), time.Hour*1)
						}
					case sql.ErrNoRows:
						log.Error().
							Err(err).
							Msg("Token not found")
						w.WriteHeader(http.StatusUnauthorized)
						w.Write([]byte(`{"code": 401, "message": "Invalid bearer token"}`))
						return
					default:
						log.Error().
							Err(err).
							Msg("Error while getting token")
						w.WriteHeader(http.StatusInternalServerError)
						w.Write([]byte(`{"code": 500, "message": "Error while getting token"}`))
						return
					}
				} else {
					// token found in cache
					gandalf_token_cache.WithLabelValues("token", "hit").Inc()
					b := bytes.Buffer{}
					b.Write([]byte(cachedTokenBytes))
					d := gob.NewDecoder(&b)
					if err := d.Decode(&token); err != nil {
						log.Error().
							Err(err).
							Str("middleware", "gandalf").
							Msg("Failed to unmarshal token stored in cache")
						w.WriteHeader(http.StatusInternalServerError)
						w.Write([]byte(`{"code": 500, "message": "Error while getting token"}`))
						return
					}
				}

				if time.Now().After(token.ValidUntil) {
					log.Error().
						Err(err).
						Msg("Token expired")
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte(`{"code": 401, "message": "Token expired"}`))
					return
				}

				// get token scopes
				var scopes []db.Scope
				// check scopes in cache
				cachedScopes, err := m.Cache.Client.Get(r.Context(), "token_to_scopes_"+givenToken).Result()
				if err != nil {
					gandalf_token_cache.WithLabelValues("scopes", "miss").Inc()
					// scopes not in cache, get from db
					dbTokenScopes, err := m.DB.GetTokenScopes(r.Context(), m.DBConn, token.ID)
					switch err {
					case nil:
						scopes = dbTokenScopes
						// save scopes in cache
						buf := &bytes.Buffer{}
						encoder := gob.NewEncoder(buf)
						if err := encoder.Encode(scopes); err != nil {
							log.Error().
								Err(err).
								Str("middleware", "gandalf").
								Msg("Failed to marshal scopes")
						} else {
							m.Cache.Client.Set(r.Context(), "token_to_scopes_"+givenToken, buf.Bytes(), time.Hour*1)
						}

					case sql.ErrNoRows:
						log.Error().
							Err(err).
							Msg("Token scopes not found")
						w.WriteHeader(http.StatusUnauthorized)
						w.Write([]byte(`{"code": 401, "message": "Invalid bearer token"}`))
						return
					default:
						log.Error().Err(err).Msg("Error while getting token scopes")
						http.Error(w, "Error while getting token scopes", http.StatusInternalServerError)
						w.Write([]byte(`{"code": 500, "message": "Error while getting token scopes"}`))
						return
					}
				} else {
					gandalf_token_cache.WithLabelValues("scopes", "hit").Inc()
					// scopes found in cache
					b := bytes.Buffer{}
					b.Write([]byte(cachedScopes))
					d := gob.NewDecoder(&b)
					if err := d.Decode(&scopes); err != nil {
						log.Error().
							Err(err).
							Str("middleware", "gandalf").
							Msg("Failed to unmarshal token scopes")
						w.WriteHeader(http.StatusInternalServerError)
						w.Write([]byte(`{"code": 500, "message": "Error while getting token scopes"}`))
						return
					}
				}

				scopesMap := map[string]bool{}
				for i := range scopes {
					scopesMap[scopes[i].Name] = true
				}

				// add scopes to request context
				r = contexthelper.SetRequestedScopesInContext(r, scopesMap)
				// check user in cache
				var username string
				cacheKey := "userid_to_username" + strconv.Itoa(int(token.UserID))
				cached_username, err := m.Cache.Client.Get(r.Context(), cacheKey).Result()
				if err != nil {
					gandalf_token_cache.WithLabelValues("user", "miss").Inc()
					// user not in cache, get from db
					usr, err := m.DB.GetUserByID(r.Context(), m.DBConn, token.UserID)
					if err != nil {
						log.Error().
							Err(err).
							Int64("userID", token.UserID).
							Msg("Failed to get user from token")
						w.WriteHeader(http.StatusInternalServerError)
						w.Write([]byte(`{"code": 500, "message": "Failed to get user from token"}`))
						return
					}
					username = usr.Username
					// store it in cache
					m.Cache.Client.Set(r.Context(), cacheKey, username, time.Minute*10)
				} else {
					gandalf_token_cache.WithLabelValues("user", "hit").Inc()
					username = cached_username
				}

				r = contexthelper.SetUsernameInContext(r, username)
				authResult = "passed_with_bearer"
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

				usr, err := m.DB.GetUserByUsername(r.Context(), m.DBConn, username)
				switch err {
				case nil:
					// exit switch
				case sql.ErrNoRows:
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
				roles, err := m.DB.GetUserRoles(r.Context(), m.DBConn, usr.ID)
				switch err {
				case nil:
					// exit switch
				case sql.ErrNoRows:
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
					scopes, err := m.DB.ListRoleScopes(r.Context(), m.DBConn, db.ListRoleScopesParams{
						ID:     roles[i].ID,
						Limit:  10000,
						Offset: 0,
					})
					switch err {
					case nil:
						// exit switch
					case sql.ErrNoRows:
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
