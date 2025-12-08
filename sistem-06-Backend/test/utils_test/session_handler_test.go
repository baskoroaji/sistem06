package utils_test

import (
	"testing"
	"time"

	"backend-sistem06.com/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/memory"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
)

// Helper function to create test app with session
func setupTestApp() (*fiber.App, *session.Store, *logrus.Logger) {
	app := fiber.New()

	// Use memory storage for testing
	storage := memory.New()
	store := session.New(session.Config{
		Storage:    storage,
		Expiration: 24 * time.Hour,
	})

	log := logrus.New()
	log.SetLevel(logrus.WarnLevel) // Reduce noise in tests

	return app, store, log
}

// Helper to create test context
func createTestContext(app *fiber.App) *fiber.Ctx {
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	return ctx
}

func TestNewSessionHandler(t *testing.T) {
	_, store, log := setupTestApp()

	handler := utils.NewSessionHandler(store, log)

	assert.NotNil(t, handler)
	assert.Equal(t, log, handler.Log)
}

func TestSetUserSession(t *testing.T) {
	app, store, log := setupTestApp()
	handler := utils.NewSessionHandler(store, log)

	tests := []struct {
		name    string
		userID  int
		email   string
		wantErr bool
	}{
		{
			name:    "Valid user session",
			userID:  1,
			email:   "user@example.com",
			wantErr: false,
		},
		{
			name:    "Valid user with different ID",
			userID:  999,
			email:   "test@example.com",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := createTestContext(app)
			defer app.ReleaseCtx(ctx)

			err := handler.SetUserSession(ctx, tt.userID, tt.email)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Verify session was set correctly
				sess, err := store.Get(ctx)
				require.NoError(t, err)

				assert.Equal(t, tt.userID, sess.Get("user_id"))
				assert.Equal(t, tt.email, sess.Get("email"))
				assert.Equal(t, true, sess.Get("authenticated"))
				assert.NotNil(t, sess.Get("created_at"))
			}
		})
	}
}

func TestGetUserID(t *testing.T) {
	app, store, log := setupTestApp()
	handler := utils.NewSessionHandler(store, log)

	t.Run("Get user ID from existing session", func(t *testing.T) {
		ctx := createTestContext(app)
		defer app.ReleaseCtx(ctx)

		// Set up session first
		expectedID := 123
		err := handler.SetUserSession(ctx, expectedID, "test@example.com")
		require.NoError(t, err)

		// Get user ID
		userID, err := handler.GetUserID(ctx)

		assert.NoError(t, err)
		assert.Equal(t, expectedID, userID)
	})

	t.Run("Get user ID from empty session", func(t *testing.T) {
		ctx := createTestContext(app)
		defer app.ReleaseCtx(ctx)

		// Don't set session
		userID, err := handler.GetUserID(ctx)

		assert.Error(t, err)
		assert.Equal(t, fiber.ErrUnauthorized, err)
		assert.Equal(t, 0, userID)
	})
}

func TestGetUserEmail(t *testing.T) {
	app, store, log := setupTestApp()
	handler := utils.NewSessionHandler(store, log)

	t.Run("Get email from existing session", func(t *testing.T) {
		ctx := createTestContext(app)
		defer app.ReleaseCtx(ctx)

		expectedEmail := "user@example.com"
		err := handler.SetUserSession(ctx, 1, expectedEmail)
		require.NoError(t, err)

		email, err := handler.GetUserEmail(ctx)

		assert.NoError(t, err)
		assert.Equal(t, expectedEmail, email)
	})

	t.Run("Get email from empty session", func(t *testing.T) {
		ctx := createTestContext(app)
		defer app.ReleaseCtx(ctx)

		email, err := handler.GetUserEmail(ctx)

		assert.Error(t, err)
		assert.Equal(t, fiber.ErrUnauthorized, err)
		assert.Equal(t, "", email)
	})

	t.Run("Get email with invalid type in session", func(t *testing.T) {
		ctx := createTestContext(app)
		defer app.ReleaseCtx(ctx)

		// Manually set invalid email type
		sess, err := store.Get(ctx)
		require.NoError(t, err)
		sess.Set("email", 12345) // Wrong type
		sess.Save()

		email, err := handler.GetUserEmail(ctx)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid session data")
		assert.Equal(t, "", email)
	})
}

func TestGetUserSession(t *testing.T) {
	app, store, log := setupTestApp()
	handler := utils.NewSessionHandler(store, log)

	t.Run("Get full session data", func(t *testing.T) {
		ctx := createTestContext(app)
		defer app.ReleaseCtx(ctx)

		expectedID := 1
		expectedEmail := "user@example.com"

		err := handler.SetUserSession(ctx, expectedID, expectedEmail)
		require.NoError(t, err)

		sessionData, err := handler.GetUserSession(ctx)

		assert.NoError(t, err)
		assert.NotNil(t, sessionData)
		assert.Equal(t, expectedID, sessionData["user_id"])
		assert.Equal(t, expectedEmail, sessionData["email"])
		assert.Equal(t, true, sessionData["authenticated"])
	})

	t.Run("Get session data when empty", func(t *testing.T) {
		ctx := createTestContext(app)
		defer app.ReleaseCtx(ctx)

		sessionData, err := handler.GetUserSession(ctx)

		assert.Error(t, err)
		assert.Equal(t, fiber.ErrUnauthorized, err)
		assert.Nil(t, sessionData)
	})

	t.Run("Get session data with missing user_id", func(t *testing.T) {
		ctx := createTestContext(app)
		defer app.ReleaseCtx(ctx)

		// Set only email
		sess, err := store.Get(ctx)
		require.NoError(t, err)
		sess.Set("email", "test@example.com")
		sess.Save()

		sessionData, err := handler.GetUserSession(ctx)

		assert.Error(t, err)
		assert.Equal(t, fiber.ErrUnauthorized, err)
		assert.Nil(t, sessionData)
	})
}

func TestIsAuthenticated(t *testing.T) {
	app, store, log := setupTestApp()
	handler := utils.NewSessionHandler(store, log)

	t.Run("User is authenticated", func(t *testing.T) {
		ctx := createTestContext(app)
		defer app.ReleaseCtx(ctx)

		err := handler.SetUserSession(ctx, 1, "user@example.com")
		require.NoError(t, err)

		isAuth := handler.IsAuthenticated(ctx)

		assert.True(t, isAuth)
	})

	t.Run("User is not authenticated - empty session", func(t *testing.T) {
		ctx := createTestContext(app)
		defer app.ReleaseCtx(ctx)

		isAuth := handler.IsAuthenticated(ctx)

		assert.False(t, isAuth)
	})

	t.Run("User is not authenticated - authenticated flag is false", func(t *testing.T) {
		ctx := createTestContext(app)
		defer app.ReleaseCtx(ctx)

		sess, err := store.Get(ctx)
		require.NoError(t, err)
		sess.Set("user_id", 1)
		sess.Set("email", "user@example.com")
		sess.Set("authenticated", false)
		sess.Save()

		isAuth := handler.IsAuthenticated(ctx)

		assert.False(t, isAuth)
	})

	t.Run("User is not authenticated - authenticated flag has wrong type", func(t *testing.T) {
		ctx := createTestContext(app)
		defer app.ReleaseCtx(ctx)

		sess, err := store.Get(ctx)
		require.NoError(t, err)
		sess.Set("user_id", 1)
		sess.Set("email", "user@example.com")
		sess.Set("authenticated", "true") // String instead of bool
		sess.Save()

		isAuth := handler.IsAuthenticated(ctx)

		assert.False(t, isAuth)
	})
}

func TestDestroySession(t *testing.T) {
	app, store, log := setupTestApp()
	handler := utils.NewSessionHandler(store, log)

	t.Run("Destroy existing session", func(t *testing.T) {
		ctx := createTestContext(app)
		defer app.ReleaseCtx(ctx)

		// Create session
		err := handler.SetUserSession(ctx, 1, "user@example.com")
		require.NoError(t, err)

		// Verify session exists
		assert.True(t, handler.IsAuthenticated(ctx))

		// Destroy session
		err = handler.DestroySession(ctx)
		assert.NoError(t, err)

		// Verify session is destroyed
		assert.False(t, handler.IsAuthenticated(ctx))
	})

	t.Run("Destroy non-existent session", func(t *testing.T) {
		ctx := createTestContext(app)
		defer app.ReleaseCtx(ctx)

		// Should not error even if session doesn't exist
		err := handler.DestroySession(ctx)
		assert.NoError(t, err)
	})
}

func TestRefreshSession(t *testing.T) {
	app, store, log := setupTestApp()
	handler := utils.NewSessionHandler(store, log)

	t.Run("Refresh existing session", func(t *testing.T) {
		ctx := createTestContext(app)
		defer app.ReleaseCtx(ctx)

		// Create session
		err := handler.SetUserSession(ctx, 1, "user@example.com")
		require.NoError(t, err)

		// Get original created_at
		sess, _ := store.Get(ctx)
		originalCreatedAt := sess.Get("created_at")

		// Wait a bit
		time.Sleep(10 * time.Millisecond)

		// Refresh session
		err = handler.RefreshSession(ctx)
		assert.NoError(t, err)

		// Verify session still exists
		assert.True(t, handler.IsAuthenticated(ctx))

		// created_at should remain the same
		sess, _ = store.Get(ctx)
		newCreatedAt := sess.Get("created_at")
		assert.Equal(t, originalCreatedAt, newCreatedAt)
	})

	t.Run("Refresh non-existent session", func(t *testing.T) {
		ctx := createTestContext(app)
		defer app.ReleaseCtx(ctx)

		// Should not error
		err := handler.RefreshSession(ctx)
		assert.NoError(t, err)
	})
}

func TestGetUserRole(t *testing.T) {
	app, store, log := setupTestApp()
	handler := utils.NewSessionHandler(store, log)

	t.Run("Get role from session", func(t *testing.T) {
		ctx := createTestContext(app)
		defer app.ReleaseCtx(ctx)

		expectedRole := "admin"

		// Set session with role
		sess, err := store.Get(ctx)
		require.NoError(t, err)
		sess.Set("role", expectedRole)
		sess.Save()

		role, err := handler.GetUserRole(ctx)

		assert.NoError(t, err)
		assert.Equal(t, expectedRole, role)
	})

	t.Run("Get role from empty session", func(t *testing.T) {
		ctx := createTestContext(app)
		defer app.ReleaseCtx(ctx)

		role, err := handler.GetUserRole(ctx)

		assert.NoError(t, err)
		assert.Equal(t, "", role)
	})

	t.Run("Get role with invalid type", func(t *testing.T) {
		ctx := createTestContext(app)
		defer app.ReleaseCtx(ctx)

		// Set invalid role type
		sess, err := store.Get(ctx)
		require.NoError(t, err)
		sess.Set("role", 12345) // Not a string
		sess.Save()

		role, err := handler.GetUserRole(ctx)

		assert.NoError(t, err)
		assert.Equal(t, "", role)
	})
}

// Integration test
func TestSessionHandlerIntegration(t *testing.T) {
	app, store, log := setupTestApp()
	handler := utils.NewSessionHandler(store, log)

	ctx := createTestContext(app)
	defer app.ReleaseCtx(ctx)

	// 1. Initially not authenticated
	assert.False(t, handler.IsAuthenticated(ctx))

	// 2. Set session
	expectedID := 42
	expectedEmail := "integration@test.com"
	err := handler.SetUserSession(ctx, expectedID, expectedEmail)
	require.NoError(t, err)

	// 3. Now authenticated
	assert.True(t, handler.IsAuthenticated(ctx))

	// 4. Can get user ID
	userID, err := handler.GetUserID(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expectedID, userID)

	// 5. Can get email
	email, err := handler.GetUserEmail(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expectedEmail, email)

	// 6. Can get full session
	sessionData, err := handler.GetUserSession(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expectedID, sessionData["user_id"])
	assert.Equal(t, expectedEmail, sessionData["email"])
	assert.Equal(t, true, sessionData["authenticated"])

	// 7. Refresh session
	err = handler.RefreshSession(ctx)
	assert.NoError(t, err)
	assert.True(t, handler.IsAuthenticated(ctx))

	// 8. Destroy session
	err = handler.DestroySession(ctx)
	assert.NoError(t, err)

	// 9. No longer authenticated
	assert.False(t, handler.IsAuthenticated(ctx))

	// 10. Cannot get user ID
	_, err = handler.GetUserID(ctx)
	assert.Error(t, err)
}

// Benchmark tests
func BenchmarkSetUserSession(b *testing.B) {
	app, store, log := setupTestApp()
	handler := utils.NewSessionHandler(store, log)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := createTestContext(app)
		handler.SetUserSession(ctx, i, "bench@test.com")
		app.ReleaseCtx(ctx)
	}
}

func BenchmarkIsAuthenticated(b *testing.B) {
	app, store, log := setupTestApp()
	handler := utils.NewSessionHandler(store, log)
	ctx := createTestContext(app)
	defer app.ReleaseCtx(ctx)

	handler.SetUserSession(ctx, 1, "bench@test.com")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handler.IsAuthenticated(ctx)
	}
}
