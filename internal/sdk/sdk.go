package sdk

import (
	"context"
	"database/sql"
	"os"

	"github.com/julienlavocat/spacetraders/internal/api"
	"github.com/julienlavocat/spacetraders/internal/utils"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Sdk struct {
	Client     *api.APIClient
	logger     zerolog.Logger
	Market     *Market
	Waypoints  *WaypointsService
	Navigation *Navigation
	Ships      map[string]*Ship
	DB         *sql.DB
}

func NewSdk() *Sdk {
	logger := log.With().Str("component", "SDK").Logger()

	db, err := sql.Open("postgres", "postgresql://spacetraders@localhost:5432/spacetraders?sslmode=disable")
	if err != nil {
		logger.Fatal().Err(err).Msg("database connection failed")
	}

	market := NewMarket(db)
	waypoints := NewWaypointsService(db)
	navigation := NewNavigation(db)

	sdk := &Sdk{
		Market:     market,
		Waypoints:  waypoints,
		logger:     logger,
		Navigation: navigation,
		DB:         db,
	}

	sdk.loadAgent()
	sdk.loadShips()

	return sdk
}

func (s *Sdk) loadAgent() {
	cfg := api.NewConfiguration()
	client := api.NewAPIClient(cfg)
	ctx := context.Background()

	content, err := os.ReadFile("./token")
	if err == nil && len(content) > 0 {
		s.logger.Info().Msg("token found, loading the agent")
		token := string(content)
		ctx = context.WithValue(ctx, api.ContextAccessToken, token)
		res := utils.RetryRequest(client.AgentsAPI.GetMyAgent(ctx).Execute, s.logger, "initial agent fetch failed")

		s.logger.Info().Interface("agent", res.Data).Msgf("agent %s loaded", res.Data.Symbol)
		cfg.AddDefaultHeader("Authorization", "Bearer "+token)
		s.Client = api.NewAPIClient(cfg)
		return
	}

	s.logger.Info().Msg("no token found, registering a new agent")

	res, _, err := client.DefaultAPI.Register(ctx).RegisterRequest(*api.NewRegisterRequest(api.FACTION_COSMIC, "JLVC")).Execute()
	if err != nil {
		s.logger.Fatal().Err(err).Msg("unable to register agent")
	}

	err = os.WriteFile("./token", []byte(res.Data.Token), 0644)
	if err != nil {
		s.logger.Fatal().Err(err).Str("token", res.Data.Token).Msg("unable to write token to file")
	}

	cfg.AddDefaultHeader("Authorization", "Bearer "+res.Data.Token)
	s.Client = api.NewAPIClient(cfg)
}

func (s *Sdk) loadShips() {
	shipsData := utils.RetryRequest(s.Client.FleetAPI.GetMyShips(context.Background()).Execute, log.Logger, "unable to retrieve ships")

	ships := make(map[string]*Ship)

	for _, ship := range shipsData.Data {
		ships[ship.Symbol] = NewShip(s.Client, ship)
	}

	s.Ships = ships
}
