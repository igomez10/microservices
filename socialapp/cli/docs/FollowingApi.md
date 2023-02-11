# FollowingApi

All URIs are relative to **

Method | HTTP request | Description
------------- | ------------- | -------------
[**getUserFollowers**](FollowingApi.md#getUserFollowers) | **GET** /v1/users/{username}/followers | Get all followers for a user



## getUserFollowers

Get all followers for a user

Get all followers for a user

### Example

```bash
socialapp-cli getUserFollowers username=value
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **username** | **string** | username of the user | [default to null]

### Return type

[**array[User]**](User.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not Applicable
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

