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

// checks if the GetScrapShip200ResponseData type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &GetScrapShip200ResponseData{}

// GetScrapShip200ResponseData struct for GetScrapShip200ResponseData
type GetScrapShip200ResponseData struct {
	Transaction ScrapTransaction `json:"transaction"`
}

type _GetScrapShip200ResponseData GetScrapShip200ResponseData

// NewGetScrapShip200ResponseData instantiates a new GetScrapShip200ResponseData object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewGetScrapShip200ResponseData(transaction ScrapTransaction) *GetScrapShip200ResponseData {
	this := GetScrapShip200ResponseData{}
	this.Transaction = transaction
	return &this
}

// NewGetScrapShip200ResponseDataWithDefaults instantiates a new GetScrapShip200ResponseData object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewGetScrapShip200ResponseDataWithDefaults() *GetScrapShip200ResponseData {
	this := GetScrapShip200ResponseData{}
	return &this
}

// GetTransaction returns the Transaction field value
func (o *GetScrapShip200ResponseData) GetTransaction() ScrapTransaction {
	if o == nil {
		var ret ScrapTransaction
		return ret
	}

	return o.Transaction
}

// GetTransactionOk returns a tuple with the Transaction field value
// and a boolean to check if the value has been set.
func (o *GetScrapShip200ResponseData) GetTransactionOk() (*ScrapTransaction, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Transaction, true
}

// SetTransaction sets field value
func (o *GetScrapShip200ResponseData) SetTransaction(v ScrapTransaction) {
	o.Transaction = v
}

func (o GetScrapShip200ResponseData) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o GetScrapShip200ResponseData) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["transaction"] = o.Transaction
	return toSerialize, nil
}

func (o *GetScrapShip200ResponseData) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"transaction",
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

	varGetScrapShip200ResponseData := _GetScrapShip200ResponseData{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varGetScrapShip200ResponseData)

	if err != nil {
		return err
	}

	*o = GetScrapShip200ResponseData(varGetScrapShip200ResponseData)

	return err
}

type NullableGetScrapShip200ResponseData struct {
	value *GetScrapShip200ResponseData
	isSet bool
}

func (v NullableGetScrapShip200ResponseData) Get() *GetScrapShip200ResponseData {
	return v.value
}

func (v *NullableGetScrapShip200ResponseData) Set(val *GetScrapShip200ResponseData) {
	v.value = val
	v.isSet = true
}

func (v NullableGetScrapShip200ResponseData) IsSet() bool {
	return v.isSet
}

func (v *NullableGetScrapShip200ResponseData) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableGetScrapShip200ResponseData(val *GetScrapShip200ResponseData) *NullableGetScrapShip200ResponseData {
	return &NullableGetScrapShip200ResponseData{value: val, isSet: true}
}

func (v NullableGetScrapShip200ResponseData) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableGetScrapShip200ResponseData) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


