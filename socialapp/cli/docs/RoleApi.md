# RoleApi

All URIs are relative to **

Method | HTTP request | Description
------------- | ------------- | -------------
[**createRole**](RoleApi.md#createRole) | **POST** /roles | Create a new role
[**deleteRole**](RoleApi.md#deleteRole) | **DELETE** /roles/{id} | Delete a role
[**getRole**](RoleApi.md#getRole) | **GET** /roles/{id} | Returns a role
[**listRoles**](RoleApi.md#listRoles) | **GET** /roles | Returns a list of roles
[**updateRole**](RoleApi.md#updateRole) | **PUT** /roles/{id} | Update a role



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
socialapp-cli listRoles  offset=value  limit=value
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **offset** | **integer** | The number of items to skip before starting to collect the result | [optional] [default to null]
 **limit** | **integer** | The numbers of roles to return | [optional] [default to null]

### Return type

[**array[Role]**](Role.md)

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

