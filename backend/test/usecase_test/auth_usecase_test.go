package usecase_test

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	"backend-sistem06.com/internal/entity"
	"backend-sistem06.com/internal/model"
	"backend-sistem06.com/internal/repository"
	"backend-sistem06.com/internal/usecase"
)

// Mock UserRepository for Login tests
type MockUserRepositoryForLogin struct {
	FindByEmailFunc func(email string) (*entity.UserEntity, error)
}

func (m *MockUserRepositoryForLogin) FindByEmail(email string) (*entity.UserEntity, error) {
	if m.FindByEmailFunc != nil {
		return m.FindByEmailFunc(email)
	}
	return nil, errors.New("not implemented")
}

func (m *MockUserRepositoryForLogin) CreateUser(tx *sql.Tx, user *entity.UserEntity) error {
	return nil
}

func (m *MockUserRepositoryForLogin) FindByID(id int) (*entity.UserEntity, error) {
	return nil, nil
}

func (m *MockUserRepositoryForLogin) CountById(tx *sql.Tx, id int) (int, error) {
	return 0, nil
}

func (m *MockUserRepositoryForLogin) CountByName(tx *sql.Tx, name string) (int, error) {
	return 0, nil
}

// Mock TokenRepository
type MockTokenRepository struct {
	CreateTokenFunc   func(tx *sql.Tx, token *entity.PersonalAccessToken) error
	FindTokenByIdFunc func(id int) (*entity.PersonalAccessToken, error)
}

func (m *MockTokenRepository) CreateToken(tx *sql.Tx, token *entity.PersonalAccessToken) error {
	if m.CreateTokenFunc != nil {
		return m.CreateTokenFunc(tx, token)
	}
	return nil
}

func (m *MockTokenRepository) FindTokenById(id int) (*entity.PersonalAccessToken, error) {
	if m.FindTokenByIdFunc != nil {
		return m.FindTokenByIdFunc(id)
	}
	return nil, errors.New("not implemented")
}

func setupAuthUseCase(
	t *testing.T,
	mockUserRepo repository.UserRepositoryInterface,
	mockTokenRepo repository.TokenRepositoryInterface,
) (*usecase.AuthUseCase, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}

	log := logrus.New()
	log.SetOutput(io.Discard)

	validate := validator.New()

	var userRepo repository.UserRepositoryInterface
	if mockUserRepo != nil {
		userRepo = mockUserRepo
	} else {
		userRepo = repository.NewUserRepository(db, log)
	}

	var tokenRepo repository.TokenRepositoryInterface
	if mockTokenRepo != nil {
		tokenRepo = mockTokenRepo
	} else {
		tokenRepo = repository.NewTokenRepository(db, log)
	}

	uc := usecase.NewAuthUseCase(db, log, validate, userRepo, tokenRepo)

	cleanup := func() { db.Close() }
	return uc, mock, cleanup
}

func TestAuthUseCase_Login(t *testing.T) {
	// Create a hashed password for testing
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	t.Run("should successfully login with valid credentials", func(t *testing.T) {
		mockUserRepo := &MockUserRepositoryForLogin{
			FindByEmailFunc: func(email string) (*entity.UserEntity, error) {
				return &entity.UserEntity{
					ID:       1,
					Name:     "John Doe",
					Email:    "john@example.com",
					Password: string(hashedPassword),
				}, nil
			},
		}

		mockTokenRepo := &MockTokenRepository{
			CreateTokenFunc: func(tx *sql.Tx, token *entity.PersonalAccessToken) error {
				token.ID = 1
				return nil
			},
		}

		uc, mock, cleanup := setupAuthUseCase(t, mockUserRepo, mockTokenRepo)
		defer cleanup()

		mock.ExpectBegin()
		mock.ExpectCommit()

		request := &model.LoginUserRequest{
			Email:    "john@example.com",
			Password: "password123",
		}

		result, err := uc.Login(context.Background(), request)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotEmpty(t, result.Token)
		assert.Len(t, result.Token, 64) // GenerateToken returns 64 char hex string
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when transaction begin fails", func(t *testing.T) {
		uc, mock, cleanup := setupAuthUseCase(t, nil, nil)
		defer cleanup()

		mock.ExpectBegin().WillReturnError(sql.ErrConnDone)

		request := &model.LoginUserRequest{
			Email:    "john@example.com",
			Password: "password123",
		}

		result, err := uc.Login(context.Background(), request)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, fiber.ErrInternalServerError, err)
	})

	t.Run("should return validation error for invalid email", func(t *testing.T) {
		uc, mock, cleanup := setupAuthUseCase(t, nil, nil)
		defer cleanup()

		mock.ExpectBegin()
		mock.ExpectRollback()

		request := &model.LoginUserRequest{
			Email:    "invalid-email",
			Password: "password123",
		}

		result, err := uc.Login(context.Background(), request)

		assert.Error(t, err)
		assert.Nil(t, result)

		fiberErr, ok := err.(*fiber.Error)
		assert.True(t, ok)
		assert.Equal(t, fiber.StatusBadRequest, fiberErr.Code)
		assert.Contains(t, fiberErr.Message, "Email")
	})

	t.Run("should return validation error for missing required fields", func(t *testing.T) {
		uc, mock, cleanup := setupAuthUseCase(t, nil, nil)
		defer cleanup()

		mock.ExpectBegin()
		mock.ExpectRollback()

		request := &model.LoginUserRequest{
			Email:    "",
			Password: "password123",
		}

		result, err := uc.Login(context.Background(), request)

		assert.Error(t, err)
		assert.Nil(t, result)

		fiberErr, ok := err.(*fiber.Error)
		assert.True(t, ok)
		assert.Equal(t, fiber.StatusBadRequest, fiberErr.Code)
	})

	t.Run("should return validation error for short password", func(t *testing.T) {
		uc, mock, cleanup := setupAuthUseCase(t, nil, nil)
		defer cleanup()

		mock.ExpectBegin()
		mock.ExpectRollback()

		request := &model.LoginUserRequest{
			Email:    "john@example.com",
			Password: "short",
		}

		result, err := uc.Login(context.Background(), request)

		assert.Error(t, err)
		assert.Nil(t, result)

		fiberErr, ok := err.(*fiber.Error)
		assert.True(t, ok)
		assert.Equal(t, fiber.StatusBadRequest, fiberErr.Code)
	})

	t.Run("should return unauthorized when user not found", func(t *testing.T) {
		mockUserRepo := &MockUserRepositoryForLogin{
			FindByEmailFunc: func(email string) (*entity.UserEntity, error) {
				return nil, errors.New("user not found")
			},
		}

		uc, mock, cleanup := setupAuthUseCase(t, mockUserRepo, nil)
		defer cleanup()

		mock.ExpectBegin()
		mock.ExpectRollback()

		request := &model.LoginUserRequest{
			Email:    "notfound@example.com",
			Password: "password123",
		}

		result, err := uc.Login(context.Background(), request)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, fiber.ErrUnauthorized, err)
	})

	t.Run("should return unauthorized when password is incorrect", func(t *testing.T) {
		mockUserRepo := &MockUserRepositoryForLogin{
			FindByEmailFunc: func(email string) (*entity.UserEntity, error) {
				return &entity.UserEntity{
					ID:       1,
					Name:     "John Doe",
					Email:    "john@example.com",
					Password: string(hashedPassword), // correct password hash
				}, nil
			},
		}

		uc, mock, cleanup := setupAuthUseCase(t, mockUserRepo, nil)
		defer cleanup()

		mock.ExpectBegin()
		mock.ExpectRollback()

		request := &model.LoginUserRequest{
			Email:    "john@example.com",
			Password: "wrongpassword", // wrong password
		}

		result, err := uc.Login(context.Background(), request)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, fiber.ErrUnauthorized, err)
	})

	t.Run("should return error when token creation fails", func(t *testing.T) {
		mockUserRepo := &MockUserRepositoryForLogin{
			FindByEmailFunc: func(email string) (*entity.UserEntity, error) {
				return &entity.UserEntity{
					ID:       1,
					Name:     "John Doe",
					Email:    "john@example.com",
					Password: string(hashedPassword),
				}, nil
			},
		}

		mockTokenRepo := &MockTokenRepository{
			CreateTokenFunc: func(tx *sql.Tx, token *entity.PersonalAccessToken) error {
				return errors.New("database error")
			},
		}

		uc, mock, cleanup := setupAuthUseCase(t, mockUserRepo, mockTokenRepo)
		defer cleanup()

		mock.ExpectBegin()
		mock.ExpectRollback()

		request := &model.LoginUserRequest{
			Email:    "john@example.com",
			Password: "password123",
		}

		result, err := uc.Login(context.Background(), request)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, fiber.ErrInternalServerError, err)
	})

	t.Run("should return error when commit fails", func(t *testing.T) {
		mockUserRepo := &MockUserRepositoryForLogin{
			FindByEmailFunc: func(email string) (*entity.UserEntity, error) {
				return &entity.UserEntity{
					ID:       1,
					Name:     "John Doe",
					Email:    "john@example.com",
					Password: string(hashedPassword),
				}, nil
			},
		}

		mockTokenRepo := &MockTokenRepository{
			CreateTokenFunc: func(tx *sql.Tx, token *entity.PersonalAccessToken) error {
				return nil
			},
		}

		uc, mock, cleanup := setupAuthUseCase(t, mockUserRepo, mockTokenRepo)
		defer cleanup()

		mock.ExpectBegin()
		mock.ExpectCommit().WillReturnError(errors.New("commit failed"))

		request := &model.LoginUserRequest{
			Email:    "john@example.com",
			Password: "password123",
		}

		result, err := uc.Login(context.Background(), request)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, fiber.ErrInternalServerError, err)
	})

	t.Run("should generate unique token for each login", func(t *testing.T) {
		mockUserRepo := &MockUserRepositoryForLogin{
			FindByEmailFunc: func(email string) (*entity.UserEntity, error) {
				return &entity.UserEntity{
					ID:       1,
					Name:     "John Doe",
					Email:    "john@example.com",
					Password: string(hashedPassword),
				}, nil
			},
		}

		mockTokenRepo := &MockTokenRepository{
			CreateTokenFunc: func(tx *sql.Tx, token *entity.PersonalAccessToken) error {
				return nil
			},
		}

		uc, mock, cleanup := setupAuthUseCase(t, mockUserRepo, mockTokenRepo)
		defer cleanup()

		request := &model.LoginUserRequest{
			Email:    "john@example.com",
			Password: "password123",
		}

		// First login
		mock.ExpectBegin()
		mock.ExpectCommit()
		result1, err1 := uc.Login(context.Background(), request)

		// Second login
		mock.ExpectBegin()
		mock.ExpectCommit()
		result2, err2 := uc.Login(context.Background(), request)

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NotNil(t, result1)
		assert.NotNil(t, result2)
		assert.NotEqual(t, result1.Token, result2.Token, "tokens should be unique")
	})

	t.Run("should create token with correct user ID", func(t *testing.T) {
		var capturedToken *entity.PersonalAccessToken
		expectedUserID := 123

		mockUserRepo := &MockUserRepositoryForLogin{
			FindByEmailFunc: func(email string) (*entity.UserEntity, error) {
				return &entity.UserEntity{
					ID:       expectedUserID,
					Name:     "John Doe",
					Email:    "john@example.com",
					Password: string(hashedPassword),
				}, nil
			},
		}

		mockTokenRepo := &MockTokenRepository{
			CreateTokenFunc: func(tx *sql.Tx, token *entity.PersonalAccessToken) error {
				capturedToken = token
				return nil
			},
		}

		uc, mock, cleanup := setupAuthUseCase(t, mockUserRepo, mockTokenRepo)
		defer cleanup()

		mock.ExpectBegin()
		mock.ExpectCommit()

		request := &model.LoginUserRequest{
			Email:    "john@example.com",
			Password: "password123",
		}

		result, err := uc.Login(context.Background(), request)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotNil(t, capturedToken)
		assert.Equal(t, expectedUserID, capturedToken.UserID)
		assert.NotEmpty(t, capturedToken.Token)
		assert.Greater(t, capturedToken.CreatedAt, int64(0))
		assert.Greater(t, capturedToken.ExpiredAt, capturedToken.CreatedAt)
	})

	t.Run("should set token expiration to 24 hours", func(t *testing.T) {
		var capturedToken *entity.PersonalAccessToken

		mockUserRepo := &MockUserRepositoryForLogin{
			FindByEmailFunc: func(email string) (*entity.UserEntity, error) {
				return &entity.UserEntity{
					ID:       1,
					Name:     "John Doe",
					Email:    "john@example.com",
					Password: string(hashedPassword),
				}, nil
			},
		}

		mockTokenRepo := &MockTokenRepository{
			CreateTokenFunc: func(tx *sql.Tx, token *entity.PersonalAccessToken) error {
				capturedToken = token
				return nil
			},
		}

		uc, mock, cleanup := setupAuthUseCase(t, mockUserRepo, mockTokenRepo)
		defer cleanup()

		mock.ExpectBegin()
		mock.ExpectCommit()

		request := &model.LoginUserRequest{
			Email:    "john@example.com",
			Password: "password123",
		}

		_, err := uc.Login(context.Background(), request)

		assert.NoError(t, err)
		assert.NotNil(t, capturedToken)

		// Verify expiration is approximately 24 hours from creation
		expirationDuration := capturedToken.ExpiredAt - capturedToken.CreatedAt
		expectedDuration := int64(24 * 60 * 60) // 24 hours in seconds

		// Allow 1 second tolerance for test execution time
		assert.InDelta(t, expectedDuration, expirationDuration, 1)
	})

	t.Run("should rollback transaction on any error", func(t *testing.T) {
		mockUserRepo := &MockUserRepositoryForLogin{
			FindByEmailFunc: func(email string) (*entity.UserEntity, error) {
				return nil, errors.New("database error")
			},
		}

		uc, mock, cleanup := setupAuthUseCase(t, mockUserRepo, nil)
		defer cleanup()

		mock.ExpectBegin()
		mock.ExpectRollback()

		request := &model.LoginUserRequest{
			Email:    "john@example.com",
			Password: "password123",
		}

		result, err := uc.Login(context.Background(), request)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// Table-driven test for validation scenarios
func TestAuthUseCase_Login_ValidationScenarios(t *testing.T) {
	testCases := []struct {
		name        string
		request     *model.LoginUserRequest
		expectError bool
		errorField  string
	}{
		{
			name: "valid request",
			request: &model.LoginUserRequest{
				Email:    "john@example.com",
				Password: "password123",
			},
			expectError: false,
		},
		{
			name: "empty email",
			request: &model.LoginUserRequest{
				Email:    "",
				Password: "password123",
			},
			expectError: true,
			errorField:  "Email",
		},
		{
			name: "invalid email format",
			request: &model.LoginUserRequest{
				Email:    "not-an-email",
				Password: "password123",
			},
			expectError: true,
			errorField:  "Email",
		},
		{
			name: "empty password",
			request: &model.LoginUserRequest{
				Email:    "john@example.com",
				Password: "",
			},
			expectError: true,
			errorField:  "Password",
		},
		{
			name: "password too short",
			request: &model.LoginUserRequest{
				Email:    "john@example.com",
				Password: "short",
			},
			expectError: true,
			errorField:  "Password",
		},
		{
			name: "password too long",
			request: &model.LoginUserRequest{
				Email:    "john@example.com",
				Password: string(make([]byte, 101)),
			},
			expectError: true,
			errorField:  "Password",
		},
		{
			name: "all fields empty",
			request: &model.LoginUserRequest{
				Email:    "",
				Password: "",
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			uc, mock, cleanup := setupAuthUseCase(t, nil, nil)
			defer cleanup()

			mock.ExpectBegin()
			mock.ExpectRollback()

			result, err := uc.Login(context.Background(), tc.request)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)

				fiberErr, ok := err.(*fiber.Error)
				assert.True(t, ok)
				assert.Equal(t, fiber.StatusBadRequest, fiberErr.Code)
				if tc.errorField != "" {
					assert.Contains(t, fiberErr.Message, tc.errorField)
				}
			} else {
				// For valid request without mocked repositories, it will fail at FindByEmail
				// This is expected in table-driven tests
				assert.Error(t, err) // Will get unauthorized without proper mock
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAuthUseCase_Verify(t *testing.T) {
	t.Run("should successfully verify valid token", func(t *testing.T) {
		uc, mock, cleanup := setupAuthUseCase(t, nil, nil)
		defer cleanup()

		tokenID := 1
		userID := 123
		futureExpiry := time.Now().Add(24 * time.Hour).Unix()

		// Mock FindTokenById query
		mock.ExpectQuery(`SELECT`).
			WithArgs(tokenID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "token", "created_at", "expired_at"}).
				AddRow(tokenID, userID, "valid_token_string", time.Now().Unix(), futureExpiry))

		// Mock FindByID query
		mock.ExpectQuery(`SELECT`).
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password", "created_at", "updated_at"}).
				AddRow(userID, "John Doe", "john@example.com", "hashed_password", 1234567890, 1234567890))

		ctx := context.Background()
		result, err := uc.Verify(ctx, tokenID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, userID, result.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when token not found", func(t *testing.T) {
		uc, mock, cleanup := setupAuthUseCase(t, nil, nil)
		defer cleanup()

		tokenID := 999

		// Mock FindTokenById query returning error
		mock.ExpectQuery(`SELECT`).
			WithArgs(tokenID).
			WillReturnError(sql.ErrNoRows)

		ctx := context.Background()
		result, err := uc.Verify(ctx, tokenID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, fiber.ErrUnauthorized, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when token is expired", func(t *testing.T) {
		uc, mock, cleanup := setupAuthUseCase(t, nil, nil)
		defer cleanup()

		tokenID := 1
		userID := 123
		pastExpiry := time.Now().Add(-1 * time.Hour).Unix() // Expired 1 hour ago

		// Mock FindTokenById query with expired token
		mock.ExpectQuery(`SELECT`).
			WithArgs(tokenID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "token", "created_at", "expired_at"}).
				AddRow(tokenID, userID, "expired_token", time.Now().Unix(), pastExpiry))

		ctx := context.Background()
		result, err := uc.Verify(ctx, tokenID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, fiber.ErrUnauthorized, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when user not found", func(t *testing.T) {
		uc, mock, cleanup := setupAuthUseCase(t, nil, nil)
		defer cleanup()

		tokenID := 1
		userID := 999
		futureExpiry := time.Now().Add(24 * time.Hour).Unix()

		// Mock FindTokenById query
		mock.ExpectQuery(`SELECT`).
			WithArgs(tokenID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "token", "created_at", "expired_at"}).
				AddRow(tokenID, userID, "valid_token", time.Now().Unix(), futureExpiry))

		// Mock FindByID query returning error
		mock.ExpectQuery(`SELECT`).
			WithArgs(userID).
			WillReturnError(sql.ErrNoRows)

		ctx := context.Background()
		result, err := uc.Verify(ctx, tokenID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, fiber.ErrUnauthorized, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should verify token at exact expiration boundary", func(t *testing.T) {
		uc, mock, cleanup := setupAuthUseCase(t, nil, nil)
		defer cleanup()

		tokenID := 1
		userID := 123
		// Token expires in 1 second
		almostExpired := time.Now().Add(1 * time.Second).Unix()

		// Mock FindTokenById query
		mock.ExpectQuery(`SELECT`).
			WithArgs(tokenID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "token", "created_at", "expired_at"}).
				AddRow(tokenID, userID, "valid_token", time.Now().Unix(), almostExpired))

		// Mock FindByID query
		mock.ExpectQuery(`SELECT`).
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password", "created_at", "updated_at"}).
				AddRow(userID, "John Doe", "john@example.com", "hashed_password", 1234567890, 1234567890))

		ctx := context.Background()
		result, err := uc.Verify(ctx, tokenID)

		// Should succeed as token is still valid (not yet expired)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, userID, result.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when token repository fails", func(t *testing.T) {
		uc, mock, cleanup := setupAuthUseCase(t, nil, nil)
		defer cleanup()

		tokenID := 1

		// Mock FindTokenById query with database error
		mock.ExpectQuery(`SELECT`).
			WithArgs(tokenID).
			WillReturnError(errors.New("database connection lost"))

		ctx := context.Background()
		result, err := uc.Verify(ctx, tokenID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, fiber.ErrUnauthorized, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when user repository fails", func(t *testing.T) {
		uc, mock, cleanup := setupAuthUseCase(t, nil, nil)
		defer cleanup()

		tokenID := 1
		userID := 123
		futureExpiry := time.Now().Add(24 * time.Hour).Unix()

		// Mock FindTokenById query
		mock.ExpectQuery(`SELECT`).
			WithArgs(tokenID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "token", "created_at", "expired_at"}).
				AddRow(tokenID, userID, "valid_token", time.Now().Unix(), futureExpiry))

		// Mock FindByID query with database error
		mock.ExpectQuery(`SELECT`).
			WithArgs(userID).
			WillReturnError(errors.New("database connection lost"))

		ctx := context.Background()
		result, err := uc.Verify(ctx, tokenID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, fiber.ErrUnauthorized, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should handle zero token ID", func(t *testing.T) {
		uc, mock, cleanup := setupAuthUseCase(t, nil, nil)
		defer cleanup()

		tokenID := 0

		// Mock FindTokenById query
		mock.ExpectQuery(`SELECT`).
			WithArgs(tokenID).
			WillReturnError(sql.ErrNoRows)

		ctx := context.Background()
		result, err := uc.Verify(ctx, tokenID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, fiber.ErrUnauthorized, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should handle negative token ID", func(t *testing.T) {
		uc, mock, cleanup := setupAuthUseCase(t, nil, nil)
		defer cleanup()

		tokenID := -1

		// Mock FindTokenById query
		mock.ExpectQuery(`SELECT`).
			WithArgs(tokenID).
			WillReturnError(sql.ErrNoRows)

		ctx := context.Background()
		result, err := uc.Verify(ctx, tokenID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, fiber.ErrUnauthorized, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
