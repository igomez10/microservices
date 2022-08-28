package integration_tests

import (
	"context"
	"fmt"
	"os"
	"socialapp/client"
	"testing"
	"time"
)

func TestListUsers(t *testing.T) {
	// Create a new user
	configuration := client.NewConfiguration()
	apiClient := client.NewAPIClient(configuration)
	resp, r, err := apiClient.UserApi.ListUsers(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `UserApi.ListUsers``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListUsers`: []User
	fmt.Fprintf(os.Stdout, "Response from `UserApi.ListUsers`: %v\n", resp)
}

func TestCreateUser(t *testing.T) {
	username := fmt.Sprintf("Test-%d", time.Now().UnixNano())
	email := fmt.Sprintf("Test-%d-@social.com", time.Now().UnixNano())
	user := *client.NewUser(username, "FirstName_example", "LastName_example", email) // User | Create a new user

	configuration := client.NewConfiguration()
	apiClient := client.NewAPIClient(configuration)
	resp, r, err := apiClient.UserApi.CreateUser(context.Background()).User(user).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `UserApi.CreateUser``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateUser`: User
	fmt.Fprintf(os.Stdout, "Response from `UserApi.CreateUser`: %v\n", resp)

	// attempt to recreate user
	errResp, errR, errErr := apiClient.UserApi.CreateUser(context.Background()).User(user).Execute()
	if errErr == nil {
		t.Fatalf("Expected error but called create user with existing user %+v, %+v", errResp, errR)
	}
}
