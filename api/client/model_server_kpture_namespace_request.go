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

// ServerKptureNamespaceRequest struct for ServerKptureNamespaceRequest
type ServerKptureNamespaceRequest struct {
	KptureName *string `json:"kptureName,omitempty"`
	KptureNamespace *string `json:"kptureNamespace,omitempty"`
}

// NewServerKptureNamespaceRequest instantiates a new ServerKptureNamespaceRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewServerKptureNamespaceRequest() *ServerKptureNamespaceRequest {
	this := ServerKptureNamespaceRequest{}
	return &this
}

// NewServerKptureNamespaceRequestWithDefaults instantiates a new ServerKptureNamespaceRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewServerKptureNamespaceRequestWithDefaults() *ServerKptureNamespaceRequest {
	this := ServerKptureNamespaceRequest{}
	return &this
}

// GetKptureName returns the KptureName field value if set, zero value otherwise.
func (o *ServerKptureNamespaceRequest) GetKptureName() string {
	if o == nil || o.KptureName == nil {
		var ret string
		return ret
	}
	return *o.KptureName
}

// GetKptureNameOk returns a tuple with the KptureName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerKptureNamespaceRequest) GetKptureNameOk() (*string, bool) {
	if o == nil || o.KptureName == nil {
		return nil, false
	}
	return o.KptureName, true
}

// HasKptureName returns a boolean if a field has been set.
func (o *ServerKptureNamespaceRequest) HasKptureName() bool {
	if o != nil && o.KptureName != nil {
		return true
	}

	return false
}

// SetKptureName gets a reference to the given string and assigns it to the KptureName field.
func (o *ServerKptureNamespaceRequest) SetKptureName(v string) {
	o.KptureName = &v
}

// GetKptureNamespace returns the KptureNamespace field value if set, zero value otherwise.
func (o *ServerKptureNamespaceRequest) GetKptureNamespace() string {
	if o == nil || o.KptureNamespace == nil {
		var ret string
		return ret
	}
	return *o.KptureNamespace
}

// GetKptureNamespaceOk returns a tuple with the KptureNamespace field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerKptureNamespaceRequest) GetKptureNamespaceOk() (*string, bool) {
	if o == nil || o.KptureNamespace == nil {
		return nil, false
	}
	return o.KptureNamespace, true
}

// HasKptureNamespace returns a boolean if a field has been set.
func (o *ServerKptureNamespaceRequest) HasKptureNamespace() bool {
	if o != nil && o.KptureNamespace != nil {
		return true
	}

	return false
}

// SetKptureNamespace gets a reference to the given string and assigns it to the KptureNamespace field.
func (o *ServerKptureNamespaceRequest) SetKptureNamespace(v string) {
	o.KptureNamespace = &v
}

func (o ServerKptureNamespaceRequest) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.KptureName != nil {
		toSerialize["kptureName"] = o.KptureName
	}
	if o.KptureNamespace != nil {
		toSerialize["kptureNamespace"] = o.KptureNamespace
	}
	return json.Marshal(toSerialize)
}

type NullableServerKptureNamespaceRequest struct {
	value *ServerKptureNamespaceRequest
	isSet bool
}

func (v NullableServerKptureNamespaceRequest) Get() *ServerKptureNamespaceRequest {
	return v.value
}

func (v *NullableServerKptureNamespaceRequest) Set(val *ServerKptureNamespaceRequest) {
	v.value = val
	v.isSet = true
}

func (v NullableServerKptureNamespaceRequest) IsSet() bool {
	return v.isSet
}

func (v *NullableServerKptureNamespaceRequest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableServerKptureNamespaceRequest(val *ServerKptureNamespaceRequest) *NullableServerKptureNamespaceRequest {
	return &NullableServerKptureNamespaceRequest{value: val, isSet: true}
}

func (v NullableServerKptureNamespaceRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableServerKptureNamespaceRequest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


