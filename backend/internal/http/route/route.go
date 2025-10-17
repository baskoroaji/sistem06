package route

import (
	"backend-sistem06.com/internal/http"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App             *fiber.App
	HelloController *http.HelloController
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
}

func (c *RouteConfig) SetupGuestRoute() {
	c.App.Get("/hello", c.HelloController.Hello)
}
