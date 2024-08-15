package rest

import (
	"github.com/gin-gonic/gin"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/julienlavocat/spacetraders/.gen/spacetraders/public/model"
	. "github.com/julienlavocat/spacetraders/.gen/spacetraders/public/table"
)

func listShips(c *gin.Context) {
	q := Ships.SELECT(Ships.AllColumns)

	var results []model.Ships
	if err := q.Query(db, &results); err != nil {
		internalServerError(c, "unable to query ships", err)
		return
	}

	ships := make([]*Ship, len(results))
	for i := range results {
		ship, err := adaptShip(results[i])
		if err != nil {
			internalServerError(c, "unable to adapt ship", err)
			return
		}
		ships[i] = ship
	}

	c.JSON(200, ships)
}

func getShipRouteToWaypoint(c *gin.Context) {
	shipId := c.Param("shipId")
	destination := c.Param("destination")

	ship, err := s.Ships.GetShip(shipId)
	if err != nil {
		notFoundError(c, "ship not found "+err.Error())
		return
	}

	route, err := s.Navigation.PlotRoute(ship.Nav.SystemSymbol, ship.Nav.WaypointSymbol, destination, ship.Fuel.Current)
	if err != nil {
		internalServerError(c, "unable to plot route to destination", err)
		return
	}
	c.JSON(200, route)
}

func getShip(c *gin.Context) {
	shipId := c.Param("shipId")

	q := Ships.SELECT(Ships.AllColumns).WHERE(Ships.ID.EQ(String(shipId)))
	var result []model.Ships
	if err := q.Query(db, &result); err != nil {
		internalServerError(c, "unable to query ships", err)
		return
	}

	if len(result) == 0 {
		notFoundError(c, "ship not found")
		return
	}

	c.JSON(200, result[0])
}
