# Documentation for Socialapp

<a name="documentation-for-api-endpoints"></a>
## Documentation for API Endpoints

All URIs are relative to *https://socialapp.gomezignacio.com*

| Class | Method | HTTP request | Description |
|------------ | ------------- | ------------- | -------------|
| *AuthenticationApi* | [**getAccessToken**](Apis/AuthenticationApi.md#getaccesstoken) | **POST** /v1/oauth/token | Get an access token |
| *CommentApi* | [**createComment**](Apis/CommentApi.md#createcomment) | **POST** /v1/comments | Create a new comment |
*CommentApi* | [**getComment**](Apis/CommentApi.md#getcomment) | **GET** /v1/comments/{id} | Returns details about a particular comment |
*CommentApi* | [**getUserComments**](Apis/CommentApi.md#getusercomments) | **GET** /v1/users/{username}/comments | Gets all comments for a user |
*CommentApi* | [**getUserFeed**](Apis/CommentApi.md#getuserfeed) | **GET** /v1/feed | Returns a users feed |
| *FollowingApi* | [**getUserFollowers**](Apis/FollowingApi.md#getuserfollowers) | **GET** /v1/users/{username}/followers | Get all followers for a user |
| *RoleApi* | [**addScopeToRole**](Apis/RoleApi.md#addscopetorole) | **POST** /v1/roles/{id}/scopes | Add a scope to a role |
*RoleApi* | [**createRole**](Apis/RoleApi.md#createrole) | **POST** /v1/roles | Create a new role |
*RoleApi* | [**deleteRole**](Apis/RoleApi.md#deleterole) | **DELETE** /v1/roles/{id} | Delete a role |
*RoleApi* | [**getRole**](Apis/RoleApi.md#getrole) | **GET** /v1/roles/{id} | Returns a role |
*RoleApi* | [**listRoles**](Apis/RoleApi.md#listroles) | **GET** /v1/roles | Returns a list of roles |
*RoleApi* | [**listScopesForRole**](Apis/RoleApi.md#listscopesforrole) | **GET** /v1/roles/{id}/scopes | Returns a list of scopes for a role |
*RoleApi* | [**removeScopeFromRole**](Apis/RoleApi.md#removescopefromrole) | **DELETE** /v1/roles/{role_id}/scopes/{scope_id} | Remove a scope from a role |
*RoleApi* | [**updateRole**](Apis/RoleApi.md#updaterole) | **PUT** /v1/roles/{id} | Update a role |
| *ScopeApi* | [**createScope**](Apis/ScopeApi.md#createscope) | **POST** /v1/scopes | Create a new scope |
*ScopeApi* | [**deleteScope**](Apis/ScopeApi.md#deletescope) | **DELETE** /v1/scopes/{id} | Delete a scope |
*ScopeApi* | [**getScope**](Apis/ScopeApi.md#getscope) | **GET** /v1/scopes/{id} | Returns a scope |
*ScopeApi* | [**listScopes**](Apis/ScopeApi.md#listscopes) | **GET** /v1/scopes | Returns a list of scopes |
*ScopeApi* | [**updateScope**](Apis/ScopeApi.md#updatescope) | **PUT** /v1/scopes/{id} | Update a scope |
| *URLApi* | [**createUrl**](Apis/URLApi.md#createurl) | **POST** /v1/urls | Create a new url |
*URLApi* | [**deleteUrl**](Apis/URLApi.md#deleteurl) | **DELETE** /v1/urls/{alias} | Delete a url |
*URLApi* | [**getUrl**](Apis/URLApi.md#geturl) | **GET** /v1/urls/{alias} | Get a url |
*URLApi* | [**getUrlData**](Apis/URLApi.md#geturldata) | **GET** /v1/urls/{alias}/data | Returns a url metadata |
| *UserApi* | [**changePassword**](Apis/UserApi.md#changepassword) | **POST** /v1/password | Change password |
*UserApi* | [**createUser**](Apis/UserApi.md#createuser) | **POST** /v1/users | Create user |
*UserApi* | [**deleteUser**](Apis/UserApi.md#deleteuser) | **DELETE** /v1/users/{username} | Deletes a particular user |
*UserApi* | [**followUser**](Apis/UserApi.md#followuser) | **POST** /v1/users/{followedUsername}/followers/{followerUsername} | Add a user as a follower |
*UserApi* | [**getFollowingUsers**](Apis/UserApi.md#getfollowingusers) | **GET** /v1/users/{username}/following | Get all followed users for a user |
*UserApi* | [**getRolesForUser**](Apis/UserApi.md#getrolesforuser) | **GET** /v1/users/{username}/roles | Get all roles for a user |
*UserApi* | [**getUserByUsername**](Apis/UserApi.md#getuserbyusername) | **GET** /v1/users/{username} | Get a particular user by username |
*UserApi* | [**getUserComments**](Apis/UserApi.md#getusercomments) | **GET** /v1/users/{username}/comments | Gets all comments for a user |
*UserApi* | [**getUserFollowers**](Apis/UserApi.md#getuserfollowers) | **GET** /v1/users/{username}/followers | Get all followers for a user |
*UserApi* | [**listUsers**](Apis/UserApi.md#listusers) | **GET** /v1/users | List users |
*UserApi* | [**resetPassword**](Apis/UserApi.md#resetpassword) | **PUT** /v1/password | Reset password |
*UserApi* | [**unfollowUser**](Apis/UserApi.md#unfollowuser) | **DELETE** /v1/users/{followedUsername}/followers/{followerUsername} | Remove a user as a follower |
*UserApi* | [**updateRolesForUser**](Apis/UserApi.md#updaterolesforuser) | **PUT** /v1/users/{username}/roles | Update all roles for a user |
*UserApi* | [**updateUser**](Apis/UserApi.md#updateuser) | **PUT** /v1/users/{username} | Update a user |


<a name="documentation-for-models"></a>
## Documentation for Models

 - [AccessToken](./Models/AccessToken.md)
 - [ChangePasswordRequest](./Models/ChangePasswordRequest.md)
 - [Comment](./Models/Comment.md)
 - [CreateUserRequest](./Models/CreateUserRequest.md)
 - [CreateUserResponse](./Models/CreateUserResponse.md)
 - [Error](./Models/Error.md)
 - [ResetPasswordRequest](./Models/ResetPasswordRequest.md)
 - [Role](./Models/Role.md)
 - [Scope](./Models/Scope.md)
 - [URL](./Models/URL.md)
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
  - socialapp.roles.scopes.create: Create scopes of a role
  - socialapp.roles.scopes.delete: Delete scopes of a role
  - socialapp.users.roles.list: List roles of a user
  - socialapp.users.roles.create: Create roles of a user
  - socialapp.users.roles.update: Update roles of a user
  - socialapp.users.roles.delete: Delete roles of a user
  - shortly.url.create: Create a url
  - shortly.url.update: Update a url
  - shortly.url.delete: Delete a url

