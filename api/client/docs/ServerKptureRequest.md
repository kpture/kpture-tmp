# ServerKptureRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AgentsRequest** | Pointer to [**[]ServerKptureRequestAgentsRequestInner**](ServerKptureRequestAgentsRequestInner.md) |  | [optional] 
**KptureName** | Pointer to **string** |  | [optional] 

## Methods

### NewServerKptureRequest

`func NewServerKptureRequest() *ServerKptureRequest`

NewServerKptureRequest instantiates a new ServerKptureRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewServerKptureRequestWithDefaults

`func NewServerKptureRequestWithDefaults() *ServerKptureRequest`

NewServerKptureRequestWithDefaults instantiates a new ServerKptureRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAgentsRequest

`func (o *ServerKptureRequest) GetAgentsRequest() []ServerKptureRequestAgentsRequestInner`

GetAgentsRequest returns the AgentsRequest field if non-nil, zero value otherwise.

### GetAgentsRequestOk

`func (o *ServerKptureRequest) GetAgentsRequestOk() (*[]ServerKptureRequestAgentsRequestInner, bool)`

GetAgentsRequestOk returns a tuple with the AgentsRequest field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAgentsRequest

`func (o *ServerKptureRequest) SetAgentsRequest(v []ServerKptureRequestAgentsRequestInner)`

SetAgentsRequest sets AgentsRequest field to given value.

### HasAgentsRequest

`func (o *ServerKptureRequest) HasAgentsRequest() bool`

HasAgentsRequest returns a boolean if a field has been set.

### GetKptureName

`func (o *ServerKptureRequest) GetKptureName() string`

GetKptureName returns the KptureName field if non-nil, zero value otherwise.

### GetKptureNameOk

`func (o *ServerKptureRequest) GetKptureNameOk() (*string, bool)`

GetKptureNameOk returns a tuple with the KptureName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKptureName

`func (o *ServerKptureRequest) SetKptureName(v string)`

SetKptureName sets KptureName field to given value.

### HasKptureName

`func (o *ServerKptureRequest) HasKptureName() bool`

HasKptureName returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


