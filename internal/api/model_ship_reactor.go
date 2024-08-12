/*
SpaceTraders API

SpaceTraders is an open-universe game and learning platform that offers a set of HTTP endpoints to control a fleet of ships and explore a multiplayer universe.  The API is documented using [OpenAPI](https://github.com/SpaceTradersAPI/api-docs). You can send your first request right here in your browser to check the status of the game server.  ```json http {   \"method\": \"GET\",   \"url\": \"https://api.spacetraders.io/v2\", } ```  Unlike a traditional game, SpaceTraders does not have a first-party client or app to play the game. Instead, you can use the API to build your own client, write a script to automate your ships, or try an app built by the community.  We have a [Discord channel](https://discord.com/invite/jh6zurdWk5) where you can share your projects, ask questions, and get help from other players.

API version: 2.0.0
Contact: joel@spacetraders.io
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package api

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// checks if the ShipReactor type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ShipReactor{}

// ShipReactor The reactor of the ship. The reactor is responsible for powering the ship's systems and weapons.
type ShipReactor struct {
	// Symbol of the reactor.
	Symbol string `json:"symbol"`
	// Name of the reactor.
	Name string `json:"name"`
	// Description of the reactor.
	Description string `json:"description"`
	// The repairable condition of a component. A value of 0 indicates the component needs significant repairs, while a value of 1 indicates the component is in near perfect condition. As the condition of a component is repaired, the overall integrity of the component decreases.
	Condition float64 `json:"condition"`
	Quality   float64 `json:"quality"`
	// The overall integrity of the component, which determines the performance of the component. A value of 0 indicates that the component is almost completely degraded, while a value of 1 indicates that the component is in near perfect condition. The integrity of the component is non-repairable, and represents permanent wear over time.
	Integrity float64 `json:"integrity"`
	// The amount of power provided by this reactor. The more power a reactor provides to the ship, the lower the cooldown it gets when using a module or mount that taxes the ship's power.
	PowerOutput  int32            `json:"powerOutput"`
	Requirements ShipRequirements `json:"requirements"`
}

type _ShipReactor ShipReactor

// NewShipReactor instantiates a new ShipReactor object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewShipReactor(symbol string, name string, description string, condition float64, integrity float64, powerOutput int32, requirements ShipRequirements) *ShipReactor {
	this := ShipReactor{}
	this.Symbol = symbol
	this.Name = name
	this.Description = description
	this.Condition = condition
	this.Integrity = integrity
	this.PowerOutput = powerOutput
	this.Requirements = requirements
	return &this
}

// NewShipReactorWithDefaults instantiates a new ShipReactor object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewShipReactorWithDefaults() *ShipReactor {
	this := ShipReactor{}
	return &this
}

// GetSymbol returns the Symbol field value
func (o *ShipReactor) GetSymbol() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Symbol
}

// GetSymbolOk returns a tuple with the Symbol field value
// and a boolean to check if the value has been set.
func (o *ShipReactor) GetSymbolOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Symbol, true
}

// SetSymbol sets field value
func (o *ShipReactor) SetSymbol(v string) {
	o.Symbol = v
}

// GetName returns the Name field value
func (o *ShipReactor) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *ShipReactor) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *ShipReactor) SetName(v string) {
	o.Name = v
}

// GetDescription returns the Description field value
func (o *ShipReactor) GetDescription() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Description
}

// GetDescriptionOk returns a tuple with the Description field value
// and a boolean to check if the value has been set.
func (o *ShipReactor) GetDescriptionOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Description, true
}

// SetDescription sets field value
func (o *ShipReactor) SetDescription(v string) {
	o.Description = v
}

// GetCondition returns the Condition field value
func (o *ShipReactor) GetCondition() float64 {
	if o == nil {
		var ret float64
		return ret
	}

	return o.Condition
}

// GetConditionOk returns a tuple with the Condition field value
// and a boolean to check if the value has been set.
func (o *ShipReactor) GetConditionOk() (*float64, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Condition, true
}

// SetCondition sets field value
func (o *ShipReactor) SetCondition(v float64) {
	o.Condition = v
}

// GetIntegrity returns the Integrity field value
func (o *ShipReactor) GetIntegrity() float64 {
	if o == nil {
		var ret float64
		return ret
	}

	return o.Integrity
}

// GetIntegrityOk returns a tuple with the Integrity field value
// and a boolean to check if the value has been set.
func (o *ShipReactor) GetIntegrityOk() (*float64, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Integrity, true
}

// SetIntegrity sets field value
func (o *ShipReactor) SetIntegrity(v float64) {
	o.Integrity = v
}

// GetPowerOutput returns the PowerOutput field value
func (o *ShipReactor) GetPowerOutput() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.PowerOutput
}

// GetPowerOutputOk returns a tuple with the PowerOutput field value
// and a boolean to check if the value has been set.
func (o *ShipReactor) GetPowerOutputOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.PowerOutput, true
}

// SetPowerOutput sets field value
func (o *ShipReactor) SetPowerOutput(v int32) {
	o.PowerOutput = v
}

// GetRequirements returns the Requirements field value
func (o *ShipReactor) GetRequirements() ShipRequirements {
	if o == nil {
		var ret ShipRequirements
		return ret
	}

	return o.Requirements
}

// GetRequirementsOk returns a tuple with the Requirements field value
// and a boolean to check if the value has been set.
func (o *ShipReactor) GetRequirementsOk() (*ShipRequirements, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Requirements, true
}

// SetRequirements sets field value
func (o *ShipReactor) SetRequirements(v ShipRequirements) {
	o.Requirements = v
}

func (o ShipReactor) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ShipReactor) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["symbol"] = o.Symbol
	toSerialize["name"] = o.Name
	toSerialize["description"] = o.Description
	toSerialize["condition"] = o.Condition
	toSerialize["integrity"] = o.Integrity
	toSerialize["powerOutput"] = o.PowerOutput
	toSerialize["requirements"] = o.Requirements
	return toSerialize, nil
}

func (o *ShipReactor) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"symbol",
		"name",
		"description",
		"condition",
		"integrity",
		"powerOutput",
		"requirements",
	}

	allProperties := make(map[string]interface{})

	err = json.Unmarshal(data, &allProperties)
	if err != nil {
		return err
	}

	for _, requiredProperty := range requiredProperties {
		if _, exists := allProperties[requiredProperty]; !exists {
			return fmt.Errorf("no value given for required property %v", requiredProperty)
		}
	}

	varShipReactor := _ShipReactor{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varShipReactor)
	if err != nil {
		return err
	}

	*o = ShipReactor(varShipReactor)

	return err
}

type NullableShipReactor struct {
	value *ShipReactor
	isSet bool
}

func (v NullableShipReactor) Get() *ShipReactor {
	return v.value
}

func (v *NullableShipReactor) Set(val *ShipReactor) {
	v.value = val
	v.isSet = true
}

func (v NullableShipReactor) IsSet() bool {
	return v.isSet
}

func (v *NullableShipReactor) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableShipReactor(val *ShipReactor) *NullableShipReactor {
	return &NullableShipReactor{value: val, isSet: true}
}

func (v NullableShipReactor) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableShipReactor) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
