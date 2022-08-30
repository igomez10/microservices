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
*CommentApi* | [**createComment**](docs/CommentApi.md#createcomment) | **POST** /comments | Create a new comment
*CommentApi* | [**getComment**](docs/CommentApi.md#getcomment) | **GET** /comments/{id} | Returns details about a particular comment
*CommentApi* | [**getUserComments**](docs/CommentApi.md#getusercomments) | **GET** /users/{username}/comments | Gets all comments for a user
*UserApi* | [**createUser**](docs/UserApi.md#createuser) | **POST** /users | Create a new user
*UserApi* | [**deleteUser**](docs/UserApi.md#deleteuser) | **DELETE** /users/{username} | Deletes a particular user
*UserApi* | [**getUserByUsername**](docs/UserApi.md#getuserbyusername) | **GET** /users/{username} | Get a particular user by username
*UserApi* | [**getUserComments**](docs/UserApi.md#getusercomments) | **GET** /users/{username}/comments | Gets all comments for a user
*UserApi* | [**listUsers**](docs/UserApi.md#listusers) | **GET** /users | Returns all the users
*UserApi* | [**updateUser**](docs/UserApi.md#updateuser) | **PUT** /users/{username} | Update a user


## Documentation For Models

 - [Comment](docs/Comment.md)
 - [Error](docs/Error.md)
 - [User](docs/User.md)


## Documentation For Authorization


## BasicAuth

- **Type**: HTTP basic authentication

## BearerAuth

- **Type**: HTTP basic authentication
