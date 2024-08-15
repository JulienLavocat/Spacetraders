package rest

import (
	"encoding/json"
	"time"

	"github.com/julienlavocat/spacetraders/.gen/spacetraders/public/model"
	"github.com/julienlavocat/spacetraders/internal/sdk"
)

type TradingFleet struct {
	StartTime time.Time                     `json:"startTime"`
	Ships     map[string]TradingShipResults `json:"ships"`
	Id        string                        `json:"id"`
	SystemId  string                        `json:"systemId"`
	Revenue   int64                         `json:"revenue"`
	Expanses  int64                         `json:"expanses"`
}

func adaptTradingFleet(fleet model.TradingFleets) (*TradingFleet, error) {
	var ships map[string]TradingShipResults
	err := json.Unmarshal([]byte(*fleet.Ships), &ships)
	if err != nil {
		return nil, err
	}

	return &TradingFleet{
		StartTime: fleet.StartTime,
		Ships:     ships,
		Id:        fleet.ID,
		SystemId:  fleet.SystemID,
		Revenue:   fleet.Revenue,
		Expanses:  fleet.Expanses,
	}, nil
}

type TradingShipResults struct {
	TradeRoute      *sdk.TradeRoute `json:"tradeRoute"`
	Step            string          `json:"step"`
	Revenue         int64           `json:"revenue"`
	Expanses        int64           `json:"expanses"`
	TradesCompleted int32           `json:"tradesCompleted"`
}
