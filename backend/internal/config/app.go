package config

import (
	"backend-sistem06.com/internal/http"
	"backend-sistem06.com/internal/http/route"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

type BootstrapConfig struct {
	App   *fiber.App
	Viper *viper.Viper
}

func Bootstrap(config *BootstrapConfig) {
	helloContoller := http.NewHelloController()

	routeConfig := route.RouteConfig{
		App:             config.App,
		HelloController: helloContoller,
	}
	routeConfig.Setup()

}
