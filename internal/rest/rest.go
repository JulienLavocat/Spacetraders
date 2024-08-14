package rest

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/julienlavocat/spacetraders/.gen/spacetraders/public/model"
	. "github.com/julienlavocat/spacetraders/.gen/spacetraders/public/table"
	"github.com/julienlavocat/spacetraders/internal/ai"
	"github.com/julienlavocat/spacetraders/internal/sdk"
	"github.com/rs/zerolog/log"
)

type RestApi struct {
	miningFleets map[string]*ai.MiningFleet
	tradingFleet map[string]*ai.TradingFleet
	s            *sdk.Sdk
	db           *sql.DB
}

func NewRestApi(s *sdk.Sdk) *RestApi {
	return &RestApi{
		s: s,
	}
}

func (r *RestApi) StartApi() {
	db, err := sql.Open("postgres", "postgresql://spacetraders:spacetraders@localhost:5432/spacetraders?sslmode=disable")
	if err != nil {
		log.Fatal().Err(err).Msg("database connection failed")
	}
	r.db = db

	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.GET("/", r.ping)
	router.GET("/ships/:shipId", r.getShip)
	router.GET("/ships/:shipId/plot/:destination", r.getShipRouteToWaypoint)
	router.GET("/fleets/trading/:fleetId", r.getTradingFleet)
	router.GET("/market/:systemId", r.listTradeRoutes)

	if err := router.Run("0.0.0.0:8080"); err != nil {
		log.Fatal().Err(err).Msg("unable to start API")
	}
}

func (r *RestApi) listTradeRoutes(c *gin.Context) {
	systemId := c.Param("systemId")
	c.JSON(200, r.s.Market.GetTradeRoutes(systemId))
}

func (r *RestApi) getTradingFleet(c *gin.Context) {
	fleetId := c.Param("fleetId")

	q := TradingFleets.SELECT(TradingFleets.AllColumns).WHERE(TradingFleets.ID.EQ(String(fleetId)))
	var result []model.TradingFleets
	if err := q.Query(r.db, &result); err != nil {
		c.JSON(500, gin.H{"message": "unable to query trading fleets", "error": err})
		return
	}

	if len(result) == 0 {
		c.JSON(404, gin.H{"message": "fleet not found"})
		return
	}

	fleet, err := adaptTradingFleet(result[0])
	if err != nil {
		c.JSON(500, gin.H{"message": "unable to parse fleet", "error": err})
	}

	c.JSON(200, fleet)
}

func (r *RestApi) getShipRouteToWaypoint(c *gin.Context) {
	shipId := c.Param("shipId")
	destination := c.Param("destination")

	ship, err := r.s.Ships.GetShip(shipId)
	if err != nil {
		c.JSON(404, gin.H{"message": "ship not found", "error": err})
		return
	}

	route, err := r.s.Navigation.PlotRoute(ship.Nav.SystemSymbol, ship.Nav.WaypointSymbol, destination, ship.Fuel.Current)
	if err != nil {
		c.String(500, err.Error())
	}
	c.JSON(200, route)
}

func (r *RestApi) getShip(c *gin.Context) {
	shipId := c.Param("shipId")

	q := Ships.SELECT(Ships.AllColumns).WHERE(Ships.ID.EQ(String(shipId)))
	var result []model.Ships
	if err := q.Query(r.db, &result); err != nil {
		c.JSON(500, gin.H{"message": "unable to query ships", "error": err})
		return
	}

	if len(result) == 0 {
		c.JSON(404, gin.H{"message": "ship not found"})
		return
	}

	c.JSON(200, result[0])
}

func (r *RestApi) ping(c *gin.Context) {
	c.JSON(200, gin.H{"ready": r.s.Ready})
}
