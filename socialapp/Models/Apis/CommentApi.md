# CommentApi

All URIs are relative to *https://socialapp.gomezignacio.com*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**createComment**](CommentApi.md#createComment) | **POST** /v1/comments | Create a new comment |
| [**getComment**](CommentApi.md#getComment) | **GET** /v1/comments/{id} | Get comment by ID |
| [**getUserFeed**](CommentApi.md#getUserFeed) | **GET** /v1/feed | Get user feed |


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
| **id** | **Long**| ID of the comment | [default to null] |

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

