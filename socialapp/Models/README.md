# Documentation for Socialapp

<a name="documentation-for-api-endpoints"></a>
## Documentation for API Endpoints

All URIs are relative to *https://socialapp.gomezignacio.com*

| Class | Method | HTTP request | Description |
|------------ | ------------- | ------------- | -------------|
| *AuthenticationApi* | [**getAccessToken**](Apis/AuthenticationApi.md#getAccessToken) | **POST** /v1/oauth/token | Get an access token |
| *CommentApi* | [**createComment**](Apis/CommentApi.md#createComment) | **POST** /v1/comments | Create a new comment |
*CommentApi* | [**getComment**](Apis/CommentApi.md#getComment) | **GET** /v1/comments/{id} | Get comment by ID |
*CommentApi* | [**getUserFeed**](Apis/CommentApi.md#getUserFeed) | **GET** /v1/feed | Get user feed |
*CommentApi* | [**likeComment**](Apis/CommentApi.md#likeComment) | **POST** /like | Like a comment |
*CommentApi* | [**searchComments**](Apis/CommentApi.md#searchComments) | **GET** /v1/comments | Search comments |
*CommentApi* | [**unlikeComment**](Apis/CommentApi.md#unlikeComment) | **DELETE** /like | Unlike a comment |
| *FollowingApi* | [**getUserFollowers**](Apis/FollowingApi.md#getUserFollowers) | **GET** /v1/users/{username}/followers | Get all followers for a user |
| *RoleApi* | [**addScopeToRole**](Apis/RoleApi.md#addScopeToRole) | **POST** /v1/roles/{id}/scopes | Add a scope to a role |
*RoleApi* | [**createRole**](Apis/RoleApi.md#createRole) | **POST** /v1/roles | Create a new role |
*RoleApi* | [**deleteRole**](Apis/RoleApi.md#deleteRole) | **DELETE** /v1/roles/{id} | Delete a role |
*RoleApi* | [**getRole**](Apis/RoleApi.md#getRole) | **GET** /v1/roles/{id} | Get role by ID |
*RoleApi* | [**listRoles**](Apis/RoleApi.md#listRoles) | **GET** /v1/roles | Returns a list of roles |
*RoleApi* | [**listScopesForRole**](Apis/RoleApi.md#listScopesForRole) | **GET** /v1/roles/{id}/scopes | List role scopes |
*RoleApi* | [**removeScopeFromRole**](Apis/RoleApi.md#removeScopeFromRole) | **DELETE** /v1/roles/{role_id}/scopes/{scope_id} | Remove a scope from a role |
*RoleApi* | [**updateRole**](Apis/RoleApi.md#updateRole) | **PUT** /v1/roles/{id} | Update a role |
| *ScopeApi* | [**createScope**](Apis/ScopeApi.md#createScope) | **POST** /v1/scopes | Create a new scope |
*ScopeApi* | [**deleteScope**](Apis/ScopeApi.md#deleteScope) | **DELETE** /v1/scopes/{id} | Delete a scope |
*ScopeApi* | [**getScope**](Apis/ScopeApi.md#getScope) | **GET** /v1/scopes/{id} | Get scope by ID |
*ScopeApi* | [**listScopes**](Apis/ScopeApi.md#listScopes) | **GET** /v1/scopes | List scopes |
*ScopeApi* | [**updateScope**](Apis/ScopeApi.md#updateScope) | **PUT** /v1/scopes/{id} | Update a scope |
| *URLApi* | [**createUrl**](Apis/URLApi.md#createUrl) | **POST** /v1/urls | Create URL alias |
*URLApi* | [**deleteUrl**](Apis/URLApi.md#deleteUrl) | **DELETE** /v1/urls/{alias} | Delete URL alias |
*URLApi* | [**getUrl**](Apis/URLApi.md#getUrl) | **GET** /v1/urls/{alias} | Redirect using URL alias |
*URLApi* | [**getUrlData**](Apis/URLApi.md#getUrlData) | **GET** /v1/urls/{alias}/data | Get URL metadata |
| *UserApi* | [**changePassword**](Apis/UserApi.md#changePassword) | **POST** /v1/password | Change password |
*UserApi* | [**createUser**](Apis/UserApi.md#createUser) | **POST** /v1/users | Create user |
*UserApi* | [**deleteUser**](Apis/UserApi.md#deleteUser) | **DELETE** /v1/users/{username} | Deletes a particular user |
*UserApi* | [**followUser**](Apis/UserApi.md#followUser) | **POST** /v1/users/{followedUsername}/followers/{followerUsername} | Add a user as a follower |
*UserApi* | [**getFollowingUsers**](Apis/UserApi.md#getFollowingUsers) | **GET** /v1/users/{username}/following | Get all followed users for a user |
*UserApi* | [**getRolesForUser**](Apis/UserApi.md#getRolesForUser) | **GET** /v1/users/{username}/roles | Get all roles for a user |
*UserApi* | [**getUserByUsername**](Apis/UserApi.md#getUserByUsername) | **GET** /v1/users/{username} | Get a particular user by username |
*UserApi* | [**getUserComments**](Apis/UserApi.md#getUserComments) | **GET** /v1/users/{username}/comments | List comments for a user |
*UserApi* | [**getUserFollowers**](Apis/UserApi.md#getUserFollowers) | **GET** /v1/users/{username}/followers | Get all followers for a user |
*UserApi* | [**listUsers**](Apis/UserApi.md#listUsers) | **GET** /v1/users | List users |
*UserApi* | [**resetPassword**](Apis/UserApi.md#resetPassword) | **PUT** /v1/password | Reset password |
*UserApi* | [**unfollowUser**](Apis/UserApi.md#unfollowUser) | **DELETE** /v1/users/{followedUsername}/followers/{followerUsername} | Remove a user as a follower |
*UserApi* | [**updateRolesForUser**](Apis/UserApi.md#updateRolesForUser) | **PUT** /v1/users/{username}/roles | Update all roles for a user |
*UserApi* | [**updateUser**](Apis/UserApi.md#updateUser) | **PUT** /v1/users/{username} | Update a user |
*UserApi* | [**welcome**](Apis/UserApi.md#welcome) | **GET** / | Get API welcome message |


<a name="documentation-for-models"></a>
## Documentation for Models

 - [AccessToken](./Models/AccessToken.md)
 - [ChangePasswordRequest](./Models/ChangePasswordRequest.md)
 - [Comment](./Models/Comment.md)
 - [CreateCommentRequest](./Models/CreateCommentRequest.md)
 - [CreateUserRequest](./Models/CreateUserRequest.md)
 - [CreateUserResponse](./Models/CreateUserResponse.md)
 - [Error](./Models/Error.md)
 - [LikeRequest](./Models/LikeRequest.md)
 - [ResetPasswordRequest](./Models/ResetPasswordRequest.md)
 - [Role](./Models/Role.md)
 - [Scope](./Models/Scope.md)
 - [URL](./Models/URL.md)
 - [User](./Models/User.md)


<a name="documentation-for-authorization"></a>
## Documentation for Authorization

<a name="BasicAuth"></a>
### BasicAuth

- **Type**: HTTP basic authentication

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
  - socialapp.comments.read: Read comments
  - socialapp.feed.read: Read user feed
  - socialapp.followers.list: List followers
  - socialapp.follower.create: Create follower relationship
  - socialapp.follower.read: Read follower information
  - socialapp.follower.delete: Delete follower relationship
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
  - socialapp.roles.scopes.list: List scopes of a role
  - socialapp.roles.scopes.create: Create scopes of a role
  - socialapp.roles.scopes.delete: Delete scopes of a role
  - socialapp.users.roles.list: List roles of a user
  - socialapp.users.roles.create: Create roles of a user
  - socialapp.users.roles.update: Update roles of a user
  - socialapp.users.roles.delete: Delete roles of a user
  - shortly.url.create: Create a url
  - shortly.url.update: Update a url
  - shortly.url.delete: Delete a url

