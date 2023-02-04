# \FollowingApi

All URIs are relative to *https://socialapp.gomezignacio.com*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetUserFollowers**](FollowingApi.md#GetUserFollowers) | **Get** /users/{username}/followers | Get all followers for a user



## GetUserFollowers

> []User GetUserFollowers(ctx, username).Execute()

Get all followers for a user



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    username := "johndoe" // string | username of the user

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.FollowingApi.GetUserFollowers(context.Background(), username).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `FollowingApi.GetUserFollowers``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetUserFollowers`: []User
    fmt.Fprintf(os.Stdout, "Response from `FollowingApi.GetUserFollowers`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**username** | **string** | username of the user | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetUserFollowersRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**[]User**](User.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

