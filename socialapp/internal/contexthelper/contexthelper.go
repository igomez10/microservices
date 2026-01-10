package contexthelper

import (
	"context"
	"log/slog"
	"net/http"
)

type contextKey string

const (
	usernameKey       contextKey = "username"
	scopesKey         contextKey = "scopes"
	requestIDKey      contextKey = "request-id"
	loggerKey         contextKey = "logger"
	requestPatternKey contextKey = "pattern"
)

func GetUsernameInContext(ctx context.Context) (string, bool) {
	username, ok := ctx.Value(usernameKey).(string)
	return username, ok
}

func SetUsernameInContext(r *http.Request, username string) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), usernameKey, username))
}

func GetRequestedScopesInContext(ctx context.Context) (map[string]bool, bool) {
	scopes, ok := ctx.Value(scopesKey).(map[string]bool)
	return scopes, ok
}

func SetRequestedScopesInContext(r *http.Request, scopes map[string]bool) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), scopesKey, scopes))
}

func GetRequestIDInContext(ctx context.Context) string {
	requestID, ok := ctx.Value(requestIDKey).(string)
	if !ok {
		slog.Error("failed to retrieve request ID from context")
		return "Request ID not found in context"
	}
	return requestID
}

func SetRequestIDInContext(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

func GetLoggerInContext(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(loggerKey).(*slog.Logger)
	if !ok || logger == nil {
		return slog.Default()
	}
	return logger
}

func SetLoggerInContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func GetRequestPatternInContext(ctx context.Context) string {
	pattern, ok := ctx.Value(requestPatternKey).(string)
	if !ok {
		slog.Error("failed to retrieve pattern from context")
		return "Pattern not found in context"
	}
	return pattern
}

func SetRequestPatternInContext(ctx context.Context, pattern string) context.Context {
	return context.WithValue(ctx, requestPatternKey, pattern)
}
