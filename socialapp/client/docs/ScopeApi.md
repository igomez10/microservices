# \ScopeAPI

All URIs are relative to *https://socialapp.gomezignacio.com*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateScope**](ScopeAPI.md#CreateScope) | **Post** /v1/scopes | Create a new scope
[**DeleteScope**](ScopeAPI.md#DeleteScope) | **Delete** /v1/scopes/{id} | Delete a scope
[**GetScope**](ScopeAPI.md#GetScope) | **Get** /v1/scopes/{id} | Returns a scope
[**ListScopes**](ScopeAPI.md#ListScopes) | **Get** /v1/scopes | Returns a list of scopes
[**UpdateScope**](ScopeAPI.md#UpdateScope) | **Put** /v1/scopes/{id} | Update a scope



## CreateScope

> Scope CreateScope(ctx).Scope(scope).Execute()

Create a new scope



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
	scope := *openapiclient.NewScope("Name_example", "Description_example") // Scope | Create a new scope

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ScopeAPI.CreateScope(context.Background()).Scope(scope).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ScopeAPI.CreateScope``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateScope`: Scope
	fmt.Fprintf(os.Stdout, "Response from `ScopeAPI.CreateScope`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateScopeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **scope** | [**Scope**](Scope.md) | Create a new scope | 

### Return type

[**Scope**](Scope.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteScope

> DeleteScope(ctx, id).Execute()

Delete a scope



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
	id := int32(123) // int32 | id of the scope

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ScopeAPI.DeleteScope(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ScopeAPI.DeleteScope``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int32** | id of the scope | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteScopeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetScope

> Scope GetScope(ctx, id).Execute()

Returns a scope



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
	id := int32(123) // int32 | The id of the scope

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ScopeAPI.GetScope(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ScopeAPI.GetScope``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetScope`: Scope
	fmt.Fprintf(os.Stdout, "Response from `ScopeAPI.GetScope`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int32** | The id of the scope | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetScopeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**Scope**](Scope.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListScopes

> []Scope ListScopes(ctx).Limit(limit).Offset(offset).Execute()

Returns a list of scopes



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
	limit := int32(56) // int32 | The numbers of scopes to return (optional) (default to 20)
	offset := int32(56) // int32 | The number of items to skip before starting to collect the result (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ScopeAPI.ListScopes(context.Background()).Limit(limit).Offset(offset).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ScopeAPI.ListScopes``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListScopes`: []Scope
	fmt.Fprintf(os.Stdout, "Response from `ScopeAPI.ListScopes`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListScopesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **limit** | **int32** | The numbers of scopes to return | [default to 20]
 **offset** | **int32** | The number of items to skip before starting to collect the result | 

### Return type

[**[]Scope**](Scope.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateScope

> Scope UpdateScope(ctx, id).Scope(scope).Execute()

Update a scope



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
	id := int32(123) // int32 | id of the scope
	scope := *openapiclient.NewScope("Name_example", "Description_example") // Scope | Update a scope (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ScopeAPI.UpdateScope(context.Background(), id).Scope(scope).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ScopeAPI.UpdateScope``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateScope`: Scope
	fmt.Fprintf(os.Stdout, "Response from `ScopeAPI.UpdateScope`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int32** | id of the scope | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateScopeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **scope** | [**Scope**](Scope.md) | Update a scope | 

### Return type

[**Scope**](Scope.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

