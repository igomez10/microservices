# CommentApi

All URIs are relative to *https://socialapp.gomezignacio.com*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**createComment**](CommentApi.md#createComment) | **POST** /v1/comments | Create a new comment |
| [**getComment**](CommentApi.md#getComment) | **GET** /v1/comments/{id} | Get comment by ID |
| [**getUserFeed**](CommentApi.md#getUserFeed) | **GET** /v1/feed | Get user feed |
| [**likeComment**](CommentApi.md#likeComment) | **POST** /like | Like a comment |
| [**searchComments**](CommentApi.md#searchComments) | **GET** /v1/comments | Search comments |
| [**unlikeComment**](CommentApi.md#unlikeComment) | **DELETE** /like | Unlike a comment |


<a name="createComment"></a>
# **createComment**
> Comment createComment(Comment)

Create a new comment

    Create a new comment.

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

Get comment by ID

    Get details for a specific comment.

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **id** | **String**| ID of the comment | [default to null] |

### Return type

[**Comment**](../Models/Comment.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="getUserFeed"></a>
# **getUserFeed**
> List getUserFeed()

Get user feed

    Return the current user&#39;s feed.

### Parameters
This endpoint does not need any parameter.

### Return type

[**List**](../Models/Comment.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="likeComment"></a>
# **likeComment**
> Comment likeComment(LikeRequest)

Like a comment

    Like a comment as the authenticated user.

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **LikeRequest** | [**LikeRequest**](../Models/LikeRequest.md)| Comment like request | |

### Return type

[**Comment**](../Models/Comment.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

<a name="searchComments"></a>
# **searchComments**
> List searchComments(username, start\_time, end\_time)

Search comments

    Search comments by author and creation time.

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **username** | **String**| Filter comments by the author&#39;s username. | [optional] [default to null] |
| **start\_time** | **Date**| Return comments created at or after this timestamp. | [optional] [default to null] |
| **end\_time** | **Date**| Return comments created at or before this timestamp. | [optional] [default to null] |

### Return type

[**List**](../Models/Comment.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="unlikeComment"></a>
# **unlikeComment**
> Comment unlikeComment(LikeRequest)

Unlike a comment

    Remove a like from a comment as the authenticated user.

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **LikeRequest** | [**LikeRequest**](../Models/LikeRequest.md)| Comment unlike request | |

### Return type

[**Comment**](../Models/Comment.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

