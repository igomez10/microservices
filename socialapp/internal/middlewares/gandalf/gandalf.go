package gandalf

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"socialapp/internal/authorizationparser"
	"socialapp/internal/contexthelper"
	"socialapp/pkg/controller/user"
	"socialapp/pkg/db"
	"strings"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var gandalf_token_cache_hits = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "gandalf_token_cache_hits",
	Help: "The total number of token cache hits",
}, []string{"cache", "status"})

var gandalf_scope_cache_hits = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "gandalf_scope_cache_hits",
	Help: "The total number of scope cache hits",
}, []string{"cache", "status"})
var gandalf_token_cache_misses = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "gandalf_token_cache_misses",
	Help: "The total number of token cache misses",
}, []string{"cache", "status"})
var gandalf_scope_cache_misses = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "gandalf_scope_cache_misses",
	Help: "The total number of scope cache misses",
}, []string{"cache", "status"})

type Middleware struct {
	DB                  db.Querier
	DBConn              *sql.DB
	Authorizationparser authorizationparser.EndpointAuthorizations
}

func (m *Middleware) Authenticate(next http.Handler) http.Handler {

	allowlistedPaths := map[string]map[string]bool{
		"/users": {
			"POST": true,
		},
		"/metrics": {
			"GET": true,
		},
		"/apispec": {
			"GET": true,
		},
	}

	tokenCache, err := bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
	if err != nil {
		panic(err)
	}

	scopesCache, err := bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
	if err != nil {
		panic(err)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := contexthelper.GetLoggerInContext(r.Context())

		// get token from header
		if allowlistedPaths[r.URL.Path] != nil && allowlistedPaths[r.URL.Path][r.Method] {
			r = contexthelper.SetRequestedScopesInContext(r, map[string]bool{})
			log.Info().
				Str("path", r.URL.Path).
				Str("method", r.Method).
				Str("middleware", "gandalf").
				Msg("Allowlisted path")

			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			givenToken := strings.TrimPrefix(authHeader, "Bearer ")
			var token db.Token
			// check token in cache
			cachedTokenBytes, err := tokenCache.Get(givenToken)
			if err != nil {
				gandalf_token_cache_misses.WithLabelValues("token", "miss").Inc()
				// token not in cache, get from db
				dbtoken, err := m.DB.GetToken(r.Context(), m.DBConn, givenToken)
				switch err {
				case nil:
					token = dbtoken
					// update cache
					tokenBytes, err := json.Marshal(token)
					if err != nil {
						log.Error().
							Err(err).
							Str("middleware", "gandalf").
							Msg("Failed to marshal token")
					} else {
						tokenCache.Set(givenToken, tokenBytes)
					}
				case sql.ErrNoRows:
					log.Error().
						Err(err).
						Msg("Token not found")
					http.Error(w, "Token not found", http.StatusUnauthorized)
					w.Write([]byte(`{"code": 401, "message": "Invalid bearer token"}`))
				default:
					log.Error().
						Err(err).
						Msg("Error while getting token")
					http.Error(w, "Error while getting token", http.StatusInternalServerError)
					w.Write([]byte(`{"code": 500, "message": "Error while getting token"}`))
				}
			} else {
				// token found in cache
				gandalf_token_cache_hits.WithLabelValues("token", "hit").Inc()
				if err := json.Unmarshal(cachedTokenBytes, &token); err != nil {
					log.Error().
						Err(err).
						Str("middleware", "gandalf").
						Msg("Failed to unmarshal token")
					w.Write([]byte(`{"code": 500, "message": "Error while getting token"}`))
					w.WriteHeader(http.StatusInternalServerError)
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
			cachedScopesBytes, err := scopesCache.Get(givenToken)
			if err != nil {
				gandalf_scope_cache_misses.WithLabelValues("scopes", "miss").Inc()
				// scopes not in cache, get from db
				dbTokenScopes, err := m.DB.GetTokenScopes(r.Context(), m.DBConn, token.ID)
				switch err {
				case nil:
					scopes = dbTokenScopes
					// update cache
					scopesBytes, err := json.Marshal(scopes)
					if err != nil {
						log.Error().
							Err(err).
							Str("middleware", "gandalf").
							Msg("Failed to marshal scopes")
					} else {
						scopesCache.Set(givenToken, scopesBytes)
					}

				case sql.ErrNoRows:
					log.Error().
						Err(err).
						Msg("Token scopes not found")
					http.Error(w, "Token scopes not found", http.StatusUnauthorized)
					w.Write([]byte(`{"code": 401, "message": "Invalid bearer token"}`))
				default:
					log.Error().Err(err).Msg("Error while getting token scopes")
					http.Error(w, "Error while getting token scopes", http.StatusInternalServerError)
					w.Write([]byte(`{"code": 500, "message": "Error while getting token scopes"}`))
				}
			} else {
				gandalf_scope_cache_hits.WithLabelValues("scopes", "hit").Inc()
				// scopes found in cache
				if err := json.Unmarshal(cachedScopesBytes, &scopes); err != nil {
					log.Error().
						Err(err).
						Str("middleware", "gandalf").
						Msg("Failed to unmarshal token scopes")
					w.Write([]byte(`{"code": 500, "message": "Error while getting token scopes"}`))
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}

			scopesMap := map[string]bool{}
			for i := range scopes {
				scopesMap[scopes[i].Name] = true
			}

			r = contexthelper.SetRequestedScopesInContext(r, scopesMap)

			usr, err := m.DB.GetUserByID(r.Context(), m.DBConn, token.UserID)
			if err != nil {
				log.Error().
					Err(err).
					Int64("userID", token.UserID).
					Msg("Failed to get user from token")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"code": 500, "message": "Failed to get user from token"}`))
				return
			}

			r = contexthelper.SetUsernameInContext(r, usr.Username)
			next.ServeHTTP(w, r)
			return
		} else if strings.HasPrefix(authHeader, "Basic ") && r.URL.Path == "/oauth/token" {
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
			default:
				log.Error().
					Err(err).
					Msg("Error while getting user")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"code": 500, "message": "Error while getting user"}`))
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
			requestedScopes := strings.Split(r.FormValue("scope"), " ")
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
			next.ServeHTTP(w, r)
			return
		}

		// no token in header
		log.Error().
			Msg("No token in header")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"code": 401, "message": "No token was provided"}`))
	})

}
