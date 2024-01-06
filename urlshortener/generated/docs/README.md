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

