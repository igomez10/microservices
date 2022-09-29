# \FollowingApi

All URIs are relative to *https://microservices.onrender.com*

Method | HTTP request | Description
------------- | ------------- | -------------
[**get_user_followers**](FollowingApi.md#get_user_followers) | **GET** /users/{username}/followers | Get all followers for a user



## get_user_followers

> Vec<crate::models::User> get_user_followers(username)
Get all followers for a user

Get all followers for a user

### Parameters


Name | Type | Description  | Required | Notes
------------- | ------------- | ------------- | ------------- | -------------
**username** | **String** | username of the user | [required] |

### Return type

[**Vec<crate::models::User>**](User.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

