package users

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/igomez10/microservices/socialapp/client"
	"github.com/rs/zerolog/log"
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
			log.Warn().
				Str("method", req.Method).
				Str("url", req.URL.String()).
				Int("attempt", attempt).
				Msgf("http retry")
		}
	}

	retryClient.HTTPClient.Transport = otelhttp.NewTransport(http.DefaultTransport)
	retryClient.RetryMax = 10
	retryClient.HTTPClient.Timeout = 15 * time.Second
	retryClient.Backoff = retryablehttp.LinearJitterBackoff
	retryClient.CheckRetry = func(ctx context.Context, resp *http.Response, err error) (bool, error) {
		// Retry on network errors or 5xx status codes
		if err != nil {
			log.Warn().
				Err(err).
				Stringer("url", resp.Request.URL).
				Str("method", resp.Request.Method).
				Str("status", resp.Status).
				Int("status_code", resp.StatusCode).
				Msg("http retry")
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
