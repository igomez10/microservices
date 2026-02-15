# URLApi

All URIs are relative to **

Method | HTTP request | Description
------------- | ------------- | -------------
[**createUrl**](URLApi.md#createUrl) | **POST** /v1/urls | Create URL alias
[**deleteUrl**](URLApi.md#deleteUrl) | **DELETE** /v1/urls/{alias} | Delete URL alias
[**getUrl**](URLApi.md#getUrl) | **GET** /v1/urls/{alias} | Redirect using URL alias
[**getUrlData**](URLApi.md#getUrlData) | **GET** /v1/urls/{alias}/data | Get URL metadata



## createUrl

Create URL alias

Create a new URL alias.

### Example

```bash
socialapp-cli createUrl
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **uRL** | [**URL**](URL.md) | Create a new url |

### Return type

[**URL**](URL.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## deleteUrl

Delete URL alias

Delete a URL alias and its metadata.

### Example

```bash
socialapp-cli deleteUrl alias=value
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **alias** | **string** | The alias of the url | [default to null]

### Return type

(empty response body)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not Applicable
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## getUrl

Redirect using URL alias

Resolve an alias and redirect to the destination URL.

### Example

```bash
socialapp-cli getUrl alias=value
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **alias** | **string** | The alias of the url | [default to null]

### Return type

(empty response body)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not Applicable
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## getUrlData

Get URL metadata

Return metadata for a URL alias.

### Example

```bash
socialapp-cli getUrlData alias=value
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **alias** | **string** | The alias of the url | [default to null]

### Return type

[**URL**](URL.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not Applicable
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

