# RoleApi

All URIs are relative to *https://socialapp.gomezignacio.com*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**addScopeToRole**](RoleApi.md#addScopeToRole) | **POST** /v1/roles/{id}/scopes | Add a scope to a role |
| [**createRole**](RoleApi.md#createRole) | **POST** /v1/roles | Create a new role |
| [**deleteRole**](RoleApi.md#deleteRole) | **DELETE** /v1/roles/{id} | Delete a role |
| [**getRole**](RoleApi.md#getRole) | **GET** /v1/roles/{id} | Get role by ID |
| [**listRoles**](RoleApi.md#listRoles) | **GET** /v1/roles | Returns a list of roles |
| [**listScopesForRole**](RoleApi.md#listScopesForRole) | **GET** /v1/roles/{id}/scopes | List role scopes |
| [**removeScopeFromRole**](RoleApi.md#removeScopeFromRole) | **DELETE** /v1/roles/{role_id}/scopes/{scope_id} | Remove a scope from a role |
| [**updateRole**](RoleApi.md#updateRole) | **PUT** /v1/roles/{id} | Update a role |


<a name="addScopeToRole"></a>
# **addScopeToRole**
> addScopeToRole(id, request\_body)

Add a scope to a role

    Add one or more scopes to a role.

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **id** | **Long**| The id of the role | [default to null] |
| **request\_body** | [**List**](../Models/string.md)| Add a scope to a role | |

### Return type

null (empty response body)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

<a name="createRole"></a>
# **createRole**
> Role createRole(Role)

Create a new role

    Create a new role.

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **Role** | [**Role**](../Models/Role.md)| Create a new role | |

### Return type

[**Role**](../Models/Role.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

<a name="deleteRole"></a>
# **deleteRole**
> deleteRole(id)

Delete a role

    Delete a role

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **id** | **Long**| id of the role | [default to null] |

### Return type

null (empty response body)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="getRole"></a>
# **getRole**
> Role getRole(id)

Get role by ID

    Get a role by ID.

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **id** | **Long**| The id of the role | [default to null] |

### Return type

[**Role**](../Models/Role.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="listRoles"></a>
# **listRoles**
> List listRoles(limit, offset)

Returns a list of roles

    List roles with offset-based pagination.

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **limit** | **Integer**| The numbers of roles to return | [optional] [default to 20] |
| **offset** | **Integer**| The number of items to skip before starting to collect the result | [optional] [default to null] |

### Return type

[**List**](../Models/Role.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="listScopesForRole"></a>
# **listScopesForRole**
> List listScopesForRole(id, limit, offset)

List role scopes

    List scopes assigned to a role.

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **id** | **Long**| The id of the role | [default to null] |
| **limit** | **Integer**| The numbers of scopes to return | [optional] [default to 20] |
| **offset** | **Integer**| The number of items to skip before starting to collect the result | [optional] [default to null] |

### Return type

[**List**](../Models/Scope.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="removeScopeFromRole"></a>
# **removeScopeFromRole**
> removeScopeFromRole(role\_id, scope\_id)

Remove a scope from a role

    Remove a scope from a role.

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **role\_id** | **Long**| The id of the role | [default to null] |
| **scope\_id** | **Long**| The id of the scope | [default to null] |

### Return type

null (empty response body)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="updateRole"></a>
# **updateRole**
> Role updateRole(id, Role)

Update a role

    Update a role by ID.

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **id** | **Long**| id of the role | [default to null] |
| **Role** | [**Role**](../Models/Role.md)| Update a role | [optional] |

### Return type

[**Role**](../Models/Role.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

