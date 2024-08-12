package ai

import (
	"container/list"
	"context"
	"time"

	. "github.com/go-jet/jet/v2/postgres"
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
	systemId          string
	probeBuyThreshold float64
}

func NewMarketProbesFleet(s *sdk.Sdk) *MarketProbesFleet {
	return &MarketProbesFleet{
		s:                 s,
		logger:            log.With().Str("component", "MarketProbesFleet").Logger(),
		waypointsQueue:    list.New(),
		shipsQueue:        list.New(),
		marketUpdateQueue: list.New(),
		probes:            make(map[string]model.MarketProbes),
		probeBuyThreshold: 0.3,
	}
}

func (m *MarketProbesFleet) BeginOperations(systemId string, updateRate time.Duration) {
	// TODO: Automate the target list :
	// - The first probe (starter one) MUST be placed at a shipyard offering probes
	// - Buy a new probe if Agents credit is > 100000
	// - Find the next marketplace where no probes are positionned
	// - Position the probe there
	// - Add the probe to the ticker for the next update

	m.systemId = systemId

	var probes []model.MarketProbes
	err := MarketProbes.SELECT(MarketProbes.AllColumns).Query(m.s.DB, &probes)
	if err != nil {
		m.logger.Fatal().Err(err).Msg("unable to retrieve probes")
	}

	for _, probe := range probes {
		m.probes[probe.Ship] = probe

		if probe.Waypoint != "" {
			m.marketUpdateQueue.PushBack(probe.Ship)
		} else {
			m.shipsQueue.PushBack(probe.Ship)
		}
	}

	waypointsQueue, err := m.getMarketplaces()
	if err != nil {
		log.Fatal().Err(err).Msgf("an error occured while retrieving marketplaces in system")
	}
	m.waypointsQueue = waypointsQueue

	timer := time.NewTicker(updateRate)
	defer timer.Stop()

	for range timer.C {
		go m.placeNextProbe()
		m.updateMarket()
	}
}

func (m *MarketProbesFleet) getMarketplaces() (*list.List, error) {
	alreadyProbed := Waypoints.INNER_JOIN(WaypointsTraits, Waypoints.ID.EQ(WaypointsTraits.WaypointID)).
		INNER_JOIN(MarketProbes, WaypointsTraits.WaypointID.EQ(MarketProbes.Waypoint)).
		SELECT(Waypoints.ID).
		WHERE(WaypointsTraits.TraitID.EQ(String(string(api.MARKETPLACE))))

	q := Waypoints.INNER_JOIN(WaypointsTraits, Waypoints.ID.EQ(WaypointsTraits.WaypointID)).
		SELECT(Waypoints.ID).
		WHERE(WaypointsTraits.TraitID.EQ(String(string(api.MARKETPLACE))).AND(Waypoints.ID.NOT_IN(alreadyProbed)))

	var results []model.Waypoints
	err := q.Query(m.s.DB, &results)
	if err != nil {
		m.logger.Error().Err(err).Msgf("unable to find markeplates in %s", m.systemId)
		return nil, err
	}

	marketplaces := list.New()

	for _, waypoint := range results {
		marketplaces.PushBack(waypoint.ID)
	}

	return marketplaces, nil
}

func (m *MarketProbesFleet) placeNextProbe() {
	// TODO: If not probes are setup, buy until credits < 25k from the first market available

	waypointEntry := m.waypointsQueue.Front()
	if waypointEntry == nil {
		return
	}

	m.tryToBuyProbe()

	shipEntry := m.shipsQueue.Front()
	if shipEntry == nil {
		return
	}

	ship := shipEntry.Value.(string)
	waypoint := waypointEntry.Value.(string)

	m.logger.Info().Str("ship", ship).Msgf("placing probe %s at marketplace %s", ship, waypoint)

	_, err := MarketProbes.UPDATE(MarketProbes.Waypoint).SET(waypoint).WHERE(MarketProbes.Ship.EQ(String(ship))).Exec(m.s.DB)
	if err != nil {
		m.logger.Error().Err(err).Msg("unable to update probe location")
	}

	utils.RetryRequest(m.s.Client.FleetAPI.OrbitShip(context.Background(), ship).Execute, m.logger, "unable to place probe %s in orbit of %s", ship, waypoint)
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

	hasShipyard, err := m.s.Market.HasShipyard(probe.Waypoint)
	if err != nil {
		m.logger.Error().Err(err).Msgf("unable to check if %s has a shipyard", probe.Waypoint)
	}
	if hasShipyard {
		res, errBody, err := utils.RetryRequestWithoutFatal(m.s.Client.SystemsAPI.GetShipyard(context.Background(), probe.System, probe.Waypoint).Execute, m.logger)
		if err != nil || errBody != nil {
			m.logger.Error().Err(err).Interface("body", errBody).Msgf("unable to get shipyard data at %s", probe.Waypoint)
		} else {
			m.s.Market.UpdateShipyard(probe.Waypoint, res.Data.Ships)
		}

	}

	m.marketUpdateQueue.Remove(shipId)
	m.marketUpdateQueue.PushBack(shipId.Value.(string))

	if len(res.Data.TradeGoods) == 0 {
		m.logger.Warn().Str("ship", shipId.Value.(string)).Msgf("no trade goods found at waypoint %s, it could be an issue with the probe fleet", res.Data.Symbol)
		return
	}

	m.s.Market.UpdateMarket(res.Data)
}

func (m *MarketProbesFleet) tryToBuyProbe() {
	if err := m.s.RefreshBalance(); err != nil {
		m.logger.Error().Err(err).Msg("unable to refresh balance")
		return
	}
	if m.shipsQueue.Len() > 0 {
		return
	}

	found, waypoint, amount, err := m.s.Market.FindLowestProductPrice(m.systemId, string(api.SHIP_PROBE))
	if err != nil {
		m.logger.Err(err).Msg("unable to query prices for probes")
		return
	}

	if !found {
		return
	}

	if float64(amount) >= float64(m.s.Balance.Load())*0.1 {
		return
	}

	m.logger.Info().Msgf("buying a probe at %s for %d", waypoint, amount)
	res, body, err := utils.RetryRequestWithoutFatal(
		m.s.Client.FleetAPI.PurchaseShip(context.Background()).PurchaseShipRequest(*api.NewPurchaseShipRequest(api.ShipType(api.SHIP_PROBE), waypoint)).Execute, m.logger,
	)
	if err != nil {
		m.logger.Error().Err(err).Interface("body", body).Msgf("unable to buy probe at %s", waypoint)
		return
	}
	m.logger.Info().Msgf("bought probe %s at %s for %d", res.Data.Ship.Symbol, waypoint, res.Data.Transaction.Price)

	_, err = MarketProbes.INSERT(MarketProbes.AllColumns).MODEL(model.MarketProbes{
		Ship:     res.Data.Ship.Symbol,
		System:   res.Data.Ship.Nav.SystemSymbol,
		Waypoint: "",
	}).Exec(m.s.DB)
	if err != nil {
		m.logger.Error().Err(err).Msgf("unable to insert probe %s", res.Data.Ship.Symbol)
	}

	m.shipsQueue.PushBack(res.Data.Ship.Symbol)
}
