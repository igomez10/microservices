// Package events defines the curated domain events emitted to the event store
// (the transactional outbox). Each event is an explicit, stable contract that is
// serialized as the event payload -- deliberately NOT the API request/response
// DTOs or database models, so the event history is decoupled from those and
// never leaks sensitive fields (e.g. passwords) into the store or downstream
// Kafka topics.
package events

// Aggregate types identify the kind of domain entity an event applies to.
const (
	AggregateTypeUser    = "User"
	AggregateTypeComment = "Comment"
	AggregateTypeRole    = "Role"
	AggregateTypeScope   = "Scope"
)

// Event types name the fact that occurred.
const (
	EventTypeUserCreated         = "UserCreated"
	EventTypeUserUpdated         = "UserUpdated"
	EventTypeUserDeleted         = "UserDeleted"
	EventTypeUserFollowed        = "UserFollowed"
	EventTypeUserUnfollowed      = "UserUnfollowed"
	EventTypeUserPasswordChanged = "UserPasswordChanged"
	EventTypeUserRolesUpdated    = "UserRolesUpdated"

	EventTypeCommentCreated = "CommentCreated"

	EventTypeRoleCreated       = "RoleCreated"
	EventTypeRoleUpdated       = "RoleUpdated"
	EventTypeRoleDeleted       = "RoleDeleted"
	EventTypeRoleScopesAdded   = "RoleScopesAdded"
	EventTypeRoleScopeRemoved  = "RoleScopeRemoved"

	EventTypeScopeCreated = "ScopeCreated"
	EventTypeScopeUpdated = "ScopeUpdated"
	EventTypeScopeDeleted = "ScopeDeleted"
)

// --- User events ---

type UserCreated struct {
	UserID    int64  `json:"user_id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type UserUpdated struct {
	UserID    int64  `json:"user_id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type UserDeleted struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
}

type UserFollowed struct {
	FollowerID       int64  `json:"follower_id"`
	FollowerUsername string `json:"follower_username"`
	FollowedID       int64  `json:"followed_id"`
	FollowedUsername string `json:"followed_username"`
}

type UserUnfollowed struct {
	FollowerID       int64  `json:"follower_id"`
	FollowerUsername string `json:"follower_username"`
	FollowedID       int64  `json:"followed_id"`
	FollowedUsername string `json:"followed_username"`
}

// UserPasswordChanged intentionally carries no password material.
type UserPasswordChanged struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
}

type UserRolesUpdated struct {
	UserID   int64    `json:"user_id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
}

// --- Comment events ---

type CommentCreated struct {
	CommentID      int64  `json:"comment_id"`
	AuthorID       int64  `json:"author_id"`
	AuthorUsername string `json:"author_username"`
	Content        string `json:"content"`
}

// --- Role events ---

type RoleCreated struct {
	RoleID      int64  `json:"role_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type RoleUpdated struct {
	RoleID      int64  `json:"role_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type RoleDeleted struct {
	RoleID int64  `json:"role_id"`
	Name   string `json:"name"`
}

type RoleScopesAdded struct {
	RoleID int64    `json:"role_id"`
	Scopes []string `json:"scopes"`
}

type RoleScopeRemoved struct {
	RoleID  int64 `json:"role_id"`
	ScopeID int64 `json:"scope_id"`
}

// --- Scope events ---

type ScopeCreated struct {
	ScopeID     int64  `json:"scope_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ScopeUpdated struct {
	ScopeID     int64  `json:"scope_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ScopeDeleted struct {
	ScopeID int64  `json:"scope_id"`
	Name    string `json:"name"`
}
