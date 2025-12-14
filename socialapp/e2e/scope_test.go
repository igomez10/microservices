package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/igomez10/microservices/socialapp/pkg/scopes"
	"github.com/igomez10/microservices/socialapp/socialappapi/openapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCreateScope_Success tests successful scope creation with authentication
func TestCreateScope_Success(t *testing.T) {
	env := setupTestEnv(t)
	defer teardownTestEnv(t, env)

	// Create a test user with admin role
	testUser := createAdminTestUser(t, env)

	// Get auth token with scope creation scope
	token := getAuthToken(t, env, testUser, scopes.SocialappScopesCreate.String())

	// Create scope request with unique name
	scopeName := fmt.Sprintf("testscope_%d", time.Now().UnixNano())
	createScopeReq := openapi.Scope{
		Name:        scopeName,
		Description: "A test scope for e2e testing",
	}

	body, err := json.Marshal(createScopeReq)
	require.NoError(t, err)

	// Make authenticated request
	resp := makeAuthenticatedRequest(t, http.MethodPost, env.BaseURL+"/v1/scopes", token, body)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	// Assert status code
	assert.Equal(t, http.StatusOK, resp.StatusCode,
		"Expected 200 OK, got %d: %s", resp.StatusCode, string(respBody))

	// Parse response
	var createdScope openapi.Scope
	err = json.Unmarshal(respBody, &createdScope)
	require.NoError(t, err, "Failed to unmarshal response: %s", string(respBody))

	// Assert response fields
	assert.NotZero(t, createdScope.Id, "Expected non-zero ID")
	assert.Equal(t, createScopeReq.Name, createdScope.Name)
	assert.Equal(t, createScopeReq.Description, createdScope.Description)
	assert.False(t, createdScope.CreatedAt.IsZero(), "Expected non-zero CreatedAt")
}

// TestCreateScope_SnowflakeIDReturned verifies that created scopes have valid snowflake IDs
func TestCreateScope_SnowflakeIDReturned(t *testing.T) {
	env := setupTestEnv(t)
	defer teardownTestEnv(t, env)

	// Create a test user with admin role
	testUser := createAdminTestUser(t, env)

	// Get auth token with scope creation scope
	token := getAuthToken(t, env, testUser, scopes.SocialappScopesCreate.String())

	scopeName := fmt.Sprintf("snowflake_scope_%d", time.Now().UnixNano())
	createScopeReq := openapi.Scope{
		Name:        scopeName,
		Description: "A scope for snowflake ID testing",
	}

	body, err := json.Marshal(createScopeReq)
	require.NoError(t, err)

	resp := makeAuthenticatedRequest(t, http.MethodPost, env.BaseURL+"/v1/scopes", token, body)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode, "Failed: %s", string(respBody))

	var createdScope openapi.Scope
	err = json.Unmarshal(respBody, &createdScope)
	require.NoError(t, err)

	// Verify the ID is a valid positive int64 (snowflake IDs are positive)
	assert.True(t, createdScope.Id > 0, "Scope ID should be positive, got %d", createdScope.Id)
	assert.True(t, createdScope.Id > 1000000, "Scope ID should be a large number, got %d", createdScope.Id)
}

// TestCreateScope_DuplicateName tests creating a scope with duplicate name
func TestCreateScope_DuplicateName(t *testing.T) {
	env := setupTestEnv(t)
	defer teardownTestEnv(t, env)

	// Create a test user with admin role
	testUser := createAdminTestUser(t, env)

	token := getAuthToken(t, env, testUser, scopes.SocialappScopesCreate.String())

	// Create first scope with unique name
	scopeName := fmt.Sprintf("duplicatescope_%d", time.Now().UnixNano())
	createScopeReq := openapi.Scope{
		Name:        scopeName,
		Description: "First scope",
	}

	body, err := json.Marshal(createScopeReq)
	require.NoError(t, err)

	firstResp := makeAuthenticatedRequest(t, http.MethodPost, env.BaseURL+"/v1/scopes", token, body)
	firstResp.Body.Close()
	require.Equal(t, http.StatusOK, firstResp.StatusCode, "First scope creation should succeed")

	// Try to create second scope with same name but different description
	duplicateScopeReq := openapi.Scope{
		Name:        scopeName, // Same name
		Description: "Second scope with same name",
	}

	duplicateBody, err := json.Marshal(duplicateScopeReq)
	require.NoError(t, err)

	duplicateResp := makeAuthenticatedRequest(t, http.MethodPost, env.BaseURL+"/v1/scopes", token, duplicateBody)
	defer duplicateResp.Body.Close()

	// Should return 409 Conflict
	assert.Equal(t, http.StatusConflict, duplicateResp.StatusCode, "Expected 409 Conflict for duplicate scope name")
}

// TestGetScope_Success tests getting a scope by ID
func TestGetScope_Success(t *testing.T) {
	env := setupTestEnv(t)
	defer teardownTestEnv(t, env)

	// Create a test user with admin role
	testUser := createAdminTestUser(t, env)

	// Get auth token with required scopes
	token := getAuthToken(t, env, testUser, scopes.SocialappScopesCreate.String(), scopes.SocialappScopesRead.String())

	// First create a scope with unique name
	scopeName := fmt.Sprintf("getscope_%d", time.Now().UnixNano())
	createScopeReq := openapi.Scope{
		Name:        scopeName,
		Description: "Scope for get test",
	}

	body, err := json.Marshal(createScopeReq)
	require.NoError(t, err)

	createResp := makeAuthenticatedRequest(t, http.MethodPost, env.BaseURL+"/v1/scopes", token, body)
	createBody, err := io.ReadAll(createResp.Body)
	createResp.Body.Close()
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, createResp.StatusCode, "Failed to create scope: %s", string(createBody))

	var createdScope openapi.Scope
	err = json.Unmarshal(createBody, &createdScope)
	require.NoError(t, err)

	// Now get the scope by ID
	getURL := fmt.Sprintf("%s/v1/scopes/%d", env.BaseURL, createdScope.Id)
	getResp := makeAuthenticatedRequest(t, http.MethodGet, getURL, token, nil)
	defer getResp.Body.Close()

	getBody, err := io.ReadAll(getResp.Body)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, getResp.StatusCode,
		"Failed to get scope: %s", string(getBody))

	var fetchedScope openapi.Scope
	err = json.Unmarshal(getBody, &fetchedScope)
	require.NoError(t, err)

	assert.Equal(t, createdScope.Id, fetchedScope.Id)
	assert.Equal(t, createdScope.Name, fetchedScope.Name)
	assert.Equal(t, createdScope.Description, fetchedScope.Description)
}

// TestGetScope_NotFound tests getting a non-existent scope
func TestGetScope_NotFound(t *testing.T) {
	env := setupTestEnv(t)
	defer teardownTestEnv(t, env)

	// Create a test user with admin role
	testUser := createAdminTestUser(t, env)

	token := getAuthToken(t, env, testUser, scopes.SocialappScopesRead.String())

	// Try to get a scope that doesn't exist
	getURL := fmt.Sprintf("%s/v1/scopes/%d", env.BaseURL, 999999999)
	resp := makeAuthenticatedRequest(t, http.MethodGet, getURL, token, nil)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "Expected 404 Not Found")
}

// TestListScopes_Success tests listing scopes with pagination
func TestListScopes_Success(t *testing.T) {
	env := setupTestEnv(t)
	defer teardownTestEnv(t, env)

	// Create a test user with admin role
	testUser := createAdminTestUser(t, env)

	token := getAuthToken(t, env, testUser, scopes.SocialappScopesList.String(), scopes.SocialappScopesCreate.String())

	// Create a few scopes
	for i := 0; i < 3; i++ {
		scopeName := fmt.Sprintf("listscope_%d_%d", time.Now().UnixNano(), i)
		createScopeReq := openapi.Scope{
			Name:        scopeName,
			Description: fmt.Sprintf("Scope number %d", i+1),
		}

		body, err := json.Marshal(createScopeReq)
		require.NoError(t, err)

		createResp := makeAuthenticatedRequest(t, http.MethodPost, env.BaseURL+"/v1/scopes", token, body)
		createResp.Body.Close()
		require.Equal(t, http.StatusOK, createResp.StatusCode, "Failed to create scope")
	}

	// List scopes with pagination
	listURL := env.BaseURL + "/v1/scopes?limit=10&offset=0"
	resp := makeAuthenticatedRequest(t, http.MethodGet, listURL, token, nil)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode,
		"Expected 200 OK, got %d: %s", resp.StatusCode, string(respBody))

	var scopesList []openapi.Scope
	err = json.Unmarshal(respBody, &scopesList)
	require.NoError(t, err, "Failed to unmarshal response: %s", string(respBody))

	// Should have at least the scopes we created
	assert.GreaterOrEqual(t, len(scopesList), 3, "Expected at least 3 scopes")
}

// TestListScopes_Empty tests listing scopes when none exist (besides seeded ones)
func TestListScopes_Empty(t *testing.T) {
	env := setupTestEnv(t)
	defer teardownTestEnv(t, env)

	// Create a test user with admin role
	testUser := createAdminTestUser(t, env)

	token := getAuthToken(t, env, testUser, scopes.SocialappScopesList.String())

	// List scopes
	listURL := env.BaseURL + "/v1/scopes?limit=10&offset=0"
	resp := makeAuthenticatedRequest(t, http.MethodGet, listURL, token, nil)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode,
		"Expected 200 OK, got %d: %s", resp.StatusCode, string(respBody))

	var scopesList []openapi.Scope
	err = json.Unmarshal(respBody, &scopesList)
	require.NoError(t, err, "Failed to unmarshal response: %s", string(respBody))

	// Should return an array (may be empty or contain seeded scopes)
	assert.NotNil(t, scopesList, "Scopes list should not be nil")
}

// TestUpdateScope_Success tests updating a scope
func TestUpdateScope_Success(t *testing.T) {
	env := setupTestEnv(t)
	defer teardownTestEnv(t, env)

	// Create a test user with admin role
	testUser := createAdminTestUser(t, env)

	token := getAuthToken(t, env, testUser, scopes.SocialappScopesCreate.String(), scopes.SocialappScopesUpdate.String(), scopes.SocialappScopesRead.String())

	// First create a scope with unique name
	scopeName := fmt.Sprintf("updatescope_%d", time.Now().UnixNano())
	createScopeReq := openapi.Scope{
		Name:        scopeName,
		Description: "Original description",
	}

	body, err := json.Marshal(createScopeReq)
	require.NoError(t, err)

	createResp := makeAuthenticatedRequest(t, http.MethodPost, env.BaseURL+"/v1/scopes", token, body)
	createBody, err := io.ReadAll(createResp.Body)
	createResp.Body.Close()
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, createResp.StatusCode, "Failed to create scope: %s", string(createBody))

	var createdScope openapi.Scope
	err = json.Unmarshal(createBody, &createdScope)
	require.NoError(t, err)

	// Update the scope with new unique name
	updatedScopeName := fmt.Sprintf("updatedscope_%d", time.Now().UnixNano())
	updateScopeReq := openapi.Scope{
		Name:        updatedScopeName,
		Description: "Updated description",
	}

	updateBody, err := json.Marshal(updateScopeReq)
	require.NoError(t, err)

	updateURL := fmt.Sprintf("%s/v1/scopes/%d", env.BaseURL, createdScope.Id)
	updateResp := makeAuthenticatedRequest(t, http.MethodPut, updateURL, token, updateBody)
	defer updateResp.Body.Close()

	updateRespBody, err := io.ReadAll(updateResp.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, updateResp.StatusCode,
		"Expected 200 OK, got %d: %s", updateResp.StatusCode, string(updateRespBody))

	var updatedScope openapi.Scope
	err = json.Unmarshal(updateRespBody, &updatedScope)
	require.NoError(t, err)

	assert.Equal(t, createdScope.Id, updatedScope.Id)
	assert.Equal(t, updateScopeReq.Name, updatedScope.Name)
	assert.Equal(t, updateScopeReq.Description, updatedScope.Description)
}

// TestUpdateScope_NotFound tests updating a non-existent scope
func TestUpdateScope_NotFound(t *testing.T) {
	env := setupTestEnv(t)
	defer teardownTestEnv(t, env)

	// Create a test user with admin role
	testUser := createAdminTestUser(t, env)

	token := getAuthToken(t, env, testUser, scopes.SocialappScopesUpdate.String())

	// Try to update a scope that doesn't exist
	updateScopeReq := openapi.Scope{
		Name:        "nonexistent",
		Description: "This should fail",
	}

	updateBody, err := json.Marshal(updateScopeReq)
	require.NoError(t, err)

	updateURL := fmt.Sprintf("%s/v1/scopes/%d", env.BaseURL, 999999999)
	updateResp := makeAuthenticatedRequest(t, http.MethodPut, updateURL, token, updateBody)
	defer updateResp.Body.Close()

	assert.Equal(t, http.StatusNotFound, updateResp.StatusCode, "Expected 404 Not Found")
}

// TestDeleteScope_Success tests deleting a scope
func TestDeleteScope_Success(t *testing.T) {
	env := setupTestEnv(t)
	defer teardownTestEnv(t, env)

	// Create a test user with admin role
	testUser := createAdminTestUser(t, env)

	token := getAuthToken(t, env, testUser, scopes.SocialappScopesCreate.String(), scopes.SocialappScopesDelete.String(), scopes.SocialappScopesRead.String())

	// First create a scope with unique name
	scopeName := fmt.Sprintf("deletescope_%d", time.Now().UnixNano())
	createScopeReq := openapi.Scope{
		Name:        scopeName,
		Description: "Scope to be deleted",
	}

	body, err := json.Marshal(createScopeReq)
	require.NoError(t, err)

	createResp := makeAuthenticatedRequest(t, http.MethodPost, env.BaseURL+"/v1/scopes", token, body)
	createBody, err := io.ReadAll(createResp.Body)
	createResp.Body.Close()
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, createResp.StatusCode, "Failed to create scope: %s", string(createBody))

	var createdScope openapi.Scope
	err = json.Unmarshal(createBody, &createdScope)
	require.NoError(t, err)

	// Delete the scope
	deleteURL := fmt.Sprintf("%s/v1/scopes/%d", env.BaseURL, createdScope.Id)
	deleteResp := makeAuthenticatedRequest(t, http.MethodDelete, deleteURL, token, nil)
	defer deleteResp.Body.Close()

	assert.Equal(t, http.StatusOK, deleteResp.StatusCode, "Expected 200 OK for delete")

	// Verify the scope no longer exists
	getResp := makeAuthenticatedRequest(t, http.MethodGet, deleteURL, token, nil)
	defer getResp.Body.Close()

	assert.Equal(t, http.StatusNotFound, getResp.StatusCode, "Expected 404 Not Found after deletion")
}

// TestDeleteScope_NotFound tests deleting a non-existent scope
func TestDeleteScope_NotFound(t *testing.T) {
	env := setupTestEnv(t)
	defer teardownTestEnv(t, env)

	// Create a test user with admin role
	testUser := createAdminTestUser(t, env)

	token := getAuthToken(t, env, testUser, scopes.SocialappScopesDelete.String())

	// Try to delete a scope that doesn't exist
	deleteURL := fmt.Sprintf("%s/v1/scopes/%d", env.BaseURL, 999999999)
	deleteResp := makeAuthenticatedRequest(t, http.MethodDelete, deleteURL, token, nil)
	defer deleteResp.Body.Close()

	assert.Equal(t, http.StatusNotFound, deleteResp.StatusCode, "Expected 404 Not Found")
}

// TestCreateScope_Unauthorized tests creating a scope without authentication
func TestCreateScope_Unauthorized(t *testing.T) {
	env := setupTestEnv(t)
	defer teardownTestEnv(t, env)

	// Create scope request
	createScopeReq := openapi.Scope{
		Name:        "testscope",
		Description: "This should fail",
	}

	body, err := json.Marshal(createScopeReq)
	require.NoError(t, err)

	// Make request without authentication
	resp, err := http.Post(env.BaseURL+"/v1/scopes", "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	// Should return 401 Unauthorized
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "Expected 401 Unauthorized")
}

// TestCreateScope_WrongScope tests creating a scope without required scope
func TestCreateScope_WrongScope(t *testing.T) {
	env := setupTestEnv(t)
	defer teardownTestEnv(t, env)

	// Create a test user with admin role
	testUser := createAdminTestUser(t, env)

	// Get auth token with wrong scope (read instead of create)
	token := getAuthToken(t, env, testUser, scopes.SocialappScopesRead.String())

	// Create scope request
	createScopeReq := openapi.Scope{
		Name:        "testscope",
		Description: "This should fail",
	}

	body, err := json.Marshal(createScopeReq)
	require.NoError(t, err)

	// Make authenticated request
	resp := makeAuthenticatedRequest(t, http.MethodPost, env.BaseURL+"/v1/scopes", token, body)
	defer resp.Body.Close()

	// Should return 403 Forbidden
	assert.Equal(t, http.StatusForbidden, resp.StatusCode, "Expected 403 Forbidden for wrong scope")
}
