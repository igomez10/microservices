# UserApi

All URIs are relative to *https://socialapp.gomezignacio.com*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**changePassword**](UserApi.md#changePassword) | **POST** /password | Change password |
| [**createUser**](UserApi.md#createUser) | **POST** /users | Create user |
| [**deleteUser**](UserApi.md#deleteUser) | **DELETE** /users/{username} | Deletes a particular user |
| [**followUser**](UserApi.md#followUser) | **POST** /users/{followedUsername}/followers/{followerUsername} | Add a user as a follower |
| [**getFollowingUsers**](UserApi.md#getFollowingUsers) | **GET** /users/{username}/following | Get all followed users for a user |
| [**getRolesForUser**](UserApi.md#getRolesForUser) | **GET** /users/{username}/roles | Get all roles for a user |
| [**getUserByUsername**](UserApi.md#getUserByUsername) | **GET** /users/{username} | Get a particular user by username |
| [**getUserComments**](UserApi.md#getUserComments) | **GET** /users/{username}/comments | Gets all comments for a user |
| [**getUserFollowers**](UserApi.md#getUserFollowers) | **GET** /users/{username}/followers | Get all followers for a user |
| [**listUsers**](UserApi.md#listUsers) | **GET** /users | List users |
| [**resetPassword**](UserApi.md#resetPassword) | **PUT** /password | Reset password |
| [**unfollowUser**](UserApi.md#unfollowUser) | **DELETE** /users/{followedUsername}/followers/{followerUsername} | Remove a user as a follower |
| [**updateRolesForUser**](UserApi.md#updateRolesForUser) | **PUT** /users/{username}/roles | Update all roles for a user |
| [**updateUser**](UserApi.md#updateUser) | **PUT** /users/{username} | Update a user |


<a name="changePassword"></a>
# **changePassword**
> User changePassword(ChangePasswordRequest)

Change password

    Change the password of the current user

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **ChangePasswordRequest** | [**ChangePasswordRequest**](../Models/ChangePasswordRequest.md)| Change password request | |

### Return type

[**User**](../Models/User.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

<a name="createUser"></a>
# **createUser**
> User createUser(CreateUserRequest)

Create user

    Create a new user in the system

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **CreateUserRequest** | [**CreateUserRequest**](../Models/CreateUserRequest.md)| Create a new user | |

### Return type

[**User**](../Models/User.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

<a name="deleteUser"></a>
# **deleteUser**
> User deleteUser(username)

Deletes a particular user

    Deletes a particular user by username

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **username** | **String**| username of the user | [default to null] |

### Return type

[**User**](../Models/User.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="followUser"></a>
# **followUser**
> followUser(followedUsername, followerUsername)

Add a user as a follower

    Add a user as a follower of another user

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **followedUsername** | **String**| username of the user | [default to null] |
| **followerUsername** | **String**| username of the follower | [default to null] |

### Return type

null (empty response body)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="getFollowingUsers"></a>
# **getFollowingUsers**
> List getFollowingUsers(username)

Get all followed users for a user

    Get all followed users for a user (users that the user is following)

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **username** | **String**| username of the user | [default to null] |

### Return type

[**List**](../Models/User.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="getRolesForUser"></a>
# **getRolesForUser**
> List getRolesForUser(username)

Get all roles for a user

    Get all roles for a user

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **username** | **String**| username of the user | [default to null] |

### Return type

[**List**](../Models/Role.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="getUserByUsername"></a>
# **getUserByUsername**
> User getUserByUsername(username)

Get a particular user by username

    Get a particular user by username

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **username** | **String**| username of the user | [default to null] |

### Return type

[**User**](../Models/User.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="getUserComments"></a>
# **getUserComments**
> List getUserComments(username, limit, offset)

Gets all comments for a user

    Gets all comments for a user

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **username** | **String**| username of the user | [default to null] |
| **limit** | **Integer**| How many items to return at one time (max 100) | [optional] [default to 20] |
| **offset** | **Integer**| The number of items to skip before starting to collect the result set | [optional] [default to null] |

### Return type

[**List**](../Models/Comment.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="getUserFollowers"></a>
# **getUserFollowers**
> List getUserFollowers(username)

Get all followers for a user

    Get all followers for a user

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **username** | **String**| username of the user | [default to null] |

### Return type

[**List**](../Models/User.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="listUsers"></a>
# **listUsers**
> List listUsers(limit, offset)

List users

    List all users in the system (paginated)

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **limit** | **Integer**| Maximum number of users to return | [optional] [default to 20] |
| **offset** | **Integer**| Pagination offset | [optional] [default to 0] |

### Return type

[**List**](../Models/User.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="resetPassword"></a>
# **resetPassword**
> User resetPassword(ResetPasswordRequest)

Reset password

    Reset the password of a user

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **ResetPasswordRequest** | [**ResetPasswordRequest**](../Models/ResetPasswordRequest.md)| Reset password | |

### Return type

[**User**](../Models/User.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

<a name="unfollowUser"></a>
# **unfollowUser**
> unfollowUser(followedUsername, followerUsername)

Remove a user as a follower

    Remove a user as a follower of another user

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **followedUsername** | **String**| username of the user | [default to null] |
| **followerUsername** | **String**| username of the follower | [default to null] |

### Return type

null (empty response body)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="updateRolesForUser"></a>
# **updateRolesForUser**
> List updateRolesForUser(username, request\_body)

Update all roles for a user

    Update all roles for a user

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **username** | **String**| username of the user | [default to null] |
| **request\_body** | [**List**](../Models/string.md)| Update all roles for a user | [optional] |

### Return type

[**List**](../Models/Role.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

<a name="updateUser"></a>
# **updateUser**
> User updateUser(username, User)

Update a user

    Update a user by username

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **username** | **String**| username of the user | [default to null] |
| **User** | [**User**](../Models/User.md)| Update a user | |

### Return type

[**User**](../Models/User.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

