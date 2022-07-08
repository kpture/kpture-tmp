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

// CaptureInfo struct for CaptureInfo
type CaptureInfo struct {
	PacketNb *int32 `json:"packetNb,omitempty"`
	Size *int32 `json:"size,omitempty"`
}

// NewCaptureInfo instantiates a new CaptureInfo object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCaptureInfo() *CaptureInfo {
	this := CaptureInfo{}
	return &this
}

// NewCaptureInfoWithDefaults instantiates a new CaptureInfo object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCaptureInfoWithDefaults() *CaptureInfo {
	this := CaptureInfo{}
	return &this
}

// GetPacketNb returns the PacketNb field value if set, zero value otherwise.
func (o *CaptureInfo) GetPacketNb() int32 {
	if o == nil || o.PacketNb == nil {
		var ret int32
		return ret
	}
	return *o.PacketNb
}

// GetPacketNbOk returns a tuple with the PacketNb field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CaptureInfo) GetPacketNbOk() (*int32, bool) {
	if o == nil || o.PacketNb == nil {
		return nil, false
	}
	return o.PacketNb, true
}

// HasPacketNb returns a boolean if a field has been set.
func (o *CaptureInfo) HasPacketNb() bool {
	if o != nil && o.PacketNb != nil {
		return true
	}

	return false
}

// SetPacketNb gets a reference to the given int32 and assigns it to the PacketNb field.
func (o *CaptureInfo) SetPacketNb(v int32) {
	o.PacketNb = &v
}

// GetSize returns the Size field value if set, zero value otherwise.
func (o *CaptureInfo) GetSize() int32 {
	if o == nil || o.Size == nil {
		var ret int32
		return ret
	}
	return *o.Size
}

// GetSizeOk returns a tuple with the Size field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CaptureInfo) GetSizeOk() (*int32, bool) {
	if o == nil || o.Size == nil {
		return nil, false
	}
	return o.Size, true
}

// HasSize returns a boolean if a field has been set.
func (o *CaptureInfo) HasSize() bool {
	if o != nil && o.Size != nil {
		return true
	}

	return false
}

// SetSize gets a reference to the given int32 and assigns it to the Size field.
func (o *CaptureInfo) SetSize(v int32) {
	o.Size = &v
}

func (o CaptureInfo) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.PacketNb != nil {
		toSerialize["packetNb"] = o.PacketNb
	}
	if o.Size != nil {
		toSerialize["size"] = o.Size
	}
	return json.Marshal(toSerialize)
}

type NullableCaptureInfo struct {
	value *CaptureInfo
	isSet bool
}

func (v NullableCaptureInfo) Get() *CaptureInfo {
	return v.value
}

func (v *NullableCaptureInfo) Set(val *CaptureInfo) {
	v.value = val
	v.isSet = true
}

func (v NullableCaptureInfo) IsSet() bool {
	return v.isSet
}

func (v *NullableCaptureInfo) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCaptureInfo(val *CaptureInfo) *NullableCaptureInfo {
	return &NullableCaptureInfo{value: val, isSet: true}
}

func (v NullableCaptureInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCaptureInfo) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

