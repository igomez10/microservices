# URLApi

All URIs are relative to *https://socialapp.gomezignacio.com*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**createUrl**](URLApi.md#createUrl) | **POST** /urls | Create a new url |
| [**deleteUrl**](URLApi.md#deleteUrl) | **DELETE** /urls/{alias} | Delete a url |
| [**getUrl**](URLApi.md#getUrl) | **GET** /urls/{alias} | Get a url |
| [**getUrlData**](URLApi.md#getUrlData) | **GET** /urls/{alias}/data | Returns a url metadata |


<a name="createUrl"></a>
# **createUrl**
> URL createUrl(URL)

Create a new url

    Returns a url

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **URL** | [**URL**](../Models/URL.md)| Create a new url | |

### Return type

[**URL**](../Models/URL.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

<a name="deleteUrl"></a>
# **deleteUrl**
> deleteUrl(alias)

Delete a url

    Delete a url

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **alias** | **String**| The alias of the url | [default to null] |

### Return type

null (empty response body)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="getUrl"></a>
# **getUrl**
> getUrl(alias)

Get a url

    Returns a url

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **alias** | **String**| The alias of the url | [default to null] |

### Return type

null (empty response body)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="getUrlData"></a>
# **getUrlData**
> URL getUrlData(alias)

Returns a url metadata

    Returns a url

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **alias** | **String**| The alias of the url | [default to null] |

### Return type

[**URL**](../Models/URL.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

