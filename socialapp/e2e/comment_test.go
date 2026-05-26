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

// TestCreateComment_Success tests successful comment creation with authentication
func TestCreateComment_Success(t *testing.T) {
	env := setupTestEnv(t)
	defer teardownTestEnv(t, env)

	// Create a test user with admin role (for comment author)
	testUser := createAdminTestUser(t, env)

	// Get auth token with comment creation scope
	token := getAuthToken(t, env, testUser, scopes.SocialappCommentsCreate.String())

	// Create comment request
	createCommentReq := openapi.CreateCommentRequest{
		Username: testUser.Username,
		Content:  "This is a test comment",
	}

	body, err := json.Marshal(createCommentReq)
	require.NoError(t, err)

	// Make authenticated request
	resp := makeAuthenticatedRequest(t, http.MethodPost, env.BaseURL+"/v1/comments", token, body)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	// Assert status code
	assert.Equal(t, http.StatusOK, resp.StatusCode,
		"Expected 200 OK, got %d: %s", resp.StatusCode, string(respBody))

	// Parse response
	var createdComment openapi.Comment
	err = json.Unmarshal(respBody, &createdComment)
	require.NoError(t, err, "Failed to unmarshal response: %s", string(respBody))

	// Assert response fields
	assert.NotZero(t, createdComment.Id, "Expected non-zero ID")
	assert.Equal(t, createCommentReq.Content, createdComment.Content)
	assert.Equal(t, createCommentReq.Username, createdComment.Username)
	assert.False(t, createdComment.CreatedAt.IsZero(), "Expected non-zero CreatedAt")
}

// TestCreateComment_SnowflakeIDReturned verifies that created comments have valid snowflake IDs
func TestCreateComment_SnowflakeIDReturned(t *testing.T) {
	env := setupTestEnv(t)
	defer teardownTestEnv(t, env)

	// Create a test user with admin role
	testUser := createAdminTestUser(t, env)

	// Get auth token with comment creation scope
	token := getAuthToken(t, env, testUser, scopes.SocialappCommentsCreate.String())

	createCommentReq := openapi.CreateCommentRequest{
		Username: testUser.Username,
		Content:  fmt.Sprintf("Snowflake test comment %d", time.Now().UnixNano()),
	}

	body, err := json.Marshal(createCommentReq)
	require.NoError(t, err)

	resp := makeAuthenticatedRequest(t, http.MethodPost, env.BaseURL+"/v1/comments", token, body)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode, "Failed: %s", string(respBody))

	var createdComment openapi.Comment
	err = json.Unmarshal(respBody, &createdComment)
	require.NoError(t, err)

	// Verify the ID is a valid positive int64 (snowflake IDs are positive).
	// Comment.Id is serialized as a JSON string on the wire (see openapi.yaml)
	// to avoid precision loss in JS clients, so we parse it back to int64 here.
	commentID, parseErr := strconv.ParseInt(createdComment.Id, 10, 64)
	require.NoError(t, parseErr, "Comment ID should be numeric, got %q", createdComment.Id)
	assert.True(t, commentID > 0, "Comment ID should be positive, got %d", commentID)
	assert.True(t, commentID > 1000000, "Comment ID should be a large number, got %d", commentID)
}

// TestGetComment_Success tests getting a comment by ID
func TestGetComment_Success(t *testing.T) {
	env := setupTestEnv(t)
	defer teardownTestEnv(t, env)

	// Create a test user with admin role
	testUser := createAdminTestUser(t, env)

	// Get auth token with required scopes
	token := getAuthToken(t, env, testUser, scopes.SocialappCommentsCreate.String(), scopes.SocialappCommentsRead.String())

	// First create a comment
	createCommentReq := openapi.CreateCommentRequest{
		Username: testUser.Username,
		Content:  "Comment for get test",
	}

	body, err := json.Marshal(createCommentReq)
	require.NoError(t, err)

	createResp := makeAuthenticatedRequest(t, http.MethodPost, env.BaseURL+"/v1/comments", token, body)
	createBody, err := io.ReadAll(createResp.Body)
	createResp.Body.Close()
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, createResp.StatusCode, "Failed to create comment: %s", string(createBody))

	var createdComment openapi.Comment
	err = json.Unmarshal(createBody, &createdComment)
	require.NoError(t, err)

	// Now get the comment by ID (Id is a string on the wire — see openapi.yaml)
	getURL := fmt.Sprintf("%s/v1/comments/%s", env.BaseURL, createdComment.Id)
	getResp := makeAuthenticatedRequest(t, http.MethodGet, getURL, token, nil)
	defer getResp.Body.Close()

	getBody, err := io.ReadAll(getResp.Body)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, getResp.StatusCode,
		"Failed to get comment: %s", string(getBody))

	var fetchedComment openapi.Comment
	err = json.Unmarshal(getBody, &fetchedComment)
	require.NoError(t, err)

	assert.Equal(t, createdComment.Id, fetchedComment.Id)
	assert.Equal(t, createdComment.Content, fetchedComment.Content)
	assert.Equal(t, createdComment.Username, fetchedComment.Username)
}

// TestGetComment_NotFound tests getting a non-existent comment
func TestGetComment_NotFound(t *testing.T) {
	env := setupTestEnv(t)
	defer teardownTestEnv(t, env)

	// Create a test user with admin role
	testUser := createAdminTestUser(t, env)

	token := getAuthToken(t, env, testUser, scopes.SocialappCommentsRead.String())

	// Try to get a comment that doesn't exist
	getURL := fmt.Sprintf("%s/v1/comments/%d", env.BaseURL, 999999999)
	resp := makeAuthenticatedRequest(t, http.MethodGet, getURL, token, nil)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "Expected 404 Not Found")
}

// TestGetUserFeed_Success tests getting user feed with comments from followed users
func TestGetUserFeed_Success(t *testing.T) {
	env := setupTestEnv(t)
	defer teardownTestEnv(t, env)

	// Create two test users
	user1 := createAdminTestUser(t, env)
	user2 := createAdminTestUser(t, env)

	// Get auth tokens
	token1 := getAuthToken(t, env, user1, scopes.SocialappCommentsCreate.String(), scopes.SocialappFeedRead.String(), scopes.SocialappFollowerCreate.String())

	// User1 follows User2
	followURL := fmt.Sprintf("%s/v1/users/%s/followers/%s", env.BaseURL, user2.Username, user1.Username)
	followResp := makeAuthenticatedRequest(t, http.MethodPost, followURL, token1, nil)
	followResp.Body.Close()
	require.Equal(t, http.StatusOK, followResp.StatusCode, "Failed to follow user")

	// User2 creates a comment
	token2 := getAuthToken(t, env, user2, scopes.SocialappCommentsCreate.String())
	createCommentReq := openapi.CreateCommentRequest{
		Username: user2.Username,
		Content:  "Comment from followed user",
	}

	body, err := json.Marshal(createCommentReq)
	require.NoError(t, err)

	createResp := makeAuthenticatedRequest(t, http.MethodPost, env.BaseURL+"/v1/comments", token2, body)
	createResp.Body.Close()
	require.Equal(t, http.StatusOK, createResp.StatusCode, "Failed to create comment")

	// User1 gets their feed
	feedResp := makeAuthenticatedRequest(t, http.MethodGet, env.BaseURL+"/v1/feed", token1, nil)
	defer feedResp.Body.Close()

	feedBody, err := io.ReadAll(feedResp.Body)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, feedResp.StatusCode,
		"Failed to get feed: %s", string(feedBody))

	var feedComments []openapi.Comment
	err = json.Unmarshal(feedBody, &feedComments)
	require.NoError(t, err, "Failed to unmarshal feed response: %s", string(feedBody))

	// Feed should contain at least one comment
	assert.GreaterOrEqual(t, len(feedComments), 1, "Feed should contain at least one comment")
}

// TestGetUserFeed_EmptyFeed tests getting feed when user has no followed users
func TestGetUserFeed_EmptyFeed(t *testing.T) {
	env := setupTestEnv(t)
	defer teardownTestEnv(t, env)

	// Create a test user with admin role
	testUser := createAdminTestUser(t, env)

	token := getAuthToken(t, env, testUser, scopes.SocialappFeedRead.String())

	// Get feed (user has no followed users)
	feedResp := makeAuthenticatedRequest(t, http.MethodGet, env.BaseURL+"/v1/feed", token, nil)
	defer feedResp.Body.Close()

	feedBody, err := io.ReadAll(feedResp.Body)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, feedResp.StatusCode,
		"Failed to get feed: %s", string(feedBody))

	var feedComments []openapi.Comment
	err = json.Unmarshal(feedBody, &feedComments)
	require.NoError(t, err, "Failed to unmarshal feed response: %s", string(feedBody))

	// Feed should be empty
	assert.Equal(t, 0, len(feedComments), "Feed should be empty when user has no followed users")
}

func TestSearchComments_ByUserAndTime(t *testing.T) {
	env := setupTestEnv(t)
	defer teardownTestEnv(t, env)

	user1 := createAdminTestUser(t, env)
	user2 := createAdminTestUser(t, env)

	token1 := getAuthToken(t, env, user1, scopes.SocialappCommentsCreate.String(), scopes.SocialappCommentsRead.String())
	token2 := getAuthToken(t, env, user2, scopes.SocialappCommentsCreate.String())

	createComment := func(token string, username string, content string) openapi.Comment {
		t.Helper()

		body, err := json.Marshal(openapi.CreateCommentRequest{
			Username: username,
			Content:  content,
		})
		require.NoError(t, err)

		resp := makeAuthenticatedRequest(t, http.MethodPost, env.BaseURL+"/v1/comments", token, body)
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode, "Failed to create comment: %s", string(respBody))

		var created openapi.Comment
		require.NoError(t, json.Unmarshal(respBody, &created))
		return created
	}

	first := createComment(token1, user1.Username, "first by user1")
	time.Sleep(10 * time.Millisecond)
	second := createComment(token2, user2.Username, "by user2")
	time.Sleep(10 * time.Millisecond)
	third := createComment(token1, user1.Username, "second by user1")

	t.Run("filters by username", func(t *testing.T) {
		searchURL := fmt.Sprintf("%s/v1/comments?username=%s", env.BaseURL, url.QueryEscape(user1.Username))
		resp := makeAuthenticatedRequest(t, http.MethodGet, searchURL, token1, nil)
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode, "search failed: %s", string(respBody))

		var comments []openapi.Comment
		require.NoError(t, json.Unmarshal(respBody, &comments))
		require.Len(t, comments, 2)
		assert.Equal(t, third.Id, comments[0].Id)
		assert.Equal(t, first.Id, comments[1].Id)
		for _, comment := range comments {
			assert.Equal(t, user1.Username, comment.Username)
		}
	})

	t.Run("filters by time window", func(t *testing.T) {
		values := url.Values{}
		values.Set("start_time", second.CreatedAt.Add(-time.Millisecond).Format(time.RFC3339Nano))
		values.Set("end_time", second.CreatedAt.Add(time.Millisecond).Format(time.RFC3339Nano))

		resp := makeAuthenticatedRequest(t, http.MethodGet, env.BaseURL+"/v1/comments?"+values.Encode(), token1, nil)
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode, "search failed: %s", string(respBody))

		var comments []openapi.Comment
		require.NoError(t, json.Unmarshal(respBody, &comments))
		require.Len(t, comments, 1)
		assert.Equal(t, second.Id, comments[0].Id)
		assert.Equal(t, user2.Username, comments[0].Username)
	})
}

func TestSearchComments_ReturnsMostRecent20(t *testing.T) {
	env := setupTestEnv(t)
	defer teardownTestEnv(t, env)

	testUser := createAdminTestUser(t, env)
	token := getAuthToken(t, env, testUser, scopes.SocialappCommentsCreate.String(), scopes.SocialappCommentsRead.String())

	createdIDs := make([]string, 0, 25)
	for i := 0; i < 25; i++ {
		body, err := json.Marshal(openapi.CreateCommentRequest{
			Username: testUser.Username,
			Content:  fmt.Sprintf("comment-%02d", i),
		})
		require.NoError(t, err)

		resp := makeAuthenticatedRequest(t, http.MethodPost, env.BaseURL+"/v1/comments", token, body)
		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode, "Failed to create comment: %s", string(respBody))

		var created openapi.Comment
		require.NoError(t, json.Unmarshal(respBody, &created))
		createdIDs = append(createdIDs, created.Id)
		time.Sleep(5 * time.Millisecond)
	}

	resp := makeAuthenticatedRequest(t, http.MethodGet, env.BaseURL+"/v1/comments", token, nil)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode, "search failed: %s", string(respBody))

	var comments []openapi.Comment
	require.NoError(t, json.Unmarshal(respBody, &comments))
	require.Len(t, comments, 20)
	assert.Equal(t, createdIDs[24], comments[0].Id)
	assert.Equal(t, createdIDs[5], comments[19].Id)
}

// TestCreateComment_Unauthorized tests creating a comment without authentication
func TestCreateComment_Unauthorized(t *testing.T) {
	env := setupTestEnv(t)
	defer teardownTestEnv(t, env)

	// Create comment request
	createCommentReq := openapi.CreateCommentRequest{
		Username: "testuser",
		Content:  "This should fail",
	}

	body, err := json.Marshal(createCommentReq)
	require.NoError(t, err)

	// Make request without authentication
	resp, err := http.Post(env.BaseURL+"/v1/comments", "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	// Should return 401 Unauthorized
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "Expected 401 Unauthorized")
}

// TestCreateComment_WrongScope tests creating a comment without required scope
func TestCreateComment_WrongScope(t *testing.T) {
	env := setupTestEnv(t)
	defer teardownTestEnv(t, env)

	// Create a test user with admin role
	testUser := createAdminTestUser(t, env)

	// Get auth token with wrong scope (read instead of create)
	token := getAuthToken(t, env, testUser, scopes.SocialappCommentsRead.String())

	// Create comment request
	createCommentReq := openapi.CreateCommentRequest{
		Username: testUser.Username,
		Content:  "This should fail",
	}

	body, err := json.Marshal(createCommentReq)
	require.NoError(t, err)

	// Make authenticated request
	resp := makeAuthenticatedRequest(t, http.MethodPost, env.BaseURL+"/v1/comments", token, body)
	defer resp.Body.Close()

	// Should return 403 Forbidden
	assert.Equal(t, http.StatusForbidden, resp.StatusCode, "Expected 403 Forbidden for wrong scope")
}
