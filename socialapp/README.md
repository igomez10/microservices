# Documentation for Socialapp

<a name="documentation-for-api-endpoints"></a>
## Documentation for API Endpoints

All URIs are relative to *https://microservices.onrender.com*

| Class | Method | HTTP request | Description |
|------------ | ------------- | ------------- | -------------|
| *AuthenticationApi* | [**getAccessToken**](Apis/AuthenticationApi.md#getaccesstoken) | **POST** /oauth/token | Get an access token |
| *CommentApi* | [**createComment**](Apis/CommentApi.md#createcomment) | **POST** /comments | Create a new comment |
*CommentApi* | [**getComment**](Apis/CommentApi.md#getcomment) | **GET** /comments/{id} | Returns details about a particular comment |
*CommentApi* | [**getUserComments**](Apis/CommentApi.md#getusercomments) | **GET** /users/{username}/comments | Gets all comments for a user |
*CommentApi* | [**getUserFeed**](Apis/CommentApi.md#getuserfeed) | **GET** /users/{username}/feed | Returns a users feed |
| *FollowingApi* | [**getUserFollowers**](Apis/FollowingApi.md#getuserfollowers) | **GET** /users/{username}/followers | Get all followers for a user |
| *RoleApi* | [**createRole**](Apis/RoleApi.md#createrole) | **POST** /roles | Create a new role |
*RoleApi* | [**deleteRole**](Apis/RoleApi.md#deleterole) | **DELETE** /roles/{id} | Delete a role |
*RoleApi* | [**getRole**](Apis/RoleApi.md#getrole) | **GET** /roles/{id} | Returns a role |
*RoleApi* | [**listRoles**](Apis/RoleApi.md#listroles) | **GET** /roles | Returns a list of roles |
*RoleApi* | [**updateRole**](Apis/RoleApi.md#updaterole) | **PUT** /roles/{id} | Update a role |
| *ScopeApi* | [**createScope**](Apis/ScopeApi.md#createscope) | **POST** /scopes | Create a new scope |
*ScopeApi* | [**deleteScope**](Apis/ScopeApi.md#deletescope) | **DELETE** /scopes/{id} | Delete a scope |
*ScopeApi* | [**getScope**](Apis/ScopeApi.md#getscope) | **GET** /scopes/{id} | Returns a scope |
*ScopeApi* | [**listScopes**](Apis/ScopeApi.md#listscopes) | **GET** /scopes | Returns a list of scopes |
*ScopeApi* | [**updateScope**](Apis/ScopeApi.md#updatescope) | **PUT** /scopes/{id} | Update a scope |
| *UserApi* | [**changePassword**](Apis/UserApi.md#changepassword) | **POST** /password | Change password |
*UserApi* | [**createUser**](Apis/UserApi.md#createuser) | **POST** /users | Create user |
*UserApi* | [**deleteUser**](Apis/UserApi.md#deleteuser) | **DELETE** /users/{username} | Deletes a particular user |
*UserApi* | [**followUser**](Apis/UserApi.md#followuser) | **POST** /users/{followedUsername}/followers/{followerUsername} | Add a user as a follower |
*UserApi* | [**getFollowingUsers**](Apis/UserApi.md#getfollowingusers) | **GET** /users/{username}/following | Get all followed users for a user |
*UserApi* | [**getUserByUsername**](Apis/UserApi.md#getuserbyusername) | **GET** /users/{username} | Get a particular user by username |
*UserApi* | [**getUserComments**](Apis/UserApi.md#getusercomments) | **GET** /users/{username}/comments | Gets all comments for a user |
*UserApi* | [**getUserFollowers**](Apis/UserApi.md#getuserfollowers) | **GET** /users/{username}/followers | Get all followers for a user |
*UserApi* | [**listUsers**](Apis/UserApi.md#listusers) | **GET** /users | List users |
*UserApi* | [**resetPassword**](Apis/UserApi.md#resetpassword) | **PUT** /password | Reset password |
*UserApi* | [**unfollowUser**](Apis/UserApi.md#unfollowuser) | **DELETE** /users/{followedUsername}/followers/{followerUsername} | Remove a user as a follower |
*UserApi* | [**updateUser**](Apis/UserApi.md#updateuser) | **PUT** /users/{username} | Update a user |


<a name="documentation-for-models"></a>
## Documentation for Models

 - [AccessToken](./Models/AccessToken.md)
 - [ChangePasswordRequest](./Models/ChangePasswordRequest.md)
 - [Comment](./Models/Comment.md)
 - [CreateUserRequest](./Models/CreateUserRequest.md)
 - [Error](./Models/Error.md)
 - [ResetPasswordRequest](./Models/ResetPasswordRequest.md)
 - [Role](./Models/Role.md)
 - [Scope](./Models/Scope.md)
 - [User](./Models/User.md)


<a name="documentation-for-authorization"></a>
## Documentation for Authorization

<a name="OAuth2"></a>
### OAuth2

- **Type**: OAuth
- **Flow**: application
- **Authorization URL**: 
- **Scopes**: 
  - socialapp.users.list: List users
  - socialapp.users.create: Create users
  - socialapp.users.update: Update users
  - socialapp.users.delete: Delete users
  - socialapp.users.read: Read a user
  - socialapp.comments.list: List comments
  - socialapp.comments.create: Create comments
  - socialapp.comments.update: Update comments
  - socialapp.comments.delete: Delete comments
  - socialapp.followers.list: List followers
  - socialapp.following.list: List following
  - socialapp.roles.list: List roles
  - socialapp.roles.create: Create roles
  - socialapp.roles.read: Read a role
  - socialapp.roles.update: Update roles
  - socialapp.roles.delete: Delete roles
  - socialapp.scopes.list: List scopes
  - socialapp.scopes.create: Create scopes
  - socialapp.scopes.read: Read a scope
  - socialapp.scopes.update: Update scopes
  - socialapp.scopes.delete: Delete scopes
  - socialapp.roles.list_scopes: List scopes of a role

