package sdk

import (
	"database/sql"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/julienlavocat/spacetraders/.gen/spacetraders/public/model"
	. "github.com/julienlavocat/spacetraders/.gen/spacetraders/public/table"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"
)

type SystemMap struct {
	graph           *simple.WeightedUndirectedGraph
	waypointsToNode map[string]int64
}

type Navigation struct {
	db     *sql.DB
	charts map[string]*SystemMap
	logger zerolog.Logger
}

func NewNavigation(db *sql.DB) *Navigation {
	return &Navigation{
		db:     db,
		charts: make(map[string]*SystemMap),
		logger: log.With().Str("component", "Navigation").Logger(),
	}
}

func newSystemMap(waypoints []model.Waypoints) *SystemMap {
	graph := simple.NewWeightedUndirectedGraph(0, 15)
	waypointsToIds := make(map[string]int64)

	for id, waypoint := range waypoints {
		graph.AddNode(simple.Node(int64(id)))
		waypointsToIds[waypoint.ID] = int64(id)
	}

	for id1, w1 := range waypoints {
		for id2, w2 := range waypoints {
			if w1.ID == w2.ID {
				continue
			}

			fuelCost := GetFuelCost(w1.X, w1.Y, w2.X, w2.Y)
			graph.SetWeightedEdge(graph.NewWeightedEdge(simple.Node(id1), simple.Node(id2), float64(fuelCost)))
		}
	}

	return &SystemMap{
		graph:           graph,
		waypointsToNode: waypointsToIds,
	}
}

func (n *Navigation) GetGraphViz(system string) string {
	systemMap, ok := n.charts[system]
	if !ok {
		var waypoints []model.Waypoints
		err := Waypoints.SELECT(Waypoints.ID, Waypoints.X, Waypoints.Y).WHERE(Waypoints.SystemID.EQ(String(system))).Query(n.db, &waypoints)
		if err != nil {
			n.logger.Fatal().Err(err).Msgf("unable to query waypoints in system %s", system)
		}

		systemMap = newSystemMap(waypoints)
		n.charts[system] = systemMap
	}

	dotViz, err := dot.Marshal(systemMap.graph, "X1-QA42", "", "")
	if err != nil {
		n.logger.Fatal().Err(err).Msgf("unable to generate dotviz for system %s", system)
	}

	return string(dotViz)
}
