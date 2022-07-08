# AgentInfo

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Errors** | Pointer to **[]string** |  | [optional] 
**Metadata** | Pointer to [**AgentMetadata**](AgentMetadata.md) |  | [optional] 
**PacketNb** | Pointer to **int32** |  | [optional] 

## Methods

### NewAgentInfo

`func NewAgentInfo() *AgentInfo`

NewAgentInfo instantiates a new AgentInfo object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAgentInfoWithDefaults

`func NewAgentInfoWithDefaults() *AgentInfo`

NewAgentInfoWithDefaults instantiates a new AgentInfo object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetErrors

`func (o *AgentInfo) GetErrors() []string`

GetErrors returns the Errors field if non-nil, zero value otherwise.

### GetErrorsOk

`func (o *AgentInfo) GetErrorsOk() (*[]string, bool)`

GetErrorsOk returns a tuple with the Errors field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetErrors

`func (o *AgentInfo) SetErrors(v []string)`

SetErrors sets Errors field to given value.

### HasErrors

`func (o *AgentInfo) HasErrors() bool`

HasErrors returns a boolean if a field has been set.

### GetMetadata

`func (o *AgentInfo) GetMetadata() AgentMetadata`

GetMetadata returns the Metadata field if non-nil, zero value otherwise.

### GetMetadataOk

`func (o *AgentInfo) GetMetadataOk() (*AgentMetadata, bool)`

GetMetadataOk returns a tuple with the Metadata field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMetadata

`func (o *AgentInfo) SetMetadata(v AgentMetadata)`

SetMetadata sets Metadata field to given value.

### HasMetadata

`func (o *AgentInfo) HasMetadata() bool`

HasMetadata returns a boolean if a field has been set.

### GetPacketNb

`func (o *AgentInfo) GetPacketNb() int32`

GetPacketNb returns the PacketNb field if non-nil, zero value otherwise.

### GetPacketNbOk

`func (o *AgentInfo) GetPacketNbOk() (*int32, bool)`

GetPacketNbOk returns a tuple with the PacketNb field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPacketNb

`func (o *AgentInfo) SetPacketNb(v int32)`

SetPacketNb sets PacketNb field to given value.

### HasPacketNb

`func (o *AgentInfo) HasPacketNb() bool`

HasPacketNb returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


