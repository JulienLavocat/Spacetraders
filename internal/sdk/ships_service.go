package sdk

import (
	"context"

	"github.com/julienlavocat/spacetraders/internal/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type ShipsService struct {
	ships  map[string]*Ship
	sdk    *Sdk
	logger zerolog.Logger
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
