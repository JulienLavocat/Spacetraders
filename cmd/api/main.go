package main

import (
	"github.com/julienlavocat/spacetraders/internal/rest"
	"github.com/julienlavocat/spacetraders/internal/sdk"
	"github.com/julienlavocat/spacetraders/internal/utils"
)

func main() {
	utils.SetupLogger()

	sdk := sdk.NewSdk(false)

	api := rest.NewRestApi(sdk)

	api.StartApi()
}
