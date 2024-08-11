package ai

import (
	"sync/atomic"
	"time"

	"github.com/julienlavocat/spacetraders/internal/sdk"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type TradingFleet struct {
	logger                   zerolog.Logger
	lastTradeRoutesUpdates   time.Time
	s                        *sdk.Sdk
	shipAvailables           chan *sdk.Ship
	systemId                 string
	ships                    []*sdk.Ship
	tradeRoutes              []*sdk.TradeRoute
	tradeRouteUpdateInterval time.Duration
	revenue                  atomic.Int32
	expanses                 atomic.Int32
}

// TODO: Elimintate ships paremeter to make the fleet autonomous
func NewTradingFleet(s *sdk.Sdk, fleetId string, systemId string, shipids []string) *TradingFleet {
	var ships []*sdk.Ship
	for _, shipId := range shipids {
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
			t.logger.Warn().Msgf("unable to find route for ship %s (fuel capacity: %d), retry in %.2fs", ship.Id, ship.Fuel.Capacity, t.tradeRouteUpdateInterval.Seconds())
			go func() {
				time.Sleep(t.tradeRouteUpdateInterval)
				t.shipAvailables <- ship
			}()
		}

		t.logger.Info().Interface("route", tradeRoute).Msgf("found trade route for ship %s", ship.Id)

		go func() {
			revenue, expanses, err := ship.FollowTradeRoute(tradeRoute)
			t.revenue.Add(revenue)
			t.expanses.Add(expanses)

			if err != nil {
				log.Error().Err(err).Msgf("ship %s failed to follow trade route", ship.Id)
			}

			t.shipAvailables <- ship
		}()
	}
}

func (t *TradingFleet) findTradeRoute(ship *sdk.Ship) *sdk.TradeRoute {
	if time.Since(t.lastTradeRoutesUpdates) >= t.tradeRouteUpdateInterval {
		t.tradeRoutes = t.s.Market.GetTradeRoutes(t.systemId)
		t.lastTradeRoutesUpdates = time.Now()
		t.logger.Info().Msg("updated trade routes")
	}

	for _, route := range t.tradeRoutes {
		if route.FuelCost > ship.Fuel.Capacity {
			continue
		}

		return route
	}

	return nil
}
