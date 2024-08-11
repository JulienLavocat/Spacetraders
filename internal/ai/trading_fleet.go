package ai

import (
	"sync/atomic"
	"time"

	"github.com/julienlavocat/spacetraders/internal/sdk"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type TradingFleet struct {
	logger                 zerolog.Logger
	lastTradeRoutesUpdates time.Time
	s                      *sdk.Sdk
	shipAvailables         chan *sdk.Ship
	systemId               string
	ships                  []*sdk.Ship
	tradeRoutes            []*sdk.TradeRoute
	updateInterval         time.Duration
	revenue                atomic.Int32
	expanses               atomic.Int32
}

// TODO: Elimintate ships paremeter to make the fleet autonomous
func NewTradingFleet(s *sdk.Sdk, fleetId string, systemId string, updateInterval time.Duration, shipIds []string) *TradingFleet {
	var ships []*sdk.Ship
	for _, shipId := range shipIds {
		ships = append(ships, s.Ships[shipId])
	}

	return &TradingFleet{
		s:              s,
		systemId:       systemId,
		ships:          ships,
		logger:         log.With().Str("component", "TradingFleet").Str("id", fleetId).Logger(),
		revenue:        atomic.Int32{},
		expanses:       atomic.Int32{},
		shipAvailables: make(chan *sdk.Ship, 100),
		updateInterval: updateInterval,
	}
}

func (t *TradingFleet) BeginOperations() {
	for _, ship := range t.ships {
		t.shipAvailables <- ship
	}

	t.logger.Info().Msg("begining operations")
	for ship := range t.shipAvailables {
		t.logger.Info().Msgf("%s is avaible", ship.Id)
		tradeRoute := t.findTradeRoute(ship)

		if tradeRoute == nil {
			t.logger.Warn().Str("ship", ship.Id).Msgf("unable to find route for ship %s (fuel capacity: %d), retry in %.2fs", ship.Id, ship.Fuel.Capacity, t.updateInterval.Seconds())
			go func() {
				time.Sleep(t.updateInterval)
				t.shipAvailables <- ship
			}()
			continue
		}

		t.logger.Info().Interface("route", tradeRoute).Str("ship", ship.Id).Msgf("found trade route for ship %s", ship.Id)

		go func() {
			revenue, expanses, err := ship.FollowTradeRoute(tradeRoute)
			t.revenue.Add(revenue)
			t.expanses.Add(expanses)

			if err != nil {
				t.logger.Error().Err(err).Str("ship", ship.Id).Msgf("ship %s failed to follow trade route", ship.Id)
			}

			t.shipAvailables <- ship
		}()
	}
}

func (t *TradingFleet) findTradeRoute(ship *sdk.Ship) *sdk.TradeRoute {
	if time.Since(t.lastTradeRoutesUpdates) >= t.updateInterval {
		t.tradeRoutes = t.s.Market.GetTradeRoutes(t.systemId)
		t.lastTradeRoutesUpdates = time.Now()
		t.logger.Info().Msg("updated trade routes")
	}

	for _, route := range t.tradeRoutes {
		if route.FuelCost > ship.Fuel.Capacity || route.EstimatedProfits < 0 {
			continue
		}

		return route
	}

	return nil
}
