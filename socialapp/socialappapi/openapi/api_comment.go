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

// CommentAPIController binds http requests to an api service and writes the service results to the http response
type CommentAPIController struct {
	service      CommentAPIServicer
	errorHandler ErrorHandler
}

// CommentAPIOption for how the controller is set up.
type CommentAPIOption func(*CommentAPIController)

// WithCommentAPIErrorHandler inject ErrorHandler into controller
func WithCommentAPIErrorHandler(h ErrorHandler) CommentAPIOption {
	return func(c *CommentAPIController) {
		c.errorHandler = h
	}
}

// NewCommentAPIController creates a default api controller
func NewCommentAPIController(s CommentAPIServicer, opts ...CommentAPIOption) Router {
	controller := &CommentAPIController{
		service:      s,
		errorHandler: DefaultErrorHandler,
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

// Routes returns all the api routes for the CommentAPIController
func (c *CommentAPIController) Routes() Routes {
	return Routes{
		"CreateComment": Route{
			strings.ToUpper("Post"),
			"/v1/comments",
			c.CreateComment,
		},
		"GetComment": Route{
			strings.ToUpper("Get"),
			"/v1/comments/{id}",
			c.GetComment,
		},
		"GetUserComments": Route{
			strings.ToUpper("Get"),
			"/v1/users/{username}/comments",
			c.GetUserComments,
		},
		"GetUserFeed": Route{
			strings.ToUpper("Get"),
			"/v1/feed",
			c.GetUserFeed,
		},
	}
}

// CreateComment - Create a new comment
func (c *CommentAPIController) CreateComment(w http.ResponseWriter, r *http.Request) {
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
	if err := AssertCommentConstraints(commentParam); err != nil {
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
func (c *CommentAPIController) GetComment(w http.ResponseWriter, r *http.Request) {
	idParam, err := parseNumericParameter[int32](
		chi.URLParam(r, "id"),
		WithRequire[int32](parseInt32),
	)
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
func (c *CommentAPIController) GetUserComments(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	usernameParam := chi.URLParam(r, "username")
	limitParam, err := parseNumericParameter[int32](
		query.Get("limit"),
		WithDefaultOrParse[int32](20, parseInt32),
		WithMinimum[int32](1),
		WithMaximum[int32](100),
	)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	offsetParam, err := parseNumericParameter[int32](
		query.Get("offset"),
		WithParse[int32](parseInt32),
	)
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
func (c *CommentAPIController) GetUserFeed(w http.ResponseWriter, r *http.Request) {
	result, err := c.service.GetUserFeed(r.Context())
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, result.Headers, w)
}
