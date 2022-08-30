# CommentApi

All URIs are relative to **

Method | HTTP request | Description
------------- | ------------- | -------------
[**createComment**](CommentApi.md#createComment) | **POST** /comments | Create a new comment
[**getComment**](CommentApi.md#getComment) | **GET** /comments/{id} | Returns details about a particular comment
[**getUserComments**](CommentApi.md#getUserComments) | **GET** /users/{username}/comments | Gets all comments for a user



## createComment

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

[BasicAuth](../README.md#BasicAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## getComment

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

