package main

import (
	"database/sql"
	"fmt"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/julienlavocat/spacetraders/.gen/spacetraders/public/model"
	. "github.com/julienlavocat/spacetraders/.gen/spacetraders/public/table"
	"github.com/julienlavocat/spacetraders/internal/sdk"
	"github.com/julienlavocat/spacetraders/internal/utils"
	"github.com/rs/zerolog/log"
)

func initWaypointsGraph(db *sql.DB) {
	var systems []model.Systems
	err := Systems.SELECT(Systems.ID, Systems.X, Systems.Y).Query(db, &systems)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to load systems")
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal().Err(err).Msg("unable to acquire transaction from db")
	}

	for i, system := range systems {
		if system.ID != "X1-NT44" {
			continue
		}
		log.Info().Msgf("creating graph of system %s %d/%d", system.ID, i, len(systems))

		var waypoints []model.Waypoints
		qW := Waypoints.
			SELECT(Waypoints.ID, Waypoints.X, Waypoints.Y, Waypoints.Gid).
			WHERE(Waypoints.SystemID.EQ(String(system.ID)))
		err := qW.Query(db, &waypoints)
		if err != nil {
			log.Fatal().Err(err).Msg("unable to load waypoints")
		}

		log.Info().Msgf("found %d waypoints having FUEL as product", len(waypoints))

		q := WaypointsGraphs.INSERT(WaypointsGraphs.AllColumns.Except(WaypointsGraphs.ID))

		var edges []model.WaypointsGraphs
		seenEdges := utils.NewSet[string]()
		for _, source := range waypoints {
			for _, target := range waypoints {
				if source.Gid == target.Gid {
					continue
				}

				// Since it's an undirected graph, we skip if the edge was already added
				// if seenEdges.Has(fmt.Sprintf("%d-%d", source.Gid, target.Gid)) || seenEdges.Has(fmt.Sprintf("%d-%d", target.Gid, source.Gid)) {
				// 	continue
				// }

				edges = append(edges, model.WaypointsGraphs{
					SystemID: system.ID,
					Source:   source.Gid,
					Target:   target.Gid,
					Cost:     sdk.GetFuelCost(source.X, source.Y, target.X, target.Y),
				})

				seenEdges.Add(fmt.Sprintf("%d-%d", source.Gid, target.Gid))
				seenEdges.Add(fmt.Sprintf("%d-%d", target.Gid, source.Gid))
			}
		}

		if len(edges) == 0 {
			log.Info().Msg("system has no waypoints")
			continue
		}

		q.MODELS(edges)

		result, err := q.Exec(tx)
		if err != nil {
			log.Fatal().Err(err).Msgf("unable to insert waypoints_edges for system %s", system.ID)
		}

		count, _ := result.RowsAffected()
		log.Info().Msgf("insert %d waypoints_graphs", count)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal().Err(err).Msg("unable to commit waypoints_graphs transation")
	}
}
