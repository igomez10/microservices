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
	"context"
	"net/http"
)

// AuthenticationApiRouter defines the required methods for binding the api requests to a responses for the AuthenticationApi
// The AuthenticationApiRouter implementation should parse necessary information from the http request,
// pass the data to a AuthenticationApiServicer to perform the required actions, then write the service results to the http response.
type AuthenticationApiRouter interface {
	GetAccessToken(http.ResponseWriter, *http.Request)
}

// CommentApiRouter defines the required methods for binding the api requests to a responses for the CommentApi
// The CommentApiRouter implementation should parse necessary information from the http request,
// pass the data to a CommentApiServicer to perform the required actions, then write the service results to the http response.
type CommentApiRouter interface {
	CreateComment(http.ResponseWriter, *http.Request)
	GetComment(http.ResponseWriter, *http.Request)
	GetUserComments(http.ResponseWriter, *http.Request)
	GetUserFeed(http.ResponseWriter, *http.Request)
}

// FollowingApiRouter defines the required methods for binding the api requests to a responses for the FollowingApi
// The FollowingApiRouter implementation should parse necessary information from the http request,
// pass the data to a FollowingApiServicer to perform the required actions, then write the service results to the http response.
type FollowingApiRouter interface {
	GetUserFollowers(http.ResponseWriter, *http.Request)
}

// RoleApiRouter defines the required methods for binding the api requests to a responses for the RoleApi
// The RoleApiRouter implementation should parse necessary information from the http request,
// pass the data to a RoleApiServicer to perform the required actions, then write the service results to the http response.
type RoleApiRouter interface {
	AddScopeToRole(http.ResponseWriter, *http.Request)
	CreateRole(http.ResponseWriter, *http.Request)
	DeleteRole(http.ResponseWriter, *http.Request)
	GetRole(http.ResponseWriter, *http.Request)
	ListRoles(http.ResponseWriter, *http.Request)
	ListScopesForRole(http.ResponseWriter, *http.Request)
	RemoveScopeFromRole(http.ResponseWriter, *http.Request)
	UpdateRole(http.ResponseWriter, *http.Request)
}

// ScopeApiRouter defines the required methods for binding the api requests to a responses for the ScopeApi
// The ScopeApiRouter implementation should parse necessary information from the http request,
// pass the data to a ScopeApiServicer to perform the required actions, then write the service results to the http response.
type ScopeApiRouter interface {
	CreateScope(http.ResponseWriter, *http.Request)
	DeleteScope(http.ResponseWriter, *http.Request)
	GetScope(http.ResponseWriter, *http.Request)
	ListScopes(http.ResponseWriter, *http.Request)
	UpdateScope(http.ResponseWriter, *http.Request)
}

// UserApiRouter defines the required methods for binding the api requests to a responses for the UserApi
// The UserApiRouter implementation should parse necessary information from the http request,
// pass the data to a UserApiServicer to perform the required actions, then write the service results to the http response.
type UserApiRouter interface {
	ChangePassword(http.ResponseWriter, *http.Request)
	CreateUser(http.ResponseWriter, *http.Request)
	DeleteUser(http.ResponseWriter, *http.Request)
	FollowUser(http.ResponseWriter, *http.Request)
	GetFollowingUsers(http.ResponseWriter, *http.Request)
	GetRolesForUser(http.ResponseWriter, *http.Request)
	GetUserByUsername(http.ResponseWriter, *http.Request)
	GetUserComments(http.ResponseWriter, *http.Request)
	GetUserFollowers(http.ResponseWriter, *http.Request)
	ListUsers(http.ResponseWriter, *http.Request)
	ResetPassword(http.ResponseWriter, *http.Request)
	UnfollowUser(http.ResponseWriter, *http.Request)
	UpdateUser(http.ResponseWriter, *http.Request)
}

// AuthenticationApiServicer defines the api actions for the AuthenticationApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type AuthenticationApiServicer interface {
	GetAccessToken(context.Context) (ImplResponse, error)
}

// CommentApiServicer defines the api actions for the CommentApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type CommentApiServicer interface {
	CreateComment(context.Context, Comment) (ImplResponse, error)
	GetComment(context.Context, int32) (ImplResponse, error)
	GetUserComments(context.Context, string, int32, int32) (ImplResponse, error)
	GetUserFeed(context.Context, string) (ImplResponse, error)
}

// FollowingApiServicer defines the api actions for the FollowingApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type FollowingApiServicer interface {
	GetUserFollowers(context.Context, string) (ImplResponse, error)
}

// RoleApiServicer defines the api actions for the RoleApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type RoleApiServicer interface {
	AddScopeToRole(context.Context, int32, []string) (ImplResponse, error)
	CreateRole(context.Context, Role) (ImplResponse, error)
	DeleteRole(context.Context, int32) (ImplResponse, error)
	GetRole(context.Context, int32) (ImplResponse, error)
	ListRoles(context.Context, int32, int32) (ImplResponse, error)
	ListScopesForRole(context.Context, int32, int32, int32) (ImplResponse, error)
	RemoveScopeFromRole(context.Context, int32, int32) (ImplResponse, error)
	UpdateRole(context.Context, int32, Role) (ImplResponse, error)
}

// ScopeApiServicer defines the api actions for the ScopeApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type ScopeApiServicer interface {
	CreateScope(context.Context, Scope) (ImplResponse, error)
	DeleteScope(context.Context, int32) (ImplResponse, error)
	GetScope(context.Context, int32) (ImplResponse, error)
	ListScopes(context.Context, int32, int32) (ImplResponse, error)
	UpdateScope(context.Context, int32, Scope) (ImplResponse, error)
}

// UserApiServicer defines the api actions for the UserApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type UserApiServicer interface {
	ChangePassword(context.Context, ChangePasswordRequest) (ImplResponse, error)
	CreateUser(context.Context, CreateUserRequest) (ImplResponse, error)
	DeleteUser(context.Context, string) (ImplResponse, error)
	FollowUser(context.Context, string, string) (ImplResponse, error)
	GetFollowingUsers(context.Context, string) (ImplResponse, error)
	GetRolesForUser(context.Context, string) (ImplResponse, error)
	GetUserByUsername(context.Context, string) (ImplResponse, error)
	GetUserComments(context.Context, string, int32, int32) (ImplResponse, error)
	GetUserFollowers(context.Context, string) (ImplResponse, error)
	ListUsers(context.Context, int32, int32) (ImplResponse, error)
	ResetPassword(context.Context, ResetPasswordRequest) (ImplResponse, error)
	UnfollowUser(context.Context, string, string) (ImplResponse, error)
	UpdateUser(context.Context, string, User) (ImplResponse, error)
}
