package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/julienlavocat/spacetraders/.gen/spacetraders/public/model"
	. "github.com/julienlavocat/spacetraders/.gen/spacetraders/public/table"
	"github.com/julienlavocat/spacetraders/internal/rest/adapters"
)

func listSystems(c *gin.Context) {
	q := Systems.SELECT(Systems.ID, Systems.Type, Systems.X, Systems.Y)
	var results []model.Systems
	if err := q.Query(db, &results); err != nil {
		internalServerError(c, "unable to query systems", err)
		return
	}

	c.JSON(http.StatusOK, adapters.AdaptSystems(results))
}

func listWaypoints(c *gin.Context) {
	systemId := c.Param("systemId")
	q := Waypoints.SELECT(Waypoints.AllColumns).WHERE(Waypoints.SystemID.EQ(String(systemId)))
	var results []model.Waypoints
	if err := q.Query(db, &results); err != nil {
		internalServerError(c, "unable to query waypoints", err)
	}

	c.JSON(http.StatusOK, adapters.AdaptWaypoints(results))
}
