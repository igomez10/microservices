# \CommentAPI

All URIs are relative to *https://socialapp.gomezignacio.com*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateComment**](CommentAPI.md#CreateComment) | **Post** /v1/comments | Create a new comment
[**GetComment**](CommentAPI.md#GetComment) | **Get** /v1/comments/{id} | Returns details about a particular comment
[**GetUserComments**](CommentAPI.md#GetUserComments) | **Get** /v1/users/{username}/comments | Gets all comments for a user
[**GetUserFeed**](CommentAPI.md#GetUserFeed) | **Get** /v1/feed | Returns a users feed



## CreateComment

> Comment CreateComment(ctx).Comment(comment).Execute()

Create a new comment



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/client"
)

func main() {
    comment := *openapiclient.NewComment("Content_example", "Username_example") // Comment | Create a new comment

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.CommentAPI.CreateComment(context.Background()).Comment(comment).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `CommentAPI.CreateComment``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `CreateComment`: Comment
    fmt.Fprintf(os.Stdout, "Response from `CommentAPI.CreateComment`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateCommentRequest struct via the builder pattern


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

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetComment

> Comment GetComment(ctx, id).Execute()

Returns details about a particular comment



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/client"
)

func main() {
    id := int32(123) // int32 | ID of the comment

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.CommentAPI.GetComment(context.Background(), id).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `CommentAPI.GetComment``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetComment`: Comment
    fmt.Fprintf(os.Stdout, "Response from `CommentAPI.GetComment`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int32** | ID of the comment | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetCommentRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**Comment**](Comment.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetUserComments

> []Comment GetUserComments(ctx, username).Limit(limit).Offset(offset).Execute()

Gets all comments for a user



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/client"
)

func main() {
    username := "johndoe" // string | username of the user
    limit := int32(20) // int32 | How many items to return at one time (max 100) (optional) (default to 20)
    offset := int32(0) // int32 | The number of items to skip before starting to collect the result set (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.CommentAPI.GetUserComments(context.Background(), username).Limit(limit).Offset(offset).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `CommentAPI.GetUserComments``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetUserComments`: []Comment
    fmt.Fprintf(os.Stdout, "Response from `CommentAPI.GetUserComments`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**username** | **string** | username of the user | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetUserCommentsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **limit** | **int32** | How many items to return at one time (max 100) | [default to 20]
 **offset** | **int32** | The number of items to skip before starting to collect the result set | 

### Return type

[**[]Comment**](Comment.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetUserFeed

> []Comment GetUserFeed(ctx).Execute()

Returns a users feed



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/client"
)

func main() {

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.CommentAPI.GetUserFeed(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `CommentAPI.GetUserFeed``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetUserFeed`: []Comment
    fmt.Fprintf(os.Stdout, "Response from `CommentAPI.GetUserFeed`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetUserFeedRequest struct via the builder pattern


### Return type

[**[]Comment**](Comment.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

