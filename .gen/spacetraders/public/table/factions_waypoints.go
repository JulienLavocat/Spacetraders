//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/postgres"
)

var FactionsWaypoints = newFactionsWaypointsTable("public", "factions_waypoints", "")

type factionsWaypointsTable struct {
	postgres.Table

	// Columns
	WaypointID postgres.ColumnString
	FactionID  postgres.ColumnString

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type FactionsWaypointsTable struct {
	factionsWaypointsTable

	EXCLUDED factionsWaypointsTable
}

// AS creates new FactionsWaypointsTable with assigned alias
func (a FactionsWaypointsTable) AS(alias string) *FactionsWaypointsTable {
	return newFactionsWaypointsTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new FactionsWaypointsTable with assigned schema name
func (a FactionsWaypointsTable) FromSchema(schemaName string) *FactionsWaypointsTable {
	return newFactionsWaypointsTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new FactionsWaypointsTable with assigned table prefix
func (a FactionsWaypointsTable) WithPrefix(prefix string) *FactionsWaypointsTable {
	return newFactionsWaypointsTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new FactionsWaypointsTable with assigned table suffix
func (a FactionsWaypointsTable) WithSuffix(suffix string) *FactionsWaypointsTable {
	return newFactionsWaypointsTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newFactionsWaypointsTable(schemaName, tableName, alias string) *FactionsWaypointsTable {
	return &FactionsWaypointsTable{
		factionsWaypointsTable: newFactionsWaypointsTableImpl(schemaName, tableName, alias),
		EXCLUDED:               newFactionsWaypointsTableImpl("", "excluded", ""),
	}
}

func newFactionsWaypointsTableImpl(schemaName, tableName, alias string) factionsWaypointsTable {
	var (
		WaypointIDColumn = postgres.StringColumn("waypoint_id")
		FactionIDColumn  = postgres.StringColumn("faction_id")
		allColumns       = postgres.ColumnList{WaypointIDColumn, FactionIDColumn}
		mutableColumns   = postgres.ColumnList{}
	)

	return factionsWaypointsTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		WaypointID: WaypointIDColumn,
		FactionID:  FactionIDColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
