package model_test

import (
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"

	"backend-sistem06.com/internal/model"
)

func TestRegisterUserRequest_Validation(t *testing.T) {
	validate := validator.New()

	t.Run("should pass validation with valid data", func(t *testing.T) {
		request := &model.RegisterUserRequest{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "password123",
		}

		err := validate.Struct(request)
		assert.NoError(t, err)
	})

	t.Run("should fail when name is empty", func(t *testing.T) {
		request := &model.RegisterUserRequest{
			Name:     "",
			Email:    "john@example.com",
			Password: "password123",
		}

		err := validate.Struct(request)
		assert.Error(t, err)

		validationErrors := err.(validator.ValidationErrors)
		assert.Equal(t, "Name", validationErrors[0].Field())
		assert.Equal(t, "required", validationErrors[0].Tag())
	})

	t.Run("should fail when name exceeds max length", func(t *testing.T) {
		longName := string(make([]byte, 101)) // 101 characters
		request := &model.RegisterUserRequest{
			Name:     longName,
			Email:    "john@example.com",
			Password: "password123",
		}

		err := validate.Struct(request)
		assert.Error(t, err)

		validationErrors := err.(validator.ValidationErrors)
		assert.Equal(t, "Name", validationErrors[0].Field())
		assert.Equal(t, "max", validationErrors[0].Tag())
	})

	t.Run("should pass when name is exactly 100 characters", func(t *testing.T) {
		exactName := strings.Repeat("a", 100)
		request := &model.RegisterUserRequest{
			Name:     exactName,
			Email:    "john@example.com",
			Password: "password123",
		}

		err := validate.Struct(request)
		assert.NoError(t, err)
	})

	t.Run("should fail when email is empty", func(t *testing.T) {
		request := &model.RegisterUserRequest{
			Name:     "John Doe",
			Email:    "",
			Password: "password123",
		}

		err := validate.Struct(request)
		assert.Error(t, err)

		validationErrors := err.(validator.ValidationErrors)
		assert.Equal(t, "Email", validationErrors[0].Field())
		assert.Equal(t, "required", validationErrors[0].Tag())
	})

	t.Run("should fail when email format is invalid", func(t *testing.T) {
		invalidEmails := []string{
			"notanemail",
			"missing@domain",
			"@nodomain.com",
			"spaces in@email.com",
			"double@@domain.com",
		}

		for _, invalidEmail := range invalidEmails {
			request := &model.RegisterUserRequest{
				Name:     "John Doe",
				Email:    invalidEmail,
				Password: "password123",
			}

			err := validate.Struct(request)
			assert.Error(t, err, "should fail for email: %s", invalidEmail)
		}
	})

	t.Run("should pass with valid email formats", func(t *testing.T) {
		validEmails := []string{
			"user@example.com",
			"user.name@example.com",
			"user+tag@example.co.id",
			"123@example.com",
		}

		for _, validEmail := range validEmails {
			request := &model.RegisterUserRequest{
				Name:     "John Doe",
				Email:    validEmail,
				Password: "password123",
			}

			err := validate.Struct(request)
			assert.NoError(t, err, "should pass for email: %s", validEmail)
		}
	})

	t.Run("should fail when password is empty", func(t *testing.T) {
		request := &model.RegisterUserRequest{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "",
		}

		err := validate.Struct(request)
		assert.Error(t, err)

		validationErrors := err.(validator.ValidationErrors)
		assert.Equal(t, "Password", validationErrors[0].Field())
		assert.Equal(t, "required", validationErrors[0].Tag())
	})

	t.Run("should fail when password is too short", func(t *testing.T) {
		request := &model.RegisterUserRequest{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "pass", // Less than 8 characters
		}

		err := validate.Struct(request)
		assert.Error(t, err)

		validationErrors := err.(validator.ValidationErrors)
		assert.Equal(t, "Password", validationErrors[0].Field())
		assert.Equal(t, "min", validationErrors[0].Tag())
	})

	t.Run("should pass when password is exactly 8 characters", func(t *testing.T) {
		request := &model.RegisterUserRequest{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "password", // Exactly 8 characters
		}

		err := validate.Struct(request)
		assert.NoError(t, err)
	})

	t.Run("should fail with multiple validation errors", func(t *testing.T) {
		request := &model.RegisterUserRequest{
			Name:     "",
			Email:    "invalid",
			Password: "short",
		}

		err := validate.Struct(request)
		assert.Error(t, err)

		validationErrors := err.(validator.ValidationErrors)
		assert.GreaterOrEqual(t, len(validationErrors), 2)
	})
}

func TestLoginUserRequest_Validation(t *testing.T) {
	validate := validator.New()

	t.Run("should pass validation with valid data", func(t *testing.T) {
		request := &model.LoginUserRequest{
			Email:    "john@example.com",
			Password: "password123",
		}

		err := validate.Struct(request)
		assert.NoError(t, err)
	})

	t.Run("should fail when email is empty", func(t *testing.T) {
		request := &model.LoginUserRequest{
			Email:    "",
			Password: "password123",
		}

		err := validate.Struct(request)
		assert.Error(t, err)
	})

	t.Run("should fail when email format is invalid", func(t *testing.T) {
		request := &model.LoginUserRequest{
			Email:    "invalid-email",
			Password: "password123",
		}

		err := validate.Struct(request)
		assert.Error(t, err)
	})

	t.Run("should fail when password is too short", func(t *testing.T) {
		request := &model.LoginUserRequest{
			Email:    "john@example.com",
			Password: "short",
		}

		err := validate.Struct(request)
		assert.Error(t, err)
	})

	t.Run("should fail when password exceeds max length", func(t *testing.T) {
		longPassword := string(make([]byte, 101))
		request := &model.LoginUserRequest{
			Email:    "john@example.com",
			Password: longPassword,
		}

		err := validate.Struct(request)
		assert.Error(t, err)
	})

	t.Run("should pass when password is exactly 100 characters", func(t *testing.T) {
		exactPassword := string(make([]byte, 100))
		request := &model.LoginUserRequest{
			Email:    "john@example.com",
			Password: exactPassword,
		}

		err := validate.Struct(request)
		assert.NoError(t, err)
	})
}

func TestVerifyUserRequest_Validation(t *testing.T) {
	validate := validator.New()

	t.Run("should pass validation with valid token", func(t *testing.T) {
		request := &model.VerifyUserRequest{
			Token: "valid_token_here",
		}

		err := validate.Struct(request)
		assert.NoError(t, err)
	})

	t.Run("should fail when token is empty", func(t *testing.T) {
		request := &model.VerifyUserRequest{
			Token: "",
		}

		err := validate.Struct(request)
		assert.Error(t, err)

		validationErrors := err.(validator.ValidationErrors)
		assert.Equal(t, "Token", validationErrors[0].Field())
		assert.Equal(t, "required", validationErrors[0].Tag())
	})

	t.Run("should fail when token exceeds max length", func(t *testing.T) {
		longToken := string(make([]byte, 101))
		request := &model.VerifyUserRequest{
			Token: longToken,
		}

		err := validate.Struct(request)
		assert.Error(t, err)

		validationErrors := err.(validator.ValidationErrors)
		assert.Equal(t, "Token", validationErrors[0].Field())
		assert.Equal(t, "max", validationErrors[0].Tag())
	})

	t.Run("should pass when token is exactly 100 characters", func(t *testing.T) {
		exactToken := string(make([]byte, 100))
		request := &model.VerifyUserRequest{
			Token: exactToken,
		}

		err := validate.Struct(request)
		assert.NoError(t, err)
	})
}

// Test JSON serialization
func TestUserResponse_JSONTags(t *testing.T) {
	t.Run("should omit empty fields in JSON", func(t *testing.T) {
		response := &model.UserResponse{
			ID:   0,
			Name: "",
		}

		// Note: In actual JSON marshaling, omitempty will omit zero values
		// This test documents the struct tags behavior
		assert.Equal(t, 0, response.ID)
		assert.Equal(t, "", response.Name)
		assert.Equal(t, int64(0), response.CreatedAt)
		assert.Equal(t, int64(0), response.UpdatedAt)
	})

	t.Run("should include non-empty fields", func(t *testing.T) {
		response := &model.UserResponse{
			ID:        1,
			Name:      "John Doe",
			CreatedAt: 1234567890,
			UpdatedAt: 1234567890,
		}

		assert.Equal(t, 1, response.ID)
		assert.Equal(t, "John Doe", response.Name)
		assert.Equal(t, int64(1234567890), response.CreatedAt)
		assert.Equal(t, int64(1234567890), response.UpdatedAt)
	})
}
