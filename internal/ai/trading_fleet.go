package ai

import (
	"encoding/json"
	"time"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/julienlavocat/spacetraders/.gen/spacetraders/public/model"
	. "github.com/julienlavocat/spacetraders/.gen/spacetraders/public/table"
	"github.com/julienlavocat/spacetraders/internal/sdk"
	"github.com/julienlavocat/spacetraders/internal/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.uber.org/atomic"
)

type TradingFleet struct {
	logger                 zerolog.Logger
	startTime              time.Time
	lastTradeRoutesUpdates time.Time
	s                      *sdk.Sdk
	shipAvailables         chan *sdk.Ship
	shipsResults           map[string]*TradingShipResults
	systemId               string
	Id                     string
	tradeRoutes            utils.Queue[*sdk.TradeRoute]
	ships                  []*sdk.Ship
	expanses               atomic.Int64
	revenue                atomic.Int64
	updateInterval         time.Duration
}

type TradingShipResults struct {
	TradeRoute      *sdk.TradeRoute `json:"tradeRoute"`
	Step            string          `json:"step"`
	Revenue         atomic.Int64    `json:"revenue"`
	Expanses        atomic.Int64    `json:"expanses"`
	TradesCompleted int32           `json:"tradesCompleted"`
}

func NewShipResults() *TradingShipResults {
	return &TradingShipResults{
		Revenue:         atomic.Int64{},
		Expanses:        atomic.Int64{},
		TradesCompleted: 0,
		Step:            "BUYING",
	}
}

// TODO: Elimintate ships paremeter to make the fleet autonomous
func NewTradingFleet(s *sdk.Sdk, fleetId string, systemId string, updateInterval time.Duration, shipIds []string) *TradingFleet {
	logger := log.With().Str("component", "TradingFleet").Str("id", fleetId).Logger()

	var ships []*sdk.Ship
	for _, shipId := range shipIds {
		ship, err := s.Ships.GetShip(shipId)
		if err != nil {
			logger.Fatal().Err(err).Msgf("ship %s not found", shipId)
		}

		ships = append(ships, ship)
	}

	return &TradingFleet{
		s:              s,
		Id:             fleetId,
		systemId:       systemId,
		ships:          ships,
		logger:         logger,
		revenue:        atomic.Int64{},
		expanses:       atomic.Int64{},
		shipAvailables: make(chan *sdk.Ship, 100),
		startTime:      time.Now().UTC(),
		shipsResults:   make(map[string]*TradingShipResults),
		updateInterval: updateInterval,
		tradeRoutes:    utils.NewQueue[*sdk.TradeRoute](),
	}
}

func (t *TradingFleet) BeginOperations() {
	t.startTime = time.Now().UTC()
	for _, ship := range t.ships {
		go func() {
			err := ship.JettisonCargo()
			if err != nil {
				t.logger.Warn().Str("ship", ship.Id).Interface("error", err).Msg("unable to jettison cargo")
			}

			t.shipsResults[ship.Id] = NewShipResults()
			t.shipAvailables <- ship
		}()
	}

	t.logger.Info().Msg("begining operations")
	for ship := range t.shipAvailables {
		t.logger.Info().Msgf("%s is avaible", ship.Id)
		tradeRoute := t.getNextTradeRoute()

		if tradeRoute == nil {
			t.logger.Warn().Str("ship", ship.Id).Msgf("unable to find route for ship %s (fuel capacity: %d), retry in %.2fs", ship.Id, ship.Fuel.Capacity, t.updateInterval.Seconds())
			go func() {
				time.Sleep(t.updateInterval)
				t.shipAvailables <- ship
			}()
			continue
		}

		t.logger.Info().Interface("route", tradeRoute).Str("ship", ship.Id).Msgf("found trade route for ship %s", ship.Id)
		t.shipsResults[ship.Id].TradeRoute = tradeRoute

		go func(route *sdk.TradeRoute) {
			results := t.shipsResults[ship.Id]

			revenue, expanses, err := ship.FollowTradeRoute(route, func(step string) {
				results.Step = step
				log.Debug().Interface("res", results).Msgf("status updated: %s", step)
				t.reportStatus()
			})
			t.revenue.Add(int64(revenue))
			t.expanses.Add(int64(expanses))

			results.Revenue.Add(int64(revenue))
			results.Expanses.Add(int64(expanses))

			if err != nil {
				t.logger.Error().Err(err).Str("ship", ship.Id).Msgf("ship %s failed to follow trade route", ship.Id)
			} else {
				results.TradesCompleted++
			}

			t.shipAvailables <- ship
		}(tradeRoute)

		t.reportStatus()
	}
}

func (t *TradingFleet) getNextTradeRoute() *sdk.TradeRoute {
	if time.Since(t.lastTradeRoutesUpdates) >= t.updateInterval {
		t.tradeRoutes.Clear()
		t.tradeRoutes.QueueAll(t.s.Market.GetTradeRoutes(t.systemId))
		t.lastTradeRoutesUpdates = time.Now()
		t.logger.Info().Msg("updated trade routes")
	}

	if !t.tradeRoutes.HasNext() {
		return nil
	}

	return t.tradeRoutes.Dequeue()
}

func (t *TradingFleet) reportStatus() {
	shipsResultsJson, err := json.Marshal(t.shipsResults)
	if err != nil {
		t.logger.Error().Err(err).Msgf("unable to marshall ships results")
		return
	}
	shipsResults := string(shipsResultsJson)

	q := TradingFleets.INSERT(TradingFleets.AllColumns).
		ON_CONFLICT(TradingFleets.ID).
		DO_UPDATE(SET(
			TradingFleets.SystemID.SET(TradingFleets.EXCLUDED.SystemID),
			TradingFleets.StartTime.SET(TradingFleets.EXCLUDED.StartTime),
			TradingFleets.Revenue.SET(TradingFleets.EXCLUDED.Revenue),
			TradingFleets.Expanses.SET(TradingFleets.EXCLUDED.Expanses),
			TradingFleets.Ships.SET(TradingFleets.EXCLUDED.Ships),
			TradingFleets.UpdatedAt.SET(NOW()),
		)).MODEL(model.TradingFleets{
		ID:        t.Id,
		SystemID:  t.systemId,
		StartTime: t.startTime,
		Revenue:   t.revenue.Load(),
		Expanses:  t.expanses.Load(),
		Ships:     &shipsResults,
	})

	_, err = q.Exec(t.s.DB)
	if err != nil {
		t.logger.Error().Err(err).Msg("unable to report status")
	}
}
