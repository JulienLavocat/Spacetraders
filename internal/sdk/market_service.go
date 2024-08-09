package sdk

import (
	"database/sql"
	"time"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/julienlavocat/spacetraders/.gen/spacetraders/public/model"
	. "github.com/julienlavocat/spacetraders/.gen/spacetraders/public/table"
	"github.com/julienlavocat/spacetraders/internal/api"
	"github.com/julienlavocat/spacetraders/internal/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Market struct {
	db     *sql.DB
	logger zerolog.Logger
}

type SellPlan struct {
	ToSell   Cargo  `json:"toSell"`
	Location string `json:"location"`
}

func NewMarket(db *sql.DB) *Market {
	return &Market{
		db:     db,
		logger: log.With().Str("component", "Market").Logger(),
	}
}

func (m *Market) CreateSellPlan(systemId string, cargo Cargo) []SellPlan {
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

// func (m *Market) CreateSellPlanV2(systemId string, cargo Cargo) []SellPlan {
// 	var products []Expression
// 	for product := range cargo {
// 		products = append(products, String(product))
// 	}
//
// 	query := Waypoints.SELECT(Waypoints.ID, WaypointsProducts.ProductID, WaypointsProducts.Sell).
// 		FROM(WaypointsProducts.INNER_JOIN(Waypoints, Waypoints.ID.EQ(WaypointsProducts.WaypointID))).
// 		WHERE(Waypoints.SystemID.EQ(String(systemId)).
// 			AND(WaypointsProducts.ProductID.IN(products...)).
// 			AND(WaypointsProducts.Export.EQ(Bool(false)))).
// 		ORDER_BY(Raw("sell").DESC())
//      return []
// }

func (m *Market) UpdateMarket(data api.Market) {
	query := WaypointsProducts.INSERT(WaypointsProducts.AllColumns).
		ON_CONFLICT(WaypointsProducts.WaypointID, WaypointsProducts.ProductID, WaypointsProducts.Export, WaypointsProducts.Import, WaypointsProducts.Exchange).
		DO_UPDATE(SET(
			WaypointsProducts.Export.SET(WaypointsProducts.EXCLUDED.Export),
			WaypointsProducts.Import.SET(WaypointsProducts.EXCLUDED.Import),
			WaypointsProducts.Exchange.SET(WaypointsProducts.EXCLUDED.Exchange),
			WaypointsProducts.Volume.SET(WaypointsProducts.EXCLUDED.Volume),
			WaypointsProducts.Supply.SET(WaypointsProducts.EXCLUDED.Supply),
			WaypointsProducts.Activity.SET(WaypointsProducts.EXCLUDED.Activity),
			WaypointsProducts.Buy.SET(WaypointsProducts.EXCLUDED.Buy),
			WaypointsProducts.Sell.SET(WaypointsProducts.EXCLUDED.Sell),
			WaypointsProducts.UpdatedAt.SET(WaypointsProducts.EXCLUDED.UpdatedAt),
		))

	updateTime := time.Now().UTC()

	for _, product := range data.TradeGoods {
		query.MODEL(model.WaypointsProducts{
			WaypointID: data.Symbol,
			ProductID:  string(product.Symbol),
			Export:     product.Type == "EXPORT",
			Import:     product.Type == "IMPORT",
			Exchange:   product.Type == "EXCHANGE",
			Volume:     &product.TradeVolume,
			Supply:     (*string)(&product.Supply),
			Activity:   (*string)(product.Activity),
			Buy:        &product.PurchasePrice,
			Sell:       &product.SellPrice,
			UpdatedAt:  &updateTime,
		})
	}

	res, err := query.Exec(m.db)
	if err != nil {
		m.logger.Println(query.DebugSql())
		m.logger.Fatal().Err(err).Msgf("unable to update products for waypoint %s", data.Symbol)
	}
	affectedRows, _ := res.RowsAffected()
	m.logger.Info().Msgf("updated %d products in waypoints %s", affectedRows, data.Symbol)
}
