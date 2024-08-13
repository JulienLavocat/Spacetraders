package sdk

import (
	"database/sql"
	"errors"
	"fmt"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Navigation struct {
	db     *sql.DB
	logger zerolog.Logger
}

type PathSegment struct {
	To      string `json:"to"`
	Fuel    int32  `json:"fuel"`
	AggCost int32  `json:"aggCost"`
}

func NewNavigation(db *sql.DB) *Navigation {
	return &Navigation{
		db:     db,
		logger: log.With().Str("component", "Navigation").Logger(),
	}
}

// Plot a route between two waypoints in a system. Navigation is guaranteed to go through a network of fuel stations (so it can safely be refueled) except for the first and last waypoints
func (n *Navigation) PlotRoute(systemId, origin, destination string, currentFuel int32) ([]PathSegment, error) {
	// TODO: Run a first path from origin -> destination using fuelConstraint, (effectively placing the ship on the network on the first step)
	// then, do another run from the first station of the network to the destination using the ship's MAX FUEL
	// the path returned will be: origin -> first station from first run -> second run)
	// If no path can't be found at the first run (to join the stations networks), drift to the nearest one
	// If not path can't be found for the next stops, get as close as possible and also drift until there
	path, err := n.pgrDjikstra(systemId, origin, destination, currentFuel)
	if err != nil {
		return nil, err
	}

	if len(path) == 0 && destination != origin {
		return nil, errors.New("no route found")
	}

	return path, err
}

func (n *Navigation) EstimateRouteFuelCost(systemId, origin, destination string, currentFuel int32) (int32, error) {
	path, err := n.PlotRoute(systemId, origin, destination, currentFuel)
	if err != nil {
		return 0, err
	}

	costs := int32(0)
	for i := range path {
		costs += path[i].Fuel
	}

	return costs, nil
}

func (n *Navigation) pgrDjikstra(systemId, origin, destination string, fuelConstraint int32) ([]PathSegment, error) {
	q := RawStatement(fmt.Sprintf(`
       SELECT id AS "PathSegment.to", cost AS "PathSegment.fuel", agg_cost AS "PathSegment.agg_cost"
FROM pgr_dijkstra('
       SELECT wg.id, source, target, cost, product_id
FROM waypoints_graphs wg
         INNER JOIN waypoints w ON w.gid = wg.source
         FULL JOIN waypoints_products wp on w.id = wp.waypoint_id
WHERE w.system_id = ''%[1]s''
  AND cost <= %[2]d
  AND (wp.product_id = ''FUEL'' OR w.id = ''%[3]s'' OR w.id = ''%[4]s'');',
                  (SELECT gid FROM waypoints WHERE id = '%[3]s'),
                  (SELECT gid FROM waypoints WHERE id = '%[4]s'))
         INNER JOIN waypoints ON gid = node
ORDER BY path_seq`, systemId, fuelConstraint, origin, destination))

	var results []PathSegment
	err := q.Query(n.db, &results)
	if err != nil {
		n.logger.Err(err).Msgf("unable to run djikstra algorithm in %s from %s to %s and fuel constraint %d", systemId, origin, destination, fuelConstraint)
		return nil, err
	}

	return results, nil
}
