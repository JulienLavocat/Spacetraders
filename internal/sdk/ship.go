package sdk

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/julienlavocat/spacetraders/internal/api"
	"github.com/julienlavocat/spacetraders/internal/utils"
	"github.com/rs/xid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Ship struct {
	logger          zerolog.Logger
	ctx             context.Context
	Fuel            api.ShipFuel
	tradeRoute      *TradeRoute
	sdk             *Sdk
	Cargo           Cargo
	cooldownUntil   *time.Time
	Nav             api.ShipNav
	Id              string
	Role            api.ShipRole
	route           []PathSegment
	refuelThreshold int32
	MaxCargo        int32
	CurrentCargo    int32
	IsDocked        bool
	IsInOrbit       bool
	IsCargoFull     bool
	HasCargo        bool
}

type TradeRouteStepCallback func(step string)

func NewShip(sdk *Sdk, ship api.Ship) *Ship {
	s := &Ship{
		Id:              ship.Symbol,
		logger:          log.With().Str("component", "Ship").Str("shipId", ship.Symbol).Logger(),
		refuelThreshold: int32(float64(ship.Fuel.Capacity) * 0.25),
		Fuel:            ship.Fuel,
		Cargo:           make(map[string]int32),
		Role:            ship.Registration.Role,
		sdk:             sdk,
	}

	s.setCargo(ship.Cargo)
	s.setNav(ship.Nav)
	s.reportStatus()

	if ship.Cooldown.HasExpiration() {
		s.setCooldown(*ship.Cooldown.Expiration)
	}

	s.logger.Info().Msg("ship loaded")

	return s
}

func (s *Ship) NavigateTo(destination string, correlationId string) int32 {
	s.ensureCooldown()
	s.logger.Info().Msgf("plotting route from %s to %s", s.Nav.WaypointSymbol, destination)
	expanses := int32(0)

	expanses += s.Refuel(correlationId)

	route, err := s.sdk.Navigation.PlotRoute(s.Nav.SystemSymbol, s.Nav.WaypointSymbol, destination, s.Fuel.Current)
	if err != nil {
		s.logger.Fatal().Err(err).Msg("unable to plot route")
	}

	s.route = route

	var stops []string
	for _, nextStop := range route {
		stops = append(stops, fmt.Sprintf("%s (%d)", nextStop.To, nextStop.Fuel))
	}

	s.logger.Info().Str("path", strings.Join(stops, " -> ")).Msgf("navigating from %s to %s", s.Nav.WaypointSymbol, destination)

	for _, nextStop := range route {
		s.ensureCooldown()

		if nextStop.To == s.Nav.WaypointSymbol {
			s.logger.Info().Msgf("ship already at %s", nextStop.To)
			continue
		}

		expanses += s.Refuel(correlationId)

		s.Orbit()

		res := utils.RetryRequest(
			s.sdk.Client.FleetAPI.NavigateShip(s.ctx, s.Id).NavigateShipRequest(*api.NewNavigateShipRequest(nextStop.To)).Execute,
			s.logger, "unable to navigate to waypoint %s", nextStop.To)

		navigationTime := s.Nav.Route.Arrival.UTC().Sub(s.Nav.Route.DepartureTime.UTC())

		s.logger.Info().
			Str("shipId", s.Id).
			Msgf("navigation will take %.2fs and consume %d fuel (current: %d/%d)", navigationTime.Seconds(), res.Data.Fuel.Consumed.Amount, res.Data.Fuel.Current, res.Data.Fuel.Capacity)

		s.setNav(res.Data.Nav)
		s.Fuel = res.Data.Fuel
		s.reportStatus()

		s.ensureCooldown()
		s.logger.Info().Str("shipId", s.Id).Msgf("navigated to %s", nextStop.To)
	}

	expanses += s.Refuel(correlationId)

	return expanses
}

func (s *Ship) Orbit() *Ship {
	s.ensureCooldown()

	if s.IsInOrbit {
		return s
	}

	s.logger.Info().Msgf("moving to orbit of %s", s.Nav.WaypointSymbol)

	res := utils.RetryRequest(s.sdk.Client.FleetAPI.OrbitShip(s.ctx, s.Id).Execute, s.logger, "unable to move to orbit")

	s.setNav(res.Data.Nav)
	s.reportStatus()

	return s
}

func (s *Ship) Dock() *Ship {
	s.ensureCooldown()

	if s.IsDocked {
		s.logger.Info().Msgf("ship already docked")
		return s
	}

	s.logger.Info().Msg("docking ship")

	res := utils.RetryRequest(s.sdk.Client.FleetAPI.DockShip(s.ctx, s.Id).Execute, s.logger, "unable to dock ship")

	s.setNav(res.Data.Nav)
	s.reportStatus()

	return s
}

func (s *Ship) Sell(plan []SellPlan, correlationId string) (int32, int32) {
	s.ensureCooldown()

	revenue := int32(0)
	expanses := int32(0)

	for _, step := range plan {
		expanses += s.NavigateTo(step.Location, correlationId)

		s.Dock()

		expanses += s.Refuel(correlationId)

		for product, amount := range step.ToSell {
			res := utils.RetryRequest(
				s.sdk.Client.FleetAPI.SellCargo(s.ctx, s.Id).SellCargoRequest(*api.NewSellCargoRequest(api.TradeSymbol(product), amount)).Execute,
				s.logger, "unable to sell %d %s at %s", amount, product, s.Nav.WaypointSymbol)

			s.setCargo(res.Data.Cargo)

			if err := s.sdk.Market.ReportTransaction(res.Data.Transaction, res.Data.Agent.Credits, correlationId); err != nil {
				s.logger.Error().Err(err).Interface("transation", res.Data.Transaction).Interface("agent", res.Data.Agent).Msgf("unable to report transation")
			}

			s.reportStatus()

			tx := res.Data.Transaction
			revenue += tx.TotalPrice
			s.logger.Info().Msgf("sold %d %s for %d (%d/u), balance is now %d", tx.Units, tx.TradeSymbol, tx.TotalPrice, tx.PricePerUnit, res.Data.Agent.Credits)
		}
	}

	return revenue, expanses
}

func (s *Ship) Refuel(correlationId string) int32 {
	s.ensureCooldown()
	s.Dock()

	if !s.sdk.Waypoints.CanRefuelAt(s.Nav.WaypointSymbol) {
		s.logger.Info().Msgf("can't refuel at %s", s.Nav.WaypointSymbol)
		return 0
	}

	res := utils.RetryRequest(s.sdk.Client.FleetAPI.RefuelShip(s.ctx, s.Id).RefuelShipRequest(*api.NewRefuelShipRequest()).Execute, s.logger, "unable to refuel ship")

	s.Fuel = res.Data.Fuel

	s.logger.Info().Msgf("refueled %d units (%d/%d) at the cost of %d (%d/unit), remaining credits is %d",
		res.Data.Transaction.Units,
		s.Fuel.Current,
		s.Fuel.Capacity,
		res.Data.Transaction.TotalPrice,
		res.Data.Transaction.PricePerUnit,
		res.Data.Agent.Credits)

	if res.Data.Transaction.Units > 0 {
		if err := s.sdk.Market.ReportTransaction(res.Data.Transaction, res.Data.Agent.Credits, correlationId); err != nil {
			s.logger.Error().Err(err).Interface("transation", res.Data.Transaction).Interface("agent", res.Data.Agent).Msgf("unable to report transation")
		}
	}

	s.reportStatus()

	return res.Data.Transaction.TotalPrice
}

func (s *Ship) Mine() *Ship {
	s.ensureCooldown()
	s.Orbit()

	res := utils.RetryRequest(s.sdk.Client.FleetAPI.ExtractResources(s.ctx, s.Id).ExtractResourcesRequest(*api.NewExtractResourcesRequest()).Execute, s.logger, "unable to mine")

	s.logger.Info().Msgf("extracted %d %s, ship cargo at %d/%d", res.Data.Extraction.Yield.Units, res.Data.Extraction.Yield.Symbol, res.Data.Cargo.Units, res.Data.Cargo.Capacity)

	s.setCargo(res.Data.Cargo)

	if res.Data.Cooldown.HasExpiration() {
		s.setCooldown(*res.Data.Cooldown.Expiration)
	}

	s.reportStatus()

	return s
}

func (s *Ship) DeliverAndFulfillContract(contract api.Contract) api.Contract {
	// TODO: Sort before iterating, just in case some destinations aren't the same to avoid unecessary trips between different goods

	s.ensureCooldown()
	corelationId := xid.New().String()

	canBeFulfilled := true
	for _, term := range contract.Terms.Deliver {

		if term.UnitsFulfilled == term.UnitsRequired {
			continue
		}

		amountInCargo, ok := s.Cargo[term.TradeSymbol]
		if !ok {
			s.logger.Info().Msgf("product %s not in cargo", term.TradeSymbol)
			canBeFulfilled = false
			continue
		}

		amount := min(amountInCargo, term.UnitsRequired-term.UnitsFulfilled)
		product := term.TradeSymbol

		s.NavigateTo(term.DestinationSymbol, corelationId)
		s.Dock()

		res := utils.RetryRequest(
			s.sdk.Client.ContractsAPI.DeliverContract(s.ctx, contract.Id).DeliverContractRequest(*api.NewDeliverContractRequest(s.Id, product, amount)).Execute,
			s.logger,
			"unable to deliver %d units of %s",
			term.UnitsRequired, term.TradeSymbol)
		s.setCargo(res.Data.Cargo)
		s.reportStatus()

		contract = res.Data.Contract

		if term.UnitsFulfilled+amount < term.UnitsRequired {
			canBeFulfilled = false
		}

		s.logger.Info().Msgf("delivered %d %s to fulfill contract %s", amount, product, contract.Id)
	}

	if canBeFulfilled {
		res := utils.RetryRequest(s.sdk.Client.ContractsAPI.FulfillContract(s.ctx, contract.Id).Execute, s.logger, "unable to fulfill contract %s", contract.Id)

		s.logger.Info().Msgf("fulfilled contract %s for faction %s +%d (%d)", res.Data.Contract.Id, res.Data.Contract.FactionSymbol, res.Data.Contract.Terms.Payment.OnFulfilled, res.Data.Agent.Credits)
		return res.Data.Contract
	}

	return contract
}

func (s *Ship) TransferPartialCargo(shipId string, maxAmount int32) {
	s.ensureCooldown()
	for product, amount := range s.Cargo {
		res := utils.RetryRequest(
			s.sdk.Client.FleetAPI.TransferCargo(s.ctx, s.Id).TransferCargoRequest(*api.NewTransferCargoRequest(api.TradeSymbol(product), min(amount, maxAmount), shipId)).Execute,
			s.logger, "unable to transfer cargo from %s to %s", s.Id, shipId)

		s.setCargo(res.Data.Cargo)
		maxAmount -= amount
		if maxAmount <= 0 {
			s.reportStatus()
			return
		}
	}
	s.reportStatus()
}

func (s *Ship) RefreshCargo() {
	res := utils.RetryRequest(s.sdk.Client.FleetAPI.GetMyShipCargo(s.ctx, s.Id).Execute, s.logger, "unable to refresh cargo")
	s.setCargo(res.Data)
	s.reportStatus()
}

func (s *Ship) JettisonCargo() error {
	s.logger.Info().Interface("cargo", s.Cargo).Msg("jettisonning cargo")

	s.ensureCooldown()
	for product, amount := range s.Cargo {
		res, errBody, err := utils.RetryRequestWithoutFatal(s.sdk.Client.FleetAPI.Jettison(s.ctx, s.Id).JettisonRequest(*api.NewJettisonRequest(api.TradeSymbol(product), amount)).Execute, s.logger)
		if err != nil {
			s.logger.Error().Err(err).Interface("body", errBody).Msgf("unable to jettison %d %s", amount, product)
			return err
		}

		s.logger.Info().Msgf("jettisonned %d %s", amount, product)
		s.setCargo(res.Data.Cargo)
	}
	s.reportStatus()

	return nil
}

func (s *Ship) Buy(product string, amount int32, correlationId string) (int32, error) {
	s.ensureCooldown()

	if amount == 0 {
		s.logger.Warn().Msgf("attempt to buy 0 %s, aborting", product)
		return 0, nil
	}

	res, errBody, err := utils.RetryRequestWithoutFatal(s.sdk.Client.FleetAPI.PurchaseCargo(s.ctx, s.Id).PurchaseCargoRequest(*api.NewPurchaseCargoRequest(api.TradeSymbol(product), amount)).Execute, s.logger)
	if err != nil || errBody != nil {
		s.logger.Err(err).Interface("body", errBody).Msgf("unable to buy goods at %s", s.Nav.WaypointSymbol)
		return 0, err
	}

	if err := s.sdk.Market.ReportTransaction(res.Data.Transaction, res.Data.Agent.Credits, correlationId); err != nil {
		s.logger.Error().Err(err).Interface("transation", res.Data.Transaction).Interface("agent", res.Data.Agent).Msgf("unable to report transation")
	}

	s.setCargo(res.Data.Cargo)
	s.reportStatus()

	tx := res.Data.Transaction
	s.logger.Info().Msgf("bought %d %s at %s for %d (%d/u), balance is now %d", tx.Units, tx.TradeSymbol, s.Nav.WaypointSymbol, tx.TotalPrice, tx.PricePerUnit, res.Data.Agent.Credits)

	return res.Data.Transaction.TotalPrice, nil
}

func (s *Ship) FollowTradeRoute(route *TradeRoute, stepCallback TradeRouteStepCallback) (int32, int32, error) {
	s.ensureCooldown()

	revenue := int32(0)
	expanses := int32(0)
	correlationId := xid.New().String()
	s.tradeRoute = route

	amount := min(s.MaxCargo-s.CurrentCargo, route.MaxAmount)
	expanses += s.NavigateTo(route.BuyAt, correlationId)
	stepCallback("BUYING")
	amount = min(amount, int32(s.sdk.Balance.Load()/int64(route.BuyPrice)))
	s.logger.Debug().Msgf("buying %d %s (availableCargo: %d, volume: %d, balance: %d, buyPrice: %d)", amount, route.Product, s.MaxCargo-s.CurrentCargo, route.MaxAmount, s.sdk.Balance.Load(), route.BuyPrice)
	if amount <= 0 {
		s.logger.Warn().Interface("cargo", s.Cargo).Int64("balance", s.sdk.Balance.Load()).Msgf("unable to buy %d %s, not enough money or cargo", amount, route.Product)
		return revenue, expanses, nil
	}
	txExpanses, err := s.Buy(route.Product, amount, correlationId)
	expanses += txExpanses
	if err != nil {
		return revenue, expanses, err
	}

	stepCallback("SELLING")
	txRevenue, txExpanses := s.Sell([]SellPlan{
		{
			ToSell: Cargo{
				route.Product: amount,
			},
			Location: route.SellAt,
		},
	}, correlationId)
	expanses += txExpanses
	revenue += txRevenue

	s.tradeRoute = nil

	return revenue, expanses, nil
}

func (s *Ship) GetSnapshot() ShipSnapshot {
	return newShipSnapshot(s)
}

func (s *Ship) setNav(data api.ShipNav) {
	s.logger.Debug().Interface("nav", data).Msg("navigation updated")

	s.Nav = data
	s.IsDocked = data.Status == api.DOCKED
	s.IsInOrbit = data.Status == api.IN_ORBIT

	if s.Nav.Route.Arrival.After(time.Now().UTC()) {
		s.setCooldown(s.Nav.Route.Arrival)
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

func (s *Ship) setCooldown(endsAt time.Time) {
	if endsAt.UTC().Before(time.Now().UTC()) {
		s.logger.Debug().Msgf("cooldown ending at %s is in the past, will not update the current cooldown timer ending at %s", endsAt, s.cooldownUntil)
		return
	}

	s.logger.Debug().Msgf("setting cooldown to %s", endsAt)
	s.cooldownUntil = &endsAt
}

func (s *Ship) ensureCooldown() {
	if s.cooldownUntil == nil {
		s.logger.Debug().Msg("ship isn't in cooldown")
		return
	}

	if s.cooldownUntil.Before(time.Now().UTC()) {
		s.logger.Debug().Msgf("cooldown ending at %s is in the past, skipping", s.cooldownUntil)
		return
	}

	sleepDuration := s.cooldownUntil.UTC().Sub(time.Now().UTC())
	s.logger.Info().Msgf("entering cooldown for %.2fs (until: %s)", sleepDuration.Seconds(), time.Now().UTC().Add(sleepDuration).String())
	time.Sleep(sleepDuration)
}

func (s *Ship) reportStatus() {
	if err := s.sdk.Ships.ReportStatus(s); err != nil {
		log.Error().Err(err).Msgf("an error occured while updating ship status %s", s.Id)
	}
}
