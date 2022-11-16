package contexthelper

import (
	"context"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func GetUsernameInContext(ctx context.Context) (string, bool) {
	username, ok := ctx.Value("username").(string)
	return username, ok
}

func SetUsernameInContext(r *http.Request, username string) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), "username", username))
}

func GetRequestedScopesInContext(ctx context.Context) (map[string]bool, bool) {
	scopes, ok := ctx.Value("scopes").(map[string]bool)
	return scopes, ok
}

func SetRequestedScopesInContext(r *http.Request, scopes map[string]bool) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), "scopes", scopes))
}

func GetRequestIDInContext(ctx context.Context) string {
	requestID, ok := ctx.Value("X-Request-ID").(string)
	if !ok {
		log.Error().Msg("failed to retrieve request ID from context")
		defaultRequestID := "Request ID not found in context"
		return defaultRequestID
	}
	return requestID
}

func SetRequestIDInContext(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, "X-Request-ID", requestID)
}

func GetLoggerInContext(ctx context.Context) zerolog.Logger {
	logger, ok := ctx.Value("logger").(zerolog.Logger)
	if !ok {
		return log.Logger
	}
	return logger
}

func SetLoggerInContext(ctx context.Context, logger zerolog.Logger) context.Context {
	return context.WithValue(ctx, "logger", logger)
}

func GetRequestPatternInContext(ctx context.Context) string {
	pattern, ok := ctx.Value("pattern").(string)
	if !ok {
		log.Error().Msg("failed to retrieve pattern from context")
		defaultPattern := "Pattern not found in context"
		return defaultPattern
	}
	return pattern
}

func SetRequestPatternInContext(ctx context.Context, pattern string) context.Context {
	return context.WithValue(ctx, "pattern", pattern)
}
