# FollowingApi

All URIs are relative to *https://socialapp.gomezignacio.com*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**getUserFollowers**](FollowingApi.md#getUserFollowers) | **GET** /v1/users/{username}/followers | Get all followers for a user |


<a name="getUserFollowers"></a>
# **getUserFollowers**
> List getUserFollowers(username)

Get all followers for a user

    Get all followers for a user

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **username** | **String**| username of the user | [default to null] |

### Return type

[**List**](../Models/User.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

