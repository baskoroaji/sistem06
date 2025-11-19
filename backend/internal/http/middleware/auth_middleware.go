package middleware

import (
	"backend-sistem06.com/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type AuthMiddleware struct {
	SessionHandler *utils.SessionHandler
	Log            *logrus.Logger
}

func NewAuthMiddleware(sessionHandler *utils.SessionHandler, log *logrus.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		SessionHandler: sessionHandler,
		Log:            log,
	}
}

func (m *AuthMiddleware) RequireGuest() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if m.SessionHandler.IsAuthenticated(c) {
			return fiber.NewError(fiber.StatusForbidden, "Already Authenticated")
		}
		return c.Next()
	}
}

func (m *AuthMiddleware) RequiredAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !m.SessionHandler.IsAuthenticated(c) {
			m.Log.Warn("Unauthorized access attempt")
			return fiber.NewError(fiber.StatusUnauthorized, "authentication required")
		}
		userID, err := m.SessionHandler.GetUserID(c)
		if err != nil {
			m.Log.Errorf("Failed to get user ID from session: %v", err)
			return fiber.ErrUnauthorized
		}

		email, err := m.SessionHandler.GetUserEmail(c)
		if err != nil {
			m.Log.Errorf("Failed to get user Email from session: %v", err)
			return fiber.ErrUnauthorized
		}

		c.Locals("user_id", userID)
		c.Locals("email", email)

		if err := m.SessionHandler.RefreshSession(c); err != nil {
			m.Log.Warnf("Failed to refresh session: %v", err)
		}

		return c.Next()
	}
}

// func (m *AuthMiddleware) RequiredRoles (role string) fiber.Handler {
// 	return func (c *fiber.Ctx) error  {
// 		if !m.SessionHandler.IsAuthenticated(c) {
// 			return fiber.ErrUnauthorized
// 		}
// 		if !m.SessionHandler.
// 	}
// }
