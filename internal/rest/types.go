package rest

import (
	"time"

	"github.com/julienlavocat/spacetraders/internal/api"
)

type ShipSnapshot struct {
	ArrivalAt    time.Time         `json:"arrivalAt"`
	Cargo        map[string]int32  `json:"cargo"`
	Departure    string            `json:"departure"`
	Id           string            `json:"id"`
	Status       api.ShipNavStatus `json:"status"`
	SystemId     string            `json:"system"`
	WaypointId   string            `json:"waypoint"`
	Destination  string            `json:"destination"`
	MaxFuel      int32             `json:"maxFuel"`
	CurrentCargo int32             `json:"currentCargo"`
	MaxCargo     int32             `json:"maxCargo"`
	CurrentFuel  int32             `json:"currentFuel"`
	IsCargoFull  bool              `json:"isCargoFull"`
	DepartedAt   time.Time         `json:"departedAt"`
}
