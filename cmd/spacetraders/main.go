package main

import (
	"io"
	"os"
	"sync"
	"time"

	"github.com/julienlavocat/spacetraders/internal/ai"
	"github.com/julienlavocat/spacetraders/internal/rest"
	"github.com/julienlavocat/spacetraders/internal/sdk"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	setupLogger()
	s := sdk.NewSdk()

	restApi := rest.NewRestApi()
	go restApi.StartApi(s)

	s.Init()

	// ship, _ := s.GetShip("JLVC-1")
	// ship.Buy("SHIP_PLATING", 2)
	// ship.NavigateTo("X1-NT44-A2")
	// ship.Sell([]sdk.SellPlan{
	// 	{
	// 		ToSell:   sdk.Cargo{"SHIP_PLATING": 11},
	// 		Location: "X1-NT44-A2",
	// 	},
	// })

	// go createMiningFleet(s, restApi)
	go createMarketProbesFleet(s)
	go createTradeFleet(s)

	// This allow the service to run forever
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}

func createMarketProbesFleet(s *sdk.Sdk) {
	probesFleet := ai.NewMarketProbesFleet(s)
	probesFleet.BeginOperations("XI-NT44", time.Second*5)
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

func createTradeFleet(s *sdk.Sdk) {
	fleet := ai.NewTradingFleet(s, "TRD_1", "X1-NT44", time.Minute, []string{"JLVC-1"})
	fleet.BeginOperations()
}

func setupLogger() {
	writers := []io.Writer{zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}}
	if _, ok := os.LookupEnv("PRODUCTION"); ok {
		logFile, _ := os.OpenFile(
			"logs.txt",
			os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666,
		)
		writers = append(writers, logFile)
	}

	log.Logger = zerolog.New(zerolog.MultiLevelWriter(writers...)).With().Timestamp().Logger()
}
