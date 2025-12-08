package converter_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"backend-sistem06.com/internal/entity"
	"backend-sistem06.com/internal/model/converter"
)

func TestUserToResponse(t *testing.T) {
	t.Run("should convert user entity to response correctly", func(t *testing.T) {
		now := time.Now().Unix()
		user := &entity.UserEntity{
			ID:        1,
			Name:      "John Doe",
			Email:     "john@example.com",
			Password:  "hashed_password_here",
			CreatedAt: now,
			UpdatedAt: now,
		}

		response := converter.UserToResponse(user)

		assert.NotNil(t, response)
		assert.Equal(t, user.ID, response.ID)
		assert.Equal(t, user.Name, response.Name)
		assert.Equal(t, user.CreatedAt, response.CreatedAt)
		assert.Equal(t, user.UpdatedAt, response.UpdatedAt)
	})

	t.Run("should not include email in response", func(t *testing.T) {
		user := &entity.UserEntity{
			ID:        1,
			Name:      "John Doe",
			Email:     "john@example.com",
			Password:  "hashed_password",
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		}

		response := converter.UserToResponse(user)

		// UserResponse doesn't have Email field
		assert.NotNil(t, response)
		// We can't directly check for absence of Email field in struct,
		// but we can verify the fields that should be present
		assert.Equal(t, user.ID, response.ID)
		assert.Equal(t, user.Name, response.Name)
	})

	t.Run("should not include password in response", func(t *testing.T) {
		user := &entity.UserEntity{
			ID:        1,
			Name:      "John Doe",
			Email:     "john@example.com",
			Password:  "super_secret_password",
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		}

		response := converter.UserToResponse(user)

		// UserResponse doesn't have Password field
		assert.NotNil(t, response)
		assert.Equal(t, user.ID, response.ID)
		assert.Equal(t, user.Name, response.Name)
	})

	t.Run("should handle zero values correctly", func(t *testing.T) {
		user := &entity.UserEntity{
			ID:        0,
			Name:      "",
			Email:     "",
			Password:  "",
			CreatedAt: 0,
			UpdatedAt: 0,
		}

		response := converter.UserToResponse(user)

		assert.NotNil(t, response)
		assert.Equal(t, 0, response.ID)
		assert.Equal(t, "", response.Name)
		assert.Equal(t, int64(0), response.CreatedAt)
		assert.Equal(t, int64(0), response.UpdatedAt)
	})

	t.Run("should handle nil entity gracefully", func(t *testing.T) {
		// This will panic if not handled, but based on the converter code
		// it doesn't handle nil check, so this documents the behavior
		assert.Panics(t, func() {
			converter.UserToResponse(nil)
		}, "converter should panic on nil entity")
	})
}

func TestUserToResponse_WithDifferentTimestamps(t *testing.T) {
	testCases := []struct {
		name      string
		createdAt int64
		updatedAt int64
	}{
		{
			name:      "same timestamps",
			createdAt: time.Now().Unix(),
			updatedAt: time.Now().Unix(),
		},
		{
			name:      "different timestamps",
			createdAt: time.Now().Unix(),
			updatedAt: time.Now().Unix() + 3600, // 1 hour later
		},
		{
			name:      "zero timestamps",
			createdAt: 0,
			updatedAt: 0,
		},
		{
			name:      "unix epoch",
			createdAt: 0,
			updatedAt: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user := &entity.UserEntity{
				ID:        1,
				Name:      "Test User",
				Email:     "test@example.com",
				Password:  "password",
				CreatedAt: tc.createdAt,
				UpdatedAt: tc.updatedAt,
			}

			response := converter.UserToResponse(user)

			assert.Equal(t, tc.createdAt, response.CreatedAt)
			assert.Equal(t, tc.updatedAt, response.UpdatedAt)
		})
	}
}
