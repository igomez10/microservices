# ScopeApi

All URIs are relative to *https://socialapp.gomezignacio.com*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**createScope**](ScopeApi.md#createScope) | **POST** /v1/scopes | Create a new scope |
| [**deleteScope**](ScopeApi.md#deleteScope) | **DELETE** /v1/scopes/{id} | Delete a scope |
| [**getScope**](ScopeApi.md#getScope) | **GET** /v1/scopes/{id} | Returns a scope |
| [**listScopes**](ScopeApi.md#listScopes) | **GET** /v1/scopes | Returns a list of scopes |
| [**updateScope**](ScopeApi.md#updateScope) | **PUT** /v1/scopes/{id} | Update a scope |


<a name="createScope"></a>
# **createScope**
> Scope createScope(Scope)

Create a new scope

    Create a new scope

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **Scope** | [**Scope**](../Models/Scope.md)| Create a new scope | |

### Return type

[**Scope**](../Models/Scope.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

<a name="deleteScope"></a>
# **deleteScope**
> deleteScope(id)

Delete a scope

    Delete a scope

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **id** | **Integer**| id of the scope | [default to null] |

### Return type

null (empty response body)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="getScope"></a>
# **getScope**
> Scope getScope(id)

Returns a scope

    Returns a scope

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **id** | **Integer**| The id of the scope | [default to null] |

### Return type

[**Scope**](../Models/Scope.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="listScopes"></a>
# **listScopes**
> List listScopes(limit, offset)

Returns a list of scopes

    Returns a list of scopes

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **limit** | **Integer**| The numbers of scopes to return | [optional] [default to 20] |
| **offset** | **Integer**| The number of items to skip before starting to collect the result | [optional] [default to null] |

### Return type

[**List**](../Models/Scope.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="updateScope"></a>
# **updateScope**
> Scope updateScope(id, Scope)

Update a scope

    Update a scope

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **id** | **Integer**| id of the scope | [default to null] |
| **Scope** | [**Scope**](../Models/Scope.md)| Update a scope | [optional] |

### Return type

[**Scope**](../Models/Scope.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

