/*
 * Socialapp
 *
 * Socialapp is a generic social network.
 *
 * API version: 1.0.0
 * Contact: ignacio.gomez.arboleda@gmail.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

// UserApiController binds http requests to an api service and writes the service results to the http response
type UserApiController struct {
	service      UserApiServicer
	errorHandler ErrorHandler
}

// UserApiOption for how the controller is set up.
type UserApiOption func(*UserApiController)

// WithUserApiErrorHandler inject ErrorHandler into controller
func WithUserApiErrorHandler(h ErrorHandler) UserApiOption {
	return func(c *UserApiController) {
		c.errorHandler = h
	}
}

// NewUserApiController creates a default api controller
func NewUserApiController(s UserApiServicer, opts ...UserApiOption) Router {
	controller := &UserApiController{
		service:      s,
		errorHandler: DefaultErrorHandler,
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

// Routes returns all the api routes for the UserApiController
func (c *UserApiController) Routes() Routes {
	return Routes{
		{
			"ChangePassword",
			strings.ToUpper("Post"),
			"/password",
			c.ChangePassword,
		},
		{
			"CreateUser",
			strings.ToUpper("Post"),
			"/users",
			c.CreateUser,
		},
		{
			"DeleteUser",
			strings.ToUpper("Delete"),
			"/users/{username}",
			c.DeleteUser,
		},
		{
			"FollowUser",
			strings.ToUpper("Post"),
			"/users/{followedUsername}/followers/{followerUsername}",
			c.FollowUser,
		},
		{
			"GetFollowingUsers",
			strings.ToUpper("Get"),
			"/users/{username}/following",
			c.GetFollowingUsers,
		},
		{
			"GetRolesForUser",
			strings.ToUpper("Get"),
			"/users/{username}/roles",
			c.GetRolesForUser,
		},
		{
			"GetUserByUsername",
			strings.ToUpper("Get"),
			"/users/{username}",
			c.GetUserByUsername,
		},
		{
			"GetUserComments",
			strings.ToUpper("Get"),
			"/users/{username}/comments",
			c.GetUserComments,
		},
		{
			"GetUserFollowers",
			strings.ToUpper("Get"),
			"/users/{username}/followers",
			c.GetUserFollowers,
		},
		{
			"ListUsers",
			strings.ToUpper("Get"),
			"/users",
			c.ListUsers,
		},
		{
			"ResetPassword",
			strings.ToUpper("Put"),
			"/password",
			c.ResetPassword,
		},
		{
			"UnfollowUser",
			strings.ToUpper("Delete"),
			"/users/{followedUsername}/followers/{followerUsername}",
			c.UnfollowUser,
		},
		{
			"UpdateUser",
			strings.ToUpper("Put"),
			"/users/{username}",
			c.UpdateUser,
		},
	}
}

// ChangePassword - Change password
func (c *UserApiController) ChangePassword(w http.ResponseWriter, r *http.Request) {
	changePasswordRequestParam := ChangePasswordRequest{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&changePasswordRequestParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertChangePasswordRequestRequired(changePasswordRequestParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.ChangePassword(r.Context(), changePasswordRequestParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// CreateUser - Create user
func (c *UserApiController) CreateUser(w http.ResponseWriter, r *http.Request) {
	createUserRequestParam := CreateUserRequest{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&createUserRequestParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertCreateUserRequestRequired(createUserRequestParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.CreateUser(r.Context(), createUserRequestParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// DeleteUser - Deletes a particular user
func (c *UserApiController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	usernameParam := chi.URLParam(r, "username")

	result, err := c.service.DeleteUser(r.Context(), usernameParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// FollowUser - Add a user as a follower
func (c *UserApiController) FollowUser(w http.ResponseWriter, r *http.Request) {
	followedUsernameParam := chi.URLParam(r, "followedUsername")

	followerUsernameParam := chi.URLParam(r, "followerUsername")

	result, err := c.service.FollowUser(r.Context(), followedUsernameParam, followerUsernameParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// GetFollowingUsers - Get all followed users for a user
func (c *UserApiController) GetFollowingUsers(w http.ResponseWriter, r *http.Request) {
	usernameParam := chi.URLParam(r, "username")

	result, err := c.service.GetFollowingUsers(r.Context(), usernameParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// GetRolesForUser - Get all roles for a user
func (c *UserApiController) GetRolesForUser(w http.ResponseWriter, r *http.Request) {
	usernameParam := chi.URLParam(r, "username")

	result, err := c.service.GetRolesForUser(r.Context(), usernameParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// GetUserByUsername - Get a particular user by username
func (c *UserApiController) GetUserByUsername(w http.ResponseWriter, r *http.Request) {
	usernameParam := chi.URLParam(r, "username")

	result, err := c.service.GetUserByUsername(r.Context(), usernameParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// GetUserComments - Gets all comments for a user
func (c *UserApiController) GetUserComments(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	usernameParam := chi.URLParam(r, "username")

	limitParam, err := parseInt32Parameter(query.Get("limit"), false)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	offsetParam, err := parseInt32Parameter(query.Get("offset"), false)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	result, err := c.service.GetUserComments(r.Context(), usernameParam, limitParam, offsetParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// GetUserFollowers - Get all followers for a user
func (c *UserApiController) GetUserFollowers(w http.ResponseWriter, r *http.Request) {
	usernameParam := chi.URLParam(r, "username")

	result, err := c.service.GetUserFollowers(r.Context(), usernameParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// ListUsers - List users
func (c *UserApiController) ListUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	limitParam, err := parseInt32Parameter(query.Get("limit"), false)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	offsetParam, err := parseInt32Parameter(query.Get("offset"), false)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	result, err := c.service.ListUsers(r.Context(), limitParam, offsetParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// ResetPassword - Reset password
func (c *UserApiController) ResetPassword(w http.ResponseWriter, r *http.Request) {
	resetPasswordRequestParam := ResetPasswordRequest{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&resetPasswordRequestParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertResetPasswordRequestRequired(resetPasswordRequestParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.ResetPassword(r.Context(), resetPasswordRequestParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// UnfollowUser - Remove a user as a follower
func (c *UserApiController) UnfollowUser(w http.ResponseWriter, r *http.Request) {
	followedUsernameParam := chi.URLParam(r, "followedUsername")

	followerUsernameParam := chi.URLParam(r, "followerUsername")

	result, err := c.service.UnfollowUser(r.Context(), followedUsernameParam, followerUsernameParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// UpdateUser - Update a user
func (c *UserApiController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	usernameParam := chi.URLParam(r, "username")

	userParam := User{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&userParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertUserRequired(userParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.UpdateUser(r.Context(), usernameParam, userParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}
