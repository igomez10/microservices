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
*AuthenticationApi* | [**getAccessToken**](docs/AuthenticationApi.md#getaccesstoken) | **POST** /oauth/token | Get an access token
*CommentApi* | [**createComment**](docs/CommentApi.md#createcomment) | **POST** /comments | Create a new comment
*CommentApi* | [**getComment**](docs/CommentApi.md#getcomment) | **GET** /comments/{id} | Returns details about a particular comment
*CommentApi* | [**getUserComments**](docs/CommentApi.md#getusercomments) | **GET** /users/{username}/comments | Gets all comments for a user
*CommentApi* | [**getUserFeed**](docs/CommentApi.md#getuserfeed) | **GET** /users/{username}/feed | Returns a users feed
*FollowingApi* | [**getUserFollowers**](docs/FollowingApi.md#getuserfollowers) | **GET** /users/{username}/followers | Get all followers for a user
*UserApi* | [**changePassword**](docs/UserApi.md#changepassword) | **POST** /password | Change password
*UserApi* | [**createUser**](docs/UserApi.md#createuser) | **POST** /users | Create user
*UserApi* | [**deleteUser**](docs/UserApi.md#deleteuser) | **DELETE** /users/{username} | Deletes a particular user
*UserApi* | [**followUser**](docs/UserApi.md#followuser) | **POST** /users/{followedUsername}/followers/{followerUsername} | Add a user as a follower
*UserApi* | [**getFollowingUsers**](docs/UserApi.md#getfollowingusers) | **GET** /users/{username}/following | Get all followed users for a user
*UserApi* | [**getUserByUsername**](docs/UserApi.md#getuserbyusername) | **GET** /users/{username} | Get a particular user by username
*UserApi* | [**getUserComments**](docs/UserApi.md#getusercomments) | **GET** /users/{username}/comments | Gets all comments for a user
*UserApi* | [**getUserFollowers**](docs/UserApi.md#getuserfollowers) | **GET** /users/{username}/followers | Get all followers for a user
*UserApi* | [**listUsers**](docs/UserApi.md#listusers) | **GET** /users | List users
*UserApi* | [**resetPassword**](docs/UserApi.md#resetpassword) | **PUT** /password | Reset password
*UserApi* | [**unfollowUser**](docs/UserApi.md#unfollowuser) | **DELETE** /users/{followedUsername}/followers/{followerUsername} | Remove a user as a follower
*UserApi* | [**updateUser**](docs/UserApi.md#updateuser) | **PUT** /users/{username} | Update a user


## Documentation For Models

 - [AccessToken](docs/AccessToken.md)
 - [ChangePasswordRequest](docs/ChangePasswordRequest.md)
 - [Comment](docs/Comment.md)
 - [CreateUserRequest](docs/CreateUserRequest.md)
 - [Error](docs/Error.md)
 - [ResetPasswordRequest](docs/ResetPasswordRequest.md)
 - [User](docs/User.md)


## Documentation For Authorization


## BasicAuth

- **Type**: HTTP basic authentication

## BearerAuth

- **Type**: HTTP basic authentication

## OAuth2


- **Type**: OAuth
- **Flow**: application
- **Token URL**: localhost:8080/oauth/token
- **Scopes**:
  - **write**: modify your data in your account
  - **read**: read your data

