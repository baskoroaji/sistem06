package pkg

import (
	"time"

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
	sess.Set("authenticated", true)
	sess.Set("created_at", time.Now().Unix())
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

func (s *SessionHandler) GetUserEmail(ctx *fiber.Ctx) (string, error) {
	sess, err := s.store.Get(ctx)
	if err != nil {
		return "", err
	}
	email := sess.Get("email")
	if email == nil {
		return "", fiber.ErrUnauthorized
	}
	emailStr, ok := email.(string)
	if !ok {
		s.Log.Errorf("Invalid email type: %T", email)
		return "", fiber.NewError(fiber.StatusInternalServerError, "invalid session data")
	}

	return emailStr, nil
}

func (s *SessionHandler) GetUserSession(ctx *fiber.Ctx) (map[string]interface{}, error) {
	sess, err := s.store.Get(ctx)
	if err != nil {
		return nil, err
	}

	userID := sess.Get("user_id")
	email := sess.Get("email")
	authenticated := sess.Get("authenticated")

	if userID == nil || email == nil {
		return nil, fiber.ErrUnauthorized
	}

	return map[string]interface{}{
		"user_id":       userID,
		"email":         email,
		"authenticated": authenticated,
	}, nil
}

func (s *SessionHandler) IsAuthenticated(ctx *fiber.Ctx) bool {
	sess, err := s.store.Get(ctx)
	if err != nil {
		return false
	}

	authenticated := sess.Get("authenticated")
	if authenticated == nil {
		return false
	}

	auth, ok := authenticated.(bool)
	return ok && auth
}

func (s *SessionHandler) DestroySession(ctx *fiber.Ctx) error {
	sess, err := s.store.Get(ctx)
	if err != nil {
		return err
	}

	if err := sess.Destroy(); err != nil {
		s.Log.Errorf("Failed to destroy session: %v", err)
		return err
	}

	s.Log.Info("Session destroyed successfully")
	return nil
}

func (s *SessionHandler) RefreshSession(ctx *fiber.Ctx) error {
	sess, err := s.store.Get(ctx)
	if err != nil {
		return err
	}

	// Just calling Save() will refresh the expiration
	if err := sess.Save(); err != nil {
		s.Log.Errorf("Failed to refresh session: %v", err)
		return err
	}

	return nil
}

func (s *SessionHandler) GetUserRole(ctx *fiber.Ctx) (string, error) {
	sess, err := s.store.Get(ctx)
	if err != nil {
		return "", err
	}

	if role, ok := sess.Get("role").(string); ok {
		return role, nil
	}

	return "", nil
}
