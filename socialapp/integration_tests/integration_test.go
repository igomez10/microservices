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

func TestListUsers(t *testing.T) {
	os.Setenv("HTTP_PROXY", "http://192.168.68.103:9091")
	os.Setenv("HTTPS_PROXY", "http://192.168.68.103:9091")

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
	// List users
	resp, r, err := apiClient.UserApi.ListUsers(context.Background()).Execute()
	if err != nil {
		t.Errorf("Error when calling `UserApi.ListUsers``: %v\n", err)
		t.Errorf("Full HTTP response: %v\n", r)
	}
	// response from `ListUsers`: []User
	t.Logf("Response from `UserApi.ListUsers`: %v\n", resp)
}

func TestCreateUser(t *testing.T) {
	os.Setenv("HTTP_PROXY", "localhost:9091")
	os.Setenv("HTTPS_PROXY", "localhost:9091")

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
	user := *client.NewUser(username, "FirstName_example", "LastName_example", email) // User | Create a new user

	// verify a user doesnt exist yet
	func() {
		_, r, err := apiClient.UserApi.GetUserByUsername(context.Background(), username).Execute()
		if err == nil {
			t.Errorf("User %s already exists: %+v", username, r)
		}
		if r.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status code %d, got %d", http.StatusNotFound, r.StatusCode)
		}
	}()

	func() {
		_, r, err := apiClient.UserApi.CreateUser(context.Background()).User(user).Execute()
		if err != nil {
			t.Errorf("Error when calling `UserApi.CreateUser`: %v\n %+v\n", err, r)
		}
		if r.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, r.StatusCode)
		}
	}()

	func() {
		resp, r, err := apiClient.UserApi.GetUserByUsername(context.Background(), username).Execute()
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
	os.Setenv("HTTP_PROXY", "localhost:9091")
	os.Setenv("HTTPS_PROXY", "localhost:9091")

	os.Setenv("HTTP_PROXY", "http://192.168.68.103:9091")
	os.Setenv("HTTPS_PROXY", "http://192.168.68.103:9091")

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
	user1 := *client.NewUser(username1, "FirstName_example", "LastName_example", email1) // User | Create a new user

	username2 := fmt.Sprintf("Test-%d2", time.Now().UnixNano())
	email2 := fmt.Sprintf("Test-%d-2@social.com", time.Now().UnixNano())
	user2 := *client.NewUser(username2, "FirstName_example", "LastName_example", email2) // User | Create a new user

	// create users
	func() {
		_, r1, err1 := apiClient.UserApi.CreateUser(context.Background()).User(user1).Execute()
		if err1 != nil {
			t.Errorf("Error when calling `UserApi.CreateUser`: %v\n %+v\n", err1, r1)
			fmt.Fprintf(os.Stderr, "Error when calling `UserApi.CreateUser``: %v\n", err1)
			fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r1)
		}

		_, r2, err2 := apiClient.UserApi.CreateUser(context.Background()).User(user2).Execute()
		if err2 != nil {
			t.Errorf("Error when calling `UserApi.CreateUser`: %v\n %+v\n", err2, r2)
			fmt.Fprintf(os.Stderr, "Error when calling `UserApi.CreateUser``: %v\n", err2)
			fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r2)
		}
	}()

	// user 1 follows user 2
	func() {
		r, err := apiClient.UserApi.FollowUser(context.Background(), username2, username1).Execute()
		if err != nil {
			t.Errorf("Error when calling `UserApi.FollowUser`: %v\n %+v\n", err, r)
			fmt.Fprintf(os.Stderr, "Error when calling `UserApi.FollowUser``: %v\n", err)
			fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
		}
	}()

	// validate user 1 follows user 2
	func() {
		followers, r, err := apiClient.UserApi.GetUserFollowers(context.Background(), username2).Execute()
		if err != nil {
			t.Errorf("Error when calling `UserApi.FollowUser`: %v\n %+v\n", err, r)
			fmt.Fprintf(os.Stderr, "Error when calling `UserApi.FollowUser``: %v\n", err)
			fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
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
		r, err := apiClient.UserApi.UnfollowUser(context.Background(), username2, username1).Execute()
		if err != nil {
			t.Errorf("Error when calling `UserApi.FollowUser`: %v\n %+v\n", err, r)
			fmt.Fprintf(os.Stderr, "Error when calling `UserApi.FollowUser``: %v\n", err)
			fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
		}
	}()

	// validate user 1 unfollows user 2
	func() {
		followers, r, err := apiClient.UserApi.GetUserFollowers(context.Background(), username2).Execute()
		if err != nil {
			t.Errorf("Error when calling `UserApi.FollowUser`: %v\n %+v\n", err, r)
			fmt.Fprintf(os.Stderr, "Error when calling `UserApi.FollowUser``: %v\n", err)
			fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
		}

		if len(followers) != 0 {
			t.Errorf("Expected 0 followers, got %d", len(followers))
		}
	}()
}
