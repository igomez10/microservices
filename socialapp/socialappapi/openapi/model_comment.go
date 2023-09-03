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
	"time"
)

type Comment struct {
	Id int64 `json:"id,omitempty"`

	Content string `json:"content"`

	LikeCount int64 `json:"like_count,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`

	Username string `json:"username"`
}

// AssertCommentRequired checks if the required fields are not zero-ed
func AssertCommentRequired(obj Comment) error {
	elements := map[string]interface{}{
		"content":  obj.Content,
		"username": obj.Username,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertCommentConstraints checks if the values respects the defined constraints
func AssertCommentConstraints(obj Comment) error {
	return nil
}
