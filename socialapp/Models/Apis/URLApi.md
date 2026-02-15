# URLApi

All URIs are relative to *https://socialapp.gomezignacio.com*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**createUrl**](URLApi.md#createUrl) | **POST** /v1/urls | Create URL alias |
| [**deleteUrl**](URLApi.md#deleteUrl) | **DELETE** /v1/urls/{alias} | Delete URL alias |
| [**getUrl**](URLApi.md#getUrl) | **GET** /v1/urls/{alias} | Redirect using URL alias |
| [**getUrlData**](URLApi.md#getUrlData) | **GET** /v1/urls/{alias}/data | Get URL metadata |


<a name="createUrl"></a>
# **createUrl**
> URL createUrl(URL)

Create URL alias

    Create a new URL alias.

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

Delete URL alias

    Delete a URL alias and its metadata.

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

Redirect using URL alias

    Resolve an alias and redirect to the destination URL.

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

Get URL metadata

    Return metadata for a URL alias.

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

