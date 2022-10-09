# RoleApi

All URIs are relative to **

Method | HTTP request | Description
------------- | ------------- | -------------
[**addScopeToRole**](RoleApi.md#addScopeToRole) | **POST** /roles/{id}/scopes | Add a scope to a role
[**createRole**](RoleApi.md#createRole) | **POST** /roles | Create a new role
[**deleteRole**](RoleApi.md#deleteRole) | **DELETE** /roles/{id} | Delete a role
[**getRole**](RoleApi.md#getRole) | **GET** /roles/{id} | Returns a role
[**listRoles**](RoleApi.md#listRoles) | **GET** /roles | Returns a list of roles
[**listScopesForRole**](RoleApi.md#listScopesForRole) | **GET** /roles/{id}/scopes | Returns a list of scopes for a role
[**removeScopeFromRole**](RoleApi.md#removeScopeFromRole) | **DELETE** /roles/{role_id}/scopes/{scope_id} | Remove a scope from a role
[**updateRole**](RoleApi.md#updateRole) | **PUT** /roles/{id} | Update a role



## addScopeToRole

Add a scope to a role

Add a scope to a role

### Example

```bash
socialapp-cli addScopeToRole id=value
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **integer** | The id of the role | [default to null]
 **requestBody** | [**array[string]**](string.md) | Add a scope to a role |

### Return type

[**array[Scope]**](Scope.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## createRole

Create a new role

Create a new role

### Example

```bash
socialapp-cli createRole
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **role** | [**Role**](Role.md) | Create a new role |

### Return type

[**Role**](Role.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## deleteRole

Delete a role

Delete a role

### Example

```bash
socialapp-cli deleteRole id=value
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **integer** | id of the role | [default to null]

### Return type

[**Role**](Role.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not Applicable
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## getRole

Returns a role

Returns a role

### Example

```bash
socialapp-cli getRole id=value
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **integer** | The id of the role | [default to null]

### Return type

[**Role**](Role.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not Applicable
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## listRoles

Returns a list of roles

Returns a list of roles

### Example

```bash
socialapp-cli listRoles  limit=value  offset=value
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **limit** | **integer** | The numbers of roles to return | [optional] [default to null]
 **offset** | **integer** | The number of items to skip before starting to collect the result | [optional] [default to null]

### Return type

[**array[Role]**](Role.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not Applicable
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## listScopesForRole

Returns a list of scopes for a role

Returns a list of scopes for a role

### Example

```bash
socialapp-cli listScopesForRole id=value  limit=value  offset=value
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **integer** | The id of the role | [default to null]
 **limit** | **integer** | The numbers of scopes to return | [optional] [default to null]
 **offset** | **integer** | The number of items to skip before starting to collect the result | [optional] [default to null]

### Return type

[**array[Scope]**](Scope.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not Applicable
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## removeScopeFromRole

Remove a scope from a role

Remove a scope from a role

### Example

```bash
socialapp-cli removeScopeFromRole role_id=value scope_id=value
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **roleId** | **integer** | The id of the role | [default to null]
 **scopeId** | **integer** | The id of the scope | [default to null]

### Return type

(empty response body)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not Applicable
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## updateRole

Update a role

Update a role

### Example

```bash
socialapp-cli updateRole id=value
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **integer** | id of the role | [default to null]
 **role** | [**Role**](Role.md) | Update a role | [optional]

### Return type

[**Role**](Role.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

