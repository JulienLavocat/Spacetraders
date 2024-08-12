package sdk

import (
	"context"
	"database/sql"
	"os"
	"sync/atomic"
	"time"

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
	Ready      bool
	Balance    atomic.Int64
}

func NewSdk() *Sdk {
	logger := log.With().Str("component", "SDK").Logger()

	db, err := sql.Open("postgres", "postgresql://spacetraders:spacetraders@localhost:5432/spacetraders?sslmode=disable")
	if err != nil {
		logger.Fatal().Err(err).Msg("database connection failed")
	}

	market := NewMarket(db, nil)
	waypoints := NewWaypointsService(db)
	navigation := NewNavigation(db)

	sdk := &Sdk{
		Market:     market,
		Waypoints:  waypoints,
		logger:     logger,
		Navigation: navigation,
		DB:         db,
		Ready:      false,
		Balance:    atomic.Int64{},
	}

	sdk.loadAgent()

	return sdk
}

func (s *Sdk) Init() {
	s.loadAgent()
	s.loadShips()
	s.Ready = true
	go s.updateAgentBalance()
}

func (s *Sdk) GetShip(id string) (*Ship, bool) {
	ship, ok := s.Ships[id]
	return ship, ok
}

func (s *Sdk) RefreshBalance() error {
	res, _, err := utils.RetryRequestWithoutFatal(s.Client.AgentsAPI.GetMyAgent(context.Background()).Execute, s.logger)
	if err != nil {
		return err
	}

	s.Balance.Swap(res.Data.Credits)
	return nil
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
		s.Balance.Swap(res.Data.Credits)
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
	// FIXME: PAGINATION
	shipsData := utils.RetryRequest(s.Client.FleetAPI.GetMyShips(context.Background()).Limit(20).Execute, log.Logger, "unable to retrieve ships")

	ships := make(map[string]*Ship)

	for _, ship := range shipsData.Data {
		ships[ship.Symbol] = NewShip(s, ship)
	}

	s.Ships = ships
}

func (s *Sdk) updateAgentBalance() {
	ticker := time.NewTicker(time.Second * 20)

	for range ticker.C {
		res, errBody, err := utils.RetryRequestWithoutFatal(s.Client.AgentsAPI.GetMyAgent(context.Background()).Execute, s.logger)
		if err != nil || errBody != nil {
			s.logger.Error().Err(err).Interface("body", errBody).Msg("unable to load agent")
			continue
		}

		s.Balance.Swap(res.Data.Credits)
	}
}
