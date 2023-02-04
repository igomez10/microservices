# URLApi

All URIs are relative to **

Method | HTTP request | Description
------------- | ------------- | -------------
[**createUrl**](URLApi.md#createUrl) | **POST** /urls | Create a new url
[**deleteUrl**](URLApi.md#deleteUrl) | **DELETE** /urls/{alias} | Delete a url
[**getUrl**](URLApi.md#getUrl) | **GET** /urls/{alias} | Get a url
[**getUrlData**](URLApi.md#getUrlData) | **GET** /urls/{alias}/data | Returns a url metadata



## createUrl

Create a new url

Returns a url

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

Delete a url

Delete a url

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

Get a url

Returns a url

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

Returns a url metadata

Returns a url

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
