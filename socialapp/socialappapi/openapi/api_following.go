/*
 * Socialapp
 *
 * Socialapp is a generic social network.
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

// FollowingApiController binds http requests to an api service and writes the service results to the http response
type FollowingApiController struct {
	service      FollowingApiServicer
	errorHandler ErrorHandler
}

// FollowingApiOption for how the controller is set up.
type FollowingApiOption func(*FollowingApiController)

// WithFollowingApiErrorHandler inject ErrorHandler into controller
func WithFollowingApiErrorHandler(h ErrorHandler) FollowingApiOption {
	return func(c *FollowingApiController) {
		c.errorHandler = h
	}
}

// NewFollowingApiController creates a default api controller
func NewFollowingApiController(s FollowingApiServicer, opts ...FollowingApiOption) Router {
	controller := &FollowingApiController{
		service:      s,
		errorHandler: DefaultErrorHandler,
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

// Routes returns all the api routes for the FollowingApiController
func (c *FollowingApiController) Routes() Routes {
	return Routes{
		{
			"GetUserFollowers",
			strings.ToUpper("Get"),
			"/users/{username}/followers",
			c.GetUserFollowers,
		},
	}
}

// GetUserFollowers - Get all followers for a user
func (c *FollowingApiController) GetUserFollowers(w http.ResponseWriter, r *http.Request) {
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
