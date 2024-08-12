package main

import (
	"time"

	"github.com/julienlavocat/spacetraders/internal/ai"
	"github.com/julienlavocat/spacetraders/internal/sdk"
	"github.com/julienlavocat/spacetraders/internal/utils"
)

func main() {
	utils.SetupLogger()

	s := sdk.NewSdk()

	probesFleet := ai.NewMarketProbesFleet(s)
	probesFleet.BeginOperations("XI-NT44", time.Second*5)
}
