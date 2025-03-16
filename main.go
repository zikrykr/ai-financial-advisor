package main

import (
	"github.com/ai-financial-advisor/cmd/rest"
	appSetup "github.com/ai-financial-advisor/cmd/setup"
	"github.com/ai-financial-advisor/config"
)

func main() {
	// config init
	config.InitConfig()

	// app setup init
	setup := appSetup.InitSetup()

	// starting REST server
	rest.StartServer(setup)
}
