package ai

import (
	"time"

	"github.com/julienlavocat/spacetraders/internal/sdk"
)

type MiningFleetSnapshot struct {
	StartTime  time.Time          `json:"startTime"`
	ShipStates map[string]string  `json:"shipStates"`
	Id         string             `json:"id"`
	MiningAt   string             `json:"miningAt"`
	Miners     []sdk.ShipSnapshot `json:"miners"`
	SellPlan   []sdk.SellPlan     `json:"sellPlan"`
	Hauler     sdk.ShipSnapshot   `json:"hauler"`
	Revenue    int32              `json:"revenue"`
	Expanses   int32              `json:"expanses"`
}

func newMiningFleetSnapshot(fleet *MiningFleet) MiningFleetSnapshot {
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
		Revenue:    fleet.revenue,
		Expanses:   fleet.expanses,
		StartTime:  fleet.startTime,
	}
}

type TradingShipResulsSnapshot struct {
	TradeRoute *sdk.TradeRoute `json:"tradeRoute"`
	Revenue    int64           `json:"revenue"`
	Expanses   int64           `json:"expanses"`
}
