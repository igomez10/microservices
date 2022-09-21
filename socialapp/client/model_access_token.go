/*
Socialapp

Socialapp is a generic social network.

API version: 1.0.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package client

import (
	"encoding/json"
	"time"
)

// AccessToken struct for AccessToken
type AccessToken struct {
	AccessToken string     `json:"access_token"`
	Scopes      []string   `json:"scopes,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}

// NewAccessToken instantiates a new AccessToken object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAccessToken(accessToken string) *AccessToken {
	this := AccessToken{}
	this.AccessToken = accessToken
	return &this
}

// NewAccessTokenWithDefaults instantiates a new AccessToken object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAccessTokenWithDefaults() *AccessToken {
	this := AccessToken{}
	return &this
}

// GetAccessToken returns the AccessToken field value
func (o *AccessToken) GetAccessToken() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.AccessToken
}

// GetAccessTokenOk returns a tuple with the AccessToken field value
// and a boolean to check if the value has been set.
func (o *AccessToken) GetAccessTokenOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.AccessToken, true
}

// SetAccessToken sets field value
func (o *AccessToken) SetAccessToken(v string) {
	o.AccessToken = v
}

// GetScopes returns the Scopes field value if set, zero value otherwise.
func (o *AccessToken) GetScopes() []string {
	if o == nil || o.Scopes == nil {
		var ret []string
		return ret
	}
	return o.Scopes
}

// GetScopesOk returns a tuple with the Scopes field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AccessToken) GetScopesOk() ([]string, bool) {
	if o == nil || o.Scopes == nil {
		return nil, false
	}
	return o.Scopes, true
}

// HasScopes returns a boolean if a field has been set.
func (o *AccessToken) HasScopes() bool {
	if o != nil && o.Scopes != nil {
		return true
	}

	return false
}

// SetScopes gets a reference to the given []string and assigns it to the Scopes field.
func (o *AccessToken) SetScopes(v []string) {
	o.Scopes = v
}

// GetExpiresAt returns the ExpiresAt field value if set, zero value otherwise.
func (o *AccessToken) GetExpiresAt() time.Time {
	if o == nil || o.ExpiresAt == nil {
		var ret time.Time
		return ret
	}
	return *o.ExpiresAt
}

// GetExpiresAtOk returns a tuple with the ExpiresAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AccessToken) GetExpiresAtOk() (*time.Time, bool) {
	if o == nil || o.ExpiresAt == nil {
		return nil, false
	}
	return o.ExpiresAt, true
}

// HasExpiresAt returns a boolean if a field has been set.
func (o *AccessToken) HasExpiresAt() bool {
	if o != nil && o.ExpiresAt != nil {
		return true
	}

	return false
}

// SetExpiresAt gets a reference to the given time.Time and assigns it to the ExpiresAt field.
func (o *AccessToken) SetExpiresAt(v time.Time) {
	o.ExpiresAt = &v
}

func (o AccessToken) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["access_token"] = o.AccessToken
	}
	if o.Scopes != nil {
		toSerialize["scopes"] = o.Scopes
	}
	if o.ExpiresAt != nil {
		toSerialize["expires_at"] = o.ExpiresAt
	}
	return json.Marshal(toSerialize)
}

type NullableAccessToken struct {
	value *AccessToken
	isSet bool
}

func (v NullableAccessToken) Get() *AccessToken {
	return v.value
}

func (v *NullableAccessToken) Set(val *AccessToken) {
	v.value = val
	v.isSet = true
}

func (v NullableAccessToken) IsSet() bool {
	return v.isSet
}

func (v *NullableAccessToken) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableAccessToken(val *AccessToken) *NullableAccessToken {
	return &NullableAccessToken{value: val, isSet: true}
}

func (v NullableAccessToken) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableAccessToken) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}