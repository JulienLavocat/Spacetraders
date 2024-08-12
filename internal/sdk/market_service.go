package sdk

import (
	"cmp"
	"database/sql"
	"slices"
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
	sdk    *Sdk
	logger zerolog.Logger
}

type SellPlan struct {
	ToSell   Cargo  `json:"toSell"`
	Location string `json:"location"`
}

type OpportunityRow struct {
	UpdatedAt time.Time
	ID        string
	Product   string
	Price     int32
	Volume    int32
	X         int32
	Y         int32
}

type TradeRoute struct {
	Product                 string
	BuyAt                   string
	SellAt                  string
	MaxAmount               int32
	SellPrice               int32
	BuyPrice                int32
	EstimatedProfits        int32
	EstimatedProfitsPerUnit int32
	FuelCost                int32
}

func NewMarket(db *sql.DB, sdk *Sdk) *Market {
	return &Market{
		db:     db,
		logger: log.With().Str("component", "Market").Logger(),
		sdk:    sdk,
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

	if len(data.TradeGoods) == 0 {
		return
	}

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
	m.logger.Info().Msgf("updated %d market products in waypoints %s", affectedRows, data.Symbol)
}

func (m *Market) UpdateShipyard(waypoint string, data []api.ShipyardShip) {
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

	if len(data) == 0 {
		return
	}

	for _, product := range data {
		query.MODEL(model.WaypointsProducts{
			WaypointID: waypoint,
			ProductID:  string(product.Type),
			Export:     true,
			Import:     false,
			Exchange:   false,
			Volume:     nil,
			Supply:     (*string)(&product.Supply),
			Activity:   (*string)(product.Activity),
			Buy:        &product.PurchasePrice,
			Sell:       nil,
			UpdatedAt:  &updateTime,
		})
	}

	res, err := query.Exec(m.db)
	if err != nil {
		m.logger.Println(query.DebugSql())
		m.logger.Fatal().Err(err).Msgf("unable to update products for waypoint %s", waypoint)
	}
	affectedRows, _ := res.RowsAffected()
	m.logger.Info().Msgf("updated %d shipyard products in waypoints %s", affectedRows, waypoint)
}

func (m *Market) GetTradeRoutes(systemId string) []*TradeRoute {
	buyQuery := WaypointsProducts.
		INNER_JOIN(Waypoints, Waypoints.ID.EQ(WaypointsProducts.WaypointID)).
		SELECT(
			Waypoints.ID.AS("OpportunityRow.id"),
			Waypoints.X.AS("OpportunityRow.x"),
			Waypoints.Y.AS("OpportunityRow.y"),
			WaypointsProducts.ProductID.AS("OpportunityRow.product"),
			WaypointsProducts.Buy.AS("OpportunityRow.price"),
			WaypointsProducts.UpdatedAt.AS("OpportunityRow.updated_at"),
			WaypointsProducts.Volume.AS("OpportunityRow.volume")).
		WHERE(Waypoints.SystemID.EQ(String(systemId)).
			AND(WaypointsProducts.Exchange.EQ(Bool(true)).
				OR(WaypointsProducts.Export.EQ(Bool(true)))).
			AND(WaypointsProducts.UpdatedAt.IS_NOT_NULL())).
		ORDER_BY(WaypointsProducts.Buy)

	var buyResults []OpportunityRow
	err := buyQuery.Query(m.db, &buyResults)
	if err != nil {
		log.Fatal().Err(err).Str("query", buyQuery.DebugSql()).Msgf("unable to query buy opportunities in %s", systemId)
	}

	sellQuery := WaypointsProducts.
		INNER_JOIN(Waypoints, Waypoints.ID.EQ(WaypointsProducts.WaypointID)).
		SELECT(
			Waypoints.ID.AS("OpportunityRow.id"),
			Waypoints.X.AS("OpportunityRow.x"),
			Waypoints.Y.AS("OpportunityRow.y"),
			WaypointsProducts.ProductID.AS("OpportunityRow.product"),
			WaypointsProducts.Sell.AS("OpportunityRow.price"),
			WaypointsProducts.UpdatedAt.AS("OpportunityRow.updated_at"),
			WaypointsProducts.Volume.AS("OpportunityRow.volume")).
		WHERE(Waypoints.SystemID.EQ(String(systemId)).
			AND(WaypointsProducts.Exchange.EQ(Bool(true)).
				OR(WaypointsProducts.Import.EQ(Bool(true)))).
			AND(WaypointsProducts.UpdatedAt.IS_NOT_NULL())).
		ORDER_BY(WaypointsProducts.Sell.DESC())

	var sellResults []OpportunityRow
	err = sellQuery.Query(m.db, &sellResults)
	if err != nil {
		log.Fatal().Err(err).Str("query", sellQuery.DebugSql()).Msgf("unable to query sell opportunities in %s", systemId)
	}

	buyOpportuinities := map[string]OpportunityRow{}
	for _, newRow := range buyResults {
		if currentRow, ok := buyOpportuinities[newRow.Product]; ok {
			if newRow.Price < currentRow.Price {
				buyOpportuinities[newRow.Product] = newRow
			}
		} else {
			buyOpportuinities[newRow.Product] = newRow
		}
	}

	sellOpportuinities := map[string]OpportunityRow{}
	for _, newRow := range sellResults {
		if currentRow, ok := sellOpportuinities[newRow.Product]; ok {
			if newRow.Price > currentRow.Price {
				sellOpportuinities[newRow.Product] = newRow
			}
		} else {
			sellOpportuinities[newRow.Product] = newRow
		}
	}

	// TODO: Calculate fuel cost for each routes
	var tradeRoutes []*TradeRoute
	for product, buyOpportunity := range buyOpportuinities {
		sellOpportunity, ok := sellOpportuinities[product]
		if !ok || sellOpportunity.ID == buyOpportunity.ID {
			continue
		}

		maxAmount := min(buyOpportunity.Volume, sellOpportunity.Volume)

		tradeRoutes = append(tradeRoutes, &TradeRoute{
			Product:                 product,
			SellPrice:               sellOpportunity.Price,
			BuyPrice:                buyOpportunity.Price,
			SellAt:                  sellOpportunity.ID,
			BuyAt:                   buyOpportunity.ID,
			MaxAmount:               maxAmount,
			EstimatedProfits:        (sellOpportunity.Price - buyOpportunity.Price) * maxAmount,
			EstimatedProfitsPerUnit: sellOpportunity.Price - buyOpportunity.Price,
			FuelCost:                GetFuelCost(buyOpportunity.X, buyOpportunity.Y, sellOpportunity.X, sellOpportunity.Y),
		})
	}

	slices.SortFunc(tradeRoutes, func(a, b *TradeRoute) int {
		return cmp.Compare(b.EstimatedProfits, a.EstimatedProfits)
	})

	return tradeRoutes
}

func (m *Market) HasShipyard(waypointId string) (bool, error) {
	q := Waypoints.INNER_JOIN(WaypointsTraits, Waypoints.ID.EQ(WaypointsTraits.WaypointID)).
		SELECT(Waypoints.ID).
		WHERE(Waypoints.ID.EQ(String(waypointId)).AND(WaypointsTraits.TraitID.EQ(String(string(api.SHIPYARD)))))

	var results []model.Waypoints
	err := q.Query(m.db, &results)
	if err != nil {
		m.logger.Err(err).Msgf("unable to check if waypoints has SHIPYARD trait at %s", waypointId)
		return false, err
	}

	return len(results) > 0, nil
}

func (m *Market) FindLowestProductPrice(systemId, product string) (bool, string, int32, error) {
	var result struct {
		Id    string
		Price int32
	}
	q := Waypoints.INNER_JOIN(WaypointsProducts, Waypoints.ID.EQ(WaypointsProducts.WaypointID)).
		SELECT(Waypoints.ID.AS("id"), WaypointsProducts.Buy.AS("price")).
		WHERE(WaypointsProducts.ProductID.EQ(String(product))).
		ORDER_BY(WaypointsProducts.Buy).LIMIT(1)

	err := q.Query(m.db, &result)
	if err != nil {
		m.logger.Error().Err(err).Msgf("unable to query market for %s at %s", product, systemId)
		return false, "", 0, err
	}

	if result.Id == "" {
		return false, "", 0, nil
	}

	return result.Id != "", result.Id, result.Price, nil
}

func (m *Market) ReportTransaction(tx api.MarketTransaction, balance int64, correlationId string) error {
	_, err := Transactions.INSERT(Transactions.AllColumns.Except(Transactions.ID, Transactions.Timestamp)).MODEL(model.Transactions{
		Waypoint:      tx.WaypointSymbol,
		Ship:          tx.ShipSymbol,
		Product:       tx.TradeSymbol,
		Type:          tx.Type,
		Amount:        tx.Units,
		TotalPrice:    tx.TotalPrice,
		PricePerUnit:  tx.PricePerUnit,
		AgentBalance:  balance,
		CorrelationID: &correlationId,
	}).Exec(m.db)

	return err
}
