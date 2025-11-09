package route

import (
	"backend-sistem06.com/internal/http"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App             *fiber.App
	HelloController *http.HelloController
	UserController  *http.UserController
	AuthController  *http.AuthController
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
}

func (c *RouteConfig) SetupGuestRoute() {
	c.App.Get("/hello", c.HelloController.Hello)
	c.App.Post("/api/register", c.UserController.Register)
	c.App.Post("/api/login", c.AuthController.Login)
}
