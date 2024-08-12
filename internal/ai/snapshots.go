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
		Revenue:    fleet.revenue,
		Expanses:   fleet.expanses,
		StartTime:  fleet.startTime,
	}
}

type TradingFleetSnapshot struct {
	StartTime    time.Time                            `json:"startTime"`
	ShipsResults map[string]TradingShipResulsSnapshot `json:"shipsResults"`
	Id           string                               `json:"id"`
	SystemId     string                               `json:"systemId"`
	Ships        []string                             `json:"ships"`
	Revenue      int64                                `json:"revenue"`
	Expanses     int64                                `json:"expanses"`
}

func newTradingFleetSnapshot(fleet *TradingFleet) TradingFleetSnapshot {
	assignedShips := make([]string, len(fleet.ships))
	for i, ship := range fleet.ships {
		assignedShips[i] = ship.Id
	}

	shipsResults := make(map[string]TradingShipResulsSnapshot)
	for shipId, results := range fleet.shipsResults {
		shipsResults[shipId] = TradingShipResulsSnapshot{
			Revenue:  results.Revenue.Load(),
			Expanses: results.Expanses.Load(),
		}
	}

	return TradingFleetSnapshot{
		StartTime:    fleet.startTime,
		Ships:        assignedShips,
		Id:           fleet.Id,
		SystemId:     fleet.systemId,
		Revenue:      fleet.revenue.Load(),
		Expanses:     fleet.expanses.Load(),
		ShipsResults: shipsResults,
	}
}

type TradingShipResulsSnapshot struct {
	Revenue  int64
	Expanses int64
}
