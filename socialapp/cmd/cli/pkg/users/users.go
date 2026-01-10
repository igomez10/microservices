package users

import (
	"context"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/igomez10/microservices/socialapp/client"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/oauth2"
)

type ConfigurationOpts func(*client.Configuration)

func WithDefaultHeader(key, value string) ConfigurationOpts {
	return func(cfg *client.Configuration) {
		cfg.AddDefaultHeader(key, value)
	}
}

func WithSkipCache() ConfigurationOpts {
	return func(cfg *client.Configuration) {
		WithDefaultHeader("Cache-Control", "no-store")(cfg)
	}
}

func NewDefaultConfiguration(opts ...ConfigurationOpts) *client.Configuration {
	cfg := client.NewConfiguration()

	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

func GetApiClient(ctx context.Context, url *url.URL) (context.Context, *client.APIClient) {
	configuration := NewDefaultConfiguration(WithSkipCache())
	httpClient := http.DefaultClient
	configuration.HTTPClient = getHTTPClient()
	configuration.Host = url.Host
	configuration.Scheme = url.Scheme

	proxyCtx := context.WithValue(ctx, oauth2.HTTPClient, httpClient)
	apiClient := client.NewAPIClient(configuration)
	return proxyCtx, apiClient
}

func getHTTPClient() *http.Client {
	// setup retryable http client
	retryClient := retryablehttp.NewClient()
	retryClient.Logger = nil
	retryClient.RequestLogHook = func(_ retryablehttp.Logger, req *http.Request, attempt int) {
		if attempt >= 1 {
			slog.Warn("http retry", "method", req.Method, "url", req.URL.String(), "attempt", attempt)
		}
	}

	retryClient.HTTPClient.Transport = otelhttp.NewTransport(http.DefaultTransport)
	retryClient.RetryMax = 10
	retryClient.HTTPClient.Timeout = 15 * time.Second
	retryClient.Backoff = retryablehttp.LinearJitterBackoff
	retryClient.CheckRetry = func(ctx context.Context, resp *http.Response, err error) (bool, error) {
		// Retry on network errors or 5xx status codes
		if err != nil {
			slog.Warn(
				"http retry",
				"error", err,
				"url", resp.Request.URL.String(),
				"method", resp.Request.Method,
				"status", resp.Status,
				"status_code", resp.StatusCode,
			)
			return true, err
		}

		if resp.StatusCode >= http.StatusInternalServerError {
			return true, nil
		}

		return false, nil
	}
	http.DefaultClient = retryClient.StandardClient()
	return http.DefaultClient
}
