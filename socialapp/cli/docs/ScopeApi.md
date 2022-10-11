# ScopeApi

All URIs are relative to **

Method | HTTP request | Description
------------- | ------------- | -------------
[**createScope**](ScopeApi.md#createScope) | **POST** /scopes | Create a new scope
[**deleteScope**](ScopeApi.md#deleteScope) | **DELETE** /scopes/{id} | Delete a scope
[**getScope**](ScopeApi.md#getScope) | **GET** /scopes/{id} | Returns a scope
[**listScopes**](ScopeApi.md#listScopes) | **GET** /scopes | Returns a list of scopes
[**updateScope**](ScopeApi.md#updateScope) | **PUT** /scopes/{id} | Update a scope



## createScope

Create a new scope

Create a new scope

### Example

```bash
socialapp-cli createScope
```

### Parameters


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

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## deleteScope

Delete a scope

Delete a scope

### Example

```bash
socialapp-cli deleteScope id=value
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **integer** | id of the scope | [default to null]

### Return type

(empty response body)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not Applicable
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## getScope

Returns a scope

Returns a scope

### Example

```bash
socialapp-cli getScope id=value
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **integer** | The id of the scope | [default to null]

### Return type

[**Scope**](Scope.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not Applicable
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## listScopes

Returns a list of scopes

Returns a list of scopes

### Example

```bash
socialapp-cli listScopes  limit=value  offset=value
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **limit** | **integer** | The numbers of scopes to return | [optional] [default to 20]
 **offset** | **integer** | The number of items to skip before starting to collect the result | [optional] [default to null]

### Return type

[**array[Scope]**](Scope.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not Applicable
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## updateScope

Update a scope

Update a scope

### Example

```bash
socialapp-cli updateScope id=value
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **integer** | id of the scope | [default to null]
 **scope** | [**Scope**](Scope.md) | Update a scope | [optional]

### Return type

[**Scope**](Scope.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

