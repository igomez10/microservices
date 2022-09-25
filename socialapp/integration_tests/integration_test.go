package integration_tests

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"socialapp/client"
	"testing"
	"time"
)

var apiClient *client.APIClient

var (
	RENDER_SERVER_URL          = 0
	LOCALHOST_SERVER_URL       = 1
	LOCALHOST_DEBUG_SERVER_URL = 2

	CONTEXT_SERVER = LOCALHOST_SERVER_URL
)

func TestListUsers(t *testing.T) {
	os.Setenv("HTTP_PROXY", "http://localhost:9091")
	os.Setenv("HTTPS_PROXY", "http://localhost:9091")

	configuration := client.NewConfiguration()
	proxyStr := "http://localhost:9091"
	proxyURL, err := url.Parse(proxyStr)
	if err != nil {
		log.Println(err)
	}

	//adding the proxy settings to the Transport object
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
	configuration.HTTPClient = &http.Client{
		Timeout:   time.Second * 10,
		Transport: transport,
	}

	apiClient = client.NewAPIClient(configuration)
	ctx := context.WithValue(context.Background(), client.ContextServerIndex, CONTEXT_SERVER)
	basicAuthCtx := context.WithValue(ctx, client.ContextBasicAuth, client.BasicAuth{
		UserName: "admin",
		Password: "admin",
	})

	// get access token
	token, r, err := apiClient.AuthenticationApi.GetAccessToken(basicAuthCtx).Execute()
	if err != nil {
		t.Errorf("Error when calling `AuthenticationApi.GetAccessToken`: %v\n %+v\n", err, r)
	}
	if r.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, r.StatusCode)
	}

	beaererCtx := context.WithValue(ctx, client.ContextAccessToken, token.AccessToken)
	// List users
	func() {
		_, r, err := apiClient.UserApi.ListUsers(beaererCtx).Execute()
		if err != nil {
			t.Errorf("Error when calling `UserApi.ListUsers``: %v\n", err)
			t.Errorf("Full HTTP response: %v\n", r)
		}
		if r.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, r.StatusCode)
		}
	}()
}

func TestCreateUser(t *testing.T) {
	os.Setenv("HTTP_PROXY", "http://localhost:9091")
	os.Setenv("HTTPS_PROXY", "http://localhost:9091")

	configuration := client.NewConfiguration()
	proxyStr := "http://localhost:9091"
	proxyURL, err := url.Parse(proxyStr)
	if err != nil {
		t.Error(err)
	}

	//adding the proxy settings to the Transport object
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
	configuration.HTTPClient = &http.Client{
		Timeout:   time.Second * 10,
		Transport: transport,
	}

	apiClient = client.NewAPIClient(configuration)

	username := fmt.Sprintf("Test-%d", time.Now().UnixNano())
	email := fmt.Sprintf("Test-%d-@social.com", time.Now().UnixNano())
	user := *client.NewCreateUserRequest(username, "password", "FirstName_example", "LastName_example", email) // User | Create a new user
	ctx := context.WithValue(context.Background(), client.ContextServerIndex, CONTEXT_SERVER)
	basicAuthCtx := context.WithValue(ctx, client.ContextBasicAuth, client.BasicAuth{
		UserName: "admin",
		Password: "admin",
	})

	// get access token
	token, r, err := apiClient.AuthenticationApi.GetAccessToken(basicAuthCtx).Execute()
	if err != nil {
		t.Errorf("Error when calling `AuthenticationApi.GetAccessToken`: %v\n %+v\n", err, r)
	}
	if r.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, r.StatusCode)
	}

	beaererCtx := context.WithValue(ctx, client.ContextAccessToken, token.AccessToken)

	// verify a user doesnt exist yet
	func() {
		_, r, err := apiClient.UserApi.GetUserByUsername(beaererCtx, username).Execute()
		if err == nil {
			t.Errorf("User %s already exists: %+v", username, r)
		}
		if r.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status code %d, got %d", http.StatusNotFound, r.StatusCode)
		}
	}()

	func() {
		_, r, err := apiClient.UserApi.CreateUser(ctx).
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
		basicAuthCtx := context.WithValue(ctx, client.ContextBasicAuth, client.BasicAuth{
			UserName: user.Username,
			Password: user.Password,
		})

		// get access token
		token, r, err := apiClient.AuthenticationApi.GetAccessToken(basicAuthCtx).Execute()
		if err != nil {
			t.Errorf("Error when calling `AuthenticationApi.GetAccessToken`: %v\n %+v\n", err, r)
		}
		if r.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, r.StatusCode)
		}

		beaererCtx := context.WithValue(ctx, client.ContextAccessToken, token.AccessToken)

		resp, r, err := apiClient.UserApi.GetUserByUsername(beaererCtx, username).Execute()
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
}

func TestFollowCycle(t *testing.T) {
	// create two users
	os.Setenv("HTTP_PROXY", "http://localhost:9091")
	os.Setenv("HTTPS_PROXY", "http://localhost:9091")

	configuration := client.NewConfiguration()
	proxyStr := "http://localhost:9091"
	proxyURL, err := url.Parse(proxyStr)
	if err != nil {
		log.Println(err)
	}

	//adding the proxy settings to the Transport object
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
	configuration.HTTPClient = &http.Client{
		Timeout:   time.Second * 10,
		Transport: transport,
	}

	apiClient = client.NewAPIClient(configuration)
	ctx := context.WithValue(context.Background(), client.ContextServerIndex, CONTEXT_SERVER)

	username1 := fmt.Sprintf("Test-%d1", time.Now().UnixNano())
	email1 := fmt.Sprintf("Test-%d-1@social.com", time.Now().UnixNano())
	user1 := *client.NewCreateUserRequest(username1, "password", "FirstName_example", "LastName_example", email1) // User | Create a new user

	username2 := fmt.Sprintf("Test-%d2", time.Now().UnixNano())
	email2 := fmt.Sprintf("Test-%d-2@social.com", time.Now().UnixNano())
	user2 := *client.NewCreateUserRequest(username2, "secretPassword", "FirstName_example", "LastName_example", email2) // User | Create a new user

	// create users
	func() {
		_, r1, err1 := apiClient.UserApi.CreateUser(ctx).
			CreateUserRequest(user1).
			Execute()
		if err1 != nil {
			t.Errorf("Error when calling `UserApi.CreateUser`: %v\n %+v\n", err1, r1)
		}

		_, r2, err2 := apiClient.UserApi.CreateUser(ctx).
			CreateUserRequest(user2).
			Execute()
		if err2 != nil {
			t.Errorf("Error when calling `UserApi.CreateUser`: %v\n %+v\n", err2, r2)
		}
	}()

	// get access token for user1
	basicAuthCtx := context.WithValue(ctx, client.ContextBasicAuth, client.BasicAuth{
		UserName: user1.Username,
		Password: user1.Password,
	})

	// get access token
	token, r, err := apiClient.AuthenticationApi.GetAccessToken(basicAuthCtx).Execute()
	if err != nil {
		t.Errorf("Error when calling `AuthenticationApi.GetAccessToken`: %v\n %+v\n", err, r)
	}
	if r.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, r.StatusCode)
	}

	beaererCtx := context.WithValue(ctx, client.ContextAccessToken, token.AccessToken)

	// user 1 follows user 2
	func() {
		r, err := apiClient.UserApi.FollowUser(beaererCtx, username2, username1).Execute()
		if err != nil {
			t.Errorf("Error when calling `UserApi.FollowUser`: %v\n %+v\n", err, r)
		}
	}()

	// validate user 1 follows user 2
	func() {
		followers, r, err := apiClient.UserApi.GetUserFollowers(beaererCtx, username2).Execute()
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
		r, err := apiClient.UserApi.UnfollowUser(beaererCtx, username2, username1).Execute()
		if err != nil {
			t.Errorf("Error when calling `UserApi.FollowUser`: %v\n %+v\n", err, r)
		}
	}()

	// validate user 1 unfollows user 2
	func() {
		followers, r, err := apiClient.UserApi.GetUserFollowers(beaererCtx, username2).Execute()
		if err != nil {
			t.Errorf("Error when calling `UserApi.FollowUser`: %v\n %+v\n", err, r)
		}

		if len(followers) != 0 {
			t.Errorf("Expected 0 followers, got %d", len(followers))
		}
	}()
}

func TestGetExpectedFeed(t *testing.T) {
	// create two users
	os.Setenv("HTTP_PROXY", "http://localhost:9091")
	os.Setenv("HTTPS_PROXY", "http://localhost:9091")

	configuration := client.NewConfiguration()
	proxyStr := "http://localhost:9091"
	proxyURL, err := url.Parse(proxyStr)
	if err != nil {
		log.Println(err)
	}

	//adding the proxy settings to the Transport object
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
	configuration.HTTPClient = &http.Client{
		Timeout:   time.Second * 10,
		Transport: transport,
	}

	apiClient = client.NewAPIClient(configuration)

	username1 := fmt.Sprintf("Test-%d1", time.Now().UnixNano())
	email1 := fmt.Sprintf("Test-%d-1@social.com", time.Now().UnixNano())
	user1 := *client.NewCreateUserRequest(username1, "password", "FirstName_example", "LastName_example", email1) // User | Create a new user

	username2 := fmt.Sprintf("Test-%d2", time.Now().UnixNano())
	email2 := fmt.Sprintf("Test-%d-2@social.com", time.Now().UnixNano())
	user2 := *client.NewCreateUserRequest(username2, "secretPassword", "FirstName_example", "LastName_example", email2) // User | Create a new user

	ctx := context.WithValue(context.Background(), client.ContextServerIndex, CONTEXT_SERVER)

	// create users
	func() {
		_, r1, err1 := apiClient.UserApi.
			CreateUser(ctx).
			CreateUserRequest(user1).
			Execute()
		if err1 != nil {
			t.Errorf("Error when calling `UserApi.CreateUser`: %v\n %+v\n", err1, r1)
		}

		_, r2, err2 := apiClient.UserApi.
			CreateUser(ctx).
			CreateUserRequest(user2).
			Execute()
		if err2 != nil {
			t.Errorf("Error when calling `UserApi.CreateUser`: %v\n %+v\n", err2, r2)
		}
	}()

	// get access token for user1
	// get access token for user1
	basicAuthCtx := context.WithValue(ctx, client.ContextBasicAuth, client.BasicAuth{
		UserName: user1.Username,
		Password: user1.Password,
	})

	// get access token
	token, r, err := apiClient.AuthenticationApi.GetAccessToken(basicAuthCtx).Execute()
	if err != nil {
		t.Errorf("Error when calling `AuthenticationApi.GetAccessToken`: %v\n %+v\n", err, r)
	}
	if r.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, r.StatusCode)
	}

	beaererCtx := context.WithValue(ctx, client.ContextAccessToken, token.AccessToken)

	// user 1 follows user 2
	func() {
		r, err := apiClient.UserApi.FollowUser(
			beaererCtx,
			username2,
			username1).
			Execute()
		if err != nil {
			t.Errorf("Error when calling `UserApi.FollowUser`: %v\n %+v\n", err, r)
		}
	}()

	// user 2 posts a comment
	func() {
		comment := *client.NewComment("Test comment", username2)
		_, r, err := apiClient.CommentApi.
			CreateComment(beaererCtx).
			Comment(comment).
			Execute()
		if err != nil {
			t.Errorf("Error when calling `UserApi.PostComment`: %v\n %+v\n", err, r)
		}
	}()

	// validate feed in user 1's feed
	func() {
		feed, r, err := apiClient.CommentApi.
			GetUserFeed(beaererCtx, username1).
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

}

func TestGetAccessToken(t *testing.T) {
	// create two users
	os.Setenv("HTTP_PROXY", "http://localhost:9091")
	os.Setenv("HTTPS_PROXY", "http://localhost:9091")

	configuration := client.NewConfiguration()

	proxyStr := "http://localhost:9091"
	proxyURL, err := url.Parse(proxyStr)
	if err != nil {
		log.Println(err)
	}

	//adding the proxy settings to the Transport object
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
	configuration.HTTPClient = &http.Client{
		Timeout:   time.Second * 10,
		Transport: transport,
	}

	apiClient = client.NewAPIClient(configuration)
	ctx := context.WithValue(context.Background(), client.ContextServerIndex, CONTEXT_SERVER)
	apiClient.AuthenticationApi.GetAccessToken(ctx)
}

func TestRegisterUserFlow(t *testing.T) {
	// create two users
	os.Setenv("HTTP_PROXY", "http://localhost:9091")
	os.Setenv("HTTPS_PROXY", "http://localhost:9091")

	configuration := client.NewConfiguration()
	proxyStr := "http://localhost:9091"
	proxyURL, err := url.Parse(proxyStr)
	if err != nil {
		log.Println(err)
	}

	// adding the proxy settings to the Transport object
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
	configuration.HTTPClient = &http.Client{
		Timeout:   time.Second * 10,
		Transport: transport,
	}
	urlContext := context.WithValue(context.Background(), client.ContextServerIndex, CONTEXT_SERVER)
	apiClient = client.NewAPIClient(configuration)

	username1 := fmt.Sprintf("Test-%d1", time.Now().UnixNano())
	password := fmt.Sprintf("Password-%d1", time.Now().UnixNano())
	createUsrReq := client.NewCreateUserRequest(username1, password, "FirstName_example", "LastName_example", username1)

	// create a user, no auth needed
	// POST /user
	// {user}
	_, res, err := apiClient.UserApi.CreateUser(urlContext).CreateUserRequest(*createUsrReq).Execute()
	if err != nil {
		t.Errorf("Error when calling `UserApi.CreateUser`: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 201, got %d", res.StatusCode)
	}

	// use basic auth to get a beaer token
	basicAuthContext := context.WithValue(urlContext, client.ContextBasicAuth, client.BasicAuth{
		UserName: username1,
		Password: password,
	})
	token, res, err := apiClient.AuthenticationApi.GetAccessToken(basicAuthContext).Execute()
	if err != nil {
		t.Errorf("Error when calling `AuthenticationApi.GetAccessToken`: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", res.StatusCode)
	}

	// use bearertoken to get user info
	func() {
		bearerTokenContext := context.WithValue(urlContext, client.ContextAccessToken, token.AccessToken)
		_, res, err := apiClient.UserApi.GetUserByUsername(bearerTokenContext, username1).Execute()
		if err != nil {
			t.Errorf("Error when calling `UserApi.GetUsers`: %v", err)
		}
		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", res.StatusCode)
		}
	}()

	// validate 401
	func() {
		user, res, err := apiClient.UserApi.GetUserByUsername(urlContext, username1).Execute()
		if err == nil {
			t.Errorf("Error when calling `UserApi.GetUsers`: %v", err)
		}
		if res.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status code 401, got %d", res.StatusCode)
		}
		if user != nil {
			t.Errorf("Expected nil user, got %v", user)
		}
	}()
}

func TestChangePassword(t *testing.T) {
	os.Setenv("HTTP_PROXY", "http://localhost:9091")
	os.Setenv("HTTPS_PROXY", "http://localhost:9091")
	configuration := client.NewConfiguration()

	proxyStr := "http://localhost:9091"
	proxyURL, err := url.Parse(proxyStr)
	if err != nil {
		log.Println(err)
	}

	// adding the proxy settings to the Transport object
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
	configuration.HTTPClient = &http.Client{
		Timeout:   time.Second * 10,
		Transport: transport,
	}
	urlContext := context.WithValue(context.Background(), client.ContextServerIndex, CONTEXT_SERVER)
	apiClient = client.NewAPIClient(configuration)

	username := fmt.Sprintf("Test-%d1", time.Now().UnixNano())
	password := fmt.Sprintf("Password-%d1", time.Now().UnixNano())
	createUsrReq := client.NewCreateUserRequest(username, password, "FirstName_example", "LastName_example", username)

	// create a user, no auth needed
	_, res, err := apiClient.UserApi.CreateUser(urlContext).CreateUserRequest(*createUsrReq).Execute()
	if err != nil {
		t.Errorf("Error when calling `UserApi.CreateUser`: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 201, got %d", res.StatusCode)
	}

	// use basic auth to get a beaer token
	basicAuthContext := context.WithValue(urlContext, client.ContextBasicAuth, client.BasicAuth{
		UserName: username,
		Password: password,
	})
	token, res, err := apiClient.AuthenticationApi.GetAccessToken(basicAuthContext).Execute()
	if err != nil {
		t.Errorf("Error when calling `AuthenticationApi.GetAccessToken`: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", res.StatusCode)
	}

	bearerTokenContext := context.WithValue(urlContext, client.ContextAccessToken, token.AccessToken)

	newPassword := password + "new"
	func() {
		changePwdReq := client.NewChangePasswordRequest(password, newPassword)
		_, res, err := apiClient.UserApi.ChangePassword(bearerTokenContext).ChangePasswordRequest(*changePwdReq).Execute()
		if err != nil {
			t.Errorf("Error when calling `UserApi.ChangePassword`: %v", err)
		}
		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", res.StatusCode)
		}
	}()

	// attempt to get token with old password, expect 401
	func() {
		basicAuthContext := context.WithValue(urlContext, client.ContextBasicAuth, client.BasicAuth{
			UserName: username,
			Password: password,
		})
		token, res, err := apiClient.AuthenticationApi.GetAccessToken(basicAuthContext).Execute()
		if err == nil {
			t.Errorf("Error when calling `AuthenticationApi.GetAccessToken`: %v", err)
		}
		if res.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status code 401, got %d", res.StatusCode)
		}
		if token != nil {
			t.Errorf("Expected nil user, got %v", token)
		}
	}()

	// attempt to get token with new password, expect 200
	func() {
		basicAuthContext := context.WithValue(urlContext, client.ContextBasicAuth, client.BasicAuth{
			UserName: username,
			Password: newPassword,
		})
		token, res, err := apiClient.AuthenticationApi.GetAccessToken(basicAuthContext).Execute()
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
