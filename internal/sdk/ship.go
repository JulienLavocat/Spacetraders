package sdk

import (
	"context"
	"time"

	"github.com/julienlavocat/spacetraders/internal/api"
	"github.com/julienlavocat/spacetraders/internal/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Ship struct {
	logger          zerolog.Logger
	ctx             context.Context
	Fuel            api.ShipFuel
	client          *api.APIClient
	Cargo           Cargo
	Nav             api.ShipNav
	Id              string
	Role            api.ShipRole
	refuelThreshold int32
	CurrentCargo    int32
	MaxCargo        int32
	IsDocked        bool
	IsInOrbit       bool
	IsCargoFull     bool
	HasCargo        bool
}

func NewShip(client *api.APIClient, ship api.Ship) *Ship {
	s := &Ship{
		Id:              ship.Symbol,
		logger:          log.With().Str("component", "Ship").Str("shipId", ship.Symbol).Logger(),
		refuelThreshold: int32(float64(ship.Fuel.Capacity) * 0.25),
		Fuel:            ship.Fuel,
		client:          client,
		Cargo:           make(map[string]int32),
		Role:            ship.Registration.Role,
	}

	s.setCargo(ship.Cargo)
	s.setNav(ship.Nav, true)

	if ship.Cooldown.RemainingSeconds > 0 {
		s.enterCooldown(time.Duration(ship.Cooldown.RemainingSeconds) * time.Second)
	}

	s.logger.Info().Msg("ship loaded")

	return s
}

func (s *Ship) NavigateTo(waypoint string) *Ship {
	s.logger.Info().Msgf("navigating from %s to %s", s.Nav.WaypointSymbol, waypoint)

	if waypoint == s.Nav.WaypointSymbol {
		s.logger.Info().Msgf("ship already at %s", waypoint)
		return s
	}

	// TODO: Ensure fuel is sufficient before navigating
	if s.Fuel.Current <= s.refuelThreshold {
		s.Refuel()
	}

	s.Orbit()

	res := utils.RetryRequest(s.client.FleetAPI.NavigateShip(s.ctx, s.Id).NavigateShipRequest(*api.NewNavigateShipRequest(waypoint)).Execute, s.logger, "unable to navigate to waypoint %s", waypoint)

	navigationTime := s.Nav.Route.Arrival.Sub(s.Nav.Route.DepartureTime)

	s.logger.Info().
		Str("shipId", s.Id).
		Msgf("navigation will take %.2fs and consume %d fuel (current: %d/%d)", navigationTime.Seconds(), res.Data.Fuel.Consumed.Amount, res.Data.Fuel.Current, res.Data.Fuel.Capacity)

	s.setNav(res.Data.Nav, false)
	s.Fuel = res.Data.Fuel

	s.logger.Info().Str("shipId", s.Id).Msgf("navigated to %s", waypoint)

	return s
}

func (s *Ship) Orbit() *Ship {
	if s.IsInOrbit {
		return s
	}

	s.logger.Info().Msg("moving to orbit")

	res := utils.RetryRequest(s.client.FleetAPI.OrbitShip(s.ctx, s.Id).Execute, s.logger, "unable to move to orbit")

	s.setNav(res.Data.Nav, false)

	s.logger.Info().Msgf("oribiting %s", s.Nav.WaypointSymbol)
	return s
}

func (s *Ship) Dock() *Ship {
	if s.IsDocked {
		s.logger.Info().Msgf("ship already docked")
		return s
	}

	s.logger.Info().Msg("docking ship")

	res := utils.RetryRequest(s.client.FleetAPI.DockShip(s.ctx, s.Id).Execute, s.logger, "unable to dock ship")

	s.setNav(res.Data.Nav, false)

	return s
}

func (s *Ship) Sell(plan []SellPlan) *Ship {
	for _, step := range plan {
		s.NavigateTo(step.Location)

		s.Dock()

		for product, amount := range step.ToSell {
			res := utils.RetryRequest(
				s.client.FleetAPI.SellCargo(s.ctx, s.Id).SellCargoRequest(*api.NewSellCargoRequest(api.TradeSymbol(product), amount)).Execute,
				s.logger, "unable to sell %d %s at %s", amount, product, s.Nav.WaypointSymbol)

			s.setCargo(res.Data.Cargo)

			tx := res.Data.Transaction
			s.logger.Info().Msgf("sold %d %s for %d (%d/u), balance is now %d", tx.Units, tx.TradeSymbol, tx.TotalPrice, tx.PricePerUnit, res.Data.Agent.Credits)
		}
	}

	return s
}

func (s *Ship) Refuel() *Ship {
	s.logger.Info().Msgf("refueling ship (%d/%d)", s.Fuel.Current, s.Fuel.Capacity)
	s.Dock()

	res := utils.RetryRequest(s.client.FleetAPI.RefuelShip(s.ctx, s.Id).RefuelShipRequest(*api.NewRefuelShipRequest()).Execute, s.logger, "unable to refuel ship")

	// TODO: Improve refuel logic to find the cheapest refuel point in the system, might be able to be extend search to n+1 systems in the future

	s.Fuel = res.Data.Fuel

	log.Info().Msgf("refueled %d units (%d/%d) at the cost of %d (%d/unit), remaining credits is %d",
		res.Data.Transaction.Units,
		s.Fuel.Current,
		s.Fuel.Capacity,
		res.Data.Transaction.TotalPrice,
		res.Data.Transaction.PricePerUnit,
		res.Data.Agent.Credits)

	return s
}

func (s *Ship) Mine() *Ship {
	s.Orbit()

	res := utils.RetryRequest(s.client.FleetAPI.ExtractResources(s.ctx, s.Id).ExtractResourcesRequest(*api.NewExtractResourcesRequest()).Execute, s.logger, "unable to mine")

	s.logger.Info().Msgf("extracted %d %s, ship cargo at %d/%d", res.Data.Extraction.Yield.Units, res.Data.Extraction.Yield.Symbol, res.Data.Cargo.Units, res.Data.Cargo.Capacity)

	s.setCargo(res.Data.Cargo)

	s.enterCooldown(time.Duration(res.Data.Cooldown.RemainingSeconds) * time.Second)

	return s
}

func (s *Ship) CountInCargo(product string) int32 {
	return s.Cargo[product]
}

func (s *Ship) DeliverAndFulfillContract(contract api.Contract) api.Contract {
	// TODO: Sort before iterating, just in case some destinations aren't the same to avoid unecessary trips between different goods

	canBeFulfilled := true
	for _, term := range contract.Terms.Deliver {

		if term.UnitsFulfilled == term.UnitsRequired {
			continue
		}

		amountInCargo, ok := s.Cargo[term.TradeSymbol]
		if !ok {
			log.Info().Msgf("product %s not in cargo", term.TradeSymbol)
			canBeFulfilled = false
			continue
		}

		amount := min(amountInCargo, term.UnitsRequired-term.UnitsFulfilled)
		product := term.TradeSymbol

		s.NavigateTo(term.DestinationSymbol).Dock()

		res := utils.RetryRequest(
			s.client.ContractsAPI.DeliverContract(s.ctx, contract.Id).DeliverContractRequest(*api.NewDeliverContractRequest(s.Id, product, amount)).Execute,
			s.logger,
			"unable to deliver %d units of %s",
			term.UnitsRequired, term.TradeSymbol)
		s.setCargo(res.Data.Cargo)

		contract = res.Data.Contract

		if term.UnitsFulfilled+amount < term.UnitsRequired {
			canBeFulfilled = false
		}

		s.logger.Info().Msgf("delivered %d %s to fulfill contract %s", amount, product, contract.Id)
	}

	if canBeFulfilled {
		res := utils.RetryRequest(s.client.ContractsAPI.FulfillContract(s.ctx, contract.Id).Execute, s.logger, "unable to fulfill contract %s", contract.Id)

		s.logger.Info().Msgf("fulfilled contract %s for faction %s +%d (%d)", res.Data.Contract.Id, res.Data.Contract.FactionSymbol, res.Data.Contract.Terms.Payment.OnFulfilled, res.Data.Agent.Credits)
		return res.Data.Contract
	}

	return contract
}

func (s *Ship) TransferPartialCargo(shipId string, maxAmount int32) {
	for product, amount := range s.Cargo {
		res := utils.RetryRequest(
			s.client.FleetAPI.TransferCargo(s.ctx, s.Id).TransferCargoRequest(*api.NewTransferCargoRequest(api.TradeSymbol(product), min(amount, maxAmount), shipId)).Execute,
			s.logger, "unable to transfer cargo from %s to %s", s.Id, shipId)

		s.setCargo(res.Data.Cargo)
		maxAmount -= amount
		if maxAmount <= 0 {
			return
		}
	}
}

func (s *Ship) RefreshCargo() {
	res := utils.RetryRequest(s.client.FleetAPI.GetMyShipCargo(s.ctx, s.Id).Execute, s.logger, "unable to refresh cargo")
	s.setCargo(res.Data)
}

func (s *Ship) GetSnapshot() ShipSnapshot {
	return newShipSnapshot(s)
}

func (s *Ship) setNav(data api.ShipNav, initial bool) {
	s.logger.Debug().Interface("nav", data).Msg("navigation updated")

	s.Nav = data
	s.IsDocked = data.Status == api.DOCKED
	s.IsInOrbit = data.Status == api.IN_ORBIT

	// If ship is a probe and it's the initial load time, we don't wait for it's cooldown as they are only moved once (for market data purposes)
	if s.Nav.Route.Arrival.After(time.Now().UTC()) && (s.Role != api.SHIP_ROLE_SATELLITE || !initial) {
		navigationTime := s.Nav.Route.Arrival.Sub(s.Nav.Route.DepartureTime)
		s.enterCooldown(navigationTime)
	}
}

func (s *Ship) setCargo(data api.ShipCargo) {
	cargo := make(map[string]int32)

	for _, product := range data.Inventory {
		cargo[string(product.Symbol)] = product.Units
	}

	s.logger.Debug().Interface("cargo", cargo).Msgf("cargo updated (%d/%d)", data.Units, data.Capacity)
	s.Cargo = cargo
	s.IsCargoFull = data.Capacity == data.Units
	s.HasCargo = data.Units > 0
	s.CurrentCargo = data.Units
	s.MaxCargo = data.Capacity
}

func (s *Ship) enterCooldown(d time.Duration) {
	s.logger.Info().Msgf("entering cooldown for %.2fs (until: %s)", d.Seconds(), time.Now().UTC().Add(d).String())
	time.Sleep(d)
}
