package sdk

import (
	"database/sql"

	. "github.com/go-jet/jet/v2/postgres"
	. "github.com/julienlavocat/spacetraders/.gen/spacetraders/public/table"
	"github.com/julienlavocat/spacetraders/internal/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Market struct {
	db     *sql.DB
	logger zerolog.Logger
}

type SellPlan struct {
	ToSell   Cargo
	Location string
}

func NewMarket(db *sql.DB) *Market {
	return &Market{
		db:     db,
		logger: log.With().Str("component", "Market").Logger(),
	}
}

func (m *Market) SellCargoTo(systemId string, cargo Cargo) []SellPlan {
	var products []Expression
	for product := range cargo {
		products = append(products, String(product))
	}

	query := Waypoints.SELECT(Waypoints.ID.AS("id"), Raw("ARRAY_AGG(product_id)").AS("products"), Raw("ARRAY_LENGTH(ARRAY_AGG(product_id), 1)").AS("count")).
		FROM(WaypointsProducts.INNER_JOIN(Waypoints, Waypoints.ID.EQ(WaypointsProducts.WaypointID))).
		WHERE(Waypoints.SystemID.EQ(String(systemId)).
			AND(WaypointsProducts.ProductID.IN(products...)).
			AND(WaypointsProducts.Export.EQ(Bool(false)))).
		GROUP_BY(Waypoints.ID).ORDER_BY(Raw("count").DESC())

	// TODO: Using market data, find the most appropriate destinations (i.e the destination where the total sell value PER UNIT (so total sell value / number of items sold) is the best)
	// For now, we just sorts by number of commodities that can be sold
	// Also, it would be ideal to take into account the amount of fuel used to get to the destination as a parameter

	var locations []struct {
		ID       string
		Products string
	}

	err := query.Query(m.db, &locations)
	if err != nil {
		log.Fatal().Err(err).Str("query", query.DebugSql()).Msg("unable to query waypoints")
	}

	log.Info().Interface("cargo", cargo).Interface("locations", locations).Msg("found these locations matching the cargo")

	var sellPlan []SellPlan
	productsSold := utils.NewSet[string]()

	for _, location := range locations {
		products := utils.NewSetFrom(utils.GetStringArrayFromSqlString(location.Products))

		toSellAtLocation := products.Difference(productsSold)

		productsSold = productsSold.Union(toSellAtLocation)

		if toSellAtLocation.Size() == 0 {
			continue
		}

		cargoToSell := make(Cargo)
		for _, v := range toSellAtLocation.Values() {
			cargoToSell[v] = cargo[v]
		}

		sellPlan = append(sellPlan, SellPlan{
			Location: location.ID,
			ToSell:   cargoToSell,
		})
	}

	log.Info().Interface("plan", sellPlan).Msg("established sell plan")

	return sellPlan
}
