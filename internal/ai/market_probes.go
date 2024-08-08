package ai

import (
	"container/list"
	"context"
	"time"

	"github.com/julienlavocat/spacetraders/.gen/spacetraders/public/model"
	. "github.com/julienlavocat/spacetraders/.gen/spacetraders/public/table"
	"github.com/julienlavocat/spacetraders/internal/api"
	"github.com/julienlavocat/spacetraders/internal/sdk"
	"github.com/julienlavocat/spacetraders/internal/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type MarketProbesFleet struct {
	s                 *sdk.Sdk
	logger            zerolog.Logger
	waypointsQueue    *list.List
	shipsQueue        *list.List
	marketUpdateQueue *list.List
	probes            map[string]model.MarketProbes
}

func NewMarketProbesFleet(s *sdk.Sdk) *MarketProbesFleet {
	return &MarketProbesFleet{
		s:                 s,
		logger:            log.With().Str("component", "MarketProbesFleet").Logger(),
		waypointsQueue:    list.New(),
		shipsQueue:        list.New(),
		marketUpdateQueue: list.New(),
		probes:            make(map[string]model.MarketProbes),
	}
}

func (m *MarketProbesFleet) BeginOperations(systemId string, updateRate time.Duration) {
	// TODO: Automate the target list :
	// - Buy a new probe if Agents credit is > 100000
	// - Find the next marketplace where no probes are positionned
	// - Position the probe there
	// - Add the probe to the ticker for the next update

	var probes []model.MarketProbes
	err := MarketProbes.SELECT(MarketProbes.AllColumns).Query(m.s.DB, &probes)
	if err != nil {
		m.logger.Fatal().Err(err).Msg("unable to retrieve probes")
	}

	for _, probe := range probes {
		m.probes[probe.Ship] = probe
		m.marketUpdateQueue.PushBack(probe.Ship)
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
	return []string{}
}

func (m *MarketProbesFleet) placeNextProbe() {
	waypointEntry := m.waypointsQueue.Front()
	if waypointEntry == nil {
		return
	}

	shipEntry := m.shipsQueue.Front()
	if shipEntry == nil {
		return
	}

	ship := shipEntry.Value.(string)
	waypoint := waypointEntry.Value.(string)

	m.logger.Info().Msgf("placing probe %s at marketplace %s", ship, waypoint)

	MarketProbes.INSERT(MarketProbes.AllColumns).MODEL(model.MarketProbes{
		Ship:     ship,
		Waypoint: waypoint,
	})

	utils.RetryRequest(
		m.s.Client.FleetAPI.NavigateShip(context.Background(), ship).NavigateShipRequest(*api.NewNavigateShipRequest(waypoint)).Execute,
		m.logger, "unable to place probe %s at location %s", ship, waypoint)

	m.shipsQueue.Remove(shipEntry)
	m.waypointsQueue.Remove(waypointEntry)
	m.marketUpdateQueue.PushBack(ship)
}

func (m *MarketProbesFleet) updateMarket() {
	shipId := m.marketUpdateQueue.Front()
	if shipId == nil {
		return
	}

	probe := m.probes[shipId.Value.(string)]

	res := utils.RetryRequest(
		m.s.Client.SystemsAPI.GetMarket(context.Background(), probe.System, probe.Waypoint).Execute,
		m.logger, "unable to fetch market using ship %s at waypoint %s", probe.Ship, probe.Waypoint)

	m.marketUpdateQueue.Remove(shipId)
	m.marketUpdateQueue.PushBack(shipId.Value.(string))

	if len(res.Data.TradeGoods) == 0 {
		m.logger.Warn().Msgf("no trade goods found at waypoint %s, it could be an issue with the probe fleet", res.Data.Symbol)
		return
	}

	m.s.Market.UpdateMarket(res.Data)
}
