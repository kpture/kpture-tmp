# \AgentsApi

All URIs are relative to *http://localhost:8080/kpture/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AgentsGet**](AgentsApi.md#AgentsGet) | **Get** /agents | Get Agents



## AgentsGet

> []AgentMetadata AgentsGet(ctx).Execute()

Get Agents



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
    resp, r, err := apiClient.AgentsApi.AgentsGet(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `AgentsApi.AgentsGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `AgentsGet`: []AgentMetadata
    fmt.Fprintf(os.Stdout, "Response from `AgentsApi.AgentsGet`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiAgentsGetRequest struct via the builder pattern


### Return type

[**[]AgentMetadata**](AgentMetadata.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

