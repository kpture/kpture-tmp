# \WiresharkApi

All URIs are relative to *http://localhost:8080/kpture/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**WiresharkHostfileGet**](WiresharkApi.md#WiresharkHostfileGet) | **Get** /wireshark/hostfile | Get hostfile



## WiresharkHostfileGet

> string WiresharkHostfileGet(ctx).Execute()

Get hostfile



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
    resp, r, err := apiClient.WiresharkApi.WiresharkHostfileGet(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `WiresharkApi.WiresharkHostfileGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `WiresharkHostfileGet`: string
    fmt.Fprintf(os.Stdout, "Response from `WiresharkApi.WiresharkHostfileGet`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiWiresharkHostfileGetRequest struct via the builder pattern


### Return type

**string**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

