package main

import (
	"context"
	"log"

	"github.com/julienlavocat/spacetraders/internal/api"
)

func main() {
	api := api.NewAPIClient(api.NewConfiguration())

	res, _, err := api.AgentsAPI.GetAgents(context.Background()).Execute()
	if err != nil {
		log.Fatal(err)
	}

	log.Print(res.Data)
}
