# \WelcomeAPI

All URIs are relative to *https://socialapp.gomezignacio.com*

Method | HTTP request | Description
------------- | ------------- | -------------
[**Welcome**](WelcomeAPI.md#Welcome) | **Get** / | Welcome to the Socialapp API



## Welcome

> string Welcome(ctx).Execute()

Welcome to the Socialapp API



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
	resp, r, err := apiClient.WelcomeAPI.Welcome(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `WelcomeAPI.Welcome``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `Welcome`: string
	fmt.Fprintf(os.Stdout, "Response from `WelcomeAPI.Welcome`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiWelcomeRequest struct via the builder pattern


### Return type

**string**

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: text/plain, application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

