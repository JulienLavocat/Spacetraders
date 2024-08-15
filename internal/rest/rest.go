package rest

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/julienlavocat/spacetraders/.gen/spacetraders/public/model"
	. "github.com/julienlavocat/spacetraders/.gen/spacetraders/public/table"
	"github.com/julienlavocat/spacetraders/internal/rest/adapters"
	"github.com/julienlavocat/spacetraders/internal/sdk"
	"github.com/rs/zerolog/log"
)

var (
	db *sql.DB
	s  *sdk.Sdk
)

type RestApi struct{}

func NewRestApi(sd *sdk.Sdk) *RestApi {
	dbConn, err := sql.Open("postgres", "postgresql://spacetraders:spacetraders@localhost:5432/spacetraders?sslmode=disable")
	if err != nil {
		log.Fatal().Err(err).Msg("database connection failed")
	}

	db = dbConn
	s = sd

	gin.SetMode(gin.ReleaseMode)

	return &RestApi{}
}

func (r *RestApi) StartApi() {
	router := gin.Default()

	router.Use(CORSMiddleware())

	router.GET("/", r.ping)
	router.GET("/ships", listShips)
	router.GET("/ships/:shipId", getShip)
	router.GET("/ships/:shipId/plot/:destination", getShipRouteToWaypoint)
	router.GET("/fleets/trading/:fleetId", r.getTradingFleet)
	router.GET("/market/:systemId", r.listTradeRoutes)
	router.GET("/transactions", listTransaction)

	if err := router.Run("0.0.0.0:8080"); err != nil {
		log.Fatal().Err(err).Msg("unable to start API")
	}
}

func (r *RestApi) listTradeRoutes(c *gin.Context) {
	systemId := c.Param("systemId")
	c.JSON(200, s.Market.GetTradeRoutes(systemId))
}

func (r *RestApi) getTradingFleet(c *gin.Context) {
	fleetId := c.Param("fleetId")

	q := TradingFleets.SELECT(TradingFleets.AllColumns).WHERE(TradingFleets.ID.EQ(String(fleetId)))
	var result []model.TradingFleets
	if err := q.Query(db, &result); err != nil {
		c.JSON(500, gin.H{"message": "unable to query trading fleets", "error": err})
		return
	}

	if len(result) == 0 {
		c.JSON(404, gin.H{"message": "fleet not found"})
		return
	}

	fleet, err := adapters.AdaptTradingFleet(result[0])
	if err != nil {
		c.JSON(500, gin.H{"message": "unable to parse fleet", "error": err})
	}

	c.JSON(200, fleet)
}

func (r *RestApi) ping(c *gin.Context) {
	c.JSON(200, gin.H{"ready": s.Ready})
}
