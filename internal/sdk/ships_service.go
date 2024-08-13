package sdk

import (
	"context"
	"database/sql"
	"encoding/json"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/julienlavocat/spacetraders/.gen/spacetraders/public/model"
	. "github.com/julienlavocat/spacetraders/.gen/spacetraders/public/table"
	"github.com/julienlavocat/spacetraders/internal/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type ShipsService struct {
	ships  map[string]*Ship
	sdk    *Sdk
	logger zerolog.Logger
	db     *sql.DB
}

func newShipService(s *Sdk) *ShipsService {
	return &ShipsService{
		ships:  make(map[string]*Ship),
		sdk:    s,
		logger: log.With().Str("component", "ShipsService").Logger(),
	}
}

func (s *ShipsService) GetShip(id string) (*Ship, error) {
	ship, ok := s.ships[id]

	if ok {
		return ship, nil
	}

	res, body, err := utils.RetryRequestWithoutFatal(s.sdk.FleetApi.GetMyShip(context.Background(), id).Execute, s.logger)
	if err != nil {
		s.logger.Error().Err(err).Interface("body", body).Msg("unable to load ship")
		return nil, err
	}

	ship = NewShip(s.sdk, res.Data)

	return ship, nil
}

func (s *ShipsService) ReportStatus(ship *Ship) error {
	cargo, err := json.Marshal(ship.Cargo)
	if err != nil {
		s.logger.Error().Err(err).Msgf("unable to marshal cargo as json for ship %s", ship.Id)
		return err
	}

	routeJson, err := json.Marshal(ship.route)
	if err != nil {
		s.logger.Error().Err(err).Msgf("unable to marshal route as json for ship %s", ship.Id)
		return err
	}
	route := string(routeJson)

	tradeRouteJson, err := json.Marshal(ship.tradeRoute)
	if err != nil {
		s.logger.Error().Err(err).Msgf("unable to marshal trade route as json for ship %s", ship.Id)
		return err
	}
	tradeRoute := string(tradeRouteJson)

	q := Ships.INSERT(Ships.AllColumns.Except(Ships.UpdatedAt)).ON_CONFLICT(Ships.ID).DO_UPDATE(SET(
		Ships.ArrivalAt.SET(Ships.EXCLUDED.ArrivalAt),
		Ships.DepartedAt.SET(Ships.EXCLUDED.DepartedAt),
		Ships.Waypoint.SET(Ships.EXCLUDED.Waypoint),
		Ships.System.SET(Ships.EXCLUDED.System),
		Ships.Status.SET(Ships.EXCLUDED.Status),
		Ships.Destination.SET(Ships.EXCLUDED.Destination),
		Ships.Origin.SET(Ships.EXCLUDED.Origin),
		Ships.MaxFuel.SET(Ships.EXCLUDED.MaxFuel),
		Ships.CurrentFuel.SET(Ships.EXCLUDED.CurrentFuel),
		Ships.MaxCargo.SET(Ships.EXCLUDED.MaxCargo),
		Ships.CurrentCargo.SET(Ships.EXCLUDED.CurrentCargo),
		Ships.CargoFull.SET(Ships.EXCLUDED.CargoFull),
		Ships.Cargo.SET(Ships.EXCLUDED.Cargo),
		Ships.Route.SET(Ships.EXCLUDED.Route),
		Ships.TradeRoute.SET(Ships.EXCLUDED.TradeRoute),
		Ships.Cooldown.SET(Ships.EXCLUDED.Cooldown),
		Ships.UpdatedAt.SET(NOW()),
	)).MODEL(model.Ships{
		ID:           ship.Id,
		Status:       string(ship.Nav.Status),
		ArrivalAt:    ship.Nav.Route.Arrival,
		DepartedAt:   ship.Nav.Route.DepartureTime,
		Waypoint:     ship.Nav.WaypointSymbol,
		System:       ship.Nav.SystemSymbol,
		Destination:  ship.Nav.Route.Destination.Symbol,
		Origin:       ship.Nav.Route.Origin.Symbol,
		CargoFull:    ship.IsCargoFull,
		MaxFuel:      ship.Fuel.Capacity,
		CurrentFuel:  ship.Fuel.Current,
		MaxCargo:     ship.MaxCargo,
		CurrentCargo: ship.CurrentCargo,
		Cargo:        string(cargo),
		Route:        &route,
		TradeRoute:   &tradeRoute,
		Cooldown:     ship.cooldownUntil,
	})

	_, err = q.Exec(s.sdk.DB)
	if err != nil {
		s.logger.Println(q.DebugSql())
		s.logger.Error().Err(err).Msgf("unable to insert report for ship %s", ship.Id)
		return err
	}

	return nil
}
