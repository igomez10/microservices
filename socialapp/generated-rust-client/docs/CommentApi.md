# \CommentApi

All URIs are relative to *https://microservices.onrender.com*

Method | HTTP request | Description
------------- | ------------- | -------------
[**create_comment**](CommentApi.md#create_comment) | **POST** /comments | Create a new comment
[**get_comment**](CommentApi.md#get_comment) | **GET** /comments/{id} | Returns details about a particular comment
[**get_user_comments**](CommentApi.md#get_user_comments) | **GET** /users/{username}/comments | Gets all comments for a user
[**get_user_feed**](CommentApi.md#get_user_feed) | **GET** /users/{username}/feed | Returns a users feed



## create_comment

> crate::models::Comment create_comment(comment)
Create a new comment

Create a new comment

### Parameters


Name | Type | Description  | Required | Notes
------------- | ------------- | ------------- | ------------- | -------------
**comment** | [**Comment**](Comment.md) | Create a new comment | [required] |

### Return type

[**crate::models::Comment**](Comment.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## get_comment

> crate::models::Comment get_comment(id)
Returns details about a particular comment

Returns details about a particular comment

### Parameters


Name | Type | Description  | Required | Notes
------------- | ------------- | ------------- | ------------- | -------------
**id** | **i32** | ID of the comment | [required] |

### Return type

[**crate::models::Comment**](Comment.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## get_user_comments

> Vec<crate::models::Comment> get_user_comments(username, limit, offset)
Gets all comments for a user

Gets all comments for a user

### Parameters


Name | Type | Description  | Required | Notes
------------- | ------------- | ------------- | ------------- | -------------
**username** | **String** | username of the user | [required] |
**limit** | Option<**i32**> | How many items to return at one time (max 100) |  |
**offset** | Option<**i32**> | The number of items to skip before starting to collect the result set |  |

### Return type

[**Vec<crate::models::Comment>**](Comment.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## get_user_feed

> Vec<crate::models::Comment> get_user_feed(username)
Returns a users feed

Returns a users feed

### Parameters


Name | Type | Description  | Required | Notes
------------- | ------------- | ------------- | ------------- | -------------
**username** | **String** | The username of the user | [required] |

### Return type

[**Vec<crate::models::Comment>**](Comment.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

