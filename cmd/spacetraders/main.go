package main

import (
	"sync"
	"time"

	"github.com/julienlavocat/spacetraders/internal/ai"
	"github.com/julienlavocat/spacetraders/internal/sdk"
	"github.com/julienlavocat/spacetraders/internal/utils"
	"github.com/rs/zerolog/log"
)

func main() {
	utils.SetupLogger()
	s := sdk.NewSdk(true)

	// go createMiningFleet(s, restApi)
	go createTradeFleet(s)

	// This allow the service to run forever
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}

func createMiningFleet(s *sdk.Sdk) {
	miners := []string{"JLVC-3", "JLVC-4", "JLVC-5"}
	hauler := "JLVC-1"
	miningFleet := ai.NewMiningFleet(s, "MNG_1", miners, hauler)

	if err := miningFleet.BeginOperations("X1-QA42"); err != nil {
		log.Fatal().Err(err).Msg("unable to begin operations for fleet MNG_1")
	}
}

func createTradeFleet(s *sdk.Sdk) {
	fleet := ai.NewTradingFleet(s, "TRD_1", "X1-NT44", time.Minute, []string{"JLVC-1", "JLVC-D", "JLVC-1C", "JLVC-1D"})
	fleet.BeginOperations()
}
