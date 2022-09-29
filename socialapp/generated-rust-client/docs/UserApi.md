# \UserApi

All URIs are relative to *https://microservices.onrender.com*

Method | HTTP request | Description
------------- | ------------- | -------------
[**change_password**](UserApi.md#change_password) | **POST** /password | Change password
[**create_user**](UserApi.md#create_user) | **POST** /users | Create user
[**delete_user**](UserApi.md#delete_user) | **DELETE** /users/{username} | Deletes a particular user
[**follow_user**](UserApi.md#follow_user) | **POST** /users/{followedUsername}/followers/{followerUsername} | Add a user as a follower
[**get_following_users**](UserApi.md#get_following_users) | **GET** /users/{username}/following | Get all followed users for a user
[**get_user_by_username**](UserApi.md#get_user_by_username) | **GET** /users/{username} | Get a particular user by username
[**get_user_comments**](UserApi.md#get_user_comments) | **GET** /users/{username}/comments | Gets all comments for a user
[**get_user_followers**](UserApi.md#get_user_followers) | **GET** /users/{username}/followers | Get all followers for a user
[**list_users**](UserApi.md#list_users) | **GET** /users | List users
[**reset_password**](UserApi.md#reset_password) | **PUT** /password | Reset password
[**unfollow_user**](UserApi.md#unfollow_user) | **DELETE** /users/{followedUsername}/followers/{followerUsername} | Remove a user as a follower
[**update_user**](UserApi.md#update_user) | **PUT** /users/{username} | Update a user



## change_password

> crate::models::User change_password(change_password_request)
Change password

Change the password of the current user

### Parameters


Name | Type | Description  | Required | Notes
------------- | ------------- | ------------- | ------------- | -------------
**change_password_request** | [**ChangePasswordRequest**](ChangePasswordRequest.md) | Change password request | [required] |

### Return type

[**crate::models::User**](User.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## create_user

> crate::models::User create_user(create_user_request)
Create user

Create a new user in the system

### Parameters


Name | Type | Description  | Required | Notes
------------- | ------------- | ------------- | ------------- | -------------
**create_user_request** | [**CreateUserRequest**](CreateUserRequest.md) | Create a new user | [required] |

### Return type

[**crate::models::User**](User.md)

### Authorization

[BasicAuth](../README.md#BasicAuth), [BearerAuth](../README.md#BearerAuth), [OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## delete_user

> crate::models::User delete_user(username)
Deletes a particular user

Deletes a particular user by username

### Parameters


Name | Type | Description  | Required | Notes
------------- | ------------- | ------------- | ------------- | -------------
**username** | **String** | username of the user | [required] |

### Return type

[**crate::models::User**](User.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## follow_user

> follow_user(followed_username, follower_username)
Add a user as a follower

Add a user as a follower of another user

### Parameters


Name | Type | Description  | Required | Notes
------------- | ------------- | ------------- | ------------- | -------------
**followed_username** | **String** | username of the user | [required] |
**follower_username** | **String** | username of the follower | [required] |

### Return type

 (empty response body)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## get_following_users

> Vec<crate::models::User> get_following_users(username)
Get all followed users for a user

Get all followed users for a user (users that the user is following)

### Parameters


Name | Type | Description  | Required | Notes
------------- | ------------- | ------------- | ------------- | -------------
**username** | **String** | username of the user | [required] |

### Return type

[**Vec<crate::models::User>**](User.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## get_user_by_username

> crate::models::User get_user_by_username(username)
Get a particular user by username

Get a particular user by username

### Parameters


Name | Type | Description  | Required | Notes
------------- | ------------- | ------------- | ------------- | -------------
**username** | **String** | username of the user | [required] |

### Return type

[**crate::models::User**](User.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## get_user_comments

> Vec<crate::models::Comment> get_user_comments(username, limit, offset)
Gets all comments for a user

Gets all comments for a user

### Parameters


Name | Type | Description  | Required | Notes
------------- | ------------- | ------------- | ------------- | -------------
**username** | **String** | username of the user | [required] |
**limit** | Option<**i32**> | How many items to return at one time (max 100) |  |
**offset** | Option<**i32**> | The number of items to skip before starting to collect the result set |  |

### Return type

[**Vec<crate::models::Comment>**](Comment.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## get_user_followers

> Vec<crate::models::User> get_user_followers(username)
Get all followers for a user

Get all followers for a user

### Parameters


Name | Type | Description  | Required | Notes
------------- | ------------- | ------------- | ------------- | -------------
**username** | **String** | username of the user | [required] |

### Return type

[**Vec<crate::models::User>**](User.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## list_users

> Vec<crate::models::User> list_users(limit, offset)
List users

List all users in the system (paginated)

### Parameters


Name | Type | Description  | Required | Notes
------------- | ------------- | ------------- | ------------- | -------------
**limit** | Option<**i32**> | Maximum number of users to return |  |[default to 20]
**offset** | Option<**i32**> | Pagination offset |  |[default to 0]

### Return type

[**Vec<crate::models::User>**](User.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## reset_password

> crate::models::User reset_password(reset_password_request)
Reset password

Reset the password of a user

### Parameters


Name | Type | Description  | Required | Notes
------------- | ------------- | ------------- | ------------- | -------------
**reset_password_request** | [**ResetPasswordRequest**](ResetPasswordRequest.md) | Reset password | [required] |

### Return type

[**crate::models::User**](User.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## unfollow_user

> unfollow_user(followed_username, follower_username)
Remove a user as a follower

Remove a user as a follower of another user

### Parameters


Name | Type | Description  | Required | Notes
------------- | ------------- | ------------- | ------------- | -------------
**followed_username** | **String** | username of the user | [required] |
**follower_username** | **String** | username of the follower | [required] |

### Return type

 (empty response body)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## update_user

> crate::models::User update_user(username, user)
Update a user

Update a user by username

### Parameters


Name | Type | Description  | Required | Notes
------------- | ------------- | ------------- | ------------- | -------------
**username** | **String** | username of the user | [required] |
**user** | Option<[**User**](User.md)> | Update a user |  |

### Return type

[**crate::models::User**](User.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

