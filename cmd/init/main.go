package main

import (
	"context"
	"database/sql"
	"math"
	"os"
	"time"

	"github.com/julienlavocat/spacetraders/.gen/spacetraders/public/model"
	. "github.com/julienlavocat/spacetraders/.gen/spacetraders/public/table"
	"github.com/julienlavocat/spacetraders/internal/api"
	"github.com/julienlavocat/spacetraders/internal/utils"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	PAGE_SIZE      = 20
	SLEEP_DURATION = time.Duration(550) * time.Millisecond
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	db, err := sql.Open("postgres", "postgresql://spacetraders@localhost:5432/spacetraders?sslmode=disable")
	if err != nil {
		log.Fatal().Err(err).Msg("database connection failed")
	}
	defer db.Close()

	insertTraits(db)
	insertProducts(db)
	insertFactions(db)
	insertModifiers(db)

	insertSystems(db)
	insertWaypoints(db)
}

func insertSystems(db *sql.DB) {
	client := api.NewAPIClient(api.NewConfiguration())

	maxPage := math.MaxInt32
	page := 1
	for page < maxPage {
		log.Info().Msgf("fetching page %d/%d, systems", page, maxPage)

		insertSystems := Systems.INSERT(Systems.AllColumns).ON_CONFLICT(Systems.ID).DO_NOTHING()
		insertFactionsSystems := FactionsSystems.INSERT(FactionsSystems.AllColumns).ON_CONFLICT(FactionsSystems.AllColumns...).DO_NOTHING()

		res := utils.RetryRequest(client.SystemsAPI.GetSystems(context.Background()).Page(int32(page)).Limit(20).Execute, log.Logger, "failed to fetch systems at page %d", page)

		maxPage = int(math.Ceil(float64(res.Meta.Total)/float64(PAGE_SIZE))) + 1
		page++

		if len(res.Data) == 0 {
			continue
		}

		hasFactions := false
		for _, system := range res.Data {
			insertSystems.MODEL(model.Systems{
				ID:       system.Symbol,
				SectorID: system.SectorSymbol,
				Type:     string(system.Type),
				X:        system.X,
				Y:        system.Y,
			})

			for _, sf := range system.Factions {
				hasFactions = true
				insertFactionsSystems.MODEL(model.FactionsSystems{
					FactionID: string(sf.Symbol),
					SystemID:  system.Symbol,
				})
			}
		}

		insertedSystems, err := insertSystems.Exec(db)
		if err != nil {
			log.Fatal().Err(err).Msg("unable to insert systems")
		}
		insertCount, _ := insertedSystems.RowsAffected()
		log.Info().Msgf("inserted %d systems", insertCount)

		if hasFactions {
			insertSystemsFactionsRes, err := insertFactionsSystems.Exec(db)
			if err != nil {
				log.Fatal().Err(err).Msg("unable to insert systems_factions")
			}
			insertedCount, _ := insertSystemsFactionsRes.RowsAffected()
			log.Info().Msgf("inserted %d systems_factions", insertedCount)
		}

	}
}

func insertWaypoints(db *sql.DB) {
	client := api.NewAPIClient(api.NewConfiguration())

	var systems []model.Systems
	if err := Systems.SELECT(Systems.ID).ORDER_BY(Systems.ID.DESC()).Query(db, &systems); err != nil {
		log.Fatal().Err(err).Msg("unable to query systems")
	}

	for i, system := range systems {

		log.Info().Msgf("fetching system %s %d/%d", system.ID, i, len(systems))

		insertWaypoint := Waypoints.INSERT(Waypoints.AllColumns).ON_CONFLICT(Waypoints.ID).DO_NOTHING()
		insertWaypointTraits := WaypointsTraits.INSERT(WaypointsTraits.AllColumns).ON_CONFLICT(WaypointsTraits.AllColumns...).DO_NOTHING()
		insertWaypointModifiers := WaypointsModifiers.INSERT(WaypointsModifiers.AllColumns).ON_CONFLICT(WaypointsModifiers.AllColumns...).DO_NOTHING()
		insertWaypointsProducts := WaypointsProducts.INSERT(WaypointsProducts.AllColumns).ON_CONFLICT(WaypointsProducts.AllColumns...).DO_NOTHING()
		insertWaypointsFactions := FactionsWaypoints.INSERT(FactionsWaypoints.AllColumns).ON_CONFLICT(FactionsWaypoints.AllColumns...).DO_NOTHING()

		maxPage := math.MaxInt
		page := 1
		hasWaypoints := false
		hasTraits := false
		hasModifiers := false
		hasProducts := false
		hasFaction := false
		for page < maxPage {
			res := utils.RetryRequest(client.SystemsAPI.GetSystemWaypoints(context.Background(), system.ID).Page(int32(page)).Limit(PAGE_SIZE).Execute, log.Logger, "failed to fetch waypoints for system %s", system.ID)

			hasWaypoints = len(res.Data) > 0
			log.Info().Msgf("fetching page %d of system %s, %d waypoints found", page, system.ID, len(res.Data))

			maxPage = int(math.Ceil(float64(res.Meta.Total)/float64(PAGE_SIZE))) + 1
			page++

			for _, waypoint := range res.Data {
				waypointModel := model.Waypoints{
					ID:                waypoint.Symbol,
					Type:              string(waypoint.Type),
					X:                 waypoint.X,
					Y:                 waypoint.Y,
					SystemID:          waypoint.SystemSymbol,
					Orbits:            waypoint.Orbits,
					UnderConstruction: waypoint.IsUnderConstruction,
				}

				if waypoint.Faction != nil {
					waypointModel.Faction = string(waypoint.Faction.Symbol)
				}

				if waypoint.Chart != nil {
					waypointModel.SubmittedOn = waypoint.Chart.SubmittedOn
					waypointModel.SubmittedBy = waypoint.Chart.SubmittedBy
				}

				insertWaypoint.MODEL(waypointModel)

				for _, trait := range waypoint.Traits {
					hasTraits = true
					insertWaypointTraits.MODEL(model.WaypointsTraits{
						WaypointID: waypoint.Symbol,
						TraitID:    string(trait.Symbol),
					})

					if trait.Symbol == api.MARKETPLACE {
						res := utils.RetryRequest(client.SystemsAPI.GetMarket(context.Background(), waypoint.SystemSymbol, waypoint.Symbol).Execute, log.Logger, "unable to fetch market for waypoint %s", waypoint.Symbol)

						if data, ok := res.GetDataOk(); ok {
							if exports, ok := data.GetExportsOk(); ok {
								for _, product := range exports {
									hasProducts = true
									insertWaypointsProducts.MODEL(model.WaypointsProducts{
										WaypointID: waypoint.Symbol,
										ProductID:  string(product.Symbol),
										Export:     true,
									})
								}
							}

							if imports, ok := data.GetImportsOk(); ok {
								for _, product := range imports {
									hasProducts = true
									insertWaypointsProducts.MODEL(model.WaypointsProducts{
										WaypointID: waypoint.Symbol,
										ProductID:  string(product.Symbol),
										Import:     true,
									})
								}
							}

							if exchange, ok := data.GetExchangeOk(); ok {
								for _, product := range exchange {
									hasProducts = true
									insertWaypointsProducts.MODEL(model.WaypointsProducts{
										WaypointID: waypoint.Symbol,
										ProductID:  string(product.Symbol),
										Exchange:   true,
									})
								}
							}
						}

					}

					if trait.Symbol == api.SHIPYARD {
						res := utils.RetryRequest(client.SystemsAPI.GetShipyard(context.Background(), waypoint.SystemSymbol, waypoint.Symbol).Execute, log.Logger, "unable to fetch shipyard for waypoint %s", waypoint.Symbol)
						for _, product := range res.Data.ShipTypes {
							hasProducts = true
							insertWaypointsProducts.MODEL(model.WaypointsProducts{
								WaypointID: waypoint.Symbol,
								ProductID:  string(product.Type),
								Export:     true,
							})
						}
					}
				}

				for _, modifier := range waypoint.Modifiers {
					hasModifiers = true
					insertWaypointModifiers.MODEL(model.WaypointsModifiers{
						ModifierID: string(modifier.Symbol),
						WaypointID: waypoint.Symbol,
					})
				}

				if faction, ok := waypoint.GetFactionOk(); ok {
					hasFaction = true
					insertWaypointsFactions.MODEL(model.FactionsWaypoints{
						WaypointID: waypoint.Symbol,
						FactionID:  string(faction.Symbol),
					})
				}
			}
		}

		if hasWaypoints {
			res, err := insertWaypoint.Exec(db)
			if err != nil {
				log.Fatal().Err(err).Msg("unable to insert waypoints")
			}
			insertedCount, _ := res.RowsAffected()
			log.Info().Msgf("inserted %d waypoints", insertedCount)
		}

		if hasTraits {
			res, err := insertWaypointTraits.Exec(db)
			if err != nil {
				log.Fatal().Err(err).Msg("unable to insert waypoints_traits")
			}
			insertedCount, _ := res.RowsAffected()
			log.Info().Msgf("inserted %d waypoints_traits", insertedCount)
		}

		if hasModifiers {
			res, err := insertWaypointModifiers.Exec(db)
			if err != nil {
				log.Fatal().Err(err).Msg("unable to insert waypoints_modifiers")
			}
			insertedCount, _ := res.RowsAffected()
			log.Info().Msgf("inserted %d waypoints_modifiers", insertedCount)
		}

		if hasProducts {
			res, err := insertWaypointsProducts.Exec(db)
			if err != nil {
				log.Fatal().Err(err).Msg("unable to insert waypoints_products")
			}
			insertedCount, _ := res.RowsAffected()
			log.Info().Msgf("inserted %d waypoints_products", insertedCount)
		}

		if hasFaction {
			res, err := insertWaypointsFactions.Exec(db)
			if err != nil {
				log.Fatal().Err(err).Msg("unable to insert factions_waypoints")
			}
			insertedCount, _ := res.RowsAffected()
			log.Info().Msgf("inserted %d factions_waypoints", insertedCount)
		}
	}
}

func insertTraits(db *sql.DB) {
	stmt := Traits.INSERT(Traits.ID).ON_CONFLICT(Traits.ID).DO_NOTHING()

	for _, trait := range api.AllowedWaypointTraitSymbolEnumValues {
		stmt.VALUES(string(trait))
	}

	if res, err := stmt.Exec(db); err != nil {
		log.Fatal().Err(err).Msg("unable to insert traits")
	} else {
		rows, _ := res.RowsAffected()
		log.Info().Msgf("inserted %d traits", rows)
	}
}

func insertFactions(db *sql.DB) {
	stmt := Factions.INSERT(Factions.ID).ON_CONFLICT(Factions.ID).DO_NOTHING()

	for _, faction := range api.AllowedFactionSymbolEnumValues {
		stmt.VALUES(string(faction))
	}

	if res, err := stmt.Exec(db); err != nil {
		log.Fatal().Err(err).Msg("unable to insert factions")
	} else {
		rows, _ := res.RowsAffected()
		log.Info().Msgf("inserted %d factions", rows)
	}
}

func insertProducts(db *sql.DB) {
	stmt := Products.INSERT(Products.ID).ON_CONFLICT(Products.ID).DO_NOTHING()

	for _, trait := range api.AllowedTradeSymbolEnumValues {
		stmt.VALUES(string(trait))
	}

	if res, err := stmt.Exec(db); err != nil {
		log.Fatal().Err(err).Msg("unable to insert products")
	} else {
		rows, _ := res.RowsAffected()
		log.Info().Msgf("inserted %d products", rows)
	}
}

func insertModifiers(db *sql.DB) {
	stmt := Modifiers.INSERT(Modifiers.ID).ON_CONFLICT(Modifiers.ID).DO_NOTHING()

	for _, modifier := range api.AllowedWaypointModifierSymbolEnumValues {
		stmt.VALUES(string(modifier))
	}

	if res, err := stmt.Exec(db); err != nil {
		log.Fatal().Err(err).Msg("unable to insert modifiers")
	} else {
		rows, _ := res.RowsAffected()
		log.Info().Msgf("inserted %d modifiers", rows)
	}
}
