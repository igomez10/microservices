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

type Url struct {
	Url string `json:"url"`

	Alias string `json:"alias"`

	CreatedAt time.Time `json:"created_at,omitempty"`

	UpdatedAt time.Time `json:"updated_at,omitempty"`

	DeletedAt time.Time `json:"deleted_at,omitempty"`
}

// AssertUrlRequired checks if the required fields are not zero-ed
func AssertUrlRequired(obj Url) error {
	elements := map[string]interface{}{
		"url":   obj.Url,
		"alias": obj.Alias,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertUrlConstraints checks if the values respects the defined constraints
func AssertUrlConstraints(obj Url) error {
	return nil
}
