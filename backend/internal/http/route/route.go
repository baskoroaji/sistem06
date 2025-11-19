package route

import (
	"backend-sistem06.com/internal/http"
	"backend-sistem06.com/internal/http/middleware"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App             *fiber.App
	HelloController *http.HelloController
	UserController  *http.UserController
	AuthController  *http.AuthController
	AuthMiddleware  *middleware.AuthMiddleware
}

func (c *RouteConfig) Setup() {
	c.AuthRoute()
}

func (c *RouteConfig) AuthRoute() {
	auth := c.App.Group("/v1/api/auth")
	// c.App.Get("/hello", c.HelloController.Hello)
	auth.Post("/register", c.UserController.Register)
	auth.Post("/login", c.AuthController.Login)
}
