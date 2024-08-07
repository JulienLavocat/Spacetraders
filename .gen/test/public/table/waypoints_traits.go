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

var WaypointsTraits = newWaypointsTraitsTable("public", "waypoints_traits", "")

type waypointsTraitsTable struct {
	postgres.Table

	// Columns
	WaypointID postgres.ColumnString
	TraitID    postgres.ColumnString

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type WaypointsTraitsTable struct {
	waypointsTraitsTable

	EXCLUDED waypointsTraitsTable
}

// AS creates new WaypointsTraitsTable with assigned alias
func (a WaypointsTraitsTable) AS(alias string) *WaypointsTraitsTable {
	return newWaypointsTraitsTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new WaypointsTraitsTable with assigned schema name
func (a WaypointsTraitsTable) FromSchema(schemaName string) *WaypointsTraitsTable {
	return newWaypointsTraitsTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new WaypointsTraitsTable with assigned table prefix
func (a WaypointsTraitsTable) WithPrefix(prefix string) *WaypointsTraitsTable {
	return newWaypointsTraitsTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new WaypointsTraitsTable with assigned table suffix
func (a WaypointsTraitsTable) WithSuffix(suffix string) *WaypointsTraitsTable {
	return newWaypointsTraitsTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newWaypointsTraitsTable(schemaName, tableName, alias string) *WaypointsTraitsTable {
	return &WaypointsTraitsTable{
		waypointsTraitsTable: newWaypointsTraitsTableImpl(schemaName, tableName, alias),
		EXCLUDED:             newWaypointsTraitsTableImpl("", "excluded", ""),
	}
}

func newWaypointsTraitsTableImpl(schemaName, tableName, alias string) waypointsTraitsTable {
	var (
		WaypointIDColumn = postgres.StringColumn("waypoint_id")
		TraitIDColumn    = postgres.StringColumn("trait_id")
		allColumns       = postgres.ColumnList{WaypointIDColumn, TraitIDColumn}
		mutableColumns   = postgres.ColumnList{}
	)

	return waypointsTraitsTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		WaypointID: WaypointIDColumn,
		TraitID:    TraitIDColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
