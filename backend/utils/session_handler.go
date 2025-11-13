package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/sirupsen/logrus"
)

type SessionHandler struct {
	store *session.Store
	Log   *logrus.Logger
}

func NewSessionHandler(store *session.Store, log *logrus.Logger) *SessionHandler {
	return &SessionHandler{
		store: store,
		Log:   log,
	}
}

func (s *SessionHandler) SetUserSession(ctx *fiber.Ctx, userID int, email string) error {
	sess, err := s.store.Get(ctx)
	if err != nil {
		return err
	}
	s.Log.Infof("Setting session for user %d with email %s", userID, email)

	sess.Set("user_id", userID)
	sess.Set("email", email)
	err = sess.Save()
	if err != nil {
		s.Log.Errorf("Failed saving session: %v", err)
	}
	return err
}

func (s *SessionHandler) GetUserID(ctx *fiber.Ctx) (int, error) {
	sess, err := s.store.Get(ctx)
	if err != nil {
		return 0, err
	}

	id := sess.Get("user_id")
	if id == nil {
		return 0, fiber.ErrUnauthorized
	}
	return id.(int), nil
}
