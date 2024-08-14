package main

import (
	"github.com/julienlavocat/spacetraders/internal/sdk"
	"github.com/rs/zerolog/log"
)

func main() {
	s := sdk.NewSdk(true)

	// api := rest.NewRestApi()
	// go api.StartApi(s)

	// ship := s.GetShip("JLVC-A")
	// plan, err := s.Navigation.PlotRoute("X1-QA42", "X1-QA42-D43", "X1-QA42-B32", 40, 40)
	// plan, err := s.Navigation.PlotRoute("X1-QA42", "X1-QA42-B7", "X1-QA42-B7", ship.Fuel.Current, ship.Fuel.Capacity)

	// wp, fuel, err := s.Navigation.FindNearestStation("X1-QA42", "X1-QA42-CX5B")
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("pathfinding error")
	// }
	// println(wp, fuel)

	// for _, step := range plan {
	// 	println(step.To, step.Fuel)
	// }

	// fuel, err := s.Navigation.FuelCostBetweenWaypoints("X1-QA42-CX5B", "X1-QA42-B32")
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("fuelcost")
	// }
	// println(fuel)

	path, err := s.Navigation.PlotRoute("X1-NT44", "X1-NT44-XZ5B", "X1-NT44-B34", 200)
	if err != nil {
		log.Fatal().Err(err).Msg("failed")
	}

	for _, v := range path {
		println(v.To, v.Fuel, v.AggCost)
	}
}
