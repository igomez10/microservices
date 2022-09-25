# UserApi

All URIs are relative to **

Method | HTTP request | Description
------------- | ------------- | -------------
[**changePassword**](UserApi.md#changePassword) | **POST** /password | Change password
[**createUser**](UserApi.md#createUser) | **POST** /users | Create a new user
[**deleteUser**](UserApi.md#deleteUser) | **DELETE** /users/{username} | Deletes a particular user
[**followUser**](UserApi.md#followUser) | **POST** /users/{followedUsername}/followers/{followerUsername} | Add a user as a follower
[**getFollowingUsers**](UserApi.md#getFollowingUsers) | **GET** /users/{username}/following | Get all followed users for a user
[**getUserByUsername**](UserApi.md#getUserByUsername) | **GET** /users/{username} | Get a particular user by username
[**getUserComments**](UserApi.md#getUserComments) | **GET** /users/{username}/comments | Gets all comments for a user
[**getUserFollowers**](UserApi.md#getUserFollowers) | **GET** /users/{username}/followers | Get all followers for a user
[**listUsers**](UserApi.md#listUsers) | **GET** /users | Returns all the users
[**resetPassword**](UserApi.md#resetPassword) | **PUT** /password | Reset password
[**unfollowUser**](UserApi.md#unfollowUser) | **DELETE** /users/{followedUsername}/followers/{followerUsername} | Remove a user as a follower
[**updateUser**](UserApi.md#updateUser) | **PUT** /users/{username} | Update a user



## changePassword

Change password

### Example

```bash
socialapp-cli changePassword
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changePasswordRequest** | [**ChangePasswordRequest**](ChangePasswordRequest.md) | Change password |

### Return type

[**User**](User.md)

### Authorization

[BasicAuth](../README.md#BasicAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## createUser

Create a new user

### Example

```bash
socialapp-cli createUser
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **createUserRequest** | [**CreateUserRequest**](CreateUserRequest.md) | Create a new user |

### Return type

[**User**](User.md)

### Authorization

[BasicAuth](../README.md#BasicAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## deleteUser

Deletes a particular user

### Example

```bash
socialapp-cli deleteUser username=value
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **username** | **string** | username of the user | [default to null]

### Return type

[**User**](User.md)

### Authorization

[BasicAuth](../README.md#BasicAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not Applicable
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## followUser

Add a user as a follower

### Example

```bash
socialapp-cli followUser followedUsername=value followerUsername=value
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **followedUsername** | **string** | username of the user | [default to null]
 **followerUsername** | **string** | username of the follower | [default to null]

### Return type

(empty response body)

### Authorization

[BasicAuth](../README.md#BasicAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not Applicable
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## getFollowingUsers

Get all followed users for a user

### Example

```bash
socialapp-cli getFollowingUsers username=value
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **username** | **string** | username of the user | [default to null]

### Return type

[**array[User]**](User.md)

### Authorization

[BasicAuth](../README.md#BasicAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not Applicable
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## getUserByUsername

Get a particular user by username

### Example

```bash
socialapp-cli getUserByUsername username=value
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **username** | **string** | username of the user | [default to null]

### Return type

[**User**](User.md)

### Authorization

[BasicAuth](../README.md#BasicAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not Applicable
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## getUserComments

Gets all comments for a user

### Example

```bash
socialapp-cli getUserComments username=value  limit=value  offset=value
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **username** | **string** | username of the user | [default to null]
 **limit** | **integer** | How many items to return at one time (max 100) | [optional] [default to null]
 **offset** | **integer** | The number of items to skip before starting to collect the result set | [optional] [default to null]

### Return type

[**array[Comment]**](Comment.md)

### Authorization

[BasicAuth](../README.md#BasicAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not Applicable
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## getUserFollowers

Get all followers for a user

### Example

```bash
socialapp-cli getUserFollowers username=value
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **username** | **string** | username of the user | [default to null]

### Return type

[**array[User]**](User.md)

### Authorization

[BasicAuth](../README.md#BasicAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not Applicable
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## listUsers

Returns all the users

### Example

```bash
socialapp-cli listUsers  limit=value  offset=value
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **limit** | **integer** | Maximum number of users to return | [optional] [default to null]
 **offset** | **integer** | Pagination offset | [optional] [default to null]

### Return type

[**array[User]**](User.md)

### Authorization

[BasicAuth](../README.md#BasicAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not Applicable
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## resetPassword

Reset password

### Example

```bash
socialapp-cli resetPassword
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **resetPasswordRequest** | [**ResetPasswordRequest**](ResetPasswordRequest.md) | Reset password |

### Return type

[**User**](User.md)

### Authorization

[BasicAuth](../README.md#BasicAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## unfollowUser

Remove a user as a follower

### Example

```bash
socialapp-cli unfollowUser followedUsername=value followerUsername=value
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **followedUsername** | **string** | username of the user | [default to null]
 **followerUsername** | **string** | username of the follower | [default to null]

### Return type

(empty response body)

### Authorization

[BasicAuth](../README.md#BasicAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not Applicable
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## updateUser

Update a user

### Example

```bash
socialapp-cli updateUser username=value
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **username** | **string** | username of the user | [default to null]
 **user** | [**User**](User.md) | Update a user | [optional]

### Return type

[**User**](User.md)

### Authorization

[BasicAuth](../README.md#BasicAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

