package contexthelper

import (
	"context"
	"net/http"
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
