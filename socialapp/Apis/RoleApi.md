# RoleApi

All URIs are relative to *https://microservices.onrender.com*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**createRole**](RoleApi.md#createRole) | **POST** /roles | Create a new role |
| [**deleteRole**](RoleApi.md#deleteRole) | **DELETE** /roles/{id} | Delete a role |
| [**getRole**](RoleApi.md#getRole) | **GET** /roles/{id} | Returns a role |
| [**listRoles**](RoleApi.md#listRoles) | **GET** /roles | Returns a list of roles |
| [**updateRole**](RoleApi.md#updateRole) | **PUT** /roles/{id} | Update a role |


<a name="createRole"></a>
# **createRole**
> Role createRole(Role)

Create a new role

    Create a new role

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
> Role deleteRole(id)

Delete a role

    Delete a role

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **id** | **Integer**| id of the role | [default to null] |

### Return type

[**Role**](../Models/Role.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="getRole"></a>
# **getRole**
> Role getRole(id)

Returns a role

    Returns a role

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **id** | **Integer**| The id of the role | [default to null] |

### Return type

[**Role**](../Models/Role.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="listRoles"></a>
# **listRoles**
> List listRoles(offset, limit)

Returns a list of roles

    Returns a list of roles

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **offset** | **Integer**| The number of items to skip before starting to collect the result | [optional] [default to null] |
| **limit** | **Integer**| The numbers of roles to return | [optional] [default to null] |

### Return type

[**List**](../Models/Role.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="updateRole"></a>
# **updateRole**
> Role updateRole(id, Role)

Update a role

    Update a role

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **id** | **Integer**| id of the role | [default to null] |
| **Role** | [**Role**](../Models/Role.md)| Update a role | [optional] |

### Return type

[**Role**](../Models/Role.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

