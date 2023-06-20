# CommentApi

All URIs are relative to *https://socialapp.gomezignacio.com*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**createComment**](CommentApi.md#createComment) | **POST** /v1/comments | Create a new comment |
| [**getComment**](CommentApi.md#getComment) | **GET** /v1/comments/{id} | Returns details about a particular comment |
| [**getUserComments**](CommentApi.md#getUserComments) | **GET** /v1/users/{username}/comments | Gets all comments for a user |
| [**getUserFeed**](CommentApi.md#getUserFeed) | **GET** /v1/feed | Returns a users feed |


<a name="createComment"></a>
# **createComment**
> Comment createComment(Comment)

Create a new comment

    Create a new comment

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **Comment** | [**Comment**](../Models/Comment.md)| Create a new comment | |

### Return type

[**Comment**](../Models/Comment.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

<a name="getComment"></a>
# **getComment**
> Comment getComment(id)

Returns details about a particular comment

    Returns details about a particular comment

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **id** | **Integer**| ID of the comment | [default to null] |

### Return type

[**Comment**](../Models/Comment.md)

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

<a name="getUserFeed"></a>
# **getUserFeed**
> List getUserFeed()

Returns a users feed

    Returns a users feed

### Parameters
This endpoint does not need any parameter.

### Return type

[**List**](../Models/Comment.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

