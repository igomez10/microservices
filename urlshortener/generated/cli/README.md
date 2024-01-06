# URL Shortener Bash client

## Overview

This is a Bash client script for accessing URL Shortener service.

The script uses cURL underneath for making all REST calls.

## Usage

```shell
# Make sure the script has executable rights
$ chmod u+x urlshortener-cli

# Print the list of operations available on the service
$ ./urlshortener-cli -h

# Print the service description
$ ./urlshortener-cli --about

# Print detailed information about specific operation
$ ./urlshortener-cli <operationId> -h

# Make GET request
./urlshortener-cli --host http://<hostname>:<port> --accept xml <operationId> <queryParam1>=<value1> <header_key1>:<header_value2>

# Make GET request using arbitrary curl options (must be passed before <operationId>) to an SSL service using username:password
urlshortener-cli -k -sS --tlsv1.2 --host https://<hostname> -u <user>:<password> --accept xml <operationId> <queryParam1>=<value1> <header_key1>:<header_value2>

# Make POST request
$ echo '<body_content>' | urlshortener-cli --host <hostname> --content-type json <operationId> -

# Make POST request with simple JSON content, e.g.:
# {
#   "key1": "value1",
#   "key2": "value2",
#   "key3": 23
# }
$ echo '<body_content>' | urlshortener-cli --host <hostname> --content-type json <operationId> key1==value1 key2=value2 key3:=23 -

# Make POST request with form data
$ urlshortener-cli --host <hostname> <operationId> key1:=value1 key2:=value2 key3:=23

# Preview the cURL command without actually executing it
$ urlshortener-cli --host http://<hostname>:<port> --dry-run <operationid>

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
source urlshortener-cli.bash-completion
```

Alternatively, the script can be copied to the `/etc/bash-completion.d` (or on OSX with Homebrew to `/usr/local/etc/bash-completion.d`):

```shell
sudo cp urlshortener-cli.bash-completion /etc/bash-completion.d/urlshortener-cli
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

In Zsh, the generated `_urlshortener-cli` Zsh completion file must be copied to one of the folders under `$FPATH` variable.

## Documentation for API Endpoints

All URIs are relative to **

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*URLApi* | [**createUrl**](docs/URLApi.md#createurl) | **POST** /v1/urls | Create a new url
*URLApi* | [**deleteUrl**](docs/URLApi.md#deleteurl) | **DELETE** /v1/urls/{alias} | Delete a url
*URLApi* | [**getUrl**](docs/URLApi.md#geturl) | **GET** /v1/urls/{alias} | Get a url
*URLApi* | [**getUrlData**](docs/URLApi.md#geturldata) | **GET** /v1/urls/{alias}/data | Returns a url metadata


## Documentation For Models

 - [Error](docs/Error.md)
 - [URL](docs/URL.md)


## Documentation For Authorization


## BasicAuth


- **Type**: HTTP basic authentication

## OAuth2


- **Type**: OAuth
- **Flow**: application
- **Token URL**: /v1/oauth/token
- **Scopes**:
  - **shortly.url.create**: Create a url
  - **shortly.url.update**: Update a url
  - **shortly.url.delete**: Delete a url

