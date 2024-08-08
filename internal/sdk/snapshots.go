package sdk

import (
	"time"

	"github.com/julienlavocat/spacetraders/internal/api"
)

type ShipSnapshot struct {
	ArrivalAt    time.Time         `json:"arrivalAt"`
	DepartedAt   time.Time         `json:"departedAt"`
	Cargo        map[string]int32  `json:"cargo"`
	WaypointId   string            `json:"waypoint"`
	Status       api.ShipNavStatus `json:"status"`
	SystemId     string            `json:"system"`
	Id           string            `json:"id"`
	Destination  string            `json:"destination"`
	Departure    string            `json:"departure"`
	MaxFuel      int32             `json:"maxFuel"`
	CurrentCargo int32             `json:"currentCargo"`
	MaxCargo     int32             `json:"maxCargo"`
	CurrentFuel  int32             `json:"currentFuel"`
	IsCargoFull  bool              `json:"isCargoFull"`
}

func newShipSnapshot(ship *Ship) ShipSnapshot {
	return ShipSnapshot{
		SystemId:     ship.Nav.SystemSymbol,
		WaypointId:   ship.Nav.WaypointSymbol,
		Id:           ship.Id,
		Cargo:        ship.Cargo,
		Status:       ship.Nav.Status,
		Destination:  ship.Nav.Route.Destination.Symbol,
		MaxFuel:      ship.Fuel.Capacity,
		CurrentFuel:  ship.Fuel.Current,
		MaxCargo:     ship.MaxCargo,
		CurrentCargo: ship.CurrentCargo,
		IsCargoFull:  ship.IsCargoFull,
		Departure:    ship.Nav.Route.Origin.Symbol,
		ArrivalAt:    ship.Nav.Route.Arrival,
		DepartedAt:   ship.Nav.Route.DepartureTime,
	}
}
