//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

// UseSchema sets a new schema name for all generated table SQL builder types. It is recommended to invoke
// this method only once at the beginning of the program.
func UseSchema(schema string) {
	Factions = Factions.FromSchema(schema)
	FactionsSystems = FactionsSystems.FromSchema(schema)
	FactionsWaypoints = FactionsWaypoints.FromSchema(schema)
	MarketProbes = MarketProbes.FromSchema(schema)
	Modifiers = Modifiers.FromSchema(schema)
	Products = Products.FromSchema(schema)
	Ships = Ships.FromSchema(schema)
	ShipsCargo = ShipsCargo.FromSchema(schema)
	Systems = Systems.FromSchema(schema)
	TradingFleets = TradingFleets.FromSchema(schema)
	Traits = Traits.FromSchema(schema)
	Transactions = Transactions.FromSchema(schema)
	Waypoints = Waypoints.FromSchema(schema)
	WaypointsGraphs = WaypointsGraphs.FromSchema(schema)
	WaypointsModifiers = WaypointsModifiers.FromSchema(schema)
	WaypointsProducts = WaypointsProducts.FromSchema(schema)
	WaypointsTraits = WaypointsTraits.FromSchema(schema)
}
