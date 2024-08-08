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

	router.GET("/ships/:shipId", func(c *gin.Context) {
		shipId := c.Param("shipId")
		c.JSON(200, s.Ships[shipId].GetSnapshot())
	})

	router.GET("/fleets/mining/:fleetId", func(c *gin.Context) {
		fleetId := c.Param("fleetId")
		c.JSON(200, r.miningFleets[fleetId].GetSnapshot())
	})

	if err := router.Run(); err != nil {
		log.Fatal().Err(err).Msg("unable to start API")
	}
}
