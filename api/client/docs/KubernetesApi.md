# \KubernetesApi

All URIs are relative to *http://localhost:8080/kpture/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**K8sNamespacesGet**](KubernetesApi.md#K8sNamespacesGet) | **Get** /k8s/namespaces | Get all kubernetes namespaces
[**K8sNamespacesNamespaceInjectPost**](KubernetesApi.md#K8sNamespacesNamespaceInjectPost) | **Post** /k8s/namespaces/{namespace}/inject | Inject annotation webhook
[**KptureK8sNamespacePost**](KubernetesApi.md#KptureK8sNamespacePost) | **Post** /kpture/k8s/namespace | Start namespace kpture
[**KptureK8sNamespacesGet**](KubernetesApi.md#KptureK8sNamespacesGet) | **Get** /kpture/k8s/namespaces | Get enabled kubernetes namespaces



## K8sNamespacesGet

> []string K8sNamespacesGet(ctx).Execute()

Get all kubernetes namespaces



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
    resp, r, err := apiClient.KubernetesApi.K8sNamespacesGet(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `KubernetesApi.K8sNamespacesGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `K8sNamespacesGet`: []string
    fmt.Fprintf(os.Stdout, "Response from `KubernetesApi.K8sNamespacesGet`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiK8sNamespacesGetRequest struct via the builder pattern


### Return type

**[]string**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## K8sNamespacesNamespaceInjectPost

> K8sNamespacesNamespaceInjectPost(ctx, namespace).Execute()

Inject annotation webhook



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
    namespace := "namespace_example" // string | namespace

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.KubernetesApi.K8sNamespacesNamespaceInjectPost(context.Background(), namespace).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `KubernetesApi.K8sNamespacesNamespaceInjectPost``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**namespace** | **string** | namespace | 

### Other Parameters

Other parameters are passed through a pointer to a apiK8sNamespacesNamespaceInjectPostRequest struct via the builder pattern


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


## KptureK8sNamespacePost

> CaptureKpture KptureK8sNamespacePost(ctx).Data(data).Execute()

Start namespace kpture



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
    data := *openapiclient.NewServerKptureNamespaceRequest() // ServerKptureNamespaceRequest | namespace for capture

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.KubernetesApi.KptureK8sNamespacePost(context.Background()).Data(data).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `KubernetesApi.KptureK8sNamespacePost``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `KptureK8sNamespacePost`: CaptureKpture
    fmt.Fprintf(os.Stdout, "Response from `KubernetesApi.KptureK8sNamespacePost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiKptureK8sNamespacePostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **data** | [**ServerKptureNamespaceRequest**](ServerKptureNamespaceRequest.md) | namespace for capture | 

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


## KptureK8sNamespacesGet

> []string KptureK8sNamespacesGet(ctx).Execute()

Get enabled kubernetes namespaces



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
    resp, r, err := apiClient.KubernetesApi.KptureK8sNamespacesGet(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `KubernetesApi.KptureK8sNamespacesGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `KptureK8sNamespacesGet`: []string
    fmt.Fprintf(os.Stdout, "Response from `KubernetesApi.KptureK8sNamespacesGet`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiKptureK8sNamespacesGetRequest struct via the builder pattern


### Return type

**[]string**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

