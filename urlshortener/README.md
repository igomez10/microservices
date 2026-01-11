# Documentation for URL Shortener

<a name="documentation-for-api-endpoints"></a>
## Documentation for API Endpoints

All URIs are relative to *https://urlshortener.gomezignacio.com*

| Class | Method | HTTP request | Description |
|------------ | ------------- | ------------- | -------------|
| *URLApi* | [**createUrl**](Apis/URLApi.md#createurl) | **POST** /v1/urls | Create a new url |
*URLApi* | [**deleteUrl**](Apis/URLApi.md#deleteurl) | **DELETE** /v1/urls/{alias} | Delete a url |
*URLApi* | [**getUrl**](Apis/URLApi.md#geturl) | **GET** /v1/urls/{alias} | Get a url |
*URLApi* | [**getUrlData**](Apis/URLApi.md#geturldata) | **GET** /v1/urls/{alias}/data | Returns a url metadata |


<a name="documentation-for-models"></a>
## Documentation for Models

 - [Error](./Models/Error.md)
 - [URL](./Models/URL.md)


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
- shortly.url.create: Create a url
- shortly.url.update: Update a url
- shortly.url.delete: Delete a url

Contract tests are also available via `make contract-test`. CI pipelines can run
`../scripts/run-contract-tests.sh` from the repo root to exercise Socialapp and
UrlShortener pact suites together.

## Contract Testing with Pact

The URL Shortener service verifies the consumer contract that Socialapp generates for `/v1/urls/{alias}/data`. Run `make contract-test` after installing the Pact CLI and native libraries to ensure the provider still meets those expectations:

```sh
go install github.com/pact-foundation/pact-go/v2@v2.4.2
pact-go -l DEBUG install --libDir /tmp
cd urlshortener
make contract-test
```

Logs appear under `contracts/logs/` and are ignored by git, while the pact file itself lives in Socialapp (`../socialapp/contracts/pacts/socialapp-urlshortener.json`).
