package sdk

import (
	. "github.com/go-jet/jet/v2/postgres"
	. "github.com/julienlavocat/spacetraders/.gen/spacetraders/public/table"
)

type Market struct{}

func (m *Market) SellCargoTo(systemId string, cargo map[string]int32) {
	// SELECT id, ARRAY_AGG(product_id) AS products FROM waypoints
	//     INNER JOIN waypoints_products ON waypoints.id = waypoints_products.waypoint_id
	//          WHERE system_id = 'X1-QA42'
	//            AND product_id IN ('ALUMINUM_ORE', 'IRON_ORE', 'COPPER_ORE', 'SILICON_CRYSTALS') AND export = false
	// GROUP BY id

	var products []Expression
	for product := range cargo {
		products = append(products, String(product))
	}

	query := Waypoints.SELECT(Waypoints.ID, Raw("ARRAY_AGG(product_id)").AS("products")).
		FROM(WaypointsProducts.INNER_JOIN(Waypoints, Waypoints.ID.EQ(WaypointsProducts.WaypointID))).
		WHERE(Waypoints.SystemID.EQ(String(systemId)).
			AND(WaypointsProducts.ProductID.IN(products...)).
			AND(WaypointsProducts.Export.EQ(Bool(false)))).
		GROUP_BY(Waypoints.ID)

	println(query.DebugSql())
}
