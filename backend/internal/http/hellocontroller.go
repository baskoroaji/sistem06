package http

import (
	"github.com/gofiber/fiber/v2"
)

type HelloController struct {
}

func NewHelloController() *HelloController {
	return &HelloController{}
}

func (h *HelloController) Hello(c *fiber.Ctx) error {
	return c.Send([]byte("hello world"))
}
