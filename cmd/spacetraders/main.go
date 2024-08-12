package main

import (
	"sync"
	"time"

	"github.com/julienlavocat/spacetraders/internal/ai"
	"github.com/julienlavocat/spacetraders/internal/rest"
	"github.com/julienlavocat/spacetraders/internal/sdk"
	"github.com/julienlavocat/spacetraders/internal/utils"
	"github.com/rs/zerolog/log"
)

func main() {
	utils.SetupLogger()
	s := sdk.NewSdk()

	restApi := rest.NewRestApi()
	go restApi.StartApi(s)

	s.Init()

	// go createMiningFleet(s, restApi)
	go createTradeFleet(s, restApi)

	// This allow the service to run forever
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}

func createMiningFleet(s *sdk.Sdk, restApi *rest.RestApi) {
	miners := []*sdk.Ship{s.Ships["JLVC-3"], s.Ships["JLVC-4"], s.Ships["JLVC-5"]}
	hauler := s.Ships["JLVC-1"]
	miningFleet := ai.NewMiningFleetCommander(s, "MNG_1", miners, hauler)
	restApi.AddMiningFleet(miningFleet)

	if err := miningFleet.BeginOperations("X1-QA42"); err != nil {
		log.Fatal().Err(err).Msg("unable to begin operations for fleet MNG_1")
	}
}

func createTradeFleet(s *sdk.Sdk, restApi *rest.RestApi) {
	fleet := ai.NewTradingFleet(s, "TRD_1", "X1-NT44", time.Minute, []string{"JLVC-1", "JLVC-D"})
	restApi.AddTradingFleet(fleet)
	fleet.BeginOperations()
}
