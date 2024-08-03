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

// checks if the PurchaseCargo201Response type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &PurchaseCargo201Response{}

// PurchaseCargo201Response 
type PurchaseCargo201Response struct {
	Data SellCargo201ResponseData `json:"data"`
}

type _PurchaseCargo201Response PurchaseCargo201Response

// NewPurchaseCargo201Response instantiates a new PurchaseCargo201Response object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPurchaseCargo201Response(data SellCargo201ResponseData) *PurchaseCargo201Response {
	this := PurchaseCargo201Response{}
	this.Data = data
	return &this
}

// NewPurchaseCargo201ResponseWithDefaults instantiates a new PurchaseCargo201Response object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPurchaseCargo201ResponseWithDefaults() *PurchaseCargo201Response {
	this := PurchaseCargo201Response{}
	return &this
}

// GetData returns the Data field value
func (o *PurchaseCargo201Response) GetData() SellCargo201ResponseData {
	if o == nil {
		var ret SellCargo201ResponseData
		return ret
	}

	return o.Data
}

// GetDataOk returns a tuple with the Data field value
// and a boolean to check if the value has been set.
func (o *PurchaseCargo201Response) GetDataOk() (*SellCargo201ResponseData, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Data, true
}

// SetData sets field value
func (o *PurchaseCargo201Response) SetData(v SellCargo201ResponseData) {
	o.Data = v
}

func (o PurchaseCargo201Response) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o PurchaseCargo201Response) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["data"] = o.Data
	return toSerialize, nil
}

func (o *PurchaseCargo201Response) UnmarshalJSON(data []byte) (err error) {
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

	varPurchaseCargo201Response := _PurchaseCargo201Response{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varPurchaseCargo201Response)

	if err != nil {
		return err
	}

	*o = PurchaseCargo201Response(varPurchaseCargo201Response)

	return err
}

type NullablePurchaseCargo201Response struct {
	value *PurchaseCargo201Response
	isSet bool
}

func (v NullablePurchaseCargo201Response) Get() *PurchaseCargo201Response {
	return v.value
}

func (v *NullablePurchaseCargo201Response) Set(val *PurchaseCargo201Response) {
	v.value = val
	v.isSet = true
}

func (v NullablePurchaseCargo201Response) IsSet() bool {
	return v.isSet
}

func (v *NullablePurchaseCargo201Response) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePurchaseCargo201Response(val *PurchaseCargo201Response) *NullablePurchaseCargo201Response {
	return &NullablePurchaseCargo201Response{value: val, isSet: true}
}

func (v NullablePurchaseCargo201Response) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePurchaseCargo201Response) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


