# CommentApi

All URIs are relative to *http://localhost:8080*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**createComment**](CommentApi.md#createComment) | **POST** /users/{username}/comments | Create a new comment |
| [**getComment**](CommentApi.md#getComment) | **GET** /comments/{id} | Returns details about a particular comment |
| [**getUserComments**](CommentApi.md#getUserComments) | **GET** /users/{username}/comments | Gets all comments for a user |


<a name="createComment"></a>
# **createComment**
> Comment createComment(username, Comment)

Create a new comment

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **username** | **String**| username of the user | [default to null] |
| **Comment** | [**Comment**](../Models/Comment.md)| Create a new comment | |

### Return type

[**Comment**](../Models/Comment.md)

### Authorization

[BasicAuth](../README.md#BasicAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

<a name="getComment"></a>
# **getComment**
> Comment getComment(id)

Returns details about a particular comment

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **id** | **Integer**| ID of the comment | [default to null] |

### Return type

[**Comment**](../Models/Comment.md)

### Authorization

[BasicAuth](../README.md#BasicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="getUserComments"></a>
# **getUserComments**
> Comment getUserComments(username)

Gets all comments for a user

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **username** | **String**| username of the user | [default to null] |

### Return type

[**Comment**](../Models/Comment.md)

### Authorization

[BasicAuth](../README.md#BasicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

