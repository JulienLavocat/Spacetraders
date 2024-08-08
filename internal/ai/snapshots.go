package ai

import (
	"github.com/julienlavocat/spacetraders/internal/sdk"
)

type MiningFleetSnapshot struct {
	ShipStates map[string]string  `json:"shipStates"`
	Id         string             `json:"id"`
	MiningAt   string             `json:"miningAt"`
	Miners     []sdk.ShipSnapshot `json:"miners"`
	SellPlan   []sdk.SellPlan     `json:"sellPlan"`
	Hauler     sdk.ShipSnapshot   `json:"hauler"`
}

func newMiningFleetSnapshot(fleet *MiningFleetCommander) MiningFleetSnapshot {
	miners := make([]sdk.ShipSnapshot, len(fleet.miners))

	i := 0
	for _, ship := range fleet.miners {
		miners[i] = ship.GetSnapshot()
		i++
	}

	return MiningFleetSnapshot{
		Id:         fleet.Id,
		ShipStates: fleet.shipStates,
		MiningAt:   fleet.target,
		Hauler:     fleet.hauler.GetSnapshot(),
		Miners:     miners,
		SellPlan:   fleet.sellPlan,
	}
}
