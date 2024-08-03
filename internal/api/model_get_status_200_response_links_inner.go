/*
SpaceTraders API

SpaceTraders is an open-universe game and learning platform that offers a set of HTTP endpoints to control a fleet of ships and explore a multiplayer universe.  The API is documented using [OpenAPI](https://github.com/SpaceTradersAPI/api-docs). You can send your first request right here in your browser to check the status of the game server.  ```json http {   \"method\": \"GET\",   \"url\": \"https://api.spacetraders.io/v2\", } ```  Unlike a traditional game, SpaceTraders does not have a first-party client or app to play the game. Instead, you can use the API to build your own client, write a script to automate your ships, or try an app built by the community.  We have a [Discord channel](https://discord.com/invite/jh6zurdWk5) where you can share your projects, ask questions, and get help from other players.   

API version: 2.0.0
Contact: joel@spacetraders.io
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package api

import (
	"encoding/json"
	"bytes"
	"fmt"
)

// checks if the GetStatus200ResponseLinksInner type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &GetStatus200ResponseLinksInner{}

// GetStatus200ResponseLinksInner struct for GetStatus200ResponseLinksInner
type GetStatus200ResponseLinksInner struct {
	Name string `json:"name"`
	Url string `json:"url"`
}

type _GetStatus200ResponseLinksInner GetStatus200ResponseLinksInner

// NewGetStatus200ResponseLinksInner instantiates a new GetStatus200ResponseLinksInner object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewGetStatus200ResponseLinksInner(name string, url string) *GetStatus200ResponseLinksInner {
	this := GetStatus200ResponseLinksInner{}
	this.Name = name
	this.Url = url
	return &this
}

// NewGetStatus200ResponseLinksInnerWithDefaults instantiates a new GetStatus200ResponseLinksInner object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewGetStatus200ResponseLinksInnerWithDefaults() *GetStatus200ResponseLinksInner {
	this := GetStatus200ResponseLinksInner{}
	return &this
}

// GetName returns the Name field value
func (o *GetStatus200ResponseLinksInner) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *GetStatus200ResponseLinksInner) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *GetStatus200ResponseLinksInner) SetName(v string) {
	o.Name = v
}

// GetUrl returns the Url field value
func (o *GetStatus200ResponseLinksInner) GetUrl() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Url
}

// GetUrlOk returns a tuple with the Url field value
// and a boolean to check if the value has been set.
func (o *GetStatus200ResponseLinksInner) GetUrlOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Url, true
}

// SetUrl sets field value
func (o *GetStatus200ResponseLinksInner) SetUrl(v string) {
	o.Url = v
}

func (o GetStatus200ResponseLinksInner) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o GetStatus200ResponseLinksInner) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["name"] = o.Name
	toSerialize["url"] = o.Url
	return toSerialize, nil
}

func (o *GetStatus200ResponseLinksInner) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"name",
		"url",
	}

	allProperties := make(map[string]interface{})

	err = json.Unmarshal(data, &allProperties)

	if err != nil {
		return err;
	}

	for _, requiredProperty := range(requiredProperties) {
		if _, exists := allProperties[requiredProperty]; !exists {
			return fmt.Errorf("no value given for required property %v", requiredProperty)
		}
	}

	varGetStatus200ResponseLinksInner := _GetStatus200ResponseLinksInner{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varGetStatus200ResponseLinksInner)

	if err != nil {
		return err
	}

	*o = GetStatus200ResponseLinksInner(varGetStatus200ResponseLinksInner)

	return err
}

type NullableGetStatus200ResponseLinksInner struct {
	value *GetStatus200ResponseLinksInner
	isSet bool
}

func (v NullableGetStatus200ResponseLinksInner) Get() *GetStatus200ResponseLinksInner {
	return v.value
}

func (v *NullableGetStatus200ResponseLinksInner) Set(val *GetStatus200ResponseLinksInner) {
	v.value = val
	v.isSet = true
}

func (v NullableGetStatus200ResponseLinksInner) IsSet() bool {
	return v.isSet
}

func (v *NullableGetStatus200ResponseLinksInner) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableGetStatus200ResponseLinksInner(val *GetStatus200ResponseLinksInner) *NullableGetStatus200ResponseLinksInner {
	return &NullableGetStatus200ResponseLinksInner{value: val, isSet: true}
}

func (v NullableGetStatus200ResponseLinksInner) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableGetStatus200ResponseLinksInner) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


