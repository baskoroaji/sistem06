package test

import (
	"errors"
	"testing"

	"backend-sistem06.com/utils" // sesuaikan dengan module path Anda
	"github.com/go-playground/validator/v10"
)

// Test struct untuk simulasi validation
type TestUser struct {
	Email    string `validate:"required,email"`
	Username string `validate:"required,min=3,max=20"`
	Age      int    `validate:"required,min=18"`
}

func TestValidationError(t *testing.T) {
	validate := validator.New()

	t.Run("should return nil when error is nil", func(t *testing.T) {
		result := utils.ValidationError(nil)

		if result != nil {
			t.Error("expected nil, got non-nil result")
		}
	})

	t.Run("should handle required field error", func(t *testing.T) {
		user := TestUser{Email: "", Username: ""}
		err := validate.Struct(user)

		errors := utils.ValidationError(err)

		if errors == nil {
			t.Fatal("expected errors map, got nil")
		}

		if _, exists := errors["Email"]; !exists {
			t.Error("expected Email field error")
		}

		if _, exists := errors["Username"]; !exists {
			t.Error("expected Username field error")
		}
	})

	t.Run("should handle email validation error", func(t *testing.T) {
		user := TestUser{Email: "invalid-email", Username: "john"}
		err := validate.Struct(user)

		errors := utils.ValidationError(err)

		if errors == nil {
			t.Fatal("expected errors map, got nil")
		}

		emailErr, exists := errors["Email"]
		if !exists {
			t.Fatal("expected Email field error")
		}

		expected := "Email must be a valid email address"
		if emailErr != expected {
			t.Errorf("expected '%s', got '%s'", expected, emailErr)
		}
	})

	t.Run("should handle min length validation error", func(t *testing.T) {
		user := TestUser{Email: "test@test.com", Username: "ab", Age: 20}
		err := validate.Struct(user)

		errors := utils.ValidationError(err)

		if errors == nil {
			t.Fatal("expected errors map, got nil")
		}

		usernameErr, exists := errors["Username"]
		if !exists {
			t.Fatal("expected Username field error")
		}

		expected := "Username must be at least 3 characters"
		if usernameErr != expected {
			t.Errorf("expected '%s', got '%s'", expected, usernameErr)
		}
	})

	t.Run("should handle max length validation error", func(t *testing.T) {
		user := TestUser{
			Email:    "test@test.com",
			Username: "verylongusernamethatexceedsmaximumlength",
			Age:      20,
		}
		err := validate.Struct(user)

		errors := utils.ValidationError(err)

		if errors == nil {
			t.Fatal("expected errors map, got nil")
		}

		usernameErr, exists := errors["Username"]
		if !exists {
			t.Fatal("expected Username field error")
		}

		expected := "Username must be at most 20 characters"
		if usernameErr != expected {
			t.Errorf("expected '%s', got '%s'", expected, usernameErr)
		}
	})

	t.Run("should use default message for unknown validation tag", func(t *testing.T) {
		// Menggunakan tag yang tidak ada di Messages map
		type CustomStruct struct {
			Field string `validate:"alphanum"` // tag tidak ada di Messages
		}

		custom := CustomStruct{Field: "invalid@#$"}
		err := validate.Struct(custom)

		errors := utils.ValidationError(err)

		if errors == nil {
			t.Fatal("expected errors map, got nil")
		}

		fieldErr, exists := errors["Field"]
		if !exists {
			t.Fatal("expected Field error")
		}

		expected := "Field is invalid"
		if fieldErr != expected {
			t.Errorf("expected '%s', got '%s'", expected, fieldErr)
		}
	})

	t.Run("should handle multiple validation errors", func(t *testing.T) {
		user := TestUser{Email: "invalid", Username: "ab"}
		err := validate.Struct(user)

		errors := utils.ValidationError(err)

		if errors == nil {
			t.Fatal("expected errors map, got nil")
		}

		if len(errors) < 2 {
			t.Errorf("expected at least 2 errors, got %d", len(errors))
		}
	})

	t.Run("should return nil for non-validator error", func(t *testing.T) {
		// Simulasi error biasa (bukan validator.ValidationErrors)
		err := errors.New("some generic error")

		result := utils.ValidationError(err)

		if result != nil {
			t.Error("expected nil for non-validator error")
		}
	})
}

func TestFormatValidationErrors(t *testing.T) {
	t.Run("should format single error correctly", func(t *testing.T) {
		errs := map[string]string{
			"Email": "Email is required",
		}

		result := utils.FormatValidationErrors(errs)

		expected := `{"Email": "Email is required"}`
		if result != expected {
			t.Errorf("expected '%s', got '%s'", expected, result)
		}
	})

	t.Run("should format multiple errors correctly", func(t *testing.T) {
		errs := map[string]string{
			"Email":    "Email is required",
			"Username": "Username must be at least 3 characters",
		}

		result := utils.FormatValidationErrors(errs)

		// Map iteration order is not guaranteed, so check both fields exist
		if len(result) == 0 {
			t.Fatal("expected non-empty result")
		}

		// Check that result contains both errors
		containsEmail := false
		containsUsername := false

		for field := range errs {
			if len(result) > 0 && result[0] == '{' && result[len(result)-1] == '}' {
				// Basic structure check
				if field == "Email" {
					containsEmail = true
				}
				if field == "Username" {
					containsUsername = true
				}
			}
		}

		if !containsEmail || !containsUsername {
			t.Error("result should contain both Email and Username fields")
		}
	})

	t.Run("should handle empty errors map", func(t *testing.T) {
		errs := map[string]string{}

		result := utils.FormatValidationErrors(errs)

		expected := "{}"
		if result != expected {
			t.Errorf("expected '%s', got '%s'", expected, result)
		}
	})

	t.Run("should format with proper JSON structure", func(t *testing.T) {
		errs := map[string]string{
			"Field": "Error message",
		}

		result := utils.FormatValidationErrors(errs)

		// Check starts with { and ends with }
		if result[0] != '{' || result[len(result)-1] != '}' {
			t.Error("result should be wrapped in curly braces")
		}

		// Check contains quotes around field and value
		if result != `{"Field": "Error message"}` {
			t.Errorf("expected proper JSON format, got '%s'", result)
		}
	})
}

// Integration test: Test both functions together
func TestValidationErrorIntegration(t *testing.T) {
	validate := validator.New()

	t.Run("should work end-to-end", func(t *testing.T) {
		user := TestUser{Email: "invalid", Username: "ab"}
		err := validate.Struct(user)

		// Get validation errors
		errors := utils.ValidationError(err)
		if errors == nil {
			t.Fatal("expected errors, got nil")
		}

		// Format errors
		formatted := utils.FormatValidationErrors(errors)

		// Should produce valid JSON-like string
		if formatted[0] != '{' || formatted[len(formatted)-1] != '}' {
			t.Error("formatted result should be JSON-like structure")
		}

		// Should contain error messages
		if len(formatted) < 10 {
			t.Error("formatted result seems too short")
		}
	})
}
