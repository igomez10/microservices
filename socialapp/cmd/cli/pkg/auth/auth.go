package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// Environment represents a CLI environment configuration
type Environment struct {
	Host          string
	TokenEndpoint string
	Scheme        string
}

// Environments maps environment names to their configurations
var Environments = map[string]Environment{
	"live": {
		Host:          "socialapp.gomezignacio.com",
		TokenEndpoint: "https://socialapp.gomezignacio.com/v1/oauth/token",
		Scheme:        "https",
	},
	"local": {
		Host:          "localhost:8086",
		TokenEndpoint: "http://localhost:8086/v1/oauth/token",
		Scheme:        "http",
	},
}

// DefaultEnv is the default environment when none is specified
const DefaultEnv = "live"

// Environment variable names
const (
	EnvVarUsername    = "SOCIALAPP_CLI_USERNAME"
	EnvVarPassword    = "SOCIALAPP_CLI_PASSWORD"
	EnvVarEnvironment = "SOCIALAPP_CLI_ENV"
)

// Token cache settings
const (
	TokenCacheDir      = ".socialapp_cli"
	TokenCacheFile     = "token"
	TokenFreshDuration = 5 * time.Minute
)

// CachedToken represents a cached OAuth token with metadata
type CachedToken struct {
	AccessToken string    `json:"access_token"`
	TokenType   string    `json:"token_type"`
	ExpiresIn   int64     `json:"expires_in"`
	CreatedAt   time.Time `json:"created_at"`
	Environment string    `json:"environment"`
	Username    string    `json:"username"`
}

// GetEnvironment returns the environment configuration for the given name
func GetEnvironment(envName string) (Environment, error) {
	env, ok := Environments[envName]
	if !ok {
		return Environment{}, fmt.Errorf("unknown environment: %s (valid: live, local)", envName)
	}
	return env, nil
}

// ResolveEnvironment determines which environment to use based on flag and env var
func ResolveEnvironment(flagValue string) string {
	if flagValue != "" {
		return flagValue
	}
	if envVar := os.Getenv(EnvVarEnvironment); envVar != "" {
		return envVar
	}
	return DefaultEnv
}

// ResolveCredentials returns username and password from flags or environment variables
func ResolveCredentials(flagUsername, flagPassword string) (string, string, error) {
	username := flagUsername
	if username == "" {
		username = os.Getenv(EnvVarUsername)
	}

	password := flagPassword
	if password == "" {
		password = os.Getenv(EnvVarPassword)
	}

	if username == "" {
		return "", "", fmt.Errorf("username is required (use --username flag or %s env var)", EnvVarUsername)
	}
	if password == "" {
		return "", "", fmt.Errorf("password is required (use --password flag or %s env var)", EnvVarPassword)
	}

	return username, password, nil
}

// getTokenCachePath returns the full path to the token cache file
func getTokenCachePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(homeDir, TokenCacheDir, TokenCacheFile), nil
}

// loadCachedToken loads the cached token from disk
func loadCachedToken() (*CachedToken, error) {
	path, err := getTokenCachePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // No cached token
		}
		return nil, fmt.Errorf("failed to read token cache: %w", err)
	}

	var token CachedToken
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, fmt.Errorf("failed to parse token cache: %w", err)
	}

	return &token, nil
}

// saveCachedToken saves the token to disk
func saveCachedToken(token *CachedToken) error {
	path, err := getTokenCachePath()
	if err != nil {
		return err
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create token cache directory: %w", err)
	}

	data, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal token: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write token cache: %w", err)
	}

	return nil
}

// isCachedTokenValid checks if the cached token is still valid
func isCachedTokenValid(cached *CachedToken, envName, username string) bool {
	if cached == nil {
		return false
	}

	// Check if environment matches
	if cached.Environment != envName {
		return false
	}

	// Check if username matches
	if cached.Username != username {
		return false
	}

	// Check if token is still fresh (less than 5 minutes old)
	age := time.Since(cached.CreatedAt)
	if age >= TokenFreshDuration {
		return false
	}

	return true
}

// cachedTokenSource wraps an oauth2.Token from cache
type cachedTokenSource struct {
	token *oauth2.Token
}

func (c *cachedTokenSource) Token() (*oauth2.Token, error) {
	return c.token, nil
}

func requestToken(ctx context.Context, tokenEndpoint, username, password string, scopes []string) (*oauth2.Token, error) {
	oauthConfig := clientcredentials.Config{
		ClientID:     username,
		ClientSecret: password,
		TokenURL:     tokenEndpoint,
		Scopes:       scopes,
	}

	tokenSource := oauthConfig.TokenSource(ctx)
	token, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to obtain token: %w", err)
	}

	return token, nil
}

func cacheToken(token *oauth2.Token, environment, username string) error {
	expiresIn := int64(0)
	if !token.Expiry.IsZero() {
		expiresIn = int64(time.Until(token.Expiry).Seconds())
	}

	cachedToken := &CachedToken{
		AccessToken: token.AccessToken,
		TokenType:   token.TokenType,
		ExpiresIn:   expiresIn,
		CreatedAt:   time.Now(),
		Environment: environment,
		Username:    username,
	}

	return saveCachedToken(cachedToken)
}

// GetFreshTokenWithTokenEndpoint fetches a new OAuth token from a specific token endpoint and updates the local cache.
func GetFreshTokenWithTokenEndpoint(ctx context.Context, cacheEnvironment, tokenEndpoint, username, password string, scopes []string) (*oauth2.Token, error) {
	token, err := requestToken(ctx, tokenEndpoint, username, password, scopes)
	if err != nil {
		return nil, err
	}

	if err := cacheToken(token, cacheEnvironment, username); err != nil {
		// Log warning but continue - caching is not critical
		fmt.Fprintf(os.Stderr, "Warning: failed to cache token: %v\n", err)
	}

	return token, nil
}

// GetFreshToken fetches a new OAuth token from the configured token endpoint and updates the local cache.
func GetFreshToken(ctx context.Context, envName, username, password string, scopes []string) (*oauth2.Token, error) {
	env, err := GetEnvironment(envName)
	if err != nil {
		return nil, err
	}

	return GetFreshTokenWithTokenEndpoint(ctx, envName, env.TokenEndpoint, username, password, scopes)
}

// GetHTTPClient returns an authenticated HTTP client with token caching
func GetHTTPClient(ctx context.Context, envName, username, password string, scopes []string) (*http.Client, error) {
	// Check for valid cached token
	cached, err := loadCachedToken()
	if err != nil {
		// Log warning but continue - we'll just get a new token
		fmt.Fprintf(os.Stderr, "Warning: failed to load cached token: %v\n", err)
	}

	if isCachedTokenValid(cached, envName, username) {
		// Use cached token
		token := &oauth2.Token{
			AccessToken: cached.AccessToken,
			TokenType:   cached.TokenType,
		}
		if cached.ExpiresIn > 0 {
			token.Expiry = cached.CreatedAt.Add(time.Duration(cached.ExpiresIn) * time.Second)
		}
		return oauth2.NewClient(ctx, &cachedTokenSource{token: token}), nil
	}

	token, err := GetFreshToken(ctx, envName, username, password, scopes)
	if err != nil {
		return nil, err
	}

	return oauth2.NewClient(ctx, oauth2.StaticTokenSource(token)), nil
}

// GetAPIClientConfig returns host and scheme for the given environment
func GetAPIClientConfig(envName string) (host, scheme string, err error) {
	env, err := GetEnvironment(envName)
	if err != nil {
		return "", "", err
	}
	return env.Host, env.Scheme, nil
}

// ParseHostURL parses a host URL and returns the host and scheme
func ParseHostURL(hostURL string) (host, scheme string, err error) {
	parsed, err := url.Parse(hostURL)
	if err != nil {
		return "", "", fmt.Errorf("error parsing host URL: %w", err)
	}
	return parsed.Host, parsed.Scheme, nil
}
