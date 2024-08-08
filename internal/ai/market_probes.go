package ai

import (
	"container/list"
	"context"
	"time"

	"github.com/julienlavocat/spacetraders/internal/sdk"
	"github.com/julienlavocat/spacetraders/internal/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type MarketProbesFleet struct {
	s                 *sdk.Sdk
	ships             map[string]*sdk.Ship
	logger            zerolog.Logger
	waypointsQueue    *list.List
	shipsQueue        *list.List
	marketUpdateQueue *list.List
}

func NewMarketProbesFleet(s *sdk.Sdk, ships []*sdk.Ship) *MarketProbesFleet {
	probes := make(map[string]*sdk.Ship)
	for _, ship := range ships {
		probes[ship.Id] = ship
	}

	return &MarketProbesFleet{
		s:                 s,
		ships:             probes,
		logger:            log.With().Str("component", "MarketProbesFleet").Logger(),
		waypointsQueue:    list.New(),
		shipsQueue:        list.New(),
		marketUpdateQueue: list.New(),
	}
}

func (m *MarketProbesFleet) BeginOperations(systemId string, updateRate time.Duration) {
	// TODO: Automate the target list :
	// - Buy a new probe if Agents credit is > 100000
	// - Find the next marketplace where no probes are positionned
	// - Position the probe there
	// - Add the probe to the ticker for the next update

	for shipId := range m.ships {
		m.shipsQueue.PushBack(shipId)
	}

	for _, marketplace := range m.getMarketplaces() {
		m.waypointsQueue.PushBack(marketplace)
	}

	timer := time.NewTicker(updateRate)
	defer timer.Stop()

	for range timer.C {
		go m.placeNextProbe()
		m.updateMarket()
	}
}

func (m *MarketProbesFleet) getMarketplaces() []string {
	// TODO: Return all marketplaces of system
	return []string{"X1-QA42-B7"}
}

func (m *MarketProbesFleet) placeNextProbe() {
	waypoint := m.waypointsQueue.Front()
	if waypoint == nil {
		return
	}

	ship := m.shipsQueue.Front()
	if ship == nil {
		return
	}

	m.logger.Info().Msgf("placing probe %s at marketplace %s", ship.Value.(string), waypoint.Value.(string))
	m.ships[ship.Value.(string)].NavigateTo(waypoint.Value.(string))

	m.shipsQueue.Remove(ship)
	m.waypointsQueue.Remove(waypoint)
	m.marketUpdateQueue.PushBack(ship.Value.(string))
}

func (m *MarketProbesFleet) updateMarket() {
	shipId := m.marketUpdateQueue.Front()
	if shipId == nil {
		return
	}

	ship := m.ships[shipId.Value.(string)]
	res := utils.RetryRequest(
		m.s.Client.SystemsAPI.GetMarket(context.Background(), ship.Nav.SystemSymbol, ship.Nav.WaypointSymbol).Execute,
		m.logger, "unable to fetch market using ship %s at waypoint %s", ship.Id, ship.Nav.WaypointSymbol)

	if len(res.Data.TradeGoods) == 0 {
		m.logger.Warn().Msgf("no trade goods found at waypoint %s, it could be an issue with the probe fleet", res.Data.Symbol)
		return
	}

	m.s.Market.UpdateMarket(res.Data)

	m.marketUpdateQueue.Remove(shipId)
	m.marketUpdateQueue.PushBack(shipId.Value.(string))
}
