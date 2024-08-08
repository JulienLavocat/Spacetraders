package sdk

import (
	"database/sql"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/julienlavocat/spacetraders/.gen/spacetraders/public/model"
	. "github.com/julienlavocat/spacetraders/.gen/spacetraders/public/table"
	"github.com/julienlavocat/spacetraders/internal/api"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type WaypointsService struct {
	db     *sql.DB
	logger zerolog.Logger
}

func NewWaypointsService(db *sql.DB) *WaypointsService {
	return &WaypointsService{
		logger: log.With().Str("component", "WaypointsService").Logger(),
		db:     db,
	}
}

func (s *WaypointsService) FindWaypointsByType(systemId string, waypointType api.WaypointType) []model.Waypoints {
	var results []model.Waypoints

	err := Waypoints.
		SELECT(Waypoints.ID).
		WHERE(Waypoints.Type.EQ(String(string(waypointType))).AND(Waypoints.SystemID.EQ(String(systemId)))).
		Query(s.db, &results)
	if err != nil {
		s.logger.Fatal().Err(err).Msgf("unable to query waypoints for system %s and type %s", systemId, string(waypointType))
	}

	return results
}

func (s *WaypointsService) FindMarketplaces(systemId string) []model.Waypoints {
	var results []model.Waypoints

	err := Waypoints.
		INNER_JOIN(WaypointsTraits, Waypoints.ID.EQ(WaypointsTraits.WaypointID)).
		SELECT(Waypoints.ID).
		WHERE(Waypoints.SystemID.EQ(String(systemId)).AND(WaypointsTraits.TraitID.EQ(String(string(api.MARKETPLACE))))).
		Query(s.db, &results)
	if err != nil {
		log.Fatal().Err(err).Msgf("unable to query marketplaces in system %s", systemId)
	}

	return results
}
