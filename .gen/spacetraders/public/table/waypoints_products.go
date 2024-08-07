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

var WaypointsProducts = newWaypointsProductsTable("public", "waypoints_products", "")

type waypointsProductsTable struct {
	postgres.Table

	// Columns
	WaypointID postgres.ColumnString
	ProductID  postgres.ColumnString
	Export     postgres.ColumnBool
	Exchange   postgres.ColumnBool
	Import     postgres.ColumnBool

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type WaypointsProductsTable struct {
	waypointsProductsTable

	EXCLUDED waypointsProductsTable
}

// AS creates new WaypointsProductsTable with assigned alias
func (a WaypointsProductsTable) AS(alias string) *WaypointsProductsTable {
	return newWaypointsProductsTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new WaypointsProductsTable with assigned schema name
func (a WaypointsProductsTable) FromSchema(schemaName string) *WaypointsProductsTable {
	return newWaypointsProductsTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new WaypointsProductsTable with assigned table prefix
func (a WaypointsProductsTable) WithPrefix(prefix string) *WaypointsProductsTable {
	return newWaypointsProductsTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new WaypointsProductsTable with assigned table suffix
func (a WaypointsProductsTable) WithSuffix(suffix string) *WaypointsProductsTable {
	return newWaypointsProductsTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newWaypointsProductsTable(schemaName, tableName, alias string) *WaypointsProductsTable {
	return &WaypointsProductsTable{
		waypointsProductsTable: newWaypointsProductsTableImpl(schemaName, tableName, alias),
		EXCLUDED:               newWaypointsProductsTableImpl("", "excluded", ""),
	}
}

func newWaypointsProductsTableImpl(schemaName, tableName, alias string) waypointsProductsTable {
	var (
		WaypointIDColumn = postgres.StringColumn("waypoint_id")
		ProductIDColumn  = postgres.StringColumn("product_id")
		ExportColumn     = postgres.BoolColumn("export")
		ExchangeColumn   = postgres.BoolColumn("exchange")
		ImportColumn     = postgres.BoolColumn("import")
		allColumns       = postgres.ColumnList{WaypointIDColumn, ProductIDColumn, ExportColumn, ExchangeColumn, ImportColumn}
		mutableColumns   = postgres.ColumnList{ExportColumn, ExchangeColumn, ImportColumn}
	)

	return waypointsProductsTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		WaypointID: WaypointIDColumn,
		ProductID:  ProductIDColumn,
		Export:     ExportColumn,
		Exchange:   ExchangeColumn,
		Import:     ImportColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}