/*
Socialapp

Socialapp is a generic social network.

API version: 1.0.0
Contact: ignacio.gomez.arboleda@gmail.com
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package client

import (
	"encoding/json"
	"time"
)

// checks if the Scope type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Scope{}

// Scope struct for Scope
type Scope struct {
	Id          *int64     `json:"id,omitempty"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
}

// NewScope instantiates a new Scope object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewScope(name string, description string) *Scope {
	this := Scope{}
	this.Name = name
	this.Description = description
	return &this
}

// NewScopeWithDefaults instantiates a new Scope object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewScopeWithDefaults() *Scope {
	this := Scope{}
	return &this
}

// GetId returns the Id field value if set, zero value otherwise.
func (o *Scope) GetId() int64 {
	if o == nil || isNil(o.Id) {
		var ret int64
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Scope) GetIdOk() (*int64, bool) {
	if o == nil || isNil(o.Id) {
		return nil, false
	}
	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *Scope) HasId() bool {
	if o != nil && !isNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given int64 and assigns it to the Id field.
func (o *Scope) SetId(v int64) {
	o.Id = &v
}

// GetName returns the Name field value
func (o *Scope) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *Scope) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *Scope) SetName(v string) {
	o.Name = v
}

// GetDescription returns the Description field value
func (o *Scope) GetDescription() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Description
}

// GetDescriptionOk returns a tuple with the Description field value
// and a boolean to check if the value has been set.
func (o *Scope) GetDescriptionOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Description, true
}

// SetDescription sets field value
func (o *Scope) SetDescription(v string) {
	o.Description = v
}

// GetCreatedAt returns the CreatedAt field value if set, zero value otherwise.
func (o *Scope) GetCreatedAt() time.Time {
	if o == nil || isNil(o.CreatedAt) {
		var ret time.Time
		return ret
	}
	return *o.CreatedAt
}

// GetCreatedAtOk returns a tuple with the CreatedAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Scope) GetCreatedAtOk() (*time.Time, bool) {
	if o == nil || isNil(o.CreatedAt) {
		return nil, false
	}
	return o.CreatedAt, true
}

// HasCreatedAt returns a boolean if a field has been set.
func (o *Scope) HasCreatedAt() bool {
	if o != nil && !isNil(o.CreatedAt) {
		return true
	}

	return false
}

// SetCreatedAt gets a reference to the given time.Time and assigns it to the CreatedAt field.
func (o *Scope) SetCreatedAt(v time.Time) {
	o.CreatedAt = &v
}

func (o Scope) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Scope) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !isNil(o.Id) {
		toSerialize["id"] = o.Id
	}
	toSerialize["name"] = o.Name
	toSerialize["description"] = o.Description
	if !isNil(o.CreatedAt) {
		toSerialize["created_at"] = o.CreatedAt
	}
	return toSerialize, nil
}

type NullableScope struct {
	value *Scope
	isSet bool
}

func (v NullableScope) Get() *Scope {
	return v.value
}

func (v *NullableScope) Set(val *Scope) {
	v.value = val
	v.isSet = true
}

func (v NullableScope) IsSet() bool {
	return v.isSet
}

func (v *NullableScope) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableScope(val *Scope) *NullableScope {
	return &NullableScope{value: val, isSet: true}
}

func (v NullableScope) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableScope) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
