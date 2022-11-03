package authorization

import (
	"fmt"
	"net/http"
	"socialapp/internal/contexthelper"
)

type Middleware struct {
	RequiredScopes map[string]bool
}

// Authorize checks if the user has the required scopes
// The scopes are expected in the context of the request
func (m *Middleware) Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := contexthelper.GetLoggerInContext(r.Context())
		// get scopes from context
		tokenScopes, ok := contexthelper.GetRequestedScopesInContext(r.Context())
		if !ok {
			log.Error().
				Msg("Failed to get token scopes from context")

			w.Write([]byte(`{"code":403,"message":"No scopes in context"}`))
			w.WriteHeader(http.StatusForbidden)
			return
		}

		// check if all required scopes are in token
		for scopeName := range m.RequiredScopes {
			if exist := tokenScopes[scopeName]; !exist {
				log.Error().
					Str("scope", scopeName).
					Str("tokenScopes", fmt.Sprintf("%v", tokenScopes)).
					Msg("Missing scope")

				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte(fmt.Sprintf(`{"code": 403, "message": "Scope %s missing from token"}`, scopeName)))
				return
			}
		}

		log.Info().
			Msg("Authorization successful")
		next.ServeHTTP(w, r)
	})
}
