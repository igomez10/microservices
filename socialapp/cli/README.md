# Socialapp Bash client

## Overview

This is a Bash client script for accessing Socialapp service.

The script uses cURL underneath for making all REST calls.

## Usage

```shell
# Make sure the script has executable rights
$ chmod u+x socialapp-cli

# Print the list of operations available on the service
$ ./socialapp-cli -h

# Print the service description
$ ./socialapp-cli --about

# Print detailed information about specific operation
$ ./socialapp-cli <operationId> -h

# Make GET request
./socialapp-cli --host http://<hostname>:<port> --accept xml <operationId> <queryParam1>=<value1> <header_key1>:<header_value2>

# Make GET request using arbitrary curl options (must be passed before <operationId>) to an SSL service using username:password
socialapp-cli -k -sS --tlsv1.2 --host https://<hostname> -u <user>:<password> --accept xml <operationId> <queryParam1>=<value1> <header_key1>:<header_value2>

# Make POST request
$ echo '<body_content>' | socialapp-cli --host <hostname> --content-type json <operationId> -

# Make POST request with simple JSON content, e.g.:
# {
#   "key1": "value1",
#   "key2": "value2",
#   "key3": 23
# }
$ echo '<body_content>' | socialapp-cli --host <hostname> --content-type json <operationId> key1==value1 key2=value2 key3:=23 -

# Make POST request with form data
$ socialapp-cli --host <hostname> <operationId> key1:=value1 key2:=value2 key3:=23

# Preview the cURL command without actually executing it
$ socialapp-cli --host http://<hostname>:<port> --dry-run <operationid>

```

## Docker image

You can easily create a Docker image containing a preconfigured environment
for using the REST Bash client including working autocompletion and short
welcome message with basic instructions, using the generated Dockerfile:

```shell
docker build -t my-rest-client .
docker run -it my-rest-client
```

By default you will be logged into a Zsh environment which has much more
advanced auto completion, but you can switch to Bash, where basic autocompletion
is also available.

## Shell completion

### Bash

The generated bash-completion script can be either directly loaded to the current Bash session using:

```shell
source socialapp-cli.bash-completion
```

Alternatively, the script can be copied to the `/etc/bash-completion.d` (or on OSX with Homebrew to `/usr/local/etc/bash-completion.d`):

```shell
sudo cp socialapp-cli.bash-completion /etc/bash-completion.d/socialapp-cli
```

#### OS X

On OSX you might need to install bash-completion using Homebrew:

```shell
brew install bash-completion
```

and add the following to the `~/.bashrc`:

```shell
if [ -f $(brew --prefix)/etc/bash_completion ]; then
  . $(brew --prefix)/etc/bash_completion
fi
```

### Zsh

In Zsh, the generated `_socialapp-cli` Zsh completion file must be copied to one of the folders under `$FPATH` variable.

## Documentation for API Endpoints

All URIs are relative to **

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*AuthenticationApi* | [**getAccessToken**](docs/AuthenticationApi.md#getaccesstoken) | **POST** /v1/oauth/token | Get an access token
*CommentApi* | [**createComment**](docs/CommentApi.md#createcomment) | **POST** /v1/comments | Create a new comment
*CommentApi* | [**getComment**](docs/CommentApi.md#getcomment) | **GET** /v1/comments/{id} | Returns details about a particular comment
*CommentApi* | [**getUserComments**](docs/CommentApi.md#getusercomments) | **GET** /v1/users/{username}/comments | Gets all comments for a user
*CommentApi* | [**getUserFeed**](docs/CommentApi.md#getuserfeed) | **GET** /v1/feed | Returns a users feed
*FollowingApi* | [**getUserFollowers**](docs/FollowingApi.md#getuserfollowers) | **GET** /v1/users/{username}/followers | Get all followers for a user
*RoleApi* | [**addScopeToRole**](docs/RoleApi.md#addscopetorole) | **POST** /v1/roles/{id}/scopes | Add a scope to a role
*RoleApi* | [**createRole**](docs/RoleApi.md#createrole) | **POST** /v1/roles | Create a new role
*RoleApi* | [**deleteRole**](docs/RoleApi.md#deleterole) | **DELETE** /v1/roles/{id} | Delete a role
*RoleApi* | [**getRole**](docs/RoleApi.md#getrole) | **GET** /v1/roles/{id} | Returns a role
*RoleApi* | [**listRoles**](docs/RoleApi.md#listroles) | **GET** /v1/roles | Returns a list of roles
*RoleApi* | [**listScopesForRole**](docs/RoleApi.md#listscopesforrole) | **GET** /v1/roles/{id}/scopes | Returns a list of scopes for a role
*RoleApi* | [**removeScopeFromRole**](docs/RoleApi.md#removescopefromrole) | **DELETE** /v1/roles/{role_id}/scopes/{scope_id} | Remove a scope from a role
*RoleApi* | [**updateRole**](docs/RoleApi.md#updaterole) | **PUT** /v1/roles/{id} | Update a role
*ScopeApi* | [**createScope**](docs/ScopeApi.md#createscope) | **POST** /v1/scopes | Create a new scope
*ScopeApi* | [**deleteScope**](docs/ScopeApi.md#deletescope) | **DELETE** /v1/scopes/{id} | Delete a scope
*ScopeApi* | [**getScope**](docs/ScopeApi.md#getscope) | **GET** /v1/scopes/{id} | Returns a scope
*ScopeApi* | [**listScopes**](docs/ScopeApi.md#listscopes) | **GET** /v1/scopes | Returns a list of scopes
*ScopeApi* | [**updateScope**](docs/ScopeApi.md#updatescope) | **PUT** /v1/scopes/{id} | Update a scope
*URLApi* | [**createUrl**](docs/URLApi.md#createurl) | **POST** /v1/urls | Create a new url
*URLApi* | [**deleteUrl**](docs/URLApi.md#deleteurl) | **DELETE** /v1/urls/{alias} | Delete a url
*URLApi* | [**getUrl**](docs/URLApi.md#geturl) | **GET** /v1/urls/{alias} | Get a url
*URLApi* | [**getUrlData**](docs/URLApi.md#geturldata) | **GET** /v1/urls/{alias}/data | Returns a url metadata
*UserApi* | [**changePassword**](docs/UserApi.md#changepassword) | **POST** /v1/password | Change password
*UserApi* | [**createUser**](docs/UserApi.md#createuser) | **POST** /v1/users | Create user
*UserApi* | [**deleteUser**](docs/UserApi.md#deleteuser) | **DELETE** /v1/users/{username} | Deletes a particular user
*UserApi* | [**followUser**](docs/UserApi.md#followuser) | **POST** /v1/users/{followedUsername}/followers/{followerUsername} | Add a user as a follower
*UserApi* | [**getFollowingUsers**](docs/UserApi.md#getfollowingusers) | **GET** /v1/users/{username}/following | Get all followed users for a user
*UserApi* | [**getRolesForUser**](docs/UserApi.md#getrolesforuser) | **GET** /v1/users/{username}/roles | Get all roles for a user
*UserApi* | [**getUserByUsername**](docs/UserApi.md#getuserbyusername) | **GET** /v1/users/{username} | Get a particular user by username
*UserApi* | [**getUserComments**](docs/UserApi.md#getusercomments) | **GET** /v1/users/{username}/comments | Gets all comments for a user
*UserApi* | [**getUserFollowers**](docs/UserApi.md#getuserfollowers) | **GET** /v1/users/{username}/followers | Get all followers for a user
*UserApi* | [**listUsers**](docs/UserApi.md#listusers) | **GET** /v1/users | List users
*UserApi* | [**resetPassword**](docs/UserApi.md#resetpassword) | **PUT** /v1/password | Reset password
*UserApi* | [**unfollowUser**](docs/UserApi.md#unfollowuser) | **DELETE** /v1/users/{followedUsername}/followers/{followerUsername} | Remove a user as a follower
*UserApi* | [**updateRolesForUser**](docs/UserApi.md#updaterolesforuser) | **PUT** /v1/users/{username}/roles | Update all roles for a user
*UserApi* | [**updateUser**](docs/UserApi.md#updateuser) | **PUT** /v1/users/{username} | Update a user


## Documentation For Models

 - [AccessToken](docs/AccessToken.md)
 - [ChangePasswordRequest](docs/ChangePasswordRequest.md)
 - [Comment](docs/Comment.md)
 - [CreateUserRequest](docs/CreateUserRequest.md)
 - [CreateUserResponse](docs/CreateUserResponse.md)
 - [Error](docs/Error.md)
 - [ResetPasswordRequest](docs/ResetPasswordRequest.md)
 - [Role](docs/Role.md)
 - [Scope](docs/Scope.md)
 - [URL](docs/URL.md)
 - [User](docs/User.md)


## Documentation For Authorization


## BasicAuth


- **Type**: HTTP basic authentication

## OAuth2


- **Type**: OAuth
- **Flow**: application
- **Token URL**: /v1/oauth/token
- **Scopes**:
  - **socialapp.users.list**: List users
  - **socialapp.users.create**: Create users
  - **socialapp.users.update**: Update users
  - **socialapp.users.delete**: Delete users
  - **socialapp.users.read**: Read a user
  - **socialapp.comments.list**: List comments
  - **socialapp.comments.create**: Create comments
  - **socialapp.comments.update**: Update comments
  - **socialapp.comments.delete**: Delete comments
  - **socialapp.followers.list**: List followers
  - **socialapp.following.list**: List following
  - **socialapp.roles.list**: List roles
  - **socialapp.roles.create**: Create roles
  - **socialapp.roles.read**: Read a role
  - **socialapp.roles.update**: Update roles
  - **socialapp.roles.delete**: Delete roles
  - **socialapp.scopes.list**: List scopes
  - **socialapp.scopes.create**: Create scopes
  - **socialapp.scopes.read**: Read a scope
  - **socialapp.scopes.update**: Update scopes
  - **socialapp.scopes.delete**: Delete scopes
  - **socialapp.roles.list_scopes**: List scopes of a role
  - **socialapp.roles.scopes.create**: Create scopes of a role
  - **socialapp.roles.scopes.delete**: Delete scopes of a role
  - **socialapp.users.roles.list**: List roles of a user
  - **socialapp.users.roles.create**: Create roles of a user
  - **socialapp.users.roles.update**: Update roles of a user
  - **socialapp.users.roles.delete**: Delete roles of a user
  - **shortly.url.create**: Create a url
  - **shortly.url.update**: Update a url
  - **shortly.url.delete**: Delete a url

