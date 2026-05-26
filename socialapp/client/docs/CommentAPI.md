# \CommentAPI

All URIs are relative to *https://socialapp.gomezignacio.com*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateComment**](CommentAPI.md#CreateComment) | **Post** /v1/comments | Create a new comment
[**GetComment**](CommentAPI.md#GetComment) | **Get** /v1/comments/{id} | Get comment by ID
[**GetUserFeed**](CommentAPI.md#GetUserFeed) | **Get** /v1/feed | Get user feed
[**LikeComment**](CommentAPI.md#LikeComment) | **Post** /like | Like a comment
[**SearchComments**](CommentAPI.md#SearchComments) | **Get** /v1/comments | Search comments
[**UnlikeComment**](CommentAPI.md#UnlikeComment) | **Delete** /like | Unlike a comment



## CreateComment

> Comment CreateComment(ctx).CreateCommentRequest(createCommentRequest).Execute()

Create a new comment



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/igomez10/microservices/socialapp/client"
)

func main() {
	createCommentRequest := *openapiclient.NewCreateCommentRequest("Content_example", "Username_example") // CreateCommentRequest | Create a new comment

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CommentAPI.CreateComment(context.Background()).CreateCommentRequest(createCommentRequest).Execute()
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
 **createCommentRequest** | [**CreateCommentRequest**](CreateCommentRequest.md) | Create a new comment | 

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

Get comment by ID



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/igomez10/microservices/socialapp/client"
)

func main() {
	id := "123" // string | ID of the comment

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
**id** | **string** | ID of the comment | 

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


## GetUserFeed

> []Comment GetUserFeed(ctx).Execute()

Get user feed



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/igomez10/microservices/socialapp/client"
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


## LikeComment

> Comment LikeComment(ctx).LikeRequest(likeRequest).Execute()

Like a comment



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/igomez10/microservices/socialapp/client"
)

func main() {
	likeRequest := *openapiclient.NewLikeRequest("CommentId_example") // LikeRequest | Comment like request

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CommentAPI.LikeComment(context.Background()).LikeRequest(likeRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CommentAPI.LikeComment``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `LikeComment`: Comment
	fmt.Fprintf(os.Stdout, "Response from `CommentAPI.LikeComment`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiLikeCommentRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **likeRequest** | [**LikeRequest**](LikeRequest.md) | Comment like request | 

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


## SearchComments

> []Comment SearchComments(ctx).Username(username).StartTime(startTime).EndTime(endTime).Execute()

Search comments



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
    "time"
	openapiclient "github.com/igomez10/microservices/socialapp/client"
)

func main() {
	username := "username_example" // string | Filter comments by the author's username. (optional)
	startTime := time.Now() // time.Time | Return comments created at or after this timestamp. (optional)
	endTime := time.Now() // time.Time | Return comments created at or before this timestamp. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CommentAPI.SearchComments(context.Background()).Username(username).StartTime(startTime).EndTime(endTime).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CommentAPI.SearchComments``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `SearchComments`: []Comment
	fmt.Fprintf(os.Stdout, "Response from `CommentAPI.SearchComments`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiSearchCommentsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **username** | **string** | Filter comments by the author&#39;s username. | 
 **startTime** | **time.Time** | Return comments created at or after this timestamp. | 
 **endTime** | **time.Time** | Return comments created at or before this timestamp. | 

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


## UnlikeComment

> Comment UnlikeComment(ctx).LikeRequest(likeRequest).Execute()

Unlike a comment



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/igomez10/microservices/socialapp/client"
)

func main() {
	likeRequest := *openapiclient.NewLikeRequest("CommentId_example") // LikeRequest | Comment unlike request

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CommentAPI.UnlikeComment(context.Background()).LikeRequest(likeRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CommentAPI.UnlikeComment``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UnlikeComment`: Comment
	fmt.Fprintf(os.Stdout, "Response from `CommentAPI.UnlikeComment`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiUnlikeCommentRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **likeRequest** | [**LikeRequest**](LikeRequest.md) | Comment unlike request | 

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

