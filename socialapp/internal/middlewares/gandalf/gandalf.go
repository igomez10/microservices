package gandalf

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net/http"
	"strconv"
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
				if strings.Contains(givenToken, ".") {
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
				} else {
					var token db.Token
					// check token in cache
					ctxDBToken, spanScopes := tracerhelper.GetTracer().Start(r.Context(), "middleware.gandalf.cache.get_token")
					cachedTokenBytes, err := m.Cache.Client.Get(ctxDBToken, "token_"+givenToken).Result()
					spanScopes.End()
					if err != nil {
						gandalf_token_cache.WithLabelValues("token", "miss").Inc()
						// token not in cache, get from db
						ctx, dbSpan := tracerhelper.GetTracer().Start(r.Context(), "middleware.gandalf.db.get_token")
						dbtoken, err := m.DB.GetToken(ctx, m.DBConn, givenToken)
						dbSpan.End()
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
								ctx, span := tracerhelper.GetTracer().Start(r.Context(), "middleware.gandalf.cache.set_token")
								m.Cache.Client.Set(ctx, "token_"+givenToken, buf.Bytes(), time.Hour*1)
								span.End()
							}
						case pgx.ErrNoRows:
							log.Error().
								Err(err).
								Msg("Token not found")
							w.WriteHeader(http.StatusUnauthorized)
							w.Write([]byte(`{"code": 401, "message": "Invalid bearer token"}`))
							return
						default:
							log.Error().
								Err(err).
								Str("error type", fmt.Sprintf("%T", err)).
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

					if time.Now().After(token.ValidUntil.Time) {
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
					ctxScopes, spanScopes := tracerhelper.GetTracer().Start(r.Context(), "middleware.gandalf.cache.get_scopes")
					cachedScopes, err := m.Cache.Client.Get(ctxScopes, "token_to_scopes_"+givenToken).Result()
					spanScopes.End()
					if err != nil {
						gandalf_token_cache.WithLabelValues("scopes", "miss").Inc()
						// scopes not in cache, get from db
						ctxscopes, spanScopes := tracerhelper.GetTracer().Start(r.Context(), "middleware.gandalf.db.get_token_scopes")
						dbTokenScopes, err := m.DB.GetTokenScopes(ctxscopes, m.DBConn, token.ID)
						spanScopes.End()
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
								ctx, setScopesSpan := tracerhelper.GetTracer().Start(r.Context(), "middleware.gandalf.cache.set_scopes")
								m.Cache.Client.Set(ctx, "token_to_scopes_"+givenToken, buf.Bytes(), time.Hour*1)
								setScopesSpan.End()
							}

						case pgx.ErrNoRows:
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
					ctx, getUserCacheSpan := tracerhelper.GetTracer().Start(r.Context(), "middleware.gandalf.cache.get_username")
					cached_username, err := m.Cache.Client.Get(ctx, cacheKey).Result()
					getUserCacheSpan.End()
					if err != nil {
						gandalf_token_cache.WithLabelValues("user", "miss").Inc()
						// user not in cache, get from db
						ctx, span := tracerhelper.GetTracer().Start(r.Context(), "middleware.gandalf.db.get_user")
						usr, err := m.DB.GetUserByID(ctx, m.DBConn, token.UserID)
						span.End()
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
						dbCtx, dbUserSpan := tracerhelper.GetTracer().Start(r.Context(), "middleware.gandalf.cache.set_username")
						m.Cache.Client.Set(dbCtx, cacheKey, username, time.Minute*10)
						dbUserSpan.End()
					} else {
						gandalf_token_cache.WithLabelValues("user", "hit").Inc()
						username = cached_username
					}

					r = contexthelper.SetUsernameInContext(r, username)
					authResult = "passed_with_bearer"
				}
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
