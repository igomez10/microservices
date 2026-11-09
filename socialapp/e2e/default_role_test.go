package e2e

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/igomez10/microservices/socialapp/pkg/scopes"
	"github.com/igomez10/microservices/socialapp/socialappapi/openapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func status(t *testing.T, resp *http.Response) int {
	t.Helper()
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)
	return resp.StatusCode
}

func jsonBody(t *testing.T, v any) []byte {
	t.Helper()
	b, err := json.Marshal(v)
	require.NoError(t, err)
	return b
}

// A plain signup gets the limited "user" role: admin scopes it asks for are dropped from the token
// (not a 401, so clients that request every scope, like the frontend, still log in).
func TestDefaultRole_SignupHasNoAdminScopes(t *testing.T) {
	env := setupTestEnv(t)
	defer teardownTestEnv(t, env)

	u := createTestUser(t, env)
	token := getAuthToken(t, env, u, scopes.SocialappUsersRead.String(), scopes.SocialappRolesCreate.String())

	resp := makeAuthenticatedRequest(t, http.MethodPost, env.BaseURL+"/v1/roles", token,
		jsonBody(t, openapi.Role{Name: "should-not-exist", Description: "x"}))
	assert.Equal(t, http.StatusForbidden, status(t, resp), "a regular user must not create roles")

	resp = makeAuthenticatedRequest(t, http.MethodGet, env.BaseURL+"/v1/users/"+u.Username, token, nil)
	assert.Equal(t, http.StatusOK, status(t, resp), "the allowed scope must still be granted")
}

// A regular user can do the everyday things: edit own profile, comment, follow as themselves.
func TestDefaultRole_UserCanUseSocialFeatures(t *testing.T) {
	env := setupTestEnv(t)
	defer teardownTestEnv(t, env)

	alice := createTestUser(t, env)
	bob := createTestUser(t, env)
	token := getAuthToken(t, env, alice,
		scopes.SocialappUsersUpdate.String(), scopes.SocialappCommentsCreate.String(),
		scopes.SocialappFollowerCreate.String(), scopes.SocialappFollowerDelete.String(), scopes.SocialappFeedRead.String())

	resp := makeAuthenticatedRequest(t, http.MethodPut, env.BaseURL+"/v1/users/"+alice.Username, token,
		jsonBody(t, openapi.User{Username: alice.Username, FirstName: "Alice", LastName: "Updated", Email: alice.Email}))
	assert.Equal(t, http.StatusOK, status(t, resp), "update own profile")

	resp = makeAuthenticatedRequest(t, http.MethodPost, env.BaseURL+"/v1/comments", token,
		jsonBody(t, openapi.CreateCommentRequest{Username: alice.Username, Content: "hello"}))
	assert.Equal(t, http.StatusOK, status(t, resp), "comment as self")

	followURL := fmt.Sprintf("%s/v1/users/%s/followers/%s", env.BaseURL, bob.Username, alice.Username)
	assert.Equal(t, http.StatusOK, status(t, makeAuthenticatedRequest(t, http.MethodPost, followURL, token, nil)), "follow as self")
	assert.Equal(t, http.StatusOK, status(t, makeAuthenticatedRequest(t, http.MethodDelete, followURL, token, nil)), "unfollow as self")
}

// Without the admin marker scope a user cannot edit another user or follow/unfollow on their behalf.
// With it (administrators only), they can.
func TestOwnership_SelfOnlyUnlessAdminMarker(t *testing.T) {
	env := setupTestEnv(t)
	defer teardownTestEnv(t, env)

	alice := createTestUser(t, env)
	bob := createTestUser(t, env)
	carol := createTestUser(t, env)
	aliceToken := getAuthToken(t, env, alice,
		scopes.SocialappUsersUpdate.String(), scopes.SocialappFollowerCreate.String(), scopes.SocialappFollowerDelete.String(),
		scopes.SocialappUsersRolesUpdate.String()) // dropped: alice is not an admin

	editBob := jsonBody(t, openapi.User{Username: bob.Username, FirstName: "Hacked", LastName: "Bob", Email: bob.Email})
	bobFollowsCarol := fmt.Sprintf("%s/v1/users/%s/followers/%s", env.BaseURL, carol.Username, bob.Username)

	assert.Equal(t, http.StatusForbidden, status(t, makeAuthenticatedRequest(t, http.MethodPut, env.BaseURL+"/v1/users/"+bob.Username, aliceToken, editBob)), "edit another user")
	assert.Equal(t, http.StatusForbidden, status(t, makeAuthenticatedRequest(t, http.MethodPost, bobFollowsCarol, aliceToken, nil)), "follow on behalf of another user")
	assert.Equal(t, http.StatusForbidden, status(t, makeAuthenticatedRequest(t, http.MethodDelete, bobFollowsCarol, aliceToken, nil)), "unfollow on behalf of another user")

	admin := createAdminTestUser(t, env)
	adminToken := getAuthToken(t, env, admin,
		scopes.SocialappUsersUpdate.String(), scopes.SocialappFollowerCreate.String(), scopes.SocialappUsersRolesUpdate.String())
	assert.Equal(t, http.StatusOK, status(t, makeAuthenticatedRequest(t, http.MethodPut, env.BaseURL+"/v1/users/"+bob.Username, adminToken, editBob)), "admin edits another user")
	assert.Equal(t, http.StatusOK, status(t, makeAuthenticatedRequest(t, http.MethodPost, bobFollowsCarol, adminToken, nil)), "admin follows on behalf of another user")
}
