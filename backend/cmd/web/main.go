package main

import (
	"fmt"

	"backend-sistem06.com/internal/config"
)

func main() {
	viperConfig := config.NewViper()
	app := config.NewFiber(viperConfig)

	config.Bootstrap(&config.BootstrapConfig{
		App:   app,
		Viper: viperConfig,
	})
	webPort := viperConfig.GetInt("web.port")
	err := app.Listen(fmt.Sprintf(":%d", webPort))
	if err != nil {
		panic(err)
	}

}
