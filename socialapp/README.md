# Documentation for Socialapp

<a name="documentation-for-api-endpoints"></a>
## Documentation for API Endpoints

All URIs are relative to *https://microservices.onrender.com*

| Class | Method | HTTP request | Description |
|------------ | ------------- | ------------- | -------------|
| *CommentApi* | [**createComment**](Apis/CommentApi.md#createcomment) | **POST** /comments | Create a new comment |
*CommentApi* | [**getComment**](Apis/CommentApi.md#getcomment) | **GET** /comments/{id} | Returns details about a particular comment |
*CommentApi* | [**getUserComments**](Apis/CommentApi.md#getusercomments) | **GET** /users/{username}/comments | Gets all comments for a user |
*CommentApi* | [**getUserFeed**](Apis/CommentApi.md#getuserfeed) | **GET** /users/{username}/feed | Returns a users feed |
| *UserApi* | [**createUser**](Apis/UserApi.md#createuser) | **POST** /users | Create a new user |
*UserApi* | [**deleteUser**](Apis/UserApi.md#deleteuser) | **DELETE** /users/{username} | Deletes a particular user |
*UserApi* | [**followUser**](Apis/UserApi.md#followuser) | **POST** /users/{followedUsername}/followers/{followerUsername} | Add a user as a follower |
*UserApi* | [**getFollowingUsers**](Apis/UserApi.md#getfollowingusers) | **GET** /users/{username}/following | Get all followed users for a user |
*UserApi* | [**getUserByUsername**](Apis/UserApi.md#getuserbyusername) | **GET** /users/{username} | Get a particular user by username |
*UserApi* | [**getUserComments**](Apis/UserApi.md#getusercomments) | **GET** /users/{username}/comments | Gets all comments for a user |
*UserApi* | [**getUserFollowers**](Apis/UserApi.md#getuserfollowers) | **GET** /users/{username}/followers | Get all followers for a user |
*UserApi* | [**listUsers**](Apis/UserApi.md#listusers) | **GET** /users | Returns all the users |
*UserApi* | [**unfollowUser**](Apis/UserApi.md#unfollowuser) | **DELETE** /users/{followedUsername}/followers/{followerUsername} | Remove a user as a follower |
*UserApi* | [**updateUser**](Apis/UserApi.md#updateuser) | **PUT** /users/{username} | Update a user |


<a name="documentation-for-models"></a>
## Documentation for Models

 - [Comment](./Models/Comment.md)
 - [Error](./Models/Error.md)
 - [User](./Models/User.md)


<a name="documentation-for-authorization"></a>
## Documentation for Authorization

<a name="BasicAuth"></a>
### BasicAuth

- **Type**: HTTP basic authentication

<a name="BearerAuth"></a>
### BearerAuth

- **Type**: HTTP basic authentication

