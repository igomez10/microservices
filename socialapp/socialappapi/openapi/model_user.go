// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

/*
 * Socialapp
 *
 * Socialapp is a generic social network.
 *
 * API version: 1.0.0
 * Contact: ignacio.gomez.arboleda@gmail.com
 */

package openapi

import (
	"time"
)

type User struct {
	Id int64 `json:"id,omitempty"`

	Username string `json:"username"`

	FirstName string `json:"first_name"`

	LastName string `json:"last_name"`

	Email string `json:"email"`

	CreatedAt time.Time `json:"created_at,omitempty"`
}

// AssertUserRequired checks if the required fields are not zero-ed
func AssertUserRequired(obj User) error {
	elements := map[string]interface{}{
		"username":   obj.Username,
		"first_name": obj.FirstName,
		"last_name":  obj.LastName,
		"email":      obj.Email,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertUserConstraints checks if the values respects the defined constraints
func AssertUserConstraints(obj User) error {
	return nil
}
