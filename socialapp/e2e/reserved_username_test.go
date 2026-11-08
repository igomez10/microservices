package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/igomez10/microservices/socialapp/pkg/scopes"
	"github.com/igomez10/microservices/socialapp/socialappapi/openapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The seeded "admin" is soft-deleted, which frees its name (usernames are not unique in the DB).
// Neither signup nor rename may claim a reserved name.
func TestReservedUsername_SignupAndRenameRejected(t *testing.T) {
	env := setupTestEnv(t)
	defer teardownTestEnv(t, env)

	for _, name := range []string{"admin", "Admin", " ROOT ", "administrator"} {
		body, err := json.Marshal(openapi.CreateUserRequest{
			Username:  name,
			Password:  "TestPassword123!",
			FirstName: "Reserved",
			LastName:  "Name",
			Email:     fmt.Sprintf("reserved_%x@example.com", name),
		})
		require.NoError(t, err)
		resp, err := http.Post(env.BaseURL+"/v1/users", "application/json", bytes.NewReader(body))
		require.NoError(t, err)
		resp.Body.Close()
		assert.Equal(t, http.StatusConflict, resp.StatusCode, "signup as %q should be rejected", name)
	}

	user := createTestUser(t, env)
	token := getAuthToken(t, env, user, scopes.SocialappUsersUpdate.String())
	body, err := json.Marshal(openapi.User{Username: "Admin", FirstName: "Test", LastName: "User", Email: user.Email})
	require.NoError(t, err)
	resp := makeAuthenticatedRequest(t, http.MethodPut, env.BaseURL+"/v1/users/"+user.Username, token, body)
	resp.Body.Close()
	assert.Equal(t, http.StatusConflict, resp.StatusCode, "rename to a reserved username should be rejected")
}

// The seed must not ship a usable default credential.
func TestSeededAdmin_CannotLogIn(t *testing.T) {
	env := setupTestEnv(t)
	defer teardownTestEnv(t, env)

	req, err := http.NewRequest(http.MethodPost, env.BaseURL+"/v1/oauth/token?scope="+scopes.SocialappUsersRead.String(), nil)
	require.NoError(t, err)
	req.SetBasicAuth("admin", "admin")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}
