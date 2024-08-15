package adapters

import (
	"encoding/json"
	"time"

	"github.com/julienlavocat/spacetraders/.gen/spacetraders/public/model"
)

type Route struct {
	To      string `json:"to"`
	Fuel    int32  `json:"fuel"`
	AggCost int32  `json:"aggFuel"`
}

type TradeRoute struct {
	Product                 string `json:"product"`
	BuyAt                   string `json:"buyAt"`
	SellAt                  string `json:"sellAt"`
	MaxAmount               int32  `json:"maxAmount"`
	BuyPrice                int32  `json:"buyPrice"`
	SellPrice               int32  `json:"sellPrice"`
	EstimatedProfits        int32  `json:"estimatedProfits"`
	EstimatedProfitsPerUnit int32  `json:"estimatedProfitsPerUnit"`
	FuelCost                int32  `json:"fuelCost"`
}

type Ship struct {
	ArrivalAt    time.Time        `json:"arrivalAt"`
	DepartedAt   time.Time        `json:"departedAt"`
	UpdatedAt    time.Time        `json:"updatedAt"`
	Cooldown     *time.Time       `json:"cooldown"`
	Cargo        map[string]int32 `json:"cargo"`
	TradeRoute   *TradeRoute      `json:"tradeRoute"`
	Id           string           `json:"id"`
	Origin       string           `json:"origin"`
	Status       string           `json:"status"`
	System       string           `json:"system"`
	Waypoint     string           `json:"waypoint"`
	Destination  string           `json:"destination"`
	Type         string           `json:"type"`
	Route        []*Route         `json:"route"`
	MaxFuel      int32            `json:"maxFuel"`
	CurrentFuel  int32            `json:"currentFuel"`
	CurrentCargo int32            `json:"currentCargo"`
	MaxCargo     int32            `json:"maxCargo"`
	IsCargoFull  bool             `json:"isCargoFull"`
}

func AdaptShip(s model.Ships) (*Ship, error) {
	var cargo map[string]int32
	err := json.Unmarshal([]byte(s.Cargo), &cargo)
	if err != nil {
		return nil, err
	}

	var route []*Route
	if s.Route != nil {
		err := json.Unmarshal([]byte(*s.Route), &route)
		if err != nil {
			return nil, err
		}
	}

	var tradeRoute *TradeRoute
	if s.TradeRoute != nil {
		err := json.Unmarshal([]byte(*s.TradeRoute), &tradeRoute)
		if err != nil {
			return nil, err
		}
	}

	return &Ship{
		ArrivalAt:    s.ArrivalAt,
		DepartedAt:   s.DepartedAt,
		Cooldown:     s.Cooldown,
		UpdatedAt:    s.UpdatedAt,
		Cargo:        cargo,
		Destination:  s.Destination,
		Id:           s.ID,
		Origin:       s.Origin,
		Status:       s.Status,
		System:       s.System,
		Waypoint:     s.Waypoint,
		MaxFuel:      s.MaxFuel,
		CurrentFuel:  s.CurrentFuel,
		IsCargoFull:  s.CargoFull,
		Route:        route,
		TradeRoute:   tradeRoute,
		CurrentCargo: s.CurrentCargo,
		MaxCargo:     s.MaxCargo,
	}, nil
}

func AdaptShips(models []model.Ships) ([]*Ship, error) {
	ships := make([]*Ship, len(models))
	for i := range models {
		ship, err := AdaptShip(models[i])
		if err != nil {
			return nil, err
		}
		ships[i] = ship
	}

	return ships, nil
}
