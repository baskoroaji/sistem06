package utils_test

import (
	"testing"

	"backend-sistem06.com/utils"
)

func TestGenerateToken(t *testing.T) {
	t.Run("should generate token with correct length", func(t *testing.T) {
		token := utils.GenerateToken()

		// 32 bytes = 64 hex characters
		expectedLength := 64
		if len(token) != expectedLength {
			t.Errorf("expected token length %d, got %d", expectedLength, len(token))
		}
	})

	t.Run("should generate different tokens on multiple calls", func(t *testing.T) {
		token1 := utils.GenerateToken()
		token2 := utils.GenerateToken()

		if token1 == token2 {
			t.Error("expected different tokens, got identical tokens")
		}
	})

	t.Run("should generate valid hex string", func(t *testing.T) {
		token := utils.GenerateToken()

		// Check if all characters are valid hex (0-9, a-f)
		for _, char := range token {
			if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f')) {
				t.Errorf("token contains invalid hex character: %c", char)
			}
		}
	})

	t.Run("should never generate empty token", func(t *testing.T) {
		token := utils.GenerateToken()

		if token == "" {
			t.Error("expected non-empty token, got empty string")
		}
	})
}

func BenchmarkGenerateToken(b *testing.B) {
	for i := 0; i < b.N; i++ {
		utils.GenerateToken()
	}
}
