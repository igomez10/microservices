# \AuthenticationAPI

All URIs are relative to *https://socialapp.gomezignacio.com*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetAccessToken**](AuthenticationAPI.md#GetAccessToken) | **Post** /v1/oauth/token | Get an access token



## GetAccessToken

> AccessToken GetAccessToken(ctx).Execute()

Get an access token



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
	resp, r, err := apiClient.AuthenticationAPI.GetAccessToken(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AuthenticationAPI.GetAccessToken``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetAccessToken`: AccessToken
	fmt.Fprintf(os.Stdout, "Response from `AuthenticationAPI.GetAccessToken`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetAccessTokenRequest struct via the builder pattern


### Return type

[**AccessToken**](AccessToken.md)

### Authorization

[BasicAuth](../README.md#BasicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

