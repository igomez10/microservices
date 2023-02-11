# CommentApi

All URIs are relative to **

Method | HTTP request | Description
------------- | ------------- | -------------
[**createComment**](CommentApi.md#createComment) | **POST** /v1/comments | Create a new comment
[**getComment**](CommentApi.md#getComment) | **GET** /v1/comments/{id} | Returns details about a particular comment
[**getUserComments**](CommentApi.md#getUserComments) | **GET** /v1/users/{username}/comments | Gets all comments for a user
[**getUserFeed**](CommentApi.md#getUserFeed) | **GET** /v1/users/{username}/feed | Returns a users feed



## createComment

Create a new comment

Create a new comment

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

Returns details about a particular comment

Returns details about a particular comment

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


## getUserComments

Gets all comments for a user

Gets all comments for a user

### Example

```bash
socialapp-cli getUserComments username=value  limit=value  offset=value
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **username** | **string** | username of the user | [default to null]
 **limit** | **integer** | How many items to return at one time (max 100) | [optional] [default to 20]
 **offset** | **integer** | The number of items to skip before starting to collect the result set | [optional] [default to null]

### Return type

[**array[Comment]**](Comment.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not Applicable
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## getUserFeed

Returns a users feed

Returns a users feed

### Example

```bash
socialapp-cli getUserFeed username=value
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **username** | **string** | The username of the user | [default to null]

### Return type

[**array[Comment]**](Comment.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not Applicable
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

