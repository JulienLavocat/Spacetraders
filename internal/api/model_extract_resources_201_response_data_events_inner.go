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
	"gopkg.in/validator.v2"
	"fmt"
)

// ExtractResources201ResponseDataEventsInner - struct for ExtractResources201ResponseDataEventsInner
type ExtractResources201ResponseDataEventsInner struct {
	ShipConditionEvent *ShipConditionEvent
}

// ShipConditionEventAsExtractResources201ResponseDataEventsInner is a convenience function that returns ShipConditionEvent wrapped in ExtractResources201ResponseDataEventsInner
func ShipConditionEventAsExtractResources201ResponseDataEventsInner(v *ShipConditionEvent) ExtractResources201ResponseDataEventsInner {
	return ExtractResources201ResponseDataEventsInner{
		ShipConditionEvent: v,
	}
}


// Unmarshal JSON data into one of the pointers in the struct
func (dst *ExtractResources201ResponseDataEventsInner) UnmarshalJSON(data []byte) error {
	var err error
	match := 0
	// try to unmarshal data into ShipConditionEvent
	err = newStrictDecoder(data).Decode(&dst.ShipConditionEvent)
	if err == nil {
		jsonShipConditionEvent, _ := json.Marshal(dst.ShipConditionEvent)
		if string(jsonShipConditionEvent) == "{}" { // empty struct
			dst.ShipConditionEvent = nil
		} else {
			if err = validator.Validate(dst.ShipConditionEvent); err != nil {
				dst.ShipConditionEvent = nil
			} else {
				match++
			}
		}
	} else {
		dst.ShipConditionEvent = nil
	}

	if match > 1 { // more than 1 match
		// reset to nil
		dst.ShipConditionEvent = nil

		return fmt.Errorf("data matches more than one schema in oneOf(ExtractResources201ResponseDataEventsInner)")
	} else if match == 1 {
		return nil // exactly one match
	} else { // no match
		return fmt.Errorf("data failed to match schemas in oneOf(ExtractResources201ResponseDataEventsInner)")
	}
}

// Marshal data from the first non-nil pointers in the struct to JSON
func (src ExtractResources201ResponseDataEventsInner) MarshalJSON() ([]byte, error) {
	if src.ShipConditionEvent != nil {
		return json.Marshal(&src.ShipConditionEvent)
	}

	return nil, nil // no data in oneOf schemas
}

// Get the actual instance
func (obj *ExtractResources201ResponseDataEventsInner) GetActualInstance() (interface{}) {
	if obj == nil {
		return nil
	}
	if obj.ShipConditionEvent != nil {
		return obj.ShipConditionEvent
	}

	// all schemas are nil
	return nil
}

type NullableExtractResources201ResponseDataEventsInner struct {
	value *ExtractResources201ResponseDataEventsInner
	isSet bool
}

func (v NullableExtractResources201ResponseDataEventsInner) Get() *ExtractResources201ResponseDataEventsInner {
	return v.value
}

func (v *NullableExtractResources201ResponseDataEventsInner) Set(val *ExtractResources201ResponseDataEventsInner) {
	v.value = val
	v.isSet = true
}

func (v NullableExtractResources201ResponseDataEventsInner) IsSet() bool {
	return v.isSet
}

func (v *NullableExtractResources201ResponseDataEventsInner) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableExtractResources201ResponseDataEventsInner(val *ExtractResources201ResponseDataEventsInner) *NullableExtractResources201ResponseDataEventsInner {
	return &NullableExtractResources201ResponseDataEventsInner{value: val, isSet: true}
}

func (v NullableExtractResources201ResponseDataEventsInner) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableExtractResources201ResponseDataEventsInner) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


