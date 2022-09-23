# CreateUserRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Username** | **string** |  | 
**Password** | **string** |  | 
**FirstName** | **string** |  | 
**LastName** | **string** |  | 
**Email** | **string** |  | 
**CreatedAt** | Pointer to **time.Time** |  | [optional] 

## Methods

### NewCreateUserRequest

`func NewCreateUserRequest(username string, password string, firstName string, lastName string, email string, ) *CreateUserRequest`

NewCreateUserRequest instantiates a new CreateUserRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreateUserRequestWithDefaults

`func NewCreateUserRequestWithDefaults() *CreateUserRequest`

NewCreateUserRequestWithDefaults instantiates a new CreateUserRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *CreateUserRequest) GetId() int64`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *CreateUserRequest) GetIdOk() (*int64, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *CreateUserRequest) SetId(v int64)`

SetId sets Id field to given value.

### HasId

`func (o *CreateUserRequest) HasId() bool`

HasId returns a boolean if a field has been set.

### GetUsername

`func (o *CreateUserRequest) GetUsername() string`

GetUsername returns the Username field if non-nil, zero value otherwise.

### GetUsernameOk

`func (o *CreateUserRequest) GetUsernameOk() (*string, bool)`

GetUsernameOk returns a tuple with the Username field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUsername

`func (o *CreateUserRequest) SetUsername(v string)`

SetUsername sets Username field to given value.


### GetPassword

`func (o *CreateUserRequest) GetPassword() string`

GetPassword returns the Password field if non-nil, zero value otherwise.

### GetPasswordOk

`func (o *CreateUserRequest) GetPasswordOk() (*string, bool)`

GetPasswordOk returns a tuple with the Password field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPassword

`func (o *CreateUserRequest) SetPassword(v string)`

SetPassword sets Password field to given value.


### GetFirstName

`func (o *CreateUserRequest) GetFirstName() string`

GetFirstName returns the FirstName field if non-nil, zero value otherwise.

### GetFirstNameOk

`func (o *CreateUserRequest) GetFirstNameOk() (*string, bool)`

GetFirstNameOk returns a tuple with the FirstName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFirstName

`func (o *CreateUserRequest) SetFirstName(v string)`

SetFirstName sets FirstName field to given value.


### GetLastName

`func (o *CreateUserRequest) GetLastName() string`

GetLastName returns the LastName field if non-nil, zero value otherwise.

### GetLastNameOk

`func (o *CreateUserRequest) GetLastNameOk() (*string, bool)`

GetLastNameOk returns a tuple with the LastName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastName

`func (o *CreateUserRequest) SetLastName(v string)`

SetLastName sets LastName field to given value.


### GetEmail

`func (o *CreateUserRequest) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *CreateUserRequest) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *CreateUserRequest) SetEmail(v string)`

SetEmail sets Email field to given value.


### GetCreatedAt

`func (o *CreateUserRequest) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *CreateUserRequest) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *CreateUserRequest) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.

### HasCreatedAt

`func (o *CreateUserRequest) HasCreatedAt() bool`

HasCreatedAt returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


