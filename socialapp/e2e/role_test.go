package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/igomez10/microservices/socialapp/pkg/scopes"
	"github.com/igomez10/microservices/socialapp/socialappapi/openapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// getAuthToken obtains a JWT token using Basic Auth with the given user credentials
func getAuthToken(t *testing.T, env *TestEnv, user TestUser, scopes ...string) string {
	t.Helper()

	// Build scope query parameter with proper URL encoding
	scopeParam := ""
	for i, scope := range scopes {
		if i > 0 {
			scopeParam += " "
		}
		scopeParam += scope
	}

	// Create request to token endpoint with URL-encoded scope parameter
	tokenURL := env.BaseURL + "/v1/oauth/token"
	if scopeParam != "" {
		tokenURL += "?scope=" + url.QueryEscape(scopeParam)
	}

	req, err := http.NewRequest(http.MethodPost, tokenURL, nil)
	require.NoError(t, err)

	// Set Basic Auth with user credentials
	req.SetBasicAuth(user.Username, user.Password)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, resp.StatusCode,
		"Failed to get auth token: %s", string(respBody))

	var tokenResp openapi.AccessToken
	err = json.Unmarshal(respBody, &tokenResp)
	require.NoError(t, err, "Failed to unmarshal token response: %s", string(respBody))

	return tokenResp.AccessToken
}

// makeAuthenticatedRequest creates an HTTP request with Bearer token authentication
func makeAuthenticatedRequest(t *testing.T, method, url, token string, body []byte) *http.Response {
	t.Helper()

	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	require.NoError(t, err)

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	return resp
}

// TestCreateRole_Success tests successful role creation with authentication
func TestCreateRole_Success(t *testing.T) {
	env := setupTestEnv(t)
	defer teardownTestEnv(t, env)

	// Create a test user with admin role
	testUser := createAdminTestUser(t, env)

	// Get auth token with role creation scope
	token := getAuthToken(t, env, testUser, scopes.SocialappRolesCreate.String())

	// Create role request with unique name
	roleName := fmt.Sprintf("testrole_%d", time.Now().UnixNano())
	createRoleReq := openapi.Role{
		Name:        roleName,
		Description: "A test role for e2e testing",
	}

	body, err := json.Marshal(createRoleReq)
	require.NoError(t, err)

	// Make authenticated request
	resp := makeAuthenticatedRequest(t, http.MethodPost, env.BaseURL+"/v1/roles", token, body)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	// Assert status code
	assert.Equal(t, http.StatusOK, resp.StatusCode,
		"Expected 200 OK, got %d: %s", resp.StatusCode, string(respBody))

	// Parse response
	var createdRole openapi.Role
	err = json.Unmarshal(respBody, &createdRole)
	require.NoError(t, err, "Failed to unmarshal response: %s", string(respBody))

	// Assert response fields
	assert.NotZero(t, createdRole.Id, "Expected non-zero ID")
	assert.Equal(t, createRoleReq.Name, createdRole.Name)
	assert.Equal(t, createRoleReq.Description, createdRole.Description)
	assert.False(t, createdRole.CreatedAt.IsZero(), "Expected non-zero CreatedAt")
}

// TestGetRole_Success tests getting a role by ID
func TestGetRole_Success(t *testing.T) {
	env := setupTestEnv(t)
	defer teardownTestEnv(t, env)

	// Create a test user with admin role
	testUser := createAdminTestUser(t, env)

	// Get auth token with required scopes
	token := getAuthToken(t, env, testUser, scopes.SocialappRolesCreate.String(), scopes.SocialappRolesRead.String())

	// First create a role with unique name
	roleName := fmt.Sprintf("getrole_%d", time.Now().UnixNano())
	createRoleReq := openapi.Role{
		Name:        roleName,
		Description: "Role for get test",
	}

	body, err := json.Marshal(createRoleReq)
	require.NoError(t, err)

	createResp := makeAuthenticatedRequest(t, http.MethodPost, env.BaseURL+"/v1/roles", token, body)
	createBody, err := io.ReadAll(createResp.Body)
	createResp.Body.Close()
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, createResp.StatusCode, "Failed to create role: %s", string(createBody))

	var createdRole openapi.Role
	err = json.Unmarshal(createBody, &createdRole)
	require.NoError(t, err)

	// Now get the role by ID
	getURL := fmt.Sprintf("%s/v1/roles/%s", env.BaseURL, createdRole.Id)
	getResp := makeAuthenticatedRequest(t, http.MethodGet, getURL, token, nil)
	defer getResp.Body.Close()

	getBody, err := io.ReadAll(getResp.Body)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, getResp.StatusCode,
		"Failed to get role: %s", string(getBody))

	var fetchedRole openapi.Role
	err = json.Unmarshal(getBody, &fetchedRole)
	require.NoError(t, err)

	assert.Equal(t, createdRole.Id, fetchedRole.Id)
	assert.Equal(t, createdRole.Name, fetchedRole.Name)
	assert.Equal(t, createdRole.Description, fetchedRole.Description)
}

// TestGetRole_NotFound tests getting a non-existent role
func TestGetRole_NotFound(t *testing.T) {
	env := setupTestEnv(t)
	defer teardownTestEnv(t, env)

	// Create a test user with admin role
	testUser := createAdminTestUser(t, env)

	token := getAuthToken(t, env, testUser, scopes.SocialappRolesRead.String())

	// Try to get a role that doesn't exist
	getURL := fmt.Sprintf("%s/v1/roles/%d", env.BaseURL, 999999)
	resp := makeAuthenticatedRequest(t, http.MethodGet, getURL, token, nil)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "Expected 404 Not Found")
}

// TestListRoles_Success tests listing roles with pagination
func TestListRoles_Success(t *testing.T) {
	env := setupTestEnv(t)
	defer teardownTestEnv(t, env)

	// Create a test user with admin role
	testUser := createAdminTestUser(t, env)

	token := getAuthToken(t, env, testUser, "socialapp.roles.list")

	// List roles (should include seeded roles from schema.sql)
	listURL := env.BaseURL + "/v1/roles?limit=10&offset=0"
	resp := makeAuthenticatedRequest(t, http.MethodGet, listURL, token, nil)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode,
		"Expected 200 OK, got %d: %s", resp.StatusCode, string(respBody))

	var roles []openapi.Role
	err = json.Unmarshal(respBody, &roles)
	require.NoError(t, err, "Failed to unmarshal response: %s", string(respBody))

	// Should have at least the seeded roles (administrator, user)
	assert.GreaterOrEqual(t, len(roles), 2, "Expected at least 2 seeded roles")

	// Verify seeded roles exist
	roleNames := make(map[string]bool)
	for _, role := range roles {
		roleNames[role.Name] = true
	}
	assert.True(t, roleNames["administrator"], "Expected 'administrator' role")
	assert.True(t, roleNames["user"], "Expected 'user' role")
}

// TestUpdateRole_Success tests updating a role
func TestUpdateRole_Success(t *testing.T) {
	env := setupTestEnv(t)
	defer teardownTestEnv(t, env)

	// Create a test user with admin role
	testUser := createAdminTestUser(t, env)

	token := getAuthToken(t, env, testUser, scopes.SocialappRolesCreate.String(), scopes.SocialappRolesUpdate.String(), scopes.SocialappRolesRead.String())

	// First create a role with unique name
	roleName := fmt.Sprintf("updaterole_%d", time.Now().UnixNano())
	createRoleReq := openapi.Role{
		Name:        roleName,
		Description: "Original description",
	}

	body, err := json.Marshal(createRoleReq)
	require.NoError(t, err)

	createResp := makeAuthenticatedRequest(t, http.MethodPost, env.BaseURL+"/v1/roles", token, body)
	createBody, err := io.ReadAll(createResp.Body)
	createResp.Body.Close()
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, createResp.StatusCode, "Failed to create role: %s", string(createBody))

	var createdRole openapi.Role
	err = json.Unmarshal(createBody, &createdRole)
	require.NoError(t, err)

	// Update the role with new unique name
	updatedRoleName := fmt.Sprintf("updatedrole_%d", time.Now().UnixNano())
	updateRoleReq := openapi.Role{
		Name:        updatedRoleName,
		Description: "Updated description",
	}

	updateBody, err := json.Marshal(updateRoleReq)
	require.NoError(t, err)

	updateURL := fmt.Sprintf("%s/v1/roles/%s", env.BaseURL, createdRole.Id)
	updateResp := makeAuthenticatedRequest(t, http.MethodPut, updateURL, token, updateBody)
	defer updateResp.Body.Close()

	updateRespBody, err := io.ReadAll(updateResp.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, updateResp.StatusCode,
		"Expected 200 OK, got %d: %s", updateResp.StatusCode, string(updateRespBody))

	var updatedRole openapi.Role
	err = json.Unmarshal(updateRespBody, &updatedRole)
	require.NoError(t, err)

	assert.Equal(t, createdRole.Id, updatedRole.Id)
	assert.Equal(t, updateRoleReq.Name, updatedRole.Name)
	assert.Equal(t, updateRoleReq.Description, updatedRole.Description)
}

// TestDeleteRole_Success tests deleting a role
func TestDeleteRole_Success(t *testing.T) {
	env := setupTestEnv(t)
	defer teardownTestEnv(t, env)

	// Create a test user with admin role
	testUser := createAdminTestUser(t, env)

	token := getAuthToken(t, env, testUser, scopes.SocialappRolesCreate.String(), scopes.SocialappRolesDelete.String(), scopes.SocialappRolesRead.String())

	// First create a role with unique name
	roleName := fmt.Sprintf("deleterole_%d", time.Now().UnixNano())
	createRoleReq := openapi.Role{
		Name:        roleName,
		Description: "Role to be deleted",
	}

	body, err := json.Marshal(createRoleReq)
	require.NoError(t, err)

	createResp := makeAuthenticatedRequest(t, http.MethodPost, env.BaseURL+"/v1/roles", token, body)
	createBody, err := io.ReadAll(createResp.Body)
	createResp.Body.Close()
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, createResp.StatusCode, "Failed to create role: %s", string(createBody))

	var createdRole openapi.Role
	err = json.Unmarshal(createBody, &createdRole)
	require.NoError(t, err)

	// Delete the role
	deleteURL := fmt.Sprintf("%s/v1/roles/%s", env.BaseURL, createdRole.Id)
	deleteResp := makeAuthenticatedRequest(t, http.MethodDelete, deleteURL, token, nil)
	defer deleteResp.Body.Close()

	assert.Equal(t, http.StatusOK, deleteResp.StatusCode, "Expected 200 OK for delete")

	// Verify the role no longer exists
	getResp := makeAuthenticatedRequest(t, http.MethodGet, deleteURL, token, nil)
	defer getResp.Body.Close()

	assert.Equal(t, http.StatusNotFound, getResp.StatusCode, "Expected 404 Not Found after deletion")
}

// TestCreateRole_MultipleRoles tests creating multiple roles successfully
func TestCreateRole_MultipleRoles(t *testing.T) {
	env := setupTestEnv(t)
	defer teardownTestEnv(t, env)

	// Create a test user with admin role
	testUser := createAdminTestUser(t, env)

	token := getAuthToken(t, env, testUser, scopes.SocialappRolesCreate.String())

	// Create multiple roles with unique names
	for i := 0; i < 3; i++ {
		roleName := fmt.Sprintf("multiplerole_%d_%d", time.Now().UnixNano(), i)
		createRoleReq := openapi.Role{
			Name:        roleName,
			Description: fmt.Sprintf("Role number %d", i+1),
		}

		body, err := json.Marshal(createRoleReq)
		require.NoError(t, err)

		resp := makeAuthenticatedRequest(t, http.MethodPost, env.BaseURL+"/v1/roles", token, body)
		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode,
			"Role %d creation failed: %s", i+1, string(respBody))

		var createdRole openapi.Role
		err = json.Unmarshal(respBody, &createdRole)
		require.NoError(t, err)

		assert.NotZero(t, createdRole.Id)
		assert.Equal(t, roleName, createdRole.Name)
	}
}

// TestCreateRole_SnowflakeIDReturned verifies that created roles have valid snowflake IDs
func TestCreateRole_SnowflakeIDReturned(t *testing.T) {
	env := setupTestEnv(t)
	defer teardownTestEnv(t, env)

	// Create a test user with admin role
	testUser := createAdminTestUser(t, env)

	// Get auth token with role creation scope
	token := getAuthToken(t, env, testUser, scopes.SocialappRolesCreate.String())

	roleName := fmt.Sprintf("snowflake_role_%d", time.Now().UnixNano())
	createRoleReq := openapi.Role{
		Name:        roleName,
		Description: "A role for snowflake ID testing",
	}

	body, err := json.Marshal(createRoleReq)
	require.NoError(t, err)

	resp := makeAuthenticatedRequest(t, http.MethodPost, env.BaseURL+"/v1/roles", token, body)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode, "Failed: %s", string(respBody))

	var createdRole openapi.Role
	err = json.Unmarshal(respBody, &createdRole)
	require.NoError(t, err)

	// Verify the ID is a valid positive int64 (snowflake IDs are positive).
	// Role.Id is serialized as a JSON string (see openapi.yaml) — parse to int64 to compare.
	roleID, parseErr := strconv.ParseInt(createdRole.Id, 10, 64)
	require.NoError(t, parseErr, "Role ID should be numeric, got %q", createdRole.Id)
	assert.True(t, roleID > 0, "Role ID should be positive, got %d", roleID)
}
