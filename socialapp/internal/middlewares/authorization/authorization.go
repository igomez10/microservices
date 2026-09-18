package authorization

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/igomez10/microservices/socialapp/internal/contexthelper"
	"github.com/igomez10/microservices/socialapp/internal/tracerhelper"
	"github.com/igomez10/microservices/socialapp/socialappapi/openapi"
)

var (
	noScopesInContextError = openapi.Error{
		Code:    http.StatusForbidden,
		Message: "No scopes in context",
	}
	unauthenticatedError = openapi.Error{
		Code:    http.StatusUnauthorized,
		Message: "No scopes in context and required scopes are not empty",
	}
)

type Middleware struct {
	RequiredScopes map[string]bool
}

// Authorize checks if the user has the required scopes
// The scopes are expected in the context of the request
func (m *Middleware) Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracerhelper.GetTracer().Start(r.Context(), "middleware.authorization")
		defer span.End()

		r = r.WithContext(ctx)

		logger := contexthelper.GetLoggerInContext(r.Context())
		// get scopes from context
		tokenScopes, ok := contexthelper.GetRequestedScopesInContext(r.Context())
		if !ok {
			logger.Error("Failed to get token scopes from context")

			writeErrorResponse(w, logger, noScopesInContextError)
			return
		}
		if len(tokenScopes) == 0 && len(m.RequiredScopes) != 0 {
			logger.Info(
				"Unauthenticated request to protected endpoint",
				"token_scopes", tokenScopes,
				"required_scopes", m.RequiredScopes,
			)
			writeErrorResponse(w, logger, unauthenticatedError)
			return
		}

		// check if all required scopes are in token
		for scopeName := range m.RequiredScopes {
			if exist := tokenScopes[scopeName]; !exist {
				logger.Info("Missing scope", "scope", scopeName, "token_scopes", tokenScopes)

				writeErrorResponse(w, logger, openapi.Error{
					Code:    http.StatusForbidden,
					Message: fmt.Sprintf("Scope %s missing from token", scopeName),
				})
				return
			}
		}

		logger.Info("Authorization successful")
		next.ServeHTTP(w, r)
	})
}

func writeErrorResponse(w http.ResponseWriter, logger *slog.Logger, response openapi.Error) {
	status := int(response.Code)
	if err := openapi.EncodeJSONResponse(response, &status, nil, w); err != nil {
		logger.Error("Failed to encode authorization error response", "error", err)
	}
}
