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

// CommentApiController binds http requests to an api service and writes the service results to the http response
type CommentApiController struct {
	service      CommentApiServicer
	errorHandler ErrorHandler
}

// CommentApiOption for how the controller is set up.
type CommentApiOption func(*CommentApiController)

// WithCommentApiErrorHandler inject ErrorHandler into controller
func WithCommentApiErrorHandler(h ErrorHandler) CommentApiOption {
	return func(c *CommentApiController) {
		c.errorHandler = h
	}
}

// NewCommentApiController creates a default api controller
func NewCommentApiController(s CommentApiServicer, opts ...CommentApiOption) Router {
	controller := &CommentApiController{
		service:      s,
		errorHandler: DefaultErrorHandler,
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

// Routes returns all the api routes for the CommentApiController
func (c *CommentApiController) Routes() Routes {
	return Routes{
		{
			"CreateComment",
			strings.ToUpper("Post"),
			"/v1/comments",
			c.CreateComment,
		},
		{
			"GetComment",
			strings.ToUpper("Get"),
			"/v1/comments/{id}",
			c.GetComment,
		},
		{
			"GetUserComments",
			strings.ToUpper("Get"),
			"/v1/users/{username}/comments",
			c.GetUserComments,
		},
		{
			"GetUserFeed",
			strings.ToUpper("Get"),
			"/v1/users/{username}/feed",
			c.GetUserFeed,
		},
	}
}

// CreateComment - Create a new comment
func (c *CommentApiController) CreateComment(w http.ResponseWriter, r *http.Request) {
	commentParam := Comment{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&commentParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertCommentRequired(commentParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.CreateComment(r.Context(), commentParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, result.Headers, w)

}

// GetComment - Returns details about a particular comment
func (c *CommentApiController) GetComment(w http.ResponseWriter, r *http.Request) {
	idParam, err := parseInt32Parameter(chi.URLParam(r, "id"), true)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}

	result, err := c.service.GetComment(r.Context(), idParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, result.Headers, w)

}

// GetUserComments - Gets all comments for a user
func (c *CommentApiController) GetUserComments(w http.ResponseWriter, r *http.Request) {
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
	EncodeJSONResponse(result.Body, &result.Code, result.Headers, w)

}

// GetUserFeed - Returns a users feed
func (c *CommentApiController) GetUserFeed(w http.ResponseWriter, r *http.Request) {
	usernameParam := chi.URLParam(r, "username")

	result, err := c.service.GetUserFeed(r.Context(), usernameParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, result.Headers, w)

}
