package sdk

import (
	"context"

	"github.com/julienlavocat/spacetraders/internal/api"
	"github.com/julienlavocat/spacetraders/internal/utils"
	"github.com/rs/zerolog/log"
)

type System struct {
	waypoints    map[string]api.Waypoint
	traits       map[api.WaypointTraitSymbol][]string
	imports      map[api.TradeSymbol][]string
	exports      map[api.TradeSymbol][]string
	exchanges    map[api.TradeSymbol][]string
	marketplaces []string
}

func NewSystem(client *api.APIClient, id string) *System {
	s := &System{
		waypoints:    make(map[string]api.Waypoint),
		traits:       make(map[api.WaypointTraitSymbol][]string),
		imports:      make(map[api.TradeSymbol][]string),
		exchanges:    make(map[api.TradeSymbol][]string),
		exports:      make(map[api.TradeSymbol][]string),
		marketplaces: make([]string, 0),
	}

	res, http, err := client.SystemsAPI.GetSystemWaypoints(context.Background(), id).Execute()
	utils.FatalIfHttpError(http, err, log.Logger, "unable to load system %s", id)

	for _, waypoint := range res.Data {
		println(waypoint.Symbol)
		s.waypoints[waypoint.Symbol] = waypoint

		for _, trait := range waypoint.Traits {
			if trait.Symbol == api.MARKETPLACE {
				s.marketplaces = append(s.marketplaces, waypoint.Symbol)

				marketRes, http, err := client.SystemsAPI.GetMarket(context.Background(), id, waypoint.Symbol).Execute()
				utils.FatalIfHttpError(http, err, log.Logger, "unable to load market data for %s", waypoint.Symbol)
				log.Info().Interface("", marketRes.Data).Msg("market")

				for _, product := range marketRes.Data.Imports {
					s.imports[product.Symbol] = append(s.imports[product.Symbol], waypoint.Symbol)
					log.Info().Interface("imports", s.imports).Msg("")
				}

				for _, product := range marketRes.Data.Exports {
					s.exports[product.Symbol] = append(s.exports[product.Symbol], waypoint.Symbol)
				}

				for _, product := range marketRes.Data.Exchange {
					s.exchanges[product.Symbol] = append(s.exchanges[product.Symbol], waypoint.Symbol)
				}
			}

			s.traits[trait.Symbol] = append(s.traits[trait.Symbol], waypoint.Symbol)
		}

	}

	return s
}

func (s *System) GetByTrait(trait api.WaypointTraitSymbol) []string {
	return s.traits[trait]
}

func (s *System) GetByImports(product api.TradeSymbol) []string {
	return s.imports[product]
}

func (s *System) GetByExports(product api.TradeSymbol) []string {
	return s.exports[product]
}

func (s *System) GetByExchange(product api.TradeSymbol) []string {
	return s.exchanges[product]
}
