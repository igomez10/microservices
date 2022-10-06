package contexthelper

import (
	"context"
	"net/http"

	"github.com/rs/zerolog/log"
)

func GetUsernameInContext(ctx context.Context) (string, bool) {
	username, ok := ctx.Value("username").(string)
	return username, ok
}

func SetUsernameInContext(r *http.Request, username string) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), "username", username))
}

func SetRequestedScopesInContext(r *http.Request, scopes map[string]bool) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), "scopes", scopes))
}

func GetRequestedScopesInContext(ctx context.Context) (map[string]bool, bool) {
	scopes, ok := ctx.Value("scopes").(map[string]bool)
	return scopes, ok
}

func GetRequestIDInContext(ctx context.Context) string {
	requestID, ok := ctx.Value("X-Request-ID").(string)
	if !ok {
		log.Error().Msg("failed to retrieve request ID from context")
		return "Request ID not found in context"
	}
	return requestID
}
