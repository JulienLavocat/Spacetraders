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

var Waypoints = newWaypointsTable("public", "waypoints", "")

type waypointsTable struct {
	postgres.Table

	// Columns
	ID                postgres.ColumnString
	SystemID          postgres.ColumnString
	X                 postgres.ColumnInteger
	Y                 postgres.ColumnInteger
	Type              postgres.ColumnString
	Faction           postgres.ColumnString
	Orbits            postgres.ColumnString
	UnderConstruction postgres.ColumnBool
	SubmittedOn       postgres.ColumnTimestampz
	SubmittedBy       postgres.ColumnString
	Geom              postgres.ColumnString
	Gid               postgres.ColumnInteger

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type WaypointsTable struct {
	waypointsTable

	EXCLUDED waypointsTable
}

// AS creates new WaypointsTable with assigned alias
func (a WaypointsTable) AS(alias string) *WaypointsTable {
	return newWaypointsTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new WaypointsTable with assigned schema name
func (a WaypointsTable) FromSchema(schemaName string) *WaypointsTable {
	return newWaypointsTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new WaypointsTable with assigned table prefix
func (a WaypointsTable) WithPrefix(prefix string) *WaypointsTable {
	return newWaypointsTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new WaypointsTable with assigned table suffix
func (a WaypointsTable) WithSuffix(suffix string) *WaypointsTable {
	return newWaypointsTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newWaypointsTable(schemaName, tableName, alias string) *WaypointsTable {
	return &WaypointsTable{
		waypointsTable: newWaypointsTableImpl(schemaName, tableName, alias),
		EXCLUDED:       newWaypointsTableImpl("", "excluded", ""),
	}
}

func newWaypointsTableImpl(schemaName, tableName, alias string) waypointsTable {
	var (
		IDColumn                = postgres.StringColumn("id")
		SystemIDColumn          = postgres.StringColumn("system_id")
		XColumn                 = postgres.IntegerColumn("x")
		YColumn                 = postgres.IntegerColumn("y")
		TypeColumn              = postgres.StringColumn("type")
		FactionColumn           = postgres.StringColumn("faction")
		OrbitsColumn            = postgres.StringColumn("orbits")
		UnderConstructionColumn = postgres.BoolColumn("under_construction")
		SubmittedOnColumn       = postgres.TimestampzColumn("submitted_on")
		SubmittedByColumn       = postgres.StringColumn("submitted_by")
		GeomColumn              = postgres.StringColumn("geom")
		GidColumn               = postgres.IntegerColumn("gid")
		allColumns              = postgres.ColumnList{IDColumn, SystemIDColumn, XColumn, YColumn, TypeColumn, FactionColumn, OrbitsColumn, UnderConstructionColumn, SubmittedOnColumn, SubmittedByColumn, GeomColumn, GidColumn}
		mutableColumns          = postgres.ColumnList{SystemIDColumn, XColumn, YColumn, TypeColumn, FactionColumn, OrbitsColumn, UnderConstructionColumn, SubmittedOnColumn, SubmittedByColumn, GeomColumn, GidColumn}
	)

	return waypointsTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:                IDColumn,
		SystemID:          SystemIDColumn,
		X:                 XColumn,
		Y:                 YColumn,
		Type:              TypeColumn,
		Faction:           FactionColumn,
		Orbits:            OrbitsColumn,
		UnderConstruction: UnderConstructionColumn,
		SubmittedOn:       SubmittedOnColumn,
		SubmittedBy:       SubmittedByColumn,
		Geom:              GeomColumn,
		Gid:               GidColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
