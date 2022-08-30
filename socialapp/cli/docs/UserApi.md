# UserApi

All URIs are relative to **

Method | HTTP request | Description
------------- | ------------- | -------------
[**createUser**](UserApi.md#createUser) | **POST** /users | Create a new user
[**deleteUser**](UserApi.md#deleteUser) | **DELETE** /users/{username} | Deletes a particular user
[**getUserByUsername**](UserApi.md#getUserByUsername) | **GET** /users/{username} | Get a particular user by username
[**getUserComments**](UserApi.md#getUserComments) | **GET** /users/{username}/comments | Gets all comments for a user
[**listUsers**](UserApi.md#listUsers) | **GET** /users | Returns all the users
[**updateUser**](UserApi.md#updateUser) | **PUT** /users/{username} | Update a user



## createUser

Create a new user

### Example

```bash
socialapp-cli createUser
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **user** | [**User**](User.md) | Create a new user |

### Return type

[**User**](User.md)

### Authorization

[BasicAuth](../README.md#BasicAuth)

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

[BasicAuth](../README.md#BasicAuth)

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

[BasicAuth](../README.md#BasicAuth)

### HTTP request headers

- **Content-Type**: Not Applicable
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## getUserComments

Gets all comments for a user

### Example

```bash
socialapp-cli getUserComments username=value
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **username** | **string** | username of the user | [default to null]

### Return type

[**Comment**](Comment.md)

### Authorization

[BasicAuth](../README.md#BasicAuth)

### HTTP request headers

- **Content-Type**: Not Applicable
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## listUsers

Returns all the users

### Example

```bash
socialapp-cli listUsers
```

### Parameters

This endpoint does not need any parameter.

### Return type

[**array[User]**](User.md)

### Authorization

[BasicAuth](../README.md#BasicAuth)

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

[BasicAuth](../README.md#BasicAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

