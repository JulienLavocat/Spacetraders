package ai

import (
	"fmt"
	"sync"
	"time"

	"github.com/julienlavocat/spacetraders/internal/api"
	"github.com/julienlavocat/spacetraders/internal/sdk"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type MiningFleetCommander struct {
	Id               string
	logger           zerolog.Logger
	s                *sdk.Sdk
	hauler           *sdk.Ship
	shipNeedsHauling chan string
	miners           map[string]*sdk.Ship
	shipStates       map[string]string
	systemId         string
	target           string
	sellPlan         []sdk.SellPlan
	revenue          int32
	expanses         int32
	startTime        time.Time
}

func NewMiningFleetCommander(s *sdk.Sdk, id string, miners []*sdk.Ship, hauler *sdk.Ship) *MiningFleetCommander {
	ships := make(map[string]*sdk.Ship)
	shipStates := make(map[string]string)
	for _, miner := range miners {
		ships[miner.Id] = miner
		shipStates[miner.Id] = "IDLE"
	}

	shipStates[hauler.Id] = "IDLE"

	fleet := &MiningFleetCommander{
		logger:     log.With().Str("component", "MiningFleetCommander").Str("id", id).Logger(),
		s:          s,
		miners:     ships,
		hauler:     hauler,
		shipStates: shipStates,
		Id:         id,
	}

	return fleet
}

func (m *MiningFleetCommander) BeginOperations(systemId string) error {
	m.systemId = systemId

	target, err := m.determineTarget()
	if err != nil {
		return err
	}

	m.logger.Info().Msgf("begining operations in system %s at waypoint %s", systemId, target)

	m.target = target
	m.moveFleetToTarget()

	m.startTime = time.Now().UTC()

	m.shipNeedsHauling = make(chan string)
	for _, miner := range m.miners {
		go m.performMiningOperations(miner)
	}
	m.performHaulingOperation()

	return nil
}

func (m *MiningFleetCommander) GetSnapshot() MiningFleetSnapshot {
	return newMiningFleetSnapshot(m)
}

func (m *MiningFleetCommander) determineTarget() (string, error) {
	// TODO: Find appropriate target based on what sells the most in the system
	m.logger.Info().Msgf("determining mining target in %s", m.systemId)

	waypoints := m.s.Waypoints.FindWaypointsByType(m.systemId, api.ENGINEERED_ASTEROID)

	if len(waypoints) == 0 {
		return "", fmt.Errorf("no target found matching system: %s and trait: %s", m.systemId, api.ENGINEERED_ASTEROID)
	}

	m.logger.Debug().Interface("waypoints", waypoints).Msg("found matching waypoints")

	return waypoints[0].ID, nil
}

func (m *MiningFleetCommander) moveFleetToTarget() {
	var wg sync.WaitGroup
	for i := range m.miners {
		m.logger.Info().Msgf("requesting mining ship %s to move to target %s", m.miners[i].Id, m.target)
		wg.Add(1)
		go func() {
			defer wg.Done()
			m.miners[i].NavigateTo(m.target)
			m.logger.Info().Msgf("miner %s arrived at target %s", m.miners[i].Id, m.target)
		}()
	}

	m.logger.Info().Msgf("requesting hauling ship %s to move to target %s", m.hauler.Id, m.target)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if m.hauler.HasCargo {
			m.logger.Info().Msg("hauler has cargo, selling before navigating to target")
			m.sellHaulerCargo()
		}

		m.hauler.NavigateTo(m.target)
		m.logger.Info().Msgf("hauler %s arrived at target %s", m.hauler.Id, m.target)
	}()

	wg.Wait()

	m.logger.Info().Msgf("all ships in fleet have been moved to target %s", m.target)
}

func (m *MiningFleetCommander) performMiningOperations(ship *sdk.Ship) {
	log.Info().Msgf("%s begining mining operations", ship.Id)
	m.shipStates[ship.Id] = "MINING"
	for !ship.IsCargoFull {
		ship.Mine()
	}

	m.shipStates[ship.Id] = "FULL"
	m.logger.Info().Msgf("%s is full, waiting for cargo transfer to %s", ship.Id, m.hauler.Id)
	m.shipNeedsHauling <- ship.Id
}

func (m *MiningFleetCommander) performHaulingOperation() {
	m.logger.Info().Msgf("%s begining hauling operations", m.hauler.Id)

	for shipId := range m.shipNeedsHauling {
		m.shipStates[m.hauler.Id] = "FILLING"
		m.logger.Info().Msgf("transfering cargo from ship %s to hauler %s", shipId, m.hauler.Id)

		ship := m.miners[shipId]

		if m.hauler.IsCargoFull {
			m.logger.Info().Msg("hauler is full, selling")
			m.sellHaulerCargo()
		}

		maxTransferAmount := m.hauler.MaxCargo
		if ship.CurrentCargo+m.hauler.CurrentCargo >= m.hauler.MaxCargo {
			maxTransferAmount = m.hauler.MaxCargo - m.hauler.CurrentCargo
		}

		ship.TransferPartialCargo(m.hauler.Id, maxTransferAmount)
		go m.performMiningOperations(ship)

		m.hauler.RefreshCargo()
		if m.hauler.IsCargoFull {
			m.logger.Info().Msg("hauler is full, selling")
			m.sellHaulerCargo()
		}
	}

	m.shipStates[m.hauler.Id] = "IDLE"
	m.logger.Info().Msg("hauling operations completed")
}

func (m *MiningFleetCommander) sellHaulerCargo() {
	m.shipStates[m.hauler.Id] = "SELLING"
	m.sellPlan = m.s.Market.CreateSellPlan(m.systemId, m.hauler.Cargo)
	revenue, expanses := m.hauler.Sell(m.sellPlan)
	expanses += m.hauler.Refuel()
	m.hauler.NavigateTo(m.target)
	m.shipStates[m.hauler.Id] = "IDLE"
	m.revenue += revenue
	m.expanses += expanses
}
