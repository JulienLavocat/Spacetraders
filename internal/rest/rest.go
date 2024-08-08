package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/julienlavocat/spacetraders/internal/sdk"
	"github.com/rs/zerolog/log"
)

func StartApi(s *sdk.Sdk) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Use(gin.Recovery())

	r.GET("/ships/:shipId", func(c *gin.Context) {
		shipId := c.Param("shipId")
		c.JSON(200, serializeShip(s.Ships[shipId]))
	})

	if err := r.Run(); err != nil {
		log.Fatal().Err(err).Msg("unable to start API")
	}
}

func serializeShip(ship *sdk.Ship) ShipSnapshot {
	return ShipSnapshot{
		SystemId:     ship.Nav.SystemSymbol,
		WaypointId:   ship.Nav.WaypointSymbol,
		Id:           ship.Id,
		Cargo:        ship.Cargo,
		Status:       ship.Nav.Status,
		Destination:  ship.Nav.Route.Destination.Symbol,
		MaxFuel:      ship.Fuel.Capacity,
		CurrentFuel:  ship.Fuel.Current,
		MaxCargo:     ship.MaxCargo,
		CurrentCargo: ship.CurrentCargo,
		IsCargoFull:  ship.IsCargoFull,
		Departure:    ship.Nav.Route.Origin.Symbol,
		ArrivalAt:    ship.Nav.Route.Arrival,
		DepartedAt:   ship.Nav.Route.DepartureTime,
	}
}
