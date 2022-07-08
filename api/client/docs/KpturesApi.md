# \KpturesApi

All URIs are relative to *http://localhost:8080/kpture/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**KpturePost**](KpturesApi.md#KpturePost) | **Post** /kpture | Start Kpture
[**KptureUuidDelete**](KpturesApi.md#KptureUuidDelete) | **Delete** /kpture/{uuid} | Delete kapture
[**KptureUuidDownloadGet**](KpturesApi.md#KptureUuidDownloadGet) | **Get** /kpture/{uuid}/download | Download kpture
[**KptureUuidGet**](KpturesApi.md#KptureUuidGet) | **Get** /kpture/{uuid} | Get kapture
[**KptureUuidStopPut**](KpturesApi.md#KptureUuidStopPut) | **Put** /kpture/{uuid}/stop | Stop Kpture
[**KpturesGet**](KpturesApi.md#KpturesGet) | **Get** /kptures | Get kaptures



## KpturePost

> CaptureKpture KpturePost(ctx).Data(data).Execute()

Start Kpture



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    data := *openapiclient.NewServerKptureRequest() // ServerKptureRequest | selected agents for capture

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.KpturesApi.KpturePost(context.Background()).Data(data).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `KpturesApi.KpturePost``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `KpturePost`: CaptureKpture
    fmt.Fprintf(os.Stdout, "Response from `KpturesApi.KpturePost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiKpturePostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **data** | [**ServerKptureRequest**](ServerKptureRequest.md) | selected agents for capture | 

### Return type

[**CaptureKpture**](CaptureKpture.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## KptureUuidDelete

> KptureUuidDelete(ctx, uuid).Execute()

Delete kapture



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    uuid := "uuid_example" // string | capture uuid

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.KpturesApi.KptureUuidDelete(context.Background(), uuid).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `KpturesApi.KptureUuidDelete``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**uuid** | **string** | capture uuid | 

### Other Parameters

Other parameters are passed through a pointer to a apiKptureUuidDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## KptureUuidDownloadGet

> KptureUuidDownloadGet(ctx, uuid).Execute()

Download kpture



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    uuid := "uuid_example" // string | capture uuid

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.KpturesApi.KptureUuidDownloadGet(context.Background(), uuid).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `KpturesApi.KptureUuidDownloadGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**uuid** | **string** | capture uuid | 

### Other Parameters

Other parameters are passed through a pointer to a apiKptureUuidDownloadGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## KptureUuidGet

> CaptureKpture KptureUuidGet(ctx, uuid).Execute()

Get kapture



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    uuid := "uuid_example" // string | capture uuid

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.KpturesApi.KptureUuidGet(context.Background(), uuid).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `KpturesApi.KptureUuidGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `KptureUuidGet`: CaptureKpture
    fmt.Fprintf(os.Stdout, "Response from `KpturesApi.KptureUuidGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**uuid** | **string** | capture uuid | 

### Other Parameters

Other parameters are passed through a pointer to a apiKptureUuidGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**CaptureKpture**](CaptureKpture.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## KptureUuidStopPut

> CaptureKpture KptureUuidStopPut(ctx, uuid).Execute()

Stop Kpture



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    uuid := "uuid_example" // string | capture uuid

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.KpturesApi.KptureUuidStopPut(context.Background(), uuid).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `KpturesApi.KptureUuidStopPut``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `KptureUuidStopPut`: CaptureKpture
    fmt.Fprintf(os.Stdout, "Response from `KpturesApi.KptureUuidStopPut`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**uuid** | **string** | capture uuid | 

### Other Parameters

Other parameters are passed through a pointer to a apiKptureUuidStopPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**CaptureKpture**](CaptureKpture.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## KpturesGet

> map[string]CaptureKpture KpturesGet(ctx).Execute()

Get kaptures



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.KpturesApi.KpturesGet(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `KpturesApi.KpturesGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `KpturesGet`: map[string]CaptureKpture
    fmt.Fprintf(os.Stdout, "Response from `KpturesApi.KpturesGet`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiKpturesGetRequest struct via the builder pattern


### Return type

[**map[string]CaptureKpture**](CaptureKpture.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

