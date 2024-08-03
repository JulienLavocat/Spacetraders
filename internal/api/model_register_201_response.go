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

// checks if the Register201Response type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Register201Response{}

// Register201Response struct for Register201Response
type Register201Response struct {
	Data Register201ResponseData `json:"data"`
}

type _Register201Response Register201Response

// NewRegister201Response instantiates a new Register201Response object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewRegister201Response(data Register201ResponseData) *Register201Response {
	this := Register201Response{}
	this.Data = data
	return &this
}

// NewRegister201ResponseWithDefaults instantiates a new Register201Response object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewRegister201ResponseWithDefaults() *Register201Response {
	this := Register201Response{}
	return &this
}

// GetData returns the Data field value
func (o *Register201Response) GetData() Register201ResponseData {
	if o == nil {
		var ret Register201ResponseData
		return ret
	}

	return o.Data
}

// GetDataOk returns a tuple with the Data field value
// and a boolean to check if the value has been set.
func (o *Register201Response) GetDataOk() (*Register201ResponseData, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Data, true
}

// SetData sets field value
func (o *Register201Response) SetData(v Register201ResponseData) {
	o.Data = v
}

func (o Register201Response) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Register201Response) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["data"] = o.Data
	return toSerialize, nil
}

func (o *Register201Response) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"data",
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

	varRegister201Response := _Register201Response{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varRegister201Response)

	if err != nil {
		return err
	}

	*o = Register201Response(varRegister201Response)

	return err
}

type NullableRegister201Response struct {
	value *Register201Response
	isSet bool
}

func (v NullableRegister201Response) Get() *Register201Response {
	return v.value
}

func (v *NullableRegister201Response) Set(val *Register201Response) {
	v.value = val
	v.isSet = true
}

func (v NullableRegister201Response) IsSet() bool {
	return v.isSet
}

func (v *NullableRegister201Response) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableRegister201Response(val *Register201Response) *NullableRegister201Response {
	return &NullableRegister201Response{value: val, isSet: true}
}

func (v NullableRegister201Response) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableRegister201Response) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


