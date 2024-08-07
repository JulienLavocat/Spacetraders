package main

import (
	"context"
	"database/sql"
	"os"
	"time"

	_ "github.com/lib/pq"

	"github.com/julienlavocat/spacetraders/internal/api"
	"github.com/julienlavocat/spacetraders/internal/sdk"
	"github.com/julienlavocat/spacetraders/internal/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	db, err := sql.Open("postgres", "postgresql://spacetraders@localhost:5432/spacetraders?sslmode=disable")
	if err != nil {
		log.Fatal().Err(err).Msg("database connection failed")
	}
	defer db.Close()

	client := loadAgent()

	agentRes := utils.RetryRequest(client.AgentsAPI.GetMyAgent(context.Background()).Execute, log.Logger, "unable to retrieve agent")
	log.Info().Interface("agent", agentRes.Data).Msg("agent loaded")

	shipsData := utils.RetryRequest(client.FleetAPI.GetMyShips(context.Background()).Execute, log.Logger, "unable to retrieve ships")

	market := sdk.NewMarket(db)
	ship := sdk.NewShip(client, shipsData.Data[0])

	asteroidRes := utils.RetryRequest(client.SystemsAPI.GetSystemWaypoints(context.Background(), ship.Nav.SystemSymbol).Type_(api.ENGINEERED_ASTEROID).Execute, log.Logger, "unable to find and engineered asteroid")

	contractsRes := utils.RetryRequest(client.ContractsAPI.GetContracts(context.Background()).Execute, log.Logger, "unable to retrieve contract")

	asteroid := asteroidRes.Data[0]
	contract := contractsRes.Data[0]

	for !ship.DeliverAndFulfillContract(contract).Fulfilled {
		if ship.IsCargoFull {
			plan := market.SellCargoTo(ship.Nav.SystemSymbol, ship.Cargo)
			for _, step := range plan {
				ship.Sell(step)
			}
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
