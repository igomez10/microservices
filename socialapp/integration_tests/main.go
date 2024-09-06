package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/igomez10/microservices/socialapp/client"
	"github.com/igomez10/microservices/socialapp/internal/tracerhelper"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

var (
	RENDER_SERVER_URL          = 0
	LOCALHOST_SERVER_URL       = 1
	LOCALHOST_DEBUG_SERVER_URL = 2

	CONTEXT_SERVER       int
	apiClient            *client.APIClient
	ENDPOINT_OAUTH_TOKEN string
	urlAgent             = flag.String("agentURL", os.Getenv("AGENT_URL"), "Agent URL \"http://localhost:4317\"")
)

const appName = "integration-tests-socialapp"
const defaultUsername = "Test-%d1"
const defaultPassword = "Password-%d1"

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

// NewDefaultConfiguration creates a new Configuration with default values, usually we want to set the same default values for all tests
func NewDefaultConfiguration(opts ...ConfigurationOpts) *client.Configuration {
	cfg := client.NewConfiguration()

	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

// add setup function
func Setup() {
	//  set the endpoint for the oauth token
	testSetup := os.Getenv("TEST_SETUP")
	// if testSetup == "" {
	// 	testSetup = "LOCALHOST_DEBUG"
	// }

	switch testSetup {
	case "LOCALHOST":
		CONTEXT_SERVER = LOCALHOST_SERVER_URL
		ENDPOINT_OAUTH_TOKEN = "http://localhost:8085/v1/oauth/token"
	case "LOCALHOST_DEBUG":
		CONTEXT_SERVER = LOCALHOST_DEBUG_SERVER_URL
		ENDPOINT_OAUTH_TOKEN = "http://localhost:8085/v1/oauth/token"
	default:
		CONTEXT_SERVER = RENDER_SERVER_URL
		ENDPOINT_OAUTH_TOKEN = "https://socialapp.gomezignacio.com/v1/oauth/token"
	}

	// with timestamp and caller
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339, NoColor: false}).
		With().Caller().Timestamp().Logger()
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

	retryClient.HTTPClient.Transport = newrelic.NewRoundTripper(otelhttp.NewTransport(http.DefaultTransport))
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

		if resp.StatusCode >= 500 {
			return true, nil
		}

		return false, nil
	}
	http.DefaultClient = retryClient.StandardClient()
	return http.DefaultClient
}

func getOuath2Context(initialContext context.Context, config clientcredentials.Config) (context.Context, error) {
	tokenSource := config.TokenSource(initialContext)
	initialContext = context.WithValue(initialContext, client.ContextOAuth2, tokenSource)

	return initialContext, nil
}

func main() {

	flag.Parse()
	log.Info().Msg("Starting integration tests")
	ctx := context.Background()
	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure(), otlptracegrpc.WithEndpointURL(*urlAgent))
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to create otlp exporter for tracing %q", *urlAgent)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("integration-tests"),
		// Add more attributes as needed
		)),
	)

	// Register the tracer provider as the global provider.
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	// set root trace
	ctx, span := tracerhelper.GetTracer().Start(ctx, "")
	defer span.End()

	if err := ListUsersLifecycle(ctx); err != nil {
		log.Error().Err(err).Msg("error ListUsersLifecycle")
	}
	if err := CreateUserLifecycle(ctx); err != nil {
		log.Error().Err(err).Msg("error CreateUserLifecycle")
	}
	if err := FollowLifeCycle(ctx); err != nil {
		log.Error().Err(err).Msg("error FollowLifeCycle")
	}
	if err := GetExpectedFeed(ctx); err != nil {
		log.Error().Err(err).Msg("error GetExpectedFeed")
	}
	if err := GetAccessToken(ctx); err != nil {
		log.Error().Err(err).Msg("error GetAccessToken")
	}
	if err := RegisterUserFlow(ctx); err != nil {
		log.Error().Err(err).Msg("error RegisterUserFlow")
	}
	if err := ChangePassword(ctx); err != nil {
		log.Error().Err(err).Msg("error ChangePassword")
	}
	if err := RoleLifecycle(ctx); err != nil {
		log.Error().Err(err).Msg("error RoleLifecycle")
	}
	if err := ScopeLifecycle(ctx); err != nil {
		log.Error().Err(err).Msg("error ScopeLifecycle")
	}
	if err := UserRoleLifeCycle(ctx); err != nil {
		log.Error().Err(err).Msg("error UserRoleLifeCycle")
	}
	if err := CacheRequestSameUser(ctx); err != nil {
		log.Error().Err(err).Msg("error CacheRequestSameUser")
	}
	if err := URLLifeCycle(ctx); err != nil {
		log.Error().Err(err).Msg("error URLLifeCycle")
	}
}

func ListUsersLifecycle(ctx context.Context) error {
	Setup()

	ctx, span := tracerhelper.GetTracer().Start(ctx, "ListUsersLifecycle")
	defer span.End()

	configuration := NewDefaultConfiguration(WithSkipCache())
	httpClient := getHTTPClient()
	configuration.HTTPClient = httpClient

	proxyCtx := context.WithValue(context.Background(), oauth2.HTTPClient, httpClient)
	proxyCtx = context.WithValue(proxyCtx, client.ContextServerIndex, CONTEXT_SERVER)

	username1 := fmt.Sprintf(defaultUsername, time.Now().UnixNano())
	password := fmt.Sprintf(defaultPassword, time.Now().UnixNano())
	apiClient = client.NewAPIClient(configuration)

	if err := func() error {
		createUsrReq := client.NewCreateUserRequest(username1, password, "FirstName_example", "LastName_example", username1)
		_, _, err := apiClient.UserAPI.CreateUser(proxyCtx).CreateUserRequest(*createUsrReq).Execute()
		if err != nil {
			return fmt.Errorf("Error creating user: %v", err)
		}
		return nil
	}(); err != nil {
		return err
	}

	conf := clientcredentials.Config{
		ClientID:     username1,
		ClientSecret: password,
		Scopes:       []string{"socialapp.users.list"},
		TokenURL:     ENDPOINT_OAUTH_TOKEN,
	}
	oauth2Ctx, err := getOuath2Context(proxyCtx, conf)
	if err != nil {
		return fmt.Errorf("Error getting oauth2 context: %v", err)
	}
	openAPICtx := context.WithValue(oauth2Ctx, client.ContextServerIndex, CONTEXT_SERVER)

	// List users
	_, r, err := apiClient.UserAPI.ListUsers(openAPICtx).Limit(10).Offset(0).Execute()
	if err != nil {
		return fmt.Errorf("Error when calling `UserAPI.ListUsers`: %v\n %+v\n", err, r)
	}
	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", r.StatusCode)
	}

	return nil
}

func CreateUserLifecycle(ctx context.Context) error {
	Setup()
	ctx, span := tracerhelper.GetTracer().Start(ctx, "CreateUserLifecycle")
	defer span.End()

	configuration := NewDefaultConfiguration(WithSkipCache())
	httpClient := getHTTPClient()
	configuration.HTTPClient = httpClient
	noAuthCtx := context.WithValue(context.Background(), oauth2.HTTPClient, httpClient)
	noAuthCtx = context.WithValue(noAuthCtx, client.ContextServerIndex, CONTEXT_SERVER)

	apiClient = client.NewAPIClient(configuration)

	username := fmt.Sprintf("Test-%d", time.Now().UnixNano())
	password := "password"
	email := fmt.Sprintf("Test-%d-@social.com", time.Now().UnixNano())
	user := *client.NewCreateUserRequest(username, "password", "FirstName_example", "LastName_example", email) // User | Create a new user

	if err := func() error {
		_, r, err := apiClient.UserAPI.CreateUser(noAuthCtx).
			CreateUserRequest(user).
			Execute()
		if err != nil {
			return fmt.Errorf("Error when calling `UserAPI.CreateUser`: %v\n %+v\n", err, r)
		}
		if r.StatusCode != http.StatusOK {
			return fmt.Errorf("Expected status code %d, got %d", http.StatusOK, r.StatusCode)
		}
		return nil
	}(); err != nil {
		return err
	}

	if err := func() error {
		conf := clientcredentials.Config{
			ClientID:     username,
			ClientSecret: password,
			Scopes:       []string{"socialapp.users.read"},
			TokenURL:     ENDPOINT_OAUTH_TOKEN,
		}
		oauth2Ctx, err := getOuath2Context(noAuthCtx, conf)
		if err != nil {
			return fmt.Errorf("Error getting oauth2 context: %v", err)
		}

		resp, r, err := apiClient.UserAPI.GetUserByUsername(oauth2Ctx, username).Execute()
		if err != nil {
			return fmt.Errorf("Error when calling `UserAPI.GetUserByUsername`: %v\n %+v\n", err, r)
		}
		if r.StatusCode != http.StatusOK {
			return fmt.Errorf("Expected status code %d, got %d", http.StatusOK, r.StatusCode)
		}
		if resp.Username != user.Username {
			return fmt.Errorf("Expected username %s, got %s", user.Username, resp.Username)
		}
		if resp.Email != user.Email {
			return fmt.Errorf("Expected email %s, got %s", user.Email, resp.Email)
		}
		if resp.FirstName != user.FirstName {
			return fmt.Errorf("Expected first name %q, got %q", user.FirstName, resp.FirstName)
		}
		if resp.LastName != user.LastName {
			return fmt.Errorf("Expected last name %q, got %q", user.LastName, resp.LastName)
		}
		return nil
	}(); err != nil {
		return err
	}

	// update user
	if err := func() error {
		conf := clientcredentials.Config{
			ClientID:     username,
			ClientSecret: password,
			Scopes:       []string{"socialapp.users.update"},
			TokenURL:     ENDPOINT_OAUTH_TOKEN,
		}
		oauth2Ctx, err := getOuath2Context(noAuthCtx, conf)
		if err != nil {
			return fmt.Errorf("Error getting oauth2 context: %v", err)
		}

		updatedFirstName := "UpdatedFirstName" + uuid.NewString()
		updatedLastName := "UpdatedLastName" + uuid.NewString()
		updatedEmail := "UpdatedEmail" + uuid.NewString() + "@social.com"
		updatedUser := client.User{
			Username:  username,
			FirstName: updatedFirstName,
			LastName:  updatedLastName,
			Email:     updatedEmail,
		}

		updateUserReq := apiClient.UserAPI.
			UpdateUser(oauth2Ctx, username).
			User(updatedUser)

		uUser, res, err := updateUserReq.Execute()
		if err != nil {
			return fmt.Errorf("Error when calling `UserAPI.UpdateUser`: %v\n %+v\n", err, res)
		}
		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("Expected status code %d, got %d", http.StatusOK, res.StatusCode)
		}
		if uUser.FirstName != updatedFirstName {
			return fmt.Errorf("Expected first name %q, got %q", updatedFirstName, uUser.FirstName)
		}
		if uUser.LastName != updatedLastName {
			return fmt.Errorf("Expected last name %q, got %q", updatedLastName, uUser.LastName)
		}
		if uUser.Email != updatedEmail {
			return fmt.Errorf("Expected email %q, got %q", updatedEmail, uUser.Email)
		}
		return nil
	}(); err != nil {
		return err
	}

	return nil
}

func FollowLifeCycle(ctx context.Context) error {
	Setup()
	ctx, span := tracerhelper.GetTracer().Start(ctx, "FollowLifeCycle")
	defer span.End()
	// create two users
	configuration := NewDefaultConfiguration(WithSkipCache())
	httpClient := getHTTPClient()
	configuration.HTTPClient = httpClient
	proxyCtx := context.WithValue(ctx, oauth2.HTTPClient, httpClient)
	proxyCtx = context.WithValue(proxyCtx, client.ContextServerIndex, CONTEXT_SERVER)

	apiClient = client.NewAPIClient(configuration)

	username1 := fmt.Sprintf(defaultUsername, time.Now().UnixNano())
	password1 := fmt.Sprintf("TestPassword-%d1", time.Now().UnixNano())
	email1 := fmt.Sprintf("Test-%d-1@social.com", time.Now().UnixNano())
	user1 := *client.NewCreateUserRequest(username1, password1, "FirstName_example", "LastName_example", email1) // User | Create a new user

	username2 := fmt.Sprintf("Test-%d2", time.Now().UnixNano())
	password2 := fmt.Sprintf("TestPassword-%d2", time.Now().UnixNano())
	email2 := fmt.Sprintf("Test-%d-2@social.com", time.Now().UnixNano())
	user2 := *client.NewCreateUserRequest(username2, password2, "FirstName_example", "LastName_example", email2) // User | Create a new user

	// create users
	if err := func() error {
		_, r1, err1 := apiClient.UserAPI.CreateUser(proxyCtx).
			CreateUserRequest(user1).
			Execute()
		if err1 != nil {
			return fmt.Errorf("Error when calling `UserAPI.CreateUser`: %v\n %+v\n", err1, r1)
		}

		_, r2, err2 := apiClient.UserAPI.CreateUser(proxyCtx).
			CreateUserRequest(user2).
			Execute()
		if err2 != nil {
			return fmt.Errorf("Error when calling `UserAPI.CreateUser`: %v\n %+v\n", err2, r2)
		}
		return nil
	}(); err != nil {
		return err
	}

	conf := clientcredentials.Config{
		ClientID:     username1,
		ClientSecret: password1,
		Scopes: []string{
			"socialapp.users.read",
			"socialapp.follower.create",
			"socialapp.follower.read",
			"socialapp.follower.delete",
		},
		TokenURL: ENDPOINT_OAUTH_TOKEN,
	}

	oauth2Ctx, err := getOuath2Context(proxyCtx, conf)

	// user 1 follows user 2
	if err := func() error {
		if err != nil {
			return fmt.Errorf("Error getting oauth2 context: %v", err)
		}
		r, err := apiClient.UserAPI.FollowUser(oauth2Ctx, username2, username1).Execute()
		if err != nil {
			return fmt.Errorf("Error when calling `UserAPI.FollowUser`: %v\n %+v\n", err, r)
		}
		return nil
	}(); err != nil {
		return err
	}

	// validate user 1 follows user 2
	if err := func() error {
		followers, r, err := apiClient.UserAPI.GetUserFollowers(oauth2Ctx, username2).Execute()
		if err != nil {
			return fmt.Errorf("Error when calling `UserAPI.FollowUser`: %v\n %+v\n", err, r)
		}
		if len(followers) != 1 {
			return fmt.Errorf("Expected 1 follower, got %d", len(followers))
		}
		if followers[0].Username != username1 {
			return fmt.Errorf("Expected follower %s, got %s", username1, followers[0].Username)
		}
		return nil
	}(); err != nil {
		return err
	}

	// user 1 unfollows user 2
	if err := func() error {
		r, err := apiClient.UserAPI.UnfollowUser(oauth2Ctx, username2, username1).Execute()
		if err != nil {
			return fmt.Errorf("Error when calling `UserAPI.FollowUser`: %v\n %+v\n", err, r)
		}
		return nil
	}(); err != nil {
		return err
	}

	// validate user 1 unfollows user 2
	if err := func() error {
		followers, r, err := apiClient.UserAPI.GetUserFollowers(oauth2Ctx, username2).Execute()
		if err != nil {
			return fmt.Errorf("Error when calling `UserAPI.FollowUser`: %v\n %+v\n", err, r)
		}

		if len(followers) != 0 {
			return fmt.Errorf("Expected 0 followers, got %d", len(followers))
		}
		return nil
	}(); err != nil {
		return err
	}
	return nil
}

func GetExpectedFeed(ctx context.Context) error {
	Setup()

	ctx, span := tracerhelper.GetTracer().Start(ctx, "GetExpectedFeed")
	defer span.End()

	// create two users
	configuration := NewDefaultConfiguration(WithSkipCache())
	httpClient := getHTTPClient()
	configuration.HTTPClient = httpClient
	proxyCtx := context.WithValue(ctx, oauth2.HTTPClient, httpClient)
	proxyCtx = context.WithValue(proxyCtx, client.ContextServerIndex, CONTEXT_SERVER)

	apiClient = client.NewAPIClient(configuration)

	username1 := fmt.Sprintf(defaultUsername, time.Now().UnixNano())
	password1 := "password"
	email1 := fmt.Sprintf("Test-%d-1@social.com", time.Now().UnixNano())
	user1 := *client.NewCreateUserRequest(username1, password1, "FirstName_example", "LastName_example", email1) // User | Create a new user

	username2 := fmt.Sprintf("Test-%d2", time.Now().UnixNano())
	password2 := "secretPassword"
	email2 := fmt.Sprintf("Test-%d-2@social.com", time.Now().UnixNano())
	user2 := *client.NewCreateUserRequest(username2, password2, "FirstName_example", "LastName_example", email2) // User | Create a new user

	// create users
	if err := func() error {
		_, r1, err1 := apiClient.UserAPI.
			CreateUser(proxyCtx).
			CreateUserRequest(user1).
			Execute()
		if err1 != nil {
			return fmt.Errorf("Error when calling `UserAPI.CreateUser`: %v\n %+v\n", err1, r1)
		}

		_, r2, err2 := apiClient.UserAPI.
			CreateUser(proxyCtx).
			CreateUserRequest(user2).
			Execute()
		if err2 != nil {
			return fmt.Errorf("Error when calling `UserAPI.CreateUser`: %v\n %+v\n", err2, r2)
		}
		return nil
	}(); err != nil {
		return err
	}

	// user 1 follows user 2
	conf1 := clientcredentials.Config{
		ClientID:     username1,
		ClientSecret: password1,
		Scopes: []string{
			"socialapp.users.read",
			"socialapp.follower.create",
			"socialapp.follower.read",
			"socialapp.follower.delete",
			"socialapp.comments.create",
			"socialapp.feed.read",
		},
		TokenURL: ENDPOINT_OAUTH_TOKEN,
	}

	oauth2Ctx1, err := getOuath2Context(proxyCtx, conf1)
	if err != nil {
		return fmt.Errorf("Error getting oauth2 context: %v", err)
	}

	if err := func() error {
		r, err := apiClient.UserAPI.FollowUser(
			oauth2Ctx1,
			username2,
			username1).
			Execute()
		if err != nil {
			return fmt.Errorf("Error when calling `UserAPI.FollowUser`: %v\n %+v\n", err, r)
		}
		return nil
	}(); err != nil {
		return err
	}

	// user 2 posts a comment
	conf2 := clientcredentials.Config{
		ClientID:     username2,
		ClientSecret: password2,
		Scopes: []string{
			"socialapp.comments.create",
			"socialapp.feed.read",
		},
		TokenURL: ENDPOINT_OAUTH_TOKEN,
	}

	oauth2Ctx2, err := getOuath2Context(proxyCtx, conf2)
	if err != nil {
		return fmt.Errorf("Error getting oauth2 context: %v", err)
	}
	if err := func() error {
		comment := *client.NewComment("Test comment", username2)
		_, r, err := apiClient.CommentAPI.
			CreateComment(oauth2Ctx2).
			Comment(comment).
			Execute()
		if err != nil {
			return fmt.Errorf("Error when calling `UserAPI.PostComment`: %v\n %+v\n", err, r)
		}
		return nil
	}(); err != nil {
		return err
	}

	// validate that comment from user 2 is in feed of user 1
	if err := func() error {
		feed, r, err := apiClient.CommentAPI.
			GetUserFeed(oauth2Ctx1).
			Execute()
		if err != nil {
			return fmt.Errorf("Error when calling `UserAPI.GetUserFeed`: %v\n %+v\n", err, r)
		}
		if len(feed) != 1 {
			return fmt.Errorf("Expected 1 post in feed, got %d", len(feed))
		}
		if feed[0].Username != username2 {
			return fmt.Errorf("Expected post from %s, got %s", username2, feed[0].Username)
		}
		return nil
	}(); err != nil {
		return err
	}

	// validate that feed from user 2 is empty
	if err := func() error {
		feed, r, err := apiClient.CommentAPI.
			GetUserFeed(oauth2Ctx2).
			Execute()
		if err != nil {
			return fmt.Errorf("Error when calling `UserAPI.GetUserFeed`: %v\n %+v\n", err, r)
		}
		if len(feed) != 0 {
			return fmt.Errorf("Expected 0 post in feed, got %d: \n %+v", len(feed), feed)
		}
		return nil
	}(); err != nil {
		return err
	}
	return nil
}

func GetAccessToken(ctx context.Context) error {
	Setup()
	ctx, span := tracerhelper.GetTracer().Start(ctx, "GetAccessToken")
	defer span.End()
	configuration := NewDefaultConfiguration(WithSkipCache())
	httpClient := getHTTPClient()
	configuration.HTTPClient = httpClient
	proxyCtx := context.WithValue(context.Background(), oauth2.HTTPClient, httpClient)
	proxyCtx = context.WithValue(proxyCtx, client.ContextServerIndex, CONTEXT_SERVER)

	apiClient = client.NewAPIClient(configuration)

	username := fmt.Sprintf(defaultUsername, time.Now().UnixNano())
	password := fmt.Sprintf(defaultPassword, time.Now().UnixNano())
	email := fmt.Sprintf("%s@example.com", username)
	createUsrReq := client.NewCreateUserRequest(username, password, "FirstName_example", "LastName_example", email)
	// create user
	if err := func() error {
		_, _, err := apiClient.UserAPI.CreateUser(proxyCtx).CreateUserRequest(*createUsrReq).Execute()
		if err != nil {
			return fmt.Errorf("Error creating user: %v", err)
		}
		return nil
	}(); err != nil {
		return err
	}
	scopes := []string{
		"socialapp.users.read",
		"socialapp.follower.create",
		"socialapp.follower.read",
		"socialapp.follower.delete",
		"socialapp.comments.create",
		"socialapp.feed.read",
	}
	credConf := clientcredentials.Config{
		ClientID:     username,
		ClientSecret: password,
		Scopes:       scopes,
		TokenURL:     ENDPOINT_OAUTH_TOKEN,
	}
	oauth2Ctx, err := getOuath2Context(proxyCtx, credConf)
	if err != nil {
		return fmt.Errorf("Error getting oauth2 context: %v", err)
	}

	// get token
	token, res, err := apiClient.AuthenticationAPI.GetAccessToken(oauth2Ctx).Execute()
	if err != nil {
		return fmt.Errorf("Error when calling `AuthenticationAPI.GetAccessToken`: %v", err)
	}
	// assert scopes are correct
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Expected status 200, got %s", res.Status)
	}

	if len(token.Scopes) != len(scopes) {
		return fmt.Errorf("Expected %d scopes, got %d", len(scopes), len(token.Scopes))
	}
	return nil
}

func RegisterUserFlow(ctx context.Context) error {
	Setup()
	ctx, span := tracerhelper.GetTracer().Start(ctx, "RegisterUserFlow")
	defer span.End()
	configuration := NewDefaultConfiguration(WithSkipCache())
	httpClient := getHTTPClient()
	configuration.HTTPClient = httpClient
	proxyCtx := context.WithValue(context.Background(), oauth2.HTTPClient, httpClient)
	proxyCtx = context.WithValue(proxyCtx, client.ContextServerIndex, CONTEXT_SERVER)

	apiClient = client.NewAPIClient(configuration)

	username1 := fmt.Sprintf(defaultUsername, time.Now().UnixNano())
	password := fmt.Sprintf(defaultPassword, time.Now().UnixNano())
	createUsrReq := client.NewCreateUserRequest(username1, password, "FirstName_example", "LastName_example", username1)

	// create a user, no auth needed
	_, res, err := apiClient.UserAPI.CreateUser(proxyCtx).CreateUserRequest(*createUsrReq).Execute()
	if err != nil {
		return fmt.Errorf("Error when calling `UserAPI.CreateUser`: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Expected status code 200, got %d", res.StatusCode)
	}

	scopes := []string{
		"socialapp.users.read",
	}
	conf := clientcredentials.Config{
		ClientID:     username1,
		ClientSecret: password,
		Scopes:       scopes,
		TokenURL:     ENDPOINT_OAUTH_TOKEN,
	}
	oauth2Ctx, err := getOuath2Context(proxyCtx, conf)
	if err != nil {
		return fmt.Errorf("Error getting oauth2 context: %v", err)
	}

	// Get user by using oauath2 token
	if err := func() error {
		_, res, err := apiClient.UserAPI.GetUserByUsername(oauth2Ctx, username1).Execute()
		if err != nil {
			return fmt.Errorf("Error when calling `UserAPI.GetUsers`: %v", err)
		}
		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("Expected status code 200, got %d", res.StatusCode)
		}
		return nil
	}(); err != nil {
		return err
	}

	// TODO API Should return 401 if no auth is provided
	// validate 403 if no auth is provided
	if err := func() error {
		user, res, err := apiClient.UserAPI.GetUserByUsername(proxyCtx, username1).Execute()
		if err == nil {
			return fmt.Errorf("Error when calling `UserAPI.GetUsers`: %v", err)
		}
		if res.StatusCode != http.StatusForbidden { // TOOD fix to 401
			return fmt.Errorf("Expected status code 401, got %d", res.StatusCode)
		}
		if user != nil {
			return fmt.Errorf("Expected nil user, got %v", user)
		}
		return nil
	}(); err != nil {
		return err
	}
	return nil
}

func ChangePassword(ctx context.Context) error {
	Setup()
	ctx, span := tracerhelper.GetTracer().Start(ctx, "ChangePassword")
	defer span.End()
	configuration := NewDefaultConfiguration(WithSkipCache())
	httpClient := getHTTPClient()
	configuration.HTTPClient = httpClient
	proxyCtx := context.WithValue(ctx, oauth2.HTTPClient, httpClient)
	proxyCtx = context.WithValue(proxyCtx, client.ContextServerIndex, CONTEXT_SERVER)
	scopes := []string{
		"socialapp.users.read",
		"socialapp.users.update",
	}

	apiClient = client.NewAPIClient(configuration)

	username := fmt.Sprintf(defaultUsername, time.Now().UnixNano())
	password := fmt.Sprintf(defaultPassword, time.Now().UnixNano())
	createUsrReq := client.NewCreateUserRequest(username, password, "FirstName_example", "LastName_example", username)

	// create a user, no auth needed
	_, res, err := apiClient.UserAPI.CreateUser(proxyCtx).CreateUserRequest(*createUsrReq).Execute()
	if err != nil {
		return fmt.Errorf("Error when calling `UserAPI.CreateUser`: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Expected status code 201, got %d", res.StatusCode)
	}

	conf := clientcredentials.Config{
		ClientID:     username,
		ClientSecret: password,
		Scopes:       scopes,
		TokenURL:     ENDPOINT_OAUTH_TOKEN,
	}
	oauth2Ctx, err := getOuath2Context(proxyCtx, conf)
	if err != nil {
		return fmt.Errorf("Error getting oauth2 context: %v", err)
	}

	newPassword := password + "new"
	if err := func() error {
		changePwdReq := client.NewChangePasswordRequest(password, newPassword)
		_, res, err := apiClient.UserAPI.ChangePassword(oauth2Ctx).ChangePasswordRequest(*changePwdReq).Execute()
		if err != nil {
			return fmt.Errorf("Error when calling `UserAPI.ChangePassword`: %v", err)
		}
		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("Expected status code 200, got %d", res.StatusCode)
		}
		return nil
	}(); err != nil {
		return err
	}

	// attempt to get token with old password, expect 401
	// func() {
	// 	token, res, err := apiClient.AuthenticationAPI.GetAccessToken(oauth2Ctx).Execute()
	// 	if err == nil {
	// 	return fmt.Errorf("Error when calling `AuthenticationAPI.GetAccessToken`: %v", err)
	// 	}
	// 	if res.StatusCode != http.StatusUnauthorized {
	// 	return fmt.Errorf("Expected status code 401, got %d", res.StatusCode)
	// 	}
	// 	if token != nil {
	// 	return fmt.Errorf("Expected nil user, got %v", token)
	// 	}
	// }()

	// attempt to get token with new password, expect 200
	if err := func() error {
		newPwdConf := clientcredentials.Config{
			ClientID:     username,
			ClientSecret: newPassword,
			Scopes:       scopes,
			TokenURL:     ENDPOINT_OAUTH_TOKEN,
		}
		newPwdOauth2Ctx, err := getOuath2Context(proxyCtx, newPwdConf)
		if err != nil {
			return fmt.Errorf("Error getting oauth2 context: %v", err)
		}
		token, res, err := apiClient.AuthenticationAPI.GetAccessToken(newPwdOauth2Ctx).Execute()
		if err != nil {
			return fmt.Errorf("Error when calling `AuthenticationAPI.GetAccessToken`: %v", err)
		}
		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("Expected status code 200, got %d", res.StatusCode)
		}
		if token == nil {
			return fmt.Errorf("Expected token, got nil")
		}
		return nil
	}(); err != nil {
		return err
	}

	return nil
}

func RoleLifecycle(ctx context.Context) error {
	Setup()
	ctx, span := tracerhelper.GetTracer().Start(ctx, "RoleLifecycle")
	defer span.End()
	configuration := NewDefaultConfiguration(WithSkipCache())
	httpClient := getHTTPClient()
	configuration.HTTPClient = httpClient
	proxyCtx := context.WithValue(context.Background(), oauth2.HTTPClient, httpClient)
	proxyCtx = context.WithValue(proxyCtx, client.ContextServerIndex, CONTEXT_SERVER)
	scopes := []string{
		"socialapp.roles.read",
		"socialapp.roles.list",
		"socialapp.roles.create",
		"socialapp.roles.update",
		"socialapp.roles.delete",
		"socialapp.scopes.create",
		"socialapp.roles.scopes.create",
		"socialapp.roles.scopes.delete",
		"socialapp.roles.scopes.list",
		"socialapp.scopes.delete",
	}

	apiClient = client.NewAPIClient(configuration)

	username := fmt.Sprintf(defaultUsername, time.Now().UnixNano())
	password := fmt.Sprintf(defaultPassword, time.Now().UnixNano())
	createUsrReq := client.NewCreateUserRequest(username, password, "FirstName_example", "LastName_example", username)

	// create a user, no auth needed
	_, res, err := apiClient.UserAPI.CreateUser(proxyCtx).CreateUserRequest(*createUsrReq).Execute()
	if err != nil {
		return fmt.Errorf("Error when calling `UserAPI.CreateUser`: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Expected status code 201, got %d", res.StatusCode)
	}

	conf := clientcredentials.Config{
		ClientID:     username,
		ClientSecret: password,
		Scopes:       scopes,
		TokenURL:     ENDPOINT_OAUTH_TOKEN,
	}
	oauth2Ctx, err := getOuath2Context(proxyCtx, conf)
	if err != nil {
		return fmt.Errorf("Error getting oauth2 context: %v", err)
	}

	// create a role
	newRole := client.NewRole(fmt.Sprintf("Test-CreateRole-%d1", time.Now().UnixNano()))
	createdRole, res, err := apiClient.RoleAPI.CreateRole(oauth2Ctx).Role(*newRole).Execute()
	if err != nil {
		return fmt.Errorf("Error when calling `RoleAPI.CreateRole`: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Expected status code 200, got %d", res.StatusCode)
	}

	// get role
	gettedRole, res, err := apiClient.RoleAPI.GetRole(oauth2Ctx, int32(*createdRole.Id)).Execute()
	if err != nil {
		return fmt.Errorf("Error when calling `RoleAPI.GetRole`: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Expected status code 200, got %d", res.StatusCode)
	}
	if gettedRole == nil {
		return fmt.Errorf("Expected role, got nil")
	}

	// add scopes to role
	// create scope
	newScope := client.NewScope(fmt.Sprintf("Test-CreateScope-%d1", time.Now().UnixNano()), "Test-CreateScope-Description")
	if err := func() error {
		createdScope, res, err := apiClient.ScopeAPI.CreateScope(oauth2Ctx).Scope(*newScope).Execute()
		if err != nil {
			return fmt.Errorf("Error when calling `ScopeAPI.CreateScope`: %v", err)
		}
		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("Expected status code 200, got %d", res.StatusCode)
		}
		// attach scope to role
		scopesToAdd := []string{newScope.Name}
		res, err = apiClient.RoleAPI.AddScopeToRole(oauth2Ctx, int32(*createdRole.Id)).RequestBody(scopesToAdd).Execute()
		if err != nil {
			return fmt.Errorf("Error when calling `RoleAPI.AddScopeToRole`: %v", err)
		}
		// verify scope is attached to role
		// get role scopes
		roleScopes, res, err := apiClient.RoleAPI.ListScopesForRole(oauth2Ctx, int32(*createdRole.Id)).Execute()
		if err != nil {
			return fmt.Errorf("Error when calling `RoleAPI.ListScopesForRole`: %v", err)
		}
		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("Expected status code 200, got %d", res.StatusCode)
		}
		if len(roleScopes) != 1 {
			return fmt.Errorf("Expected 1 scope, got %d", len(roleScopes))
		}
		if roleScopes[0].Name != newScope.Name {
			return fmt.Errorf("Expected scope name %s, got %s", newScope.Name, roleScopes[0].Name)
		}
		// remove scope from role
		res, err = apiClient.RoleAPI.RemoveScopeFromRole(oauth2Ctx, int32(*createdRole.Id), int32(*createdScope.Id)).Execute()
		if err != nil {
			return fmt.Errorf("Error when calling `RoleAPI.RemoveScopeFromRole`: %v", err)
		}
		if res.StatusCode != http.StatusNoContent {
			return fmt.Errorf("Expected status code 204, got %d", res.StatusCode)
		}
		// verify scope is removed from role
		// get role scopes
		roleScopes, res, err = apiClient.RoleAPI.ListScopesForRole(oauth2Ctx, int32(*createdRole.Id)).Execute()
		if err != nil {
			return fmt.Errorf("Error when calling `RoleAPI.ListScopesForRole`: %v", err)
		}
		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("Expected status code 200, got %d", res.StatusCode)
		}
		if len(roleScopes) != 0 {
			return fmt.Errorf("Expected 0 scopes, got %d", len(roleScopes))
		}
		// detach scope from role
		res, err = apiClient.RoleAPI.RemoveScopeFromRole(oauth2Ctx, int32(*createdRole.Id), int32(*createdScope.Id)).Execute()
		if err != nil {
			return fmt.Errorf("Error when calling `RoleAPI.RemoveScopeFromRole`: %v", err)
		}
		if res.StatusCode != http.StatusNoContent {
			return fmt.Errorf("Expected status code 204, got %d", res.StatusCode)
		}
		// verify scope is detached from role
		// get role scopes
		roleScopes, res, err = apiClient.RoleAPI.ListScopesForRole(oauth2Ctx, int32(*createdRole.Id)).Execute()
		if err != nil {
			return fmt.Errorf("Error when calling `RoleAPI.ListScopesForRole`: %v", err)
		}
		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("Expected status code 200, got %d", res.StatusCode)
		}
		if len(roleScopes) != 0 {
			return fmt.Errorf("Expected 0 scopes, got %d", len(roleScopes))
		}

		// delete scope
		res, err = apiClient.ScopeAPI.DeleteScope(oauth2Ctx, int32(*createdScope.Id)).Execute()
		if err != nil {
			return fmt.Errorf("Error when calling `ScopeAPI.DeleteScope`: %v", err)
		}
		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("Expected status code 200, got %d", res.StatusCode)
		}
		return nil
	}(); err != nil {
		return err
	}

	updatedName := fmt.Sprintf("Test-UpdateRole-%d", time.Now().UnixNano())
	updatedDesc := fmt.Sprintf("Test-UpdateRole-Description-%d", time.Now().UnixNano())
	newRole.Description = &updatedDesc
	newRole.Name = updatedName

	// update role
	updatedRole, res, err := apiClient.RoleAPI.UpdateRole(oauth2Ctx, int32(*createdRole.Id)).Role(*newRole).Execute()
	if err != nil {
		return fmt.Errorf("Error when calling `RoleAPI.UpdateRole`: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Expected status code 200, got %d", res.StatusCode)
	}
	if updatedRole == nil {
		return fmt.Errorf("Expected role, got nil")
	}
	if updatedRole.Name != updatedName {
		return fmt.Errorf("Expected name %s, got %s", updatedName, updatedRole.Name)
	}
	if *updatedRole.Description != updatedDesc {
		return fmt.Errorf("Expected description %s, got %s", updatedDesc, *updatedRole.Description)
	}

	// get role again to check if updated
	gettedRole, res, err = apiClient.RoleAPI.GetRole(oauth2Ctx, int32(*createdRole.Id)).Execute()
	if err != nil {
		return fmt.Errorf("Error when calling `RoleAPI.GetRole`: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Expected status code 200, got %d", res.StatusCode)
	}
	if gettedRole == nil {
		return fmt.Errorf("Expected role, got nil")
	}
	if gettedRole.Name != updatedName {
		return fmt.Errorf("Expected name %s, got %s", updatedName, gettedRole.Name)
	}
	if *gettedRole.Description != updatedDesc {
		return fmt.Errorf("Expected description %s, got %s", updatedDesc, *gettedRole.Description)
	}

	// delete role
	if err := func() error {
		res, err := apiClient.RoleAPI.DeleteRole(oauth2Ctx, int32(*createdRole.Id)).Execute()
		if err != nil {
			return fmt.Errorf("Error when calling `RoleAPI.DeleteRole`: %v", err)
		}
		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("Expected status code 200, got %d", res.StatusCode)
		}
		return nil
	}(); err != nil {
		return err
	}

	// try to get deleted role
	gettedDeletedRole, res, err := apiClient.RoleAPI.GetRole(oauth2Ctx, int32(*createdRole.Id)).Execute()
	if err == nil {
		return fmt.Errorf("Expected error when calling `RoleAPI.GetRole`, got nil")
	}
	if res.StatusCode != http.StatusNotFound {
		return fmt.Errorf("Expected status code 404, got %d", res.StatusCode)
	}
	if gettedDeletedRole != nil {
		return fmt.Errorf("Expected nil, got %v", gettedRole)
	}
	return nil
}

func ScopeLifecycle(ctx context.Context) error {
	Setup()
	ctx, span := tracerhelper.GetTracer().Start(ctx, "ScopeLifecycle")
	defer span.End()
	configuration := NewDefaultConfiguration(WithSkipCache())
	httpClient := getHTTPClient()
	configuration.HTTPClient = httpClient
	proxyCtx := context.WithValue(context.Background(), oauth2.HTTPClient, httpClient)
	proxyCtx = context.WithValue(proxyCtx, client.ContextServerIndex, CONTEXT_SERVER)
	scopes := []string{
		"socialapp.scopes.read",
		"socialapp.scopes.list",
		"socialapp.scopes.create",
		"socialapp.scopes.update",
		"socialapp.scopes.delete",
	}

	apiClient = client.NewAPIClient(configuration)

	username := fmt.Sprintf(defaultUsername, time.Now().UnixNano())
	password := fmt.Sprintf(defaultPassword, time.Now().UnixNano())
	createUsrReq := client.NewCreateUserRequest(username, password, "FirstName_example", "LastName_example", username)

	// create a user, no auth needed
	_, res, err := apiClient.UserAPI.CreateUser(proxyCtx).CreateUserRequest(*createUsrReq).Execute()
	if err != nil {
		return fmt.Errorf("Error when calling `UserAPI.CreateUser`: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Expected status code 201, got %d", res.StatusCode)
	}

	conf := clientcredentials.Config{
		ClientID:     username,
		ClientSecret: password,
		Scopes:       scopes,
		TokenURL:     ENDPOINT_OAUTH_TOKEN,
	}
	oauth2Ctx, err := getOuath2Context(proxyCtx, conf)
	if err != nil {
		return fmt.Errorf("Error getting oauth2 context: %v", err)
	}

	newScope := client.NewScope(fmt.Sprintf("Test-CreateScope-%d1", time.Now().UnixNano()), "Test-CreateScope-Description1")
	// create a scope

	createdScope, res, err := apiClient.ScopeAPI.CreateScope(oauth2Ctx).Scope(*newScope).Execute()
	if err != nil {
		return fmt.Errorf("Error when calling `ScopeAPI.CreateScope`: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Expected status code 200, got %d", res.StatusCode)
	}

	// get scope
	gettedScope, res, err := apiClient.ScopeAPI.GetScope(oauth2Ctx, int32(*createdScope.Id)).Execute()
	if err != nil {
		return fmt.Errorf("Error when calling `ScopeAPI.GetScope`: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Expected status code 200, got %d", res.StatusCode)
	}
	if gettedScope == nil {
		return fmt.Errorf("Expected scope, got nil")
	}

	updatedName := fmt.Sprintf("Test-UpdateScope-%d", time.Now().UnixNano())
	updatedDesc := fmt.Sprintf("Test-UpdateScope-Description-%d", time.Now().UnixNano())
	newScope.Description = updatedDesc
	newScope.Name = updatedName

	// update scope
	updatedScope, res, err := apiClient.ScopeAPI.UpdateScope(oauth2Ctx, int32(*createdScope.Id)).Scope(*newScope).Execute()
	if err != nil {
		return fmt.Errorf("Error when calling `ScopeAPI.UpdateScope`: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Expected status code 200, got %d", res.StatusCode)
	}
	if updatedScope == nil {
		return fmt.Errorf("Expected scope, got nil")
	}
	if updatedScope.Name != updatedName {
		return fmt.Errorf("Expected name %s, got %s", updatedName, updatedScope.Name)
	}
	if updatedScope.Description != updatedDesc {
		return fmt.Errorf("Expected description %s, got %s", updatedDesc, updatedScope.Description)
	}

	// get scope again to check if updated
	gettedScope, res, err = apiClient.ScopeAPI.GetScope(oauth2Ctx, int32(*createdScope.Id)).Execute()
	if err != nil {
		return fmt.Errorf("Error when calling `ScopeAPI.GetScope`: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Expected status code 200, got %d", res.StatusCode)
	}
	if gettedScope == nil {
		return fmt.Errorf("Expected scope, got nil")
	}
	if gettedScope.Name != updatedName {
		return fmt.Errorf("Expected name %s, got %s", updatedName, gettedScope.Name)
	}
	if gettedScope.Description != updatedDesc {
		return fmt.Errorf("Expected description %s, got %s", updatedDesc, gettedScope.Description)
	}

	// delete scope
	if err := func() error {
		res, err := apiClient.ScopeAPI.DeleteScope(oauth2Ctx, int32(*createdScope.Id)).Execute()
		if err != nil {
			return fmt.Errorf("Error when calling `ScopeAPI.DeleteScope`: %v", err)
		}
		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("Expected status code 204, got %d", res.StatusCode)
		}
		return nil
	}(); err != nil {
		return err
	}

	// try to get deleted scope
	gettedDeletedScope, res, err := apiClient.ScopeAPI.GetScope(oauth2Ctx, int32(*createdScope.Id)).Execute()
	if err == nil {
		return fmt.Errorf("Expected error when calling `ScopeAPI.GetScope`, got nil")
	}
	if res.StatusCode != http.StatusNotFound {
		return fmt.Errorf("Expected status code 404, got %d", res.StatusCode)
	}
	if gettedDeletedScope != nil {
		return fmt.Errorf("Expected nil, got %v", gettedScope)
	}
	return nil
}

func UserRoleLifeCycle(ctx context.Context) (err error) {
	Setup()
	ctx, span := tracerhelper.GetTracer().Start(ctx, "UserRoleLifeCycle")
	defer span.End()
	configuration := NewDefaultConfiguration(WithSkipCache())
	httpClient := getHTTPClient()
	configuration.HTTPClient = httpClient
	proxyCtx := context.WithValue(context.Background(), oauth2.HTTPClient, httpClient)
	proxyCtx = context.WithValue(proxyCtx, client.ContextServerIndex, CONTEXT_SERVER)
	scopes := []string{
		"socialapp.scopes.read",
		"socialapp.scopes.list",
		"socialapp.scopes.create",
		"socialapp.scopes.update",
		"socialapp.scopes.delete",
		"socialapp.roles.create",
		"socialapp.users.roles.update",
		"socialapp.users.roles.list",
		"socialapp.roles.delete",
	}

	apiClient = client.NewAPIClient(configuration)

	username := fmt.Sprintf(defaultUsername, time.Now().UnixNano())
	password := fmt.Sprintf(defaultPassword, time.Now().UnixNano())
	createUsrReq := client.NewCreateUserRequest(username, password, "FirstName_example", "LastName_example", username)

	// create a user, no auth needed
	_, res, err := apiClient.UserAPI.CreateUser(proxyCtx).CreateUserRequest(*createUsrReq).Execute()
	if err != nil {
		return fmt.Errorf("Error when calling `UserAPI.CreateUser`: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Expected status code 201, got %d", res.StatusCode)
	}

	conf := clientcredentials.Config{
		ClientID:     username,
		ClientSecret: password,
		Scopes:       scopes,
		TokenURL:     ENDPOINT_OAUTH_TOKEN,
	}
	oauth2Ctx, err := getOuath2Context(proxyCtx, conf)
	if err != nil {
		return fmt.Errorf("Error getting oauth2 context: %v", err)
	}

	newScope := client.NewScope(fmt.Sprintf("Test-CreateScope-%d1", time.Now().UnixNano()), "Test-CreateScope-Description1")

	// create a scope
	cretedScope, res, err := apiClient.ScopeAPI.CreateScope(oauth2Ctx).Scope(*newScope).Execute()
	if err != nil {
		return fmt.Errorf("Error when calling `ScopeAPI.CreateScope`: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Expected status code 200, got %d", res.StatusCode)
	}
	defer func() {
		// delete scope
		res, err := apiClient.ScopeAPI.DeleteScope(oauth2Ctx, int32(*cretedScope.Id)).Execute()
		if err != nil {
			err = fmt.Errorf("Error when calling `ScopeAPI.DeleteScope`: %v", err)
		}
		if res.StatusCode != http.StatusOK {
			err = fmt.Errorf("Expected status code 204, got %d", res.StatusCode)
		}
	}()

	// create role
	newRole := client.NewRole(fmt.Sprintf("Test-CreateRole-%d1", time.Now().UnixNano()))
	createdRole, res, err := apiClient.RoleAPI.CreateRole(oauth2Ctx).Role(*newRole).Execute()
	if err != nil {
		return fmt.Errorf("Error when calling `RoleAPI.CreateRole`: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Expected status code 200, got %d", res.StatusCode)
	}

	// attach role to user
	_, res, err = apiClient.UserAPI.UpdateRolesForUser(oauth2Ctx, username).RequestBody([]string{createdRole.Name}).Execute()
	if err != nil {
		return fmt.Errorf("Error when calling `UserAPI.UpdateRolesForUser`: %v", err)
	}

	defer func() {
		// delete role
		res, err := apiClient.RoleAPI.DeleteRole(oauth2Ctx, int32(*createdRole.Id)).Execute()
		if err != nil {
			err = fmt.Errorf("Error when calling `RoleAPI.DeleteRole`: %v", err)
		}
		if res.StatusCode != http.StatusOK {
			err = fmt.Errorf("Expected status code 204, got %d", res.StatusCode)
		}
	}()

	// List user roles
	roles, res, err := apiClient.UserAPI.GetRolesForUser(oauth2Ctx, username).Execute()
	if err != nil {
		return fmt.Errorf("Error when calling `UserAPI.GetRolesForUser`: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Expected status code 200, got %d", res.StatusCode)
	}
	// verify role is attached to user, default roles are already attached
	foundRole := false
	for i := range roles {
		if roles[i].Name == createdRole.Name {
			foundRole = true
			break
		}
	}
	if !foundRole {
		return fmt.Errorf("Expected to find role %s, got %v", createdRole.Name, roles)
	}

	return nil
}

func CacheRequestSameUser(ctx context.Context) error {
	Setup()
	ctx, span := tracerhelper.GetTracer().Start(ctx, "CacheRequestSameUser")
	defer span.End()
	configuration := client.NewConfiguration()
	httpClient := getHTTPClient()
	configuration.HTTPClient = httpClient

	proxyCtx := context.WithValue(context.Background(), oauth2.HTTPClient, httpClient)
	proxyCtx = context.WithValue(proxyCtx, client.ContextServerIndex, CONTEXT_SERVER)

	username1 := fmt.Sprintf(defaultUsername, time.Now().UnixNano())
	password := fmt.Sprintf(defaultPassword, time.Now().UnixNano())
	apiClient = client.NewAPIClient(configuration)
	if err := func() error {
		createUsrReq := client.NewCreateUserRequest(username1, password, "FirstName_example", "LastName_example", username1)
		_, _, err := apiClient.UserAPI.CreateUser(proxyCtx).CreateUserRequest(*createUsrReq).Execute()
		if err != nil {
			return fmt.Errorf("Error creating user: %v", err)
		}
		return nil
	}(); err != nil {
		return err
	}

	conf := clientcredentials.Config{
		ClientID:     username1,
		ClientSecret: password,
		Scopes: []string{
			"socialapp.users.list",
			"socialapp.users.read",
			"socialapp.feed.read",
			"socialapp.comments.read",
		},
		TokenURL: ENDPOINT_OAUTH_TOKEN,
	}
	oauth2Ctx, err := getOuath2Context(proxyCtx, conf)
	if err != nil {
		return fmt.Errorf("Error getting oauth2 context: %v", err)
	}
	openAPICtx := context.WithValue(oauth2Ctx, client.ContextServerIndex, CONTEXT_SERVER)

	// List 100 users, different offset on every execution
	offset := time.Now().UnixNano() % 5000
	listedUsers, r, err := apiClient.UserAPI.
		ListUsers(openAPICtx).
		Limit(5).
		Offset(int32(offset)).
		Execute()

	if err != nil {
		return fmt.Errorf("Error when calling `UserAPI.ListUsers`: %v %v", err, r)
	}
	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("Expected status code %d, got %d", http.StatusOK, r.StatusCode)
	}

	// get user info 5 times
	for _, currentUser := range listedUsers {
		for i := 0; i < 50; i++ {
			_, r, err = apiClient.UserAPI.GetUserByUsername(openAPICtx, currentUser.Username).Execute()
			if err != nil {
				return fmt.Errorf("Error when calling `UserAPI.GetUser`: %v %v", err, r)

			}
			if r.StatusCode != http.StatusOK {
				return fmt.Errorf("Expected status code %d, got %d", http.StatusOK, r.StatusCode)
			}

			// get user comments
			_, r, err = apiClient.UserAPI.GetUserComments(openAPICtx, currentUser.Username).Execute()
			if err != nil {
				return fmt.Errorf("Error when calling `UserAPI.GetUserComments`: %v %v", err, r)

			}
			if r.StatusCode != http.StatusOK {
				return fmt.Errorf("Expected status code %d, got %d", http.StatusOK, r.StatusCode)
			}

			// get feed
			_, r, err = apiClient.CommentAPI.GetUserFeed(openAPICtx).Execute()
			if err != nil {
				return fmt.Errorf("Error when calling `CommentAPI.GetUserFeed`: %v %v", err, r)

			}
			if r.StatusCode != http.StatusOK {
				return fmt.Errorf("Expected status code %d, got %d", http.StatusOK, r.StatusCode)
			}

			// t.Logf("%q: %d/100\n", currentUser.Username, i)
		}
	}
	return nil
}

func URLLifeCycle(ctx context.Context) error {
	Setup()
	ctx, span := tracerhelper.GetTracer().Start(ctx, "URLLifeCycle")
	defer span.End()
	configuration := NewDefaultConfiguration(WithSkipCache())
	httpClient := getHTTPClient()
	configuration.HTTPClient = httpClient

	proxyCtx := context.WithValue(context.Background(), oauth2.HTTPClient, httpClient)
	proxyCtx = context.WithValue(proxyCtx, client.ContextServerIndex, CONTEXT_SERVER)

	username1 := fmt.Sprintf(defaultUsername, time.Now().UnixNano())
	password := fmt.Sprintf(defaultPassword, time.Now().UnixNano())
	apiClient = client.NewAPIClient(configuration)
	if err := func() error {
		createUsrReq := client.NewCreateUserRequest(username1, password, "FirstName_example", "LastName_example", username1)
		_, _, err := apiClient.UserAPI.CreateUser(proxyCtx).CreateUserRequest(*createUsrReq).Execute()
		if err != nil {
			return fmt.Errorf("Error creating user: %v", err)
		}
		return nil
	}(); err != nil {
		return err
	}

	conf := clientcredentials.Config{
		ClientID:     username1,
		ClientSecret: password,
		Scopes:       []string{"shortly.url.create", "shortly.url.delete"},
		TokenURL:     ENDPOINT_OAUTH_TOKEN,
	}

	oauth2Ctx, err := getOuath2Context(proxyCtx, conf)
	if err != nil {
		return fmt.Errorf("Error getting oauth2 context: %v", err)
	}
	openAPICtx := context.WithValue(oauth2Ctx, client.ContextServerIndex, CONTEXT_SERVER)

	// create url
	newURL := client.NewURL("https://www.google.com", fmt.Sprintf("%d", time.Now().Unix()))
	_, r, err := apiClient.URLAPI.CreateUrl(openAPICtx).URL(*newURL).Execute()
	if err != nil {
		return fmt.Errorf("Error when calling `URLAPI.CreateURL`: %v %v", err, r)
	}

	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("Expected status code %d, got %d", http.StatusOK, r.StatusCode)
	}

	// get url
	getUrlRes, err := apiClient.URLAPI.GetUrl(proxyCtx, newURL.Alias).Execute()
	if err != nil {
		return fmt.Errorf("Error when calling `URLAPI.GetURL`: %v %v", err, getUrlRes)
	}

	if getUrlRes.StatusCode != http.StatusOK {
		return fmt.Errorf("Expected status code %d, got %d", http.StatusOK, getUrlRes.StatusCode)
	}

	// delete url
	deleteUrlRes, err := apiClient.URLAPI.DeleteUrl(openAPICtx, newURL.Alias).Execute()
	if err != nil {
		return fmt.Errorf("Error when calling `URLAPI.DeleteURL`: %v %v", err, deleteUrlRes)
	}

	if deleteUrlRes.StatusCode != http.StatusOK {
		return fmt.Errorf("Expected status code %d, got %d", http.StatusOK, deleteUrlRes.StatusCode)
	}

	// get url
	getUrlRes, err = apiClient.URLAPI.GetUrl(proxyCtx, newURL.Alias).Execute()
	if err == nil {
		return fmt.Errorf("Expected error when calling `URLAPI.GetURL`, got none")
	}
	if getUrlRes.StatusCode != http.StatusNotFound {
		return fmt.Errorf("Expected status code %d, got %d", http.StatusNotFound, getUrlRes.StatusCode)
	}

	return nil
}
