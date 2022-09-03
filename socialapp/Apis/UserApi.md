# UserApi

All URIs are relative to *http://localhost:8080*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**createUser**](UserApi.md#createUser) | **POST** /users | Create a new user |
| [**deleteUser**](UserApi.md#deleteUser) | **DELETE** /users/{username} | Deletes a particular user |
| [**followUser**](UserApi.md#followUser) | **POST** /users/{followedUsername}/followers/{followerUsername} | Add a user as a follower |
| [**getFollowingUsers**](UserApi.md#getFollowingUsers) | **GET** /users/{username}/following | Get all followed users for a user |
| [**getUserByUsername**](UserApi.md#getUserByUsername) | **GET** /users/{username} | Get a particular user by username |
| [**getUserComments**](UserApi.md#getUserComments) | **GET** /users/{username}/comments | Gets all comments for a user |
| [**getUserFollowers**](UserApi.md#getUserFollowers) | **GET** /users/{username}/followers | Get all followers for a user |
| [**listUsers**](UserApi.md#listUsers) | **GET** /users | Returns all the users |
| [**unfollowUser**](UserApi.md#unfollowUser) | **DELETE** /users/{followedUsername}/followers/{followerUsername} | Remove a user as a follower |
| [**updateUser**](UserApi.md#updateUser) | **PUT** /users/{username} | Update a user |


<a name="createUser"></a>
# **createUser**
> User createUser(User)

Create a new user

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **User** | [**User**](../Models/User.md)| Create a new user | |

### Return type

[**User**](../Models/User.md)

### Authorization

[BasicAuth](../README.md#BasicAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

<a name="deleteUser"></a>
# **deleteUser**
> User deleteUser(username)

Deletes a particular user

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **username** | **String**| username of the user | [default to null] |

### Return type

[**User**](../Models/User.md)

### Authorization

[BasicAuth](../README.md#BasicAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="followUser"></a>
# **followUser**
> followUser(followedUsername, followerUsername)

Add a user as a follower

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **followedUsername** | **String**| username of the user | [default to null] |
| **followerUsername** | **String**| username of the follower | [default to null] |

### Return type

null (empty response body)

### Authorization

[BasicAuth](../README.md#BasicAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="getFollowingUsers"></a>
# **getFollowingUsers**
> List getFollowingUsers(username)

Get all followed users for a user

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **username** | **String**| username of the user | [default to null] |

### Return type

[**List**](../Models/User.md)

### Authorization

[BasicAuth](../README.md#BasicAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="getUserByUsername"></a>
# **getUserByUsername**
> User getUserByUsername(username)

Get a particular user by username

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **username** | **String**| username of the user | [default to null] |

### Return type

[**User**](../Models/User.md)

### Authorization

[BasicAuth](../README.md#BasicAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="getUserComments"></a>
# **getUserComments**
> List getUserComments(username)

Gets all comments for a user

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **username** | **String**| username of the user | [default to null] |

### Return type

[**List**](../Models/Comment.md)

### Authorization

[BasicAuth](../README.md#BasicAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="getUserFollowers"></a>
# **getUserFollowers**
> List getUserFollowers(username)

Get all followers for a user

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **username** | **String**| username of the user | [default to null] |

### Return type

[**List**](../Models/User.md)

### Authorization

[BasicAuth](../README.md#BasicAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="listUsers"></a>
# **listUsers**
> List listUsers()

Returns all the users

### Parameters
This endpoint does not need any parameter.

### Return type

[**List**](../Models/User.md)

### Authorization

[BasicAuth](../README.md#BasicAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="unfollowUser"></a>
# **unfollowUser**
> unfollowUser(followedUsername, followerUsername)

Remove a user as a follower

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **followedUsername** | **String**| username of the user | [default to null] |
| **followerUsername** | **String**| username of the follower | [default to null] |

### Return type

null (empty response body)

### Authorization

[BasicAuth](../README.md#BasicAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="updateUser"></a>
# **updateUser**
> User updateUser(username, User)

Update a user

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **username** | **String**| username of the user | [default to null] |
| **User** | [**User**](../Models/User.md)| Update a user | [optional] |

### Return type

[**User**](../Models/User.md)

### Authorization

[BasicAuth](../README.md#BasicAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

