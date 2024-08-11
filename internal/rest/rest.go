package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/julienlavocat/spacetraders/internal/ai"
	"github.com/julienlavocat/spacetraders/internal/sdk"
	"github.com/rs/zerolog/log"
)

type RestApi struct {
	miningFleets map[string]*ai.MiningFleetCommander
}

func NewRestApi() *RestApi {
	return &RestApi{
		miningFleets: make(map[string]*ai.MiningFleetCommander),
	}
}

func (r *RestApi) AddMiningFleet(fleet *ai.MiningFleetCommander) {
	r.miningFleets[fleet.Id] = fleet
}

func (r *RestApi) StartApi(s *sdk.Sdk) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(gin.Recovery())

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"ready": s.Ready})
	})

	router.GET("/ships/:shipId", func(c *gin.Context) {
		shipId := c.Param("shipId")

		ship, ok := s.Ships[shipId]
		if !ok {
			c.JSON(404, gin.H{"error": "ship not found"})
			return
		}

		c.JSON(200, ship.GetSnapshot())
	})

	router.GET("/ships/:shipId/plot/:destination", func(c *gin.Context) {
		shipId := c.Param("shipId")
		destination := c.Param("destination")

		ship, ok := s.GetShip(shipId)
		if !ok {
			c.JSON(404, gin.H{"error": "ship not found"})
			return
		}

		route, err := s.Navigation.PlotRoute(ship.Nav.SystemSymbol, ship.Nav.WaypointSymbol, destination, ship.Fuel.Current)
		if err != nil {
			c.String(500, err.Error())
		}
		c.JSON(200, route)
	})

	router.GET("/fleets/mining/:fleetId", func(c *gin.Context) {
		fleetId := c.Param("fleetId")

		fleet, ok := r.miningFleets[fleetId]
		if !ok {
			c.JSON(404, gin.H{"error": "fleet not found"})
			return
		}

		c.JSON(200, fleet.GetSnapshot())
	})

	router.GET("/market/:systemId", func(c *gin.Context) {
		systemId := c.Param("systemId")
		c.JSON(200, s.Market.GetTradeRoutes(systemId))
	})

	if err := router.Run("0.0.0.0:8080"); err != nil {
		log.Fatal().Err(err).Msg("unable to start API")
	}
}
