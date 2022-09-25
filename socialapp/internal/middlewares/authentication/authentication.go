package authentication

import (
	"context"
	"database/sql"
	"net/http"
	"socialapp/pkg/controller/user"
	"socialapp/pkg/db"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

type Middleware struct {
	DB     db.Querier
	DBConn *sql.DB
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

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get token from header
		if allowlistedPaths[r.URL.Path] != nil && allowlistedPaths[r.URL.Path][r.Method] {
			log.Info().Msg("Allowlisted path")
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			// check givenToken in DB
			givenToken := strings.TrimPrefix(authHeader, "Bearer ")
			token, err := m.DB.GetToken(r.Context(), m.DBConn, givenToken)
			if err != nil {
				log.Error().Err(err).Msg("Failed to get token")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			if time.Now().After(token.ValidUntil) {
				log.Error().Err(err).Msg("Token expired")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			usr, err := m.DB.GetUserByID(r.Context(), m.DBConn, token.UserID)
			if err != nil {
				log.Error().Err(err).Msg("Failed to get user")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), "username", usr.Username))
			next.ServeHTTP(w, r)
			return
		} else if strings.HasPrefix(authHeader, "Basic ") && r.URL.Path == "/oauth/token" {

			// check grant type is client_credentials
			if r.FormValue("grant_type") != "client_credentials" {
				log.Error().Msg("Grant type is not client_credentials")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Grant type is not client_credentials"))
				return
			}

			username, password, ok := r.BasicAuth()
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"message": "invalid basic auth format"}`))
				return
			}

			usr, err := m.DB.GetUserByUsername(r.Context(), m.DBConn, username)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"message": "invalid username or password"}`))
				return
			}

			encryptedPassword := user.EncryptPassword(password, usr.Salt)
			if encryptedPassword != usr.HashedPassword {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"message": "invalid username or password"}`))
				return
			}
			// passed authentication
			r = r.WithContext(context.WithValue(r.Context(), "username", username))
			next.ServeHTTP(w, r)
			return
		}

		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
	})
}

func (m *Middleware) GetTokenFromDB(ctx context.Context, token string) (db.Token, error) {
	dbToken, err := m.DB.GetToken(ctx, m.DBConn, token)
	if err != nil {
		return db.Token{}, err
	}

	return dbToken, nil
}
