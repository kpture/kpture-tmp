/*
Kpture-backend

Kpture Backend server

API version: 0.1
Contact: kpture.git@gmail.com
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package client

import (
	"encoding/json"
)

// AgentInfo struct for AgentInfo
type AgentInfo struct {
	Errors []string `json:"errors,omitempty"`
	Metadata *AgentMetadata `json:"metadata,omitempty"`
	PacketNb *int32 `json:"packetNb,omitempty"`
}

// NewAgentInfo instantiates a new AgentInfo object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAgentInfo() *AgentInfo {
	this := AgentInfo{}
	return &this
}

// NewAgentInfoWithDefaults instantiates a new AgentInfo object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAgentInfoWithDefaults() *AgentInfo {
	this := AgentInfo{}
	return &this
}

// GetErrors returns the Errors field value if set, zero value otherwise.
func (o *AgentInfo) GetErrors() []string {
	if o == nil || o.Errors == nil {
		var ret []string
		return ret
	}
	return o.Errors
}

// GetErrorsOk returns a tuple with the Errors field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AgentInfo) GetErrorsOk() ([]string, bool) {
	if o == nil || o.Errors == nil {
		return nil, false
	}
	return o.Errors, true
}

// HasErrors returns a boolean if a field has been set.
func (o *AgentInfo) HasErrors() bool {
	if o != nil && o.Errors != nil {
		return true
	}

	return false
}

// SetErrors gets a reference to the given []string and assigns it to the Errors field.
func (o *AgentInfo) SetErrors(v []string) {
	o.Errors = v
}

// GetMetadata returns the Metadata field value if set, zero value otherwise.
func (o *AgentInfo) GetMetadata() AgentMetadata {
	if o == nil || o.Metadata == nil {
		var ret AgentMetadata
		return ret
	}
	return *o.Metadata
}

// GetMetadataOk returns a tuple with the Metadata field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AgentInfo) GetMetadataOk() (*AgentMetadata, bool) {
	if o == nil || o.Metadata == nil {
		return nil, false
	}
	return o.Metadata, true
}

// HasMetadata returns a boolean if a field has been set.
func (o *AgentInfo) HasMetadata() bool {
	if o != nil && o.Metadata != nil {
		return true
	}

	return false
}

// SetMetadata gets a reference to the given AgentMetadata and assigns it to the Metadata field.
func (o *AgentInfo) SetMetadata(v AgentMetadata) {
	o.Metadata = &v
}

// GetPacketNb returns the PacketNb field value if set, zero value otherwise.
func (o *AgentInfo) GetPacketNb() int32 {
	if o == nil || o.PacketNb == nil {
		var ret int32
		return ret
	}
	return *o.PacketNb
}

// GetPacketNbOk returns a tuple with the PacketNb field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AgentInfo) GetPacketNbOk() (*int32, bool) {
	if o == nil || o.PacketNb == nil {
		return nil, false
	}
	return o.PacketNb, true
}

// HasPacketNb returns a boolean if a field has been set.
func (o *AgentInfo) HasPacketNb() bool {
	if o != nil && o.PacketNb != nil {
		return true
	}

	return false
}

// SetPacketNb gets a reference to the given int32 and assigns it to the PacketNb field.
func (o *AgentInfo) SetPacketNb(v int32) {
	o.PacketNb = &v
}

func (o AgentInfo) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Errors != nil {
		toSerialize["errors"] = o.Errors
	}
	if o.Metadata != nil {
		toSerialize["metadata"] = o.Metadata
	}
	if o.PacketNb != nil {
		toSerialize["packetNb"] = o.PacketNb
	}
	return json.Marshal(toSerialize)
}

type NullableAgentInfo struct {
	value *AgentInfo
	isSet bool
}

func (v NullableAgentInfo) Get() *AgentInfo {
	return v.value
}

func (v *NullableAgentInfo) Set(val *AgentInfo) {
	v.value = val
	v.isSet = true
}

func (v NullableAgentInfo) IsSet() bool {
	return v.isSet
}

func (v *NullableAgentInfo) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableAgentInfo(val *AgentInfo) *NullableAgentInfo {
	return &NullableAgentInfo{value: val, isSet: true}
}

func (v NullableAgentInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableAgentInfo) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

