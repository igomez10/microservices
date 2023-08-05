package integration_tests

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/igomez10/microservices/socialapp/client"
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
)

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
}

func TestMain(m *testing.M) {
	// run tests
	// add jitter
	if os.Getenv("ADD_TEST_JITTER") != "" {
		jitterInSeconds := uuid.New().ID() % 60
		log.Printf("Adding test jitter of %d seconds", jitterInSeconds)
		time.Sleep(time.Duration(jitterInSeconds))
	}

	Setup()
	code := m.Run()
	os.Exit(code)
}

func getHTTPClient() *http.Client {
	if os.Getenv("USE_PROXY") == "true" {
		proxyStr := "http://localhost:9091"
		proxyURL, err := url.Parse(proxyStr)
		if err != nil {
			return http.DefaultClient
		}
		// Setup http client with proxy to capture traffic
		httpClient := &http.Client{
			Timeout: time.Second * 10,
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			},
		}

		return httpClient
	}

	return http.DefaultClient
}

func getOuath2Context(initialContext context.Context, config clientcredentials.Config) (context.Context, error) {
	tokenSource := config.TokenSource(initialContext)
	initialContext = context.WithValue(initialContext, client.ContextOAuth2, tokenSource)

	return initialContext, nil
}

func TestListUsers(t *testing.T) {
	Setup()
	configuration := client.NewConfiguration()
	httpClient := getHTTPClient()
	configuration.HTTPClient = httpClient

	proxyCtx := context.WithValue(context.Background(), oauth2.HTTPClient, httpClient)
	proxyCtx = context.WithValue(proxyCtx, client.ContextServerIndex, CONTEXT_SERVER)

	username1 := fmt.Sprintf("Test-%d1", time.Now().UnixNano())
	password := fmt.Sprintf("Password-%d1", time.Now().UnixNano())
	apiClient = client.NewAPIClient(configuration)
	//skip cache for tests
	configuration.DefaultHeader["Cache-Control"] = "no-store"
	func() {
		createUsrReq := client.NewCreateUserRequest(username1, password, "FirstName_example", "LastName_example", username1)
		_, _, err := apiClient.UserApi.CreateUser(proxyCtx).CreateUserRequest(*createUsrReq).Execute()
		if err != nil {
			t.Fatalf("Error creating user: %v", err)
		}
	}()

	conf := clientcredentials.Config{
		ClientID:     username1,
		ClientSecret: password,
		Scopes:       []string{"socialapp.users.list"},
		TokenURL:     ENDPOINT_OAUTH_TOKEN,
	}
	oauth2Ctx, err := getOuath2Context(proxyCtx, conf)
	if err != nil {
		t.Fatalf("Error getting oauth2 context: %v", err)
	}
	openAPICtx := context.WithValue(oauth2Ctx, client.ContextServerIndex, CONTEXT_SERVER)

	// List users
	_, r, err := apiClient.UserApi.ListUsers(openAPICtx).Limit(10).Offset(0).Execute()
	if err != nil {
		t.Errorf("Error when calling `UserApi.ListUsers`: %v\n", err)
		t.Errorf("Full HTTP response: %v\n", r)
	}
	if r.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, r.StatusCode)
	}
}

func TestCreateUser(t *testing.T) {
	Setup()
	configuration := client.NewConfiguration()
	//skip cache for tests
	configuration.DefaultHeader["Cache-Control"] = "no-store"
	httpClient := getHTTPClient()
	configuration.HTTPClient = httpClient
	noAuthCtx := context.WithValue(context.Background(), oauth2.HTTPClient, httpClient)
	noAuthCtx = context.WithValue(noAuthCtx, client.ContextServerIndex, CONTEXT_SERVER)

	apiClient = client.NewAPIClient(configuration)

	username := fmt.Sprintf("Test-%d", time.Now().UnixNano())
	password := "password"
	email := fmt.Sprintf("Test-%d-@social.com", time.Now().UnixNano())
	user := *client.NewCreateUserRequest(username, "password", "FirstName_example", "LastName_example", email) // User | Create a new user

	func() {
		_, r, err := apiClient.UserApi.CreateUser(noAuthCtx).
			CreateUserRequest(user).
			Execute()
		if err != nil {
			t.Errorf("Error when calling `UserApi.CreateUser`: %v\n %+v\n", err, r)
		}
		if r.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, r.StatusCode)
		}
	}()

	func() {
		conf := clientcredentials.Config{
			ClientID:     username,
			ClientSecret: password,
			Scopes:       []string{"socialapp.users.read"},
			TokenURL:     ENDPOINT_OAUTH_TOKEN,
		}
		oauth2Ctx, err := getOuath2Context(noAuthCtx, conf)
		if err != nil {
			t.Fatalf("Error getting oauth2 context: %v", err)
		}

		resp, r, err := apiClient.UserApi.GetUserByUsername(oauth2Ctx, username).Execute()
		if err != nil {
			t.Errorf("Error when calling `UserApi.GetUserByUsername`: %v\n %+v\n", err, r)
		}
		if r.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, r.StatusCode)
		}
		if resp.Username != user.Username {
			t.Errorf("Expected username %s, got %s", user.Username, resp.Username)
		}
		if resp.Email != user.Email {
			t.Errorf("Expected email %s, got %s", user.Email, resp.Email)
		}
		if resp.FirstName != user.FirstName {
			t.Errorf("Expected first name %q, got %q", user.FirstName, resp.FirstName)
		}
		if resp.LastName != user.LastName {
			t.Errorf("Expected last name %q, got %q", user.LastName, resp.LastName)
		}
	}()

	// update user
	func() {
		conf := clientcredentials.Config{
			ClientID:     username,
			ClientSecret: password,
			Scopes:       []string{"socialapp.users.update"},
			TokenURL:     ENDPOINT_OAUTH_TOKEN,
		}
		oauth2Ctx, err := getOuath2Context(noAuthCtx, conf)
		if err != nil {
			t.Fatalf("Error getting oauth2 context: %v", err)
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

		updateUserReq := apiClient.UserApi.
			UpdateUser(oauth2Ctx, username).
			User(updatedUser)

		uUser, res, err := updateUserReq.Execute()
		if err != nil {
			t.Fatalf("Error when calling `UserApi.UpdateUser`: %v\n %+v\n", err, res)
		}
		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, res.StatusCode)
		}
		if uUser.FirstName != updatedFirstName {
			t.Errorf("Expected first name %q, got %q", updatedFirstName, uUser.FirstName)
		}
		if uUser.LastName != updatedLastName {
			t.Errorf("Expected last name %q, got %q", updatedLastName, uUser.LastName)
		}
		if uUser.Email != updatedEmail {
			t.Errorf("Expected email %q, got %q", updatedEmail, uUser.Email)
		}
	}()

}

func TestFollowCycle(t *testing.T) {
	Setup()
	// create two users
	configuration := client.NewConfiguration()
	//skip cache for tests
	configuration.DefaultHeader["Cache-Control"] = "no-store"
	httpClient := getHTTPClient()
	configuration.HTTPClient = httpClient
	proxyCtx := context.WithValue(context.Background(), oauth2.HTTPClient, httpClient)
	proxyCtx = context.WithValue(proxyCtx, client.ContextServerIndex, CONTEXT_SERVER)

	apiClient = client.NewAPIClient(configuration)

	username1 := fmt.Sprintf("Test-%d1", time.Now().UnixNano())
	password1 := fmt.Sprintf("TestPassword-%d1", time.Now().UnixNano())
	email1 := fmt.Sprintf("Test-%d-1@social.com", time.Now().UnixNano())
	user1 := *client.NewCreateUserRequest(username1, password1, "FirstName_example", "LastName_example", email1) // User | Create a new user

	username2 := fmt.Sprintf("Test-%d2", time.Now().UnixNano())
	password2 := fmt.Sprintf("TestPassword-%d2", time.Now().UnixNano())
	email2 := fmt.Sprintf("Test-%d-2@social.com", time.Now().UnixNano())
	user2 := *client.NewCreateUserRequest(username2, password2, "FirstName_example", "LastName_example", email2) // User | Create a new user

	// create users
	func() {
		_, r1, err1 := apiClient.UserApi.CreateUser(proxyCtx).
			CreateUserRequest(user1).
			Execute()
		if err1 != nil {
			t.Errorf("Error when calling `UserApi.CreateUser`: %v\n %+v\n", err1, r1)
		}

		_, r2, err2 := apiClient.UserApi.CreateUser(proxyCtx).
			CreateUserRequest(user2).
			Execute()
		if err2 != nil {
			t.Errorf("Error when calling `UserApi.CreateUser`: %v\n %+v\n", err2, r2)
		}
	}()

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
	func() {

		if err != nil {
			t.Fatalf("Error getting oauth2 context: %v", err)
		}
		r, err := apiClient.UserApi.FollowUser(oauth2Ctx, username2, username1).Execute()
		if err != nil {
			t.Errorf("Error when calling `UserApi.FollowUser`: %v\n %+v\n", err, r)
		}
	}()

	// validate user 1 follows user 2
	func() {
		followers, r, err := apiClient.UserApi.GetUserFollowers(oauth2Ctx, username2).Execute()
		if err != nil {
			t.Errorf("Error when calling `UserApi.FollowUser`: %v\n %+v\n", err, r)
		}
		if len(followers) != 1 {
			t.Errorf("Expected 1 follower, got %d", len(followers))
		}
		if followers[0].Username != username1 {
			t.Errorf("Expected follower %s, got %s", username1, followers[0].Username)
		}
	}()

	// user 1 unfollows user 2
	func() {
		r, err := apiClient.UserApi.UnfollowUser(oauth2Ctx, username2, username1).Execute()
		if err != nil {
			t.Errorf("Error when calling `UserApi.FollowUser`: %v\n %+v\n", err, r)
		}
	}()

	// validate user 1 unfollows user 2
	func() {
		followers, r, err := apiClient.UserApi.GetUserFollowers(oauth2Ctx, username2).Execute()
		if err != nil {
			t.Errorf("Error when calling `UserApi.FollowUser`: %v\n %+v\n", err, r)
		}

		if len(followers) != 0 {
			t.Errorf("Expected 0 followers, got %d", len(followers))
		}
	}()
}

func TestGetExpectedFeed(t *testing.T) {
	Setup()
	// create two users
	configuration := client.NewConfiguration()
	//skip cache for tests
	configuration.DefaultHeader["Cache-Control"] = "no-store"
	httpClient := getHTTPClient()
	configuration.HTTPClient = httpClient
	proxyCtx := context.WithValue(context.Background(), oauth2.HTTPClient, httpClient)
	proxyCtx = context.WithValue(proxyCtx, client.ContextServerIndex, CONTEXT_SERVER)

	apiClient = client.NewAPIClient(configuration)

	username1 := fmt.Sprintf("Test-%d1", time.Now().UnixNano())
	password1 := "password"
	email1 := fmt.Sprintf("Test-%d-1@social.com", time.Now().UnixNano())
	user1 := *client.NewCreateUserRequest(username1, password1, "FirstName_example", "LastName_example", email1) // User | Create a new user

	username2 := fmt.Sprintf("Test-%d2", time.Now().UnixNano())
	password2 := "secretPassword"
	email2 := fmt.Sprintf("Test-%d-2@social.com", time.Now().UnixNano())
	user2 := *client.NewCreateUserRequest(username2, password2, "FirstName_example", "LastName_example", email2) // User | Create a new user

	// create users
	func() {
		_, r1, err1 := apiClient.UserApi.
			CreateUser(proxyCtx).
			CreateUserRequest(user1).
			Execute()
		if err1 != nil {
			t.Errorf("Error when calling `UserApi.CreateUser`: %v\n %+v\n", err1, r1)
		}

		_, r2, err2 := apiClient.UserApi.
			CreateUser(proxyCtx).
			CreateUserRequest(user2).
			Execute()
		if err2 != nil {
			t.Errorf("Error when calling `UserApi.CreateUser`: %v\n %+v\n", err2, r2)
		}
	}()

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
		t.Fatalf("Error getting oauth2 context: %v", err)
	}

	func() {
		r, err := apiClient.UserApi.FollowUser(
			oauth2Ctx1,
			username2,
			username1).
			Execute()
		if err != nil {
			t.Errorf("Error when calling `UserApi.FollowUser`: %v\n %+v\n", err, r)
		}
	}()

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
		t.Fatalf("Error getting oauth2 context: %v", err)
	}
	func() {
		comment := *client.NewComment("Test comment", username2)
		_, r, err := apiClient.CommentApi.
			CreateComment(oauth2Ctx2).
			Comment(comment).
			Execute()
		if err != nil {
			t.Errorf("Error when calling `UserApi.PostComment`: %v\n %+v\n", err, r)
		}
	}()

	// validate that comment from user 2 is in feed of user 1
	func() {
		feed, r, err := apiClient.CommentApi.
			GetUserFeed(oauth2Ctx1).
			Execute()
		if err != nil {
			t.Errorf("Error when calling `UserApi.GetUserFeed`: %v\n %+v\n", err, r)
		}
		if len(feed) != 1 {
			t.Errorf("Expected 1 post in feed, got %d", len(feed))
		}
		if feed[0].Username != username2 {
			t.Errorf("Expected post from %s, got %s", username2, feed[0].Username)
		}
	}()

	// validate that feed from user 2 is empty
	func() {
		feed, r, err := apiClient.CommentApi.
			GetUserFeed(oauth2Ctx2).
			Execute()
		if err != nil {
			t.Errorf("Error when calling `UserApi.GetUserFeed`: %v\n %+v\n", err, r)
		}
		if len(feed) != 0 {
			t.Errorf("Expected 0 post in feed, got %d: \n %+v", len(feed), feed)
		}
	}()
}

func TestGetAccessToken(t *testing.T) {
	Setup()
	// create two users
	configuration := client.NewConfiguration()
	//skip cache for tests
	configuration.DefaultHeader["Cache-Control"] = "no-store"
	httpClient := getHTTPClient()
	configuration.HTTPClient = httpClient
	proxyCtx := context.WithValue(context.Background(), oauth2.HTTPClient, httpClient)
	proxyCtx = context.WithValue(proxyCtx, client.ContextServerIndex, CONTEXT_SERVER)

	apiClient = client.NewAPIClient(configuration)

	username1 := fmt.Sprintf("Test-%d1", time.Now().UnixNano())
	password := fmt.Sprintf("Password-%d1", time.Now().UnixNano())
	createUsrReq := client.NewCreateUserRequest(username1, password, "FirstName_example", "LastName_example", username1)
	func() {
		_, _, err := apiClient.UserApi.CreateUser(proxyCtx).CreateUserRequest(*createUsrReq).Execute()
		if err != nil {
			t.Fatalf("Error creating user: %v", err)
		}
	}()
	scopes := []string{
		"socialapp.users.read",
		"socialapp.follower.create",
		"socialapp.follower.read",
		"socialapp.follower.delete",
		"socialapp.comments.create",
		"socialapp.feed.read",
	}
	conf := clientcredentials.Config{
		ClientID:     username1,
		ClientSecret: password,
		Scopes:       scopes,
		TokenURL:     ENDPOINT_OAUTH_TOKEN,
	}
	oauth2Ctx, err := getOuath2Context(proxyCtx, conf)
	if err != nil {
		t.Fatalf("Error getting oauth2 context: %v", err)
	}

	token, res, err := apiClient.AuthenticationApi.GetAccessToken(oauth2Ctx).Execute()
	if err != nil {
		t.Errorf("Error when calling `AuthenticationApi.GetAccessToken`: %v", err)
	}
	// assert scopes are correct
	if res.Status != "200 OK" {
		t.Errorf("Expected status 200, got %s", res.Status)
	}

	if len(token.Scopes) != len(scopes) {
		t.Errorf("Expected %d scopes, got %d", len(scopes), len(token.Scopes))
		t.Log(cmp.Diff(scopes, token.Scopes))
	}
}

func TestRegisterUserFlow(t *testing.T) {
	Setup()
	configuration := client.NewConfiguration()
	//skip cache for tests
	configuration.DefaultHeader["Cache-Control"] = "no-store"
	httpClient := getHTTPClient()
	configuration.HTTPClient = httpClient
	proxyCtx := context.WithValue(context.Background(), oauth2.HTTPClient, httpClient)
	proxyCtx = context.WithValue(proxyCtx, client.ContextServerIndex, CONTEXT_SERVER)

	apiClient = client.NewAPIClient(configuration)

	username1 := fmt.Sprintf("Test-%d1", time.Now().UnixNano())
	password := fmt.Sprintf("Password-%d1", time.Now().UnixNano())
	createUsrReq := client.NewCreateUserRequest(username1, password, "FirstName_example", "LastName_example", username1)

	// create a user, no auth needed
	// POST /user
	// {user}
	_, res, err := apiClient.UserApi.CreateUser(proxyCtx).CreateUserRequest(*createUsrReq).Execute()
	if err != nil {
		t.Errorf("Error when calling `UserApi.CreateUser`: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 201, got %d", res.StatusCode)
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
		t.Fatalf("Error getting oauth2 context: %v", err)
	}

	// Get user by using oauath2 token
	func() {
		_, res, err := apiClient.UserApi.GetUserByUsername(oauth2Ctx, username1).Execute()
		if err != nil {
			t.Errorf("Error when calling `UserApi.GetUsers`: %v", err)
		}
		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", res.StatusCode)
		}
	}()

	// TODO API Should return 401 if no auth is provided
	// validate 403 if no auth is provided
	func() {
		user, res, err := apiClient.UserApi.GetUserByUsername(proxyCtx, username1).Execute()
		if err == nil {
			t.Errorf("Error when calling `UserApi.GetUsers`: %v", err)
		}
		if res.StatusCode != http.StatusForbidden { // TOOD fix to 401
			t.Errorf("Expected status code 401, got %d", res.StatusCode)
		}
		if user != nil {
			t.Errorf("Expected nil user, got %v", user)
		}
	}()
}

func TestChangePassword(t *testing.T) {
	Setup()
	configuration := client.NewConfiguration()
	//skip cache for tests
	configuration.DefaultHeader["Cache-Control"] = "no-store"
	httpClient := getHTTPClient()
	configuration.HTTPClient = httpClient
	proxyCtx := context.WithValue(context.Background(), oauth2.HTTPClient, httpClient)
	proxyCtx = context.WithValue(proxyCtx, client.ContextServerIndex, CONTEXT_SERVER)
	scopes := []string{
		"socialapp.users.read",
		"socialapp.users.update",
	}

	apiClient = client.NewAPIClient(configuration)

	username := fmt.Sprintf("Test-%d1", time.Now().UnixNano())
	password := fmt.Sprintf("Password-%d1", time.Now().UnixNano())
	createUsrReq := client.NewCreateUserRequest(username, password, "FirstName_example", "LastName_example", username)

	// create a user, no auth needed
	_, res, err := apiClient.UserApi.CreateUser(proxyCtx).CreateUserRequest(*createUsrReq).Execute()
	if err != nil {
		t.Errorf("Error when calling `UserApi.CreateUser`: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 201, got %d", res.StatusCode)
	}

	conf := clientcredentials.Config{
		ClientID:     username,
		ClientSecret: password,
		Scopes:       scopes,
		TokenURL:     ENDPOINT_OAUTH_TOKEN,
	}
	oauth2Ctx, err := getOuath2Context(proxyCtx, conf)
	if err != nil {
		t.Fatalf("Error getting oauth2 context: %v", err)
	}

	newPassword := password + "new"
	func() {
		changePwdReq := client.NewChangePasswordRequest(password, newPassword)
		_, res, err := apiClient.UserApi.ChangePassword(oauth2Ctx).ChangePasswordRequest(*changePwdReq).Execute()
		if err != nil {
			t.Fatalf("Error when calling `UserApi.ChangePassword`: %v", err)
		}
		if res.StatusCode != http.StatusOK {
			t.Fatalf("Expected status code 200, got %d", res.StatusCode)
		}
	}()

	// attempt to get token with old password, expect 401
	// func() {
	// 	token, res, err := apiClient.AuthenticationApi.GetAccessToken(oauth2Ctx).Execute()
	// 	if err == nil {
	// 		t.Errorf("Error when calling `AuthenticationApi.GetAccessToken`: %v", err)
	// 	}
	// 	if res.StatusCode != http.StatusUnauthorized {
	// 		t.Errorf("Expected status code 401, got %d", res.StatusCode)
	// 	}
	// 	if token != nil {
	// 		t.Errorf("Expected nil user, got %v", token)
	// 	}
	// }()

	// attempt to get token with new password, expect 200
	func() {
		newPwdConf := clientcredentials.Config{
			ClientID:     username,
			ClientSecret: newPassword,
			Scopes:       scopes,
			TokenURL:     ENDPOINT_OAUTH_TOKEN,
		}
		newPwdOauth2Ctx, err := getOuath2Context(proxyCtx, newPwdConf)
		if err != nil {
			t.Fatalf("Error getting oauth2 context: %v", err)
		}
		token, res, err := apiClient.AuthenticationApi.GetAccessToken(newPwdOauth2Ctx).Execute()
		if err != nil {
			t.Errorf("Error when calling `AuthenticationApi.GetAccessToken`: %v", err)
		}
		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", res.StatusCode)
		}
		if token == nil {
			t.Errorf("Expected token, got nil")
		}
	}()
}

func TestRoleLifecycle(t *testing.T) {
	Setup()
	configuration := client.NewConfiguration()
	//skip cache for tests
	configuration.DefaultHeader["Cache-Control"] = "no-store"
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

	username := fmt.Sprintf("Test-%d1", time.Now().UnixNano())
	password := fmt.Sprintf("Password-%d1", time.Now().UnixNano())
	createUsrReq := client.NewCreateUserRequest(username, password, "FirstName_example", "LastName_example", username)

	// create a user, no auth needed
	_, res, err := apiClient.UserApi.CreateUser(proxyCtx).CreateUserRequest(*createUsrReq).Execute()
	if err != nil {
		t.Errorf("Error when calling `UserApi.CreateUser`: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 201, got %d", res.StatusCode)
	}

	conf := clientcredentials.Config{
		ClientID:     username,
		ClientSecret: password,
		Scopes:       scopes,
		TokenURL:     ENDPOINT_OAUTH_TOKEN,
	}
	oauth2Ctx, err := getOuath2Context(proxyCtx, conf)
	if err != nil {
		t.Fatalf("Error getting oauth2 context: %v", err)
	}

	// create a role
	newRole := client.NewRole(fmt.Sprintf("Test-CreateRole-%d1", time.Now().UnixNano()))
	createdRole, res, err := apiClient.RoleApi.CreateRole(oauth2Ctx).Role(*newRole).Execute()
	if err != nil {
		t.Fatalf("Error when calling `RoleApi.CreateRole`: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", res.StatusCode)
	}

	// get role
	gettedRole, res, err := apiClient.RoleApi.GetRole(oauth2Ctx, int32(*createdRole.Id)).Execute()
	if err != nil {
		t.Fatalf("Error when calling `RoleApi.GetRole`: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", res.StatusCode)
	}
	if gettedRole == nil {
		t.Fatalf("Expected role, got nil")
	}

	// add scopes to role
	// create scope
	newScope := client.NewScope(fmt.Sprintf("Test-CreateScope-%d1", time.Now().UnixNano()), "Test-CreateScope-Description")
	func() {
		createdScope, res, err := apiClient.ScopeApi.CreateScope(oauth2Ctx).Scope(*newScope).Execute()
		if err != nil {
			t.Fatalf("Error when calling `ScopeApi.CreateScope`: %v", err)
		}
		if res.StatusCode != http.StatusOK {
			t.Fatalf("Expected status code 200, got %d", res.StatusCode)
		}
		// attach scope to role
		scopesToAdd := []string{newScope.Name}
		res, err = apiClient.RoleApi.AddScopeToRole(oauth2Ctx, int32(*createdRole.Id)).RequestBody(scopesToAdd).Execute()
		if err != nil {
			t.Fatalf("Error when calling `RoleApi.AddScopeToRole`: %v", err)
		}
		// verify scope is attached to role
		// get role scopes
		roleScopes, res, err := apiClient.RoleApi.ListScopesForRole(oauth2Ctx, int32(*createdRole.Id)).Execute()
		if err != nil {
			t.Fatalf("Error when calling `RoleApi.ListScopesForRole`: %v", err)
		}
		if res.StatusCode != http.StatusOK {
			t.Fatalf("Expected status code 200, got %d", res.StatusCode)
		}
		if len(roleScopes) != 1 {
			t.Fatalf("Expected 1 scope, got %d", len(roleScopes))
		}
		if roleScopes[0].Name != newScope.Name {
			t.Fatalf("Expected scope name %s, got %s", newScope.Name, roleScopes[0].Name)
		}
		// remove scope from role
		res, err = apiClient.RoleApi.RemoveScopeFromRole(oauth2Ctx, int32(*createdRole.Id), int32(*createdScope.Id)).Execute()
		if err != nil {
			t.Fatalf("Error when calling `RoleApi.RemoveScopeFromRole`: %v", err)
		}
		if res.StatusCode != http.StatusNoContent {
			t.Fatalf("Expected status code 204, got %d", res.StatusCode)
		}
		// verify scope is removed from role
		// get role scopes
		roleScopes, res, err = apiClient.RoleApi.ListScopesForRole(oauth2Ctx, int32(*createdRole.Id)).Execute()
		if err != nil {
			t.Fatalf("Error when calling `RoleApi.ListScopesForRole`: %v", err)
		}
		if res.StatusCode != http.StatusOK {
			t.Fatalf("Expected status code 200, got %d", res.StatusCode)
		}
		if len(roleScopes) != 0 {
			t.Fatalf("Expected 0 scopes, got %d", len(roleScopes))
		}
		// detach scope from role
		res, err = apiClient.RoleApi.RemoveScopeFromRole(oauth2Ctx, int32(*createdRole.Id), int32(*createdScope.Id)).Execute()
		if err != nil {
			t.Fatalf("Error when calling `RoleApi.RemoveScopeFromRole`: %v", err)
		}
		if res.StatusCode != http.StatusNoContent {
			t.Fatalf("Expected status code 204, got %d", res.StatusCode)
		}
		// verify scope is detached from role
		// get role scopes
		roleScopes, res, err = apiClient.RoleApi.ListScopesForRole(oauth2Ctx, int32(*createdRole.Id)).Execute()
		if err != nil {
			t.Fatalf("Error when calling `RoleApi.ListScopesForRole`: %v", err)
		}
		if res.StatusCode != http.StatusOK {
			t.Fatalf("Expected status code 200, got %d", res.StatusCode)
		}
		if len(roleScopes) != 0 {
			t.Fatalf("Expected 0 scopes, got %d", len(roleScopes))
		}

		// delete scope
		res, err = apiClient.ScopeApi.DeleteScope(oauth2Ctx, int32(*createdScope.Id)).Execute()
		if err != nil {
			t.Fatalf("Error when calling `ScopeApi.DeleteScope`: %v", err)
		}
		if res.StatusCode != http.StatusOK {
			t.Fatalf("Expected status code 200, got %d", res.StatusCode)
		}
	}()

	updatedName := fmt.Sprintf("Test-UpdateRole-%d", time.Now().UnixNano())
	updatedDesc := fmt.Sprintf("Test-UpdateRole-Description-%d", time.Now().UnixNano())
	newRole.Description = &updatedDesc
	newRole.Name = updatedName

	// update role
	updatedRole, res, err := apiClient.RoleApi.UpdateRole(oauth2Ctx, int32(*createdRole.Id)).Role(*newRole).Execute()
	if err != nil {
		t.Fatalf("Error when calling `RoleApi.UpdateRole`: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", res.StatusCode)
	}
	if updatedRole == nil {
		t.Fatalf("Expected role, got nil")
	}
	if updatedRole.Name != updatedName {
		t.Fatalf("Expected name %s, got %s", updatedName, updatedRole.Name)
	}
	if *updatedRole.Description != updatedDesc {
		t.Fatalf("Expected description %s, got %s", updatedDesc, *updatedRole.Description)
	}

	// get role again to check if updated
	gettedRole, res, err = apiClient.RoleApi.GetRole(oauth2Ctx, int32(*createdRole.Id)).Execute()
	if err != nil {
		t.Fatalf("Error when calling `RoleApi.GetRole`: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", res.StatusCode)
	}
	if gettedRole == nil {
		t.Fatalf("Expected role, got nil")
	}
	if gettedRole.Name != updatedName {
		t.Fatalf("Expected name %s, got %s", updatedName, gettedRole.Name)
	}
	if *gettedRole.Description != updatedDesc {
		t.Fatalf("Expected description %s, got %s", updatedDesc, *gettedRole.Description)
	}

	// delete role
	func() {
		res, err := apiClient.RoleApi.DeleteRole(oauth2Ctx, int32(*createdRole.Id)).Execute()
		if err != nil {
			t.Fatalf("Error when calling `RoleApi.DeleteRole`: %v", err)
		}
		if res.StatusCode != http.StatusOK {
			t.Fatalf("Expected status code 200, got %d", res.StatusCode)
		}
	}()

	// try to get deleted role
	gettedDeletedRole, res, err := apiClient.RoleApi.GetRole(oauth2Ctx, int32(*createdRole.Id)).Execute()
	if err == nil {
		t.Fatalf("Expected error when calling `RoleApi.GetRole`, got nil")
	}
	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("Expected status code 404, got %d", res.StatusCode)
	}
	if gettedDeletedRole != nil {
		t.Fatalf("Expected nil, got %v", gettedRole)
	}
}

func TestScopeLifeCycle(t *testing.T) {
	Setup()
	configuration := client.NewConfiguration()
	//skip cache for tests
	configuration.DefaultHeader["Cache-Control"] = "no-store"
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

	username := fmt.Sprintf("Test-%d1", time.Now().UnixNano())
	password := fmt.Sprintf("Password-%d1", time.Now().UnixNano())
	createUsrReq := client.NewCreateUserRequest(username, password, "FirstName_example", "LastName_example", username)

	// create a user, no auth needed
	_, res, err := apiClient.UserApi.CreateUser(proxyCtx).CreateUserRequest(*createUsrReq).Execute()
	if err != nil {
		t.Errorf("Error when calling `UserApi.CreateUser`: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 201, got %d", res.StatusCode)
	}

	conf := clientcredentials.Config{
		ClientID:     username,
		ClientSecret: password,
		Scopes:       scopes,
		TokenURL:     ENDPOINT_OAUTH_TOKEN,
	}
	oauth2Ctx, err := getOuath2Context(proxyCtx, conf)
	if err != nil {
		t.Fatalf("Error getting oauth2 context: %v", err)
	}

	newScope := client.NewScope(fmt.Sprintf("Test-CreateScope-%d1", time.Now().UnixNano()), "Test-CreateScope-Description1")
	// create a scope

	createdScope, res, err := apiClient.ScopeApi.CreateScope(oauth2Ctx).Scope(*newScope).Execute()
	if err != nil {
		t.Fatalf("Error when calling `ScopeApi.CreateScope`: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", res.StatusCode)
	}

	// get scope
	gettedScope, res, err := apiClient.ScopeApi.GetScope(oauth2Ctx, int32(*createdScope.Id)).Execute()
	if err != nil {
		t.Fatalf("Error when calling `ScopeApi.GetScope`: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", res.StatusCode)
	}
	if gettedScope == nil {
		t.Fatalf("Expected scope, got nil")
	}

	updatedName := fmt.Sprintf("Test-UpdateScope-%d", time.Now().UnixNano())
	updatedDesc := fmt.Sprintf("Test-UpdateScope-Description-%d", time.Now().UnixNano())
	newScope.Description = updatedDesc
	newScope.Name = updatedName

	// update scope
	updatedScope, res, err := apiClient.ScopeApi.UpdateScope(oauth2Ctx, int32(*createdScope.Id)).Scope(*newScope).Execute()
	if err != nil {
		t.Fatalf("Error when calling `ScopeApi.UpdateScope`: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", res.StatusCode)
	}
	if updatedScope == nil {
		t.Fatalf("Expected scope, got nil")
	}
	if updatedScope.Name != updatedName {
		t.Fatalf("Expected name %s, got %s", updatedName, updatedScope.Name)
	}
	if updatedScope.Description != updatedDesc {
		t.Fatalf("Expected description %s, got %s", updatedDesc, updatedScope.Description)
	}

	// get scope again to check if updated
	gettedScope, res, err = apiClient.ScopeApi.GetScope(oauth2Ctx, int32(*createdScope.Id)).Execute()
	if err != nil {
		t.Fatalf("Error when calling `ScopeApi.GetScope`: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", res.StatusCode)
	}
	if gettedScope == nil {
		t.Fatalf("Expected scope, got nil")
	}
	if gettedScope.Name != updatedName {
		t.Fatalf("Expected name %s, got %s", updatedName, gettedScope.Name)
	}
	if gettedScope.Description != updatedDesc {
		t.Fatalf("Expected description %s, got %s", updatedDesc, gettedScope.Description)
	}

	// delete scope
	func() {
		res, err := apiClient.ScopeApi.DeleteScope(oauth2Ctx, int32(*createdScope.Id)).Execute()
		if err != nil {
			t.Fatalf("Error when calling `ScopeApi.DeleteScope`: %v", err)
		}
		if res.StatusCode != http.StatusOK {
			t.Fatalf("Expected status code 204, got %d", res.StatusCode)
		}
	}()

	// try to get deleted scope
	gettedDeletedScope, res, err := apiClient.ScopeApi.GetScope(oauth2Ctx, int32(*createdScope.Id)).Execute()
	if err == nil {
		t.Fatalf("Expected error when calling `ScopeApi.GetScope`, got nil")
	}
	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("Expected status code 404, got %d", res.StatusCode)
	}
	if gettedDeletedScope != nil {
		t.Fatalf("Expected nil, got %v", gettedScope)
	}
}

func TestUserRoleLifeCycle(t *testing.T) {
	Setup()
	configuration := client.NewConfiguration()
	//skip cache for tests
	configuration.DefaultHeader["Cache-Control"] = "no-store"
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

	username := fmt.Sprintf("Test-%d1", time.Now().UnixNano())
	password := fmt.Sprintf("Password-%d1", time.Now().UnixNano())
	createUsrReq := client.NewCreateUserRequest(username, password, "FirstName_example", "LastName_example", username)

	// create a user, no auth needed
	_, res, err := apiClient.UserApi.CreateUser(proxyCtx).CreateUserRequest(*createUsrReq).Execute()
	if err != nil {
		t.Errorf("Error when calling `UserApi.CreateUser`: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 201, got %d", res.StatusCode)
	}

	conf := clientcredentials.Config{
		ClientID:     username,
		ClientSecret: password,
		Scopes:       scopes,
		TokenURL:     ENDPOINT_OAUTH_TOKEN,
	}
	oauth2Ctx, err := getOuath2Context(proxyCtx, conf)
	if err != nil {
		t.Fatalf("Error getting oauth2 context: %v", err)
	}

	newScope := client.NewScope(fmt.Sprintf("Test-CreateScope-%d1", time.Now().UnixNano()), "Test-CreateScope-Description1")

	// create a scope
	cretedScope, res, err := apiClient.ScopeApi.CreateScope(oauth2Ctx).Scope(*newScope).Execute()
	if err != nil {
		t.Fatalf("Error when calling `ScopeApi.CreateScope`: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", res.StatusCode)
	}
	t.Cleanup(func() {
		// delete scope
		res, err := apiClient.ScopeApi.DeleteScope(oauth2Ctx, int32(*cretedScope.Id)).Execute()
		if err != nil {
			t.Fatalf("Error when calling `ScopeApi.DeleteScope`: %v", err)
		}
		if res.StatusCode != http.StatusOK {
			t.Fatalf("Expected status code 204, got %d", res.StatusCode)
		}
	})

	// create role
	newRole := client.NewRole(fmt.Sprintf("Test-CreateRole-%d1", time.Now().UnixNano()))
	createdRole, res, err := apiClient.RoleApi.CreateRole(oauth2Ctx).Role(*newRole).Execute()
	if err != nil {
		t.Fatalf("Error when calling `RoleApi.CreateRole`: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", res.StatusCode)
	}

	// attach role to user
	_, res, err = apiClient.UserApi.UpdateRolesForUser(oauth2Ctx, username).RequestBody([]string{createdRole.Name}).Execute()
	if err != nil {
		t.Fatalf("Error when calling `UserApi.UpdateRolesForUser`: %v", err)
	}

	t.Cleanup(func() {
		// delete role
		res, err := apiClient.RoleApi.DeleteRole(oauth2Ctx, int32(*createdRole.Id)).Execute()
		if err != nil {
			t.Fatalf("Error when calling `RoleApi.DeleteRole`: %v", err)
		}
		if res.StatusCode != http.StatusOK {
			t.Fatalf("Expected status code 204, got %d", res.StatusCode)
		}
	})

	// List user roles
	roles, res, err := apiClient.UserApi.GetRolesForUser(oauth2Ctx, username).Execute()
	if err != nil {
		t.Fatalf("Error when calling `UserApi.GetRolesForUser`: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", res.StatusCode)
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
		t.Fatalf("Expected to find role %s, got %v", createdRole.Name, roles)
	}
}

func TestCacheRequestSameUser(t *testing.T) {
	Setup()
	configuration := client.NewConfiguration()
	httpClient := getHTTPClient()
	configuration.HTTPClient = httpClient

	proxyCtx := context.WithValue(context.Background(), oauth2.HTTPClient, httpClient)
	proxyCtx = context.WithValue(proxyCtx, client.ContextServerIndex, CONTEXT_SERVER)

	username1 := fmt.Sprintf("Test-%d1", time.Now().UnixNano())
	password := fmt.Sprintf("Password-%d1", time.Now().UnixNano())
	apiClient = client.NewAPIClient(configuration)
	func() {
		createUsrReq := client.NewCreateUserRequest(username1, password, "FirstName_example", "LastName_example", username1)
		_, _, err := apiClient.UserApi.CreateUser(proxyCtx).CreateUserRequest(*createUsrReq).Execute()
		if err != nil {
			t.Fatalf("Error creating user: %v", err)
		}
	}()

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
		t.Fatalf("Error getting oauth2 context: %v", err)
	}
	openAPICtx := context.WithValue(oauth2Ctx, client.ContextServerIndex, CONTEXT_SERVER)

	// List 100 users, different offset on every execution
	offset := time.Now().UnixNano() % 5000
	listedUsers, r, err := apiClient.UserApi.
		ListUsers(openAPICtx).
		Limit(5).
		Offset(int32(offset)).
		Execute()

	if err != nil {
		t.Errorf("Error when calling `UserApi.ListUsers`: %v\n", err)
		t.Errorf("Full HTTP response: %v\n", r)
	}
	if r.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, r.StatusCode)
	}

	// get user info 5 times
	for _, currentUser := range listedUsers {
		for i := 0; i < 50; i++ {
			_, r, err = apiClient.UserApi.GetUserByUsername(openAPICtx, currentUser.Username).Execute()
			if err != nil {
				t.Errorf("Error when calling `UserApi.GetUser`: %v\n", err)
				t.Errorf("Full HTTP response: %v", r)
			}
			if r.StatusCode != http.StatusOK {
				t.Errorf("Expected status code %d, got %d", http.StatusOK, r.StatusCode)
			}

			// get user comments
			_, r, err = apiClient.CommentApi.GetUserComments(openAPICtx, currentUser.Username).Execute()
			if err != nil {
				t.Errorf("Error when calling `CommentApi.GetUserComments`: %v\n", err)
				t.Errorf("Full HTTP response: %v", r)
			}
			if r.StatusCode != http.StatusOK {
				t.Errorf("Expected status code %d, got %d", http.StatusOK, r.StatusCode)
			}

			// get feed
			_, r, err = apiClient.CommentApi.GetUserFeed(openAPICtx).Execute()
			if err != nil {
				t.Errorf("Error when calling `CommentApi.GetUserFeed`: %v\n", err)
				t.Errorf("Full HTTP response: %v", r)
			}
			if r.StatusCode != http.StatusOK {
				t.Errorf("Expected status code %d, got %d", http.StatusOK, r.StatusCode)
			}

			// t.Logf("%q: %d/100\n", currentUser.Username, i)
		}
	}
}

func TestURLLifeCycle(t *testing.T) {
	Setup()
	configuration := client.NewConfiguration()
	httpClient := getHTTPClient()
	configuration.HTTPClient = httpClient

	proxyCtx := context.WithValue(context.Background(), oauth2.HTTPClient, httpClient)
	proxyCtx = context.WithValue(proxyCtx, client.ContextServerIndex, CONTEXT_SERVER)

	username1 := fmt.Sprintf("Test-%d1", time.Now().UnixNano())
	password := fmt.Sprintf("Password-%d1", time.Now().UnixNano())
	apiClient = client.NewAPIClient(configuration)
	func() {
		createUsrReq := client.NewCreateUserRequest(username1, password, "FirstName_example", "LastName_example", username1)
		_, _, err := apiClient.UserApi.CreateUser(proxyCtx).CreateUserRequest(*createUsrReq).Execute()
		if err != nil {
			t.Fatalf("Error creating user: %v", err)
		}
	}()

	conf := clientcredentials.Config{
		ClientID:     username1,
		ClientSecret: password,
		Scopes:       []string{"shortly.url.create", "shortly.url.delete"},
		TokenURL:     ENDPOINT_OAUTH_TOKEN,
	}

	oauth2Ctx, err := getOuath2Context(proxyCtx, conf)
	if err != nil {
		t.Fatalf("Error getting oauth2 context: %v", err)
	}
	openAPICtx := context.WithValue(oauth2Ctx, client.ContextServerIndex, CONTEXT_SERVER)

	// create url
	newURL := client.NewURL("https://www.google.com", "google")
	_, r, err := apiClient.URLApi.CreateUrl(openAPICtx).URL(*newURL).Execute()
	if err != nil {
		t.Errorf("Error when calling `URLApi.CreateURL`: %v\n", err)
		t.Errorf("Full HTTP response: %v ", r)
	}

	if r.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, r.StatusCode)
	}

	// get url
	getUrlRes, err := apiClient.URLApi.GetUrl(proxyCtx, "google").Execute()
	if err != nil {
		t.Errorf("Error when calling `URLApi.GetURL`: %v\n", err)
		t.Errorf("Full HTTP response: %v ", r)
		t.Fatalf("Error getting url: %v", err)
	}

	if getUrlRes.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, getUrlRes.StatusCode)
	}

	// delete url
	deleteUrlRes, err := apiClient.URLApi.DeleteUrl(openAPICtx, "google").Execute()
	if err != nil {
		t.Errorf("Error when calling `URLApi.DeleteURL`: %v\n", err)
		t.Errorf("Full HTTP response: %v ", r)
		t.Fatalf("Error deleting url: %v", err)
	}

	if deleteUrlRes.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, deleteUrlRes.StatusCode)
	}

	// get url
	getUrlRes, err = apiClient.URLApi.GetUrl(proxyCtx, "google").Execute()
	if err == nil {
		t.Errorf("Expected error when calling `URLApi.GetURL`, got none")
	}
	if getUrlRes.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, getUrlRes.StatusCode)
	}

}
