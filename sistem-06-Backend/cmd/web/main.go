package main

import (
	"fmt"

	"backend-sistem06.com/internal/config"
)

func main() {
	viperConfig := config.NewViper()
	app := config.NewFiber(viperConfig)
	log := config.NewLogger(viperConfig)
	db := config.NewPostgres(viperConfig, log)
	validate := config.NewValidator(viperConfig)
	sessionStore := config.NewSession(viperConfig, db, log)

	config.Bootstrap(&config.BootstrapConfig{
		App:       app,
		Config:    viperConfig,
		Validator: validate,
		DB:        db,
		Log:       log,
		Session:   sessionStore,
	})
	webPort := viperConfig.GetInt("web.port")
	err := app.Listen(fmt.Sprintf(":%d", webPort))
	if err != nil {
		panic(err)
	}

}
