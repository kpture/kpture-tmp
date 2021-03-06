# Go API client for client

Kpture Backend server

## Overview
This API client was generated by the [OpenAPI Generator](https://openapi-generator.tech) project.  By using the [OpenAPI-spec](https://www.openapis.org/) from a remote server, you can easily generate an API client.

- API version: 0.1
- Package version: 1.0.0
- Build package: org.openapitools.codegen.languages.GoClientCodegen
For more information, please visit [http://www.kpture.io](http://www.kpture.io)

## Installation

Install the following dependencies:

```shell
go get github.com/stretchr/testify/assert
go get golang.org/x/oauth2
go get golang.org/x/net/context
```

Put the package under your project folder and add the following in import:

```golang
import client "github.com/GIT_USER_ID/GIT_REPO_ID"
```

To use a proxy, set the environment variable `HTTP_PROXY`:

```golang
os.Setenv("HTTP_PROXY", "http://proxy_name:proxy_port")
```

## Configuration of Server URL

Default configuration comes with `Servers` field that contains server objects as defined in the OpenAPI specification.

### Select Server Configuration

For using other server than the one defined on index 0 set context value `sw.ContextServerIndex` of type `int`.

```golang
ctx := context.WithValue(context.Background(), client.ContextServerIndex, 1)
```

### Templated Server URL

Templated server URL is formatted using default variables from configuration or from context value `sw.ContextServerVariables` of type `map[string]string`.

```golang
ctx := context.WithValue(context.Background(), client.ContextServerVariables, map[string]string{
	"basePath": "v2",
})
```

Note, enum values are always validated and all unused variables are silently ignored.

### URLs Configuration per Operation

Each operation can use different server URL defined using `OperationServers` map in the `Configuration`.
An operation is uniquely identified by `"{classname}Service.{nickname}"` string.
Similar rules for overriding default operation server index and variables applies by using `sw.ContextOperationServerIndices` and `sw.ContextOperationServerVariables` context maps.

```
ctx := context.WithValue(context.Background(), client.ContextOperationServerIndices, map[string]int{
	"{classname}Service.{nickname}": 2,
})
ctx = context.WithValue(context.Background(), client.ContextOperationServerVariables, map[string]map[string]string{
	"{classname}Service.{nickname}": {
		"port": "8443",
	},
})
```

## Documentation for API Endpoints

All URIs are relative to *http://localhost:8080/kpture/api/v1*

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*AgentsApi* | [**AgentsGet**](docs/AgentsApi.md#agentsget) | **Get** /agents | Get Agents
*KpturesApi* | [**KpturePost**](docs/KpturesApi.md#kpturepost) | **Post** /kpture | Start Kpture
*KpturesApi* | [**KptureUuidDelete**](docs/KpturesApi.md#kptureuuiddelete) | **Delete** /kpture/{uuid} | Delete kapture
*KpturesApi* | [**KptureUuidDownloadGet**](docs/KpturesApi.md#kptureuuiddownloadget) | **Get** /kpture/{uuid}/download | Download kpture
*KpturesApi* | [**KptureUuidGet**](docs/KpturesApi.md#kptureuuidget) | **Get** /kpture/{uuid} | Get kapture
*KpturesApi* | [**KptureUuidStopPut**](docs/KpturesApi.md#kptureuuidstopput) | **Put** /kpture/{uuid}/stop | Stop Kpture
*KpturesApi* | [**KpturesGet**](docs/KpturesApi.md#kpturesget) | **Get** /kptures | Get kaptures
*KubernetesApi* | [**K8sNamespacesGet**](docs/KubernetesApi.md#k8snamespacesget) | **Get** /k8s/namespaces | Get all kubernetes namespaces
*KubernetesApi* | [**K8sNamespacesNamespaceInjectPost**](docs/KubernetesApi.md#k8snamespacesnamespaceinjectpost) | **Post** /k8s/namespaces/{namespace}/inject | Inject annotation webhook
*KubernetesApi* | [**KptureK8sNamespacePost**](docs/KubernetesApi.md#kpturek8snamespacepost) | **Post** /kpture/k8s/namespace | Start namespace kpture
*KubernetesApi* | [**KptureK8sNamespacesGet**](docs/KubernetesApi.md#kpturek8snamespacesget) | **Get** /kpture/k8s/namespaces | Get enabled kubernetes namespaces
*ProfilesApi* | [**ProfileProfileNameDelete**](docs/ProfilesApi.md#profileprofilenamedelete) | **Delete** /profile/{profileName} | Delete profile
*ProfilesApi* | [**ProfileProfileNamePost**](docs/ProfilesApi.md#profileprofilenamepost) | **Post** /profile/{profileName} | Create profile
*ProfilesApi* | [**ProfilesGet**](docs/ProfilesApi.md#profilesget) | **Get** /profiles | Get profiles
*WiresharkApi* | [**WiresharkHostfileGet**](docs/WiresharkApi.md#wiresharkhostfileget) | **Get** /wireshark/hostfile | Get hostfile


## Documentation For Models

 - [AgentInfo](docs/AgentInfo.md)
 - [AgentMetadata](docs/AgentMetadata.md)
 - [CaptureInfo](docs/CaptureInfo.md)
 - [CaptureKpture](docs/CaptureKpture.md)
 - [ServerKptureNamespaceRequest](docs/ServerKptureNamespaceRequest.md)
 - [ServerKptureRequest](docs/ServerKptureRequest.md)
 - [ServerKptureRequestAgentsRequestInner](docs/ServerKptureRequestAgentsRequestInner.md)
 - [ServerServerError](docs/ServerServerError.md)


## Documentation For Authorization



### BasicAuth

- **Type**: HTTP basic authentication

Example

```golang
auth := context.WithValue(context.Background(), sw.ContextBasicAuth, sw.BasicAuth{
    UserName: "username",
    Password: "password",
})
r, err := client.Service.Operation(auth, args)
```


## Documentation for Utility Methods

Due to the fact that model structure members are all pointers, this package contains
a number of utility functions to easily obtain pointers to values of basic types.
Each of these functions takes a value of the given basic type and returns a pointer to it:

* `PtrBool`
* `PtrInt`
* `PtrInt32`
* `PtrInt64`
* `PtrFloat`
* `PtrFloat32`
* `PtrFloat64`
* `PtrString`
* `PtrTime`

## Author

kpture.git@gmail.com

