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

type CreateUserResponse struct {
	Id int64 `json:"id"`

	Username string `json:"username"`

	FirstName string `json:"first_name"`

	LastName string `json:"last_name"`

	Email string `json:"email"`

	CreatedAt time.Time `json:"created_at"`
}

// AssertCreateUserResponseRequired checks if the required fields are not zero-ed
func AssertCreateUserResponseRequired(obj CreateUserResponse) error {
	elements := map[string]interface{}{
		"id":         obj.Id,
		"username":   obj.Username,
		"first_name": obj.FirstName,
		"last_name":  obj.LastName,
		"email":      obj.Email,
		"created_at": obj.CreatedAt,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertRecurseCreateUserResponseRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of CreateUserResponse (e.g. [][]CreateUserResponse), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseCreateUserResponseRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aCreateUserResponse, ok := obj.(CreateUserResponse)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertCreateUserResponseRequired(aCreateUserResponse)
	})
}