# CommentApi

All URIs are relative to **

Method | HTTP request | Description
------------- | ------------- | -------------
[**createComment**](CommentApi.md#createComment) | **POST** /v1/comments | Create a new comment
[**getComment**](CommentApi.md#getComment) | **GET** /v1/comments/{id} | Get comment by ID
[**getUserFeed**](CommentApi.md#getUserFeed) | **GET** /v1/feed | Get user feed
[**searchComments**](CommentApi.md#searchComments) | **GET** /v1/comments | Search comments



## createComment

Create a new comment

Create a new comment.

### Example

```bash
socialapp-cli createComment
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **comment** | [**Comment**](Comment.md) | Create a new comment |

### Return type

[**Comment**](Comment.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## getComment

Get comment by ID

Get details for a specific comment.

### Example

```bash
socialapp-cli getComment id=value
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **integer** | ID of the comment | [default to null]

### Return type

[**Comment**](Comment.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not Applicable
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## getUserFeed

Get user feed

Return the current user's feed.

### Example

```bash
socialapp-cli getUserFeed
```

### Parameters

This endpoint does not need any parameter.

### Return type

[**array[Comment]**](Comment.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not Applicable
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## searchComments

Search comments

Search comments by author and creation time.

### Example

```bash
socialapp-cli searchComments  username=value  start_time=value  end_time=value
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **username** | **string** | Filter comments by the author's username. | [optional] [default to null]
 **startTime** | **string** | Return comments created at or after this timestamp. | [optional] [default to null]
 **endTime** | **string** | Return comments created at or before this timestamp. | [optional] [default to null]

### Return type

[**array[Comment]**](Comment.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not Applicable
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

