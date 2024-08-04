package main

import (
	"context"
	"os"
	"time"

	"github.com/julienlavocat/spacetraders/internal/api"
	"github.com/julienlavocat/spacetraders/internal/sdk"
	"github.com/julienlavocat/spacetraders/internal/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	client := loadAgent()

	agentRes, http, err := client.AgentsAPI.GetMyAgent(context.Background()).Execute()
	utils.FatalIfHttpError(http, err, log.Logger, "unable to retrieve agent")
	log.Info().Interface("agent", agentRes.Data).Msg("agent loaded")

	shipsData, http, err := client.FleetAPI.GetMyShips(context.Background()).Execute()
	utils.FatalIfHttpError(http, err, log.Logger, "unable to retrieve ships")

	ship := sdk.NewShip(client, shipsData.Data[0])

	system := sdk.NewSystem(client, ship.Nav.SystemSymbol)
	log.Info().Interface("waypoints", system.GetByImports(api.COPPER_ORE)).Msg("Found copper ore locations")

	asteroidRes, http, err := client.SystemsAPI.GetSystemWaypoints(context.Background(), ship.Nav.SystemSymbol).Type_(api.ENGINEERED_ASTEROID).Execute()
	utils.FatalIfHttpError(http, err, log.Logger, "unable to find an engineered asteroid")

	contractsRes, http, err := client.ContractsAPI.GetContracts(context.Background()).Execute()
	utils.FatalIfHttpError(http, err, log.Logger, "unable to retrieve contracts")

	asteroid := asteroidRes.Data[0]
	contract := contractsRes.Data[0]

	for !ship.DeliverAndFulfillContract(contract).Fulfilled {
		if ship.IsCargoFull {
			ship.NavigateTo("X1-QA42-H51").SellFullCargo()
		}

		ship.NavigateTo(asteroid.Symbol).Refuel().Mine()
	}
}

func loadAgent() *api.APIClient {
	cfg := api.NewConfiguration()
	client := api.NewAPIClient(cfg)
	ctx := context.Background()

	content, err := os.ReadFile("./token")
	if err == nil && len(content) > 0 {
		token := string(content)
		ctx = context.WithValue(ctx, api.ContextAccessToken, token)
		_, _, err := client.AgentsAPI.GetMyAgent(ctx).Execute()
		if err != nil {
			log.Fatal().Err(err).Msg("initial agent fetch failed")
		}

		cfg.AddDefaultHeader("Authorization", "Bearer "+token)
		return api.NewAPIClient(cfg)
	}

	res, _, err := client.DefaultAPI.Register(ctx).RegisterRequest(*api.NewRegisterRequest(api.FACTION_COSMIC, "JLVC")).Execute()
	if err != nil {
		log.Fatal().Err(err).Msg("unable to register agent")
	}

	err = os.WriteFile("./token", []byte(res.Data.Token), 0644)
	if err != nil {
		log.Fatal().Err(err).Str("token", res.Data.Token).Msg("unable to write token to file")
	}

	cfg.AddDefaultHeader("Authorization", "Bearer "+res.Data.Token)
	return api.NewAPIClient(cfg)
}
