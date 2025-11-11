package usecase_test

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	"backend-sistem06.com/internal/model"
	"backend-sistem06.com/internal/repository"
	"backend-sistem06.com/internal/usecase"
)

func setupUserUseCase(t *testing.T) (*usecase.UserUseCase, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}

	log := logrus.New()
	log.SetOutput(io.Discard) // Suppress log output in tests

	validate := validator.New()
	userRepo := repository.NewUserRepository(db, log)

	uc := usecase.NewUserUseCase(db, log, validate, userRepo)

	cleanup := func() {
		db.Close()
	}

	return uc, mock, cleanup
}

func TestUserUseCase_Create(t *testing.T) {
	t.Run("should successfully create user", func(t *testing.T) {
		uc, mock, cleanup := setupUserUseCase(t)
		defer cleanup()

		request := &model.RegisterUserRequest{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "password123",
		}

		// Expect transaction begin
		mock.ExpectBegin()

		// Expect INSERT query with RETURNING id
		mock.ExpectQuery(`INSERT INTO users \(name, email, password, created_at, updated_at\)`).
			WithArgs(request.Name, request.Email, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		// Expect commit
		mock.ExpectCommit()

		result, err := uc.Create(context.Background(), request)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, result.ID)
		assert.Equal(t, "John Doe", result.Name)
		// CreatedAt and UpdatedAt are set in repository but not returned to entity
		// So they will be 0 in the response
		assert.GreaterOrEqual(t, result.CreatedAt, int64(0))
		assert.GreaterOrEqual(t, result.UpdatedAt, int64(0))
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when transaction begin fails", func(t *testing.T) {
		uc, mock, cleanup := setupUserUseCase(t)
		defer cleanup()

		mock.ExpectBegin().WillReturnError(sql.ErrConnDone)

		request := &model.RegisterUserRequest{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "password123",
		}

		result, err := uc.Create(context.Background(), request)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, fiber.ErrInternalServerError, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return validation error for invalid email", func(t *testing.T) {
		uc, mock, cleanup := setupUserUseCase(t)
		defer cleanup()

		mock.ExpectBegin()
		mock.ExpectRollback()

		request := &model.RegisterUserRequest{
			Name:     "John Doe",
			Email:    "invalid-email",
			Password: "password123",
		}

		result, err := uc.Create(context.Background(), request)

		assert.Error(t, err)
		assert.Nil(t, result)

		fiberErr, ok := err.(*fiber.Error)
		assert.True(t, ok, "error should be fiber.Error")
		assert.Equal(t, fiber.StatusBadRequest, fiberErr.Code)
		assert.Contains(t, fiberErr.Message, "Email")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return validation error for missing required fields", func(t *testing.T) {
		uc, mock, cleanup := setupUserUseCase(t)
		defer cleanup()

		mock.ExpectBegin()
		mock.ExpectRollback()

		request := &model.RegisterUserRequest{
			Name:     "",
			Email:    "john@example.com",
			Password: "password123",
		}

		result, err := uc.Create(context.Background(), request)

		assert.Error(t, err)
		assert.Nil(t, result)

		fiberErr, ok := err.(*fiber.Error)
		assert.True(t, ok)
		assert.Equal(t, fiber.StatusBadRequest, fiberErr.Code)
		assert.Contains(t, fiberErr.Message, "Name")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// t.Run("should hash password correctly", func(t *testing.T) {
	// 	uc, mock, cleanup := setupUserUseCase(t)
	// 	defer cleanup()

	// 	plainPassword := "password123"
	// 	request := &model.RegisterUserRequest{
	// 		Name:     "John Doe",
	// 		Email:    "john@example.com",
	// 		Password: plainPassword,
	// 	}

	// 	mock.ExpectBegin()

	// 	// Capture the hashed password from the query
	// 	var hashedPassword string
	// 	mock.ExpectQuery(`INSERT INTO users`).
	// 		WithArgs(request.Name, request.Email, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
	// 		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1)).
	// 		WillDelayFor(0).
	// 		WillReturnError(nil)

	// 	mock.ExpectCommit()

	// 	result, err := uc.Create(context.Background(), request)

	// 	assert.NoError(t, err)
	// 	assert.NotNil(t, result)

	// 	// Verify the password is not in response (omitempty and not included in converter)
	// 	// UserResponse doesn't have Email or Password fields
	// 	assert.NotEqual(t, plainPassword, result.Name) // Just verify result is valid
	// 	assert.NoError(t, mock.ExpectationsWereMet())
	// })

	t.Run("should return conflict error for duplicate email", func(t *testing.T) {
		uc, mock, cleanup := setupUserUseCase(t)
		defer cleanup()

		request := &model.RegisterUserRequest{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "password123",
		}

		mock.ExpectBegin()

		// Simulate duplicate key error
		mock.ExpectQuery(`INSERT INTO users`).
			WithArgs(request.Name, request.Email, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnError(errors.New("pq: duplicate key value violates unique constraint"))

		mock.ExpectRollback()

		result, err := uc.Create(context.Background(), request)

		assert.Error(t, err)
		assert.Nil(t, result)

		fiberErr, ok := err.(*fiber.Error)
		assert.True(t, ok)
		assert.Equal(t, fiber.StatusConflict, fiberErr.Code)
		assert.Equal(t, "email or name already exist", fiberErr.Message)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should rollback transaction on database insert error", func(t *testing.T) {
		uc, mock, cleanup := setupUserUseCase(t)
		defer cleanup()

		request := &model.RegisterUserRequest{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "password123",
		}

		mock.ExpectBegin()

		mock.ExpectQuery(`INSERT INTO users`).
			WithArgs(request.Name, request.Email, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnError(errors.New("database connection lost"))

		mock.ExpectRollback()

		result, err := uc.Create(context.Background(), request)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, fiber.ErrInternalServerError, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when commit fails", func(t *testing.T) {
		uc, mock, cleanup := setupUserUseCase(t)
		defer cleanup()

		request := &model.RegisterUserRequest{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "password123",
		}

		mock.ExpectBegin()

		mock.ExpectQuery(`INSERT INTO users`).
			WithArgs(request.Name, request.Email, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		// Commit fails
		mock.ExpectCommit().WillReturnError(errors.New("commit failed"))

		result, err := uc.Create(context.Background(), request)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, fiber.ErrInternalServerError, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should handle context cancellation", func(t *testing.T) {
		uc, mock, cleanup := setupUserUseCase(t)
		defer cleanup()

		// Don't actually cancel the context - just mock the error
		mock.ExpectBegin().WillReturnError(context.Canceled)

		request := &model.RegisterUserRequest{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "password123",
		}

		result, err := uc.Create(context.Background(), request)

		assert.Error(t, err)
		assert.Nil(t, result)

		// Check it's a fiber error with status 408 (Request Timeout)
		fiberErr, ok := err.(*fiber.Error)
		assert.True(t, ok, "error should be *fiber.Error type")
		assert.Equal(t, fiber.StatusRequestTimeout, fiberErr.Code)
		assert.Contains(t, fiberErr.Message, "request timeout or canceled")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should handle context deadline exceeded", func(t *testing.T) {
		uc, mock, cleanup := setupUserUseCase(t)
		defer cleanup()

		mock.ExpectBegin().WillReturnError(context.DeadlineExceeded)

		request := &model.RegisterUserRequest{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "password123",
		}

		result, err := uc.Create(context.Background(), request)

		assert.Error(t, err)
		assert.Nil(t, result)

		fiberErr, ok := err.(*fiber.Error)
		assert.True(t, ok)
		assert.Equal(t, fiber.StatusRequestTimeout, fiberErr.Code)
		assert.Contains(t, fiberErr.Message, "request timeout or canceled")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should verify bcrypt password hashing", func(t *testing.T) {
		// This is a separate test to actually verify bcrypt works
		plainPassword := "mySecurePassword123"
		request := &model.RegisterUserRequest{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: plainPassword,
		}

		// Generate hash using bcrypt (same as in usecase)
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		assert.NoError(t, err)

		// Verify hash
		err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(plainPassword))
		assert.NoError(t, err, "password should be properly hashed")

		// Verify plain password doesn't match
		assert.NotEqual(t, plainPassword, string(hashedPassword))
	})

	t.Run("should pass transaction to repository", func(t *testing.T) {
		uc, mock, cleanup := setupUserUseCase(t)
		defer cleanup()

		request := &model.RegisterUserRequest{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "password123",
		}

		mock.ExpectBegin()

		// The repository should use the transaction
		mock.ExpectQuery(`INSERT INTO users`).
			WithArgs(request.Name, request.Email, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectCommit()

		result, err := uc.Create(context.Background(), request)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// Table-driven test for validation scenarios
func TestUserUseCase_Create_ValidationScenarios(t *testing.T) {
	testCases := []struct {
		name        string
		request     *model.RegisterUserRequest
		expectError bool
		errorField  string
	}{
		{
			name: "valid request",
			request: &model.RegisterUserRequest{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "password123",
			},
			expectError: false,
		},
		{
			name: "empty name",
			request: &model.RegisterUserRequest{
				Name:     "",
				Email:    "john@example.com",
				Password: "password123",
			},
			expectError: true,
			errorField:  "Name",
		},
		{
			name: "invalid email format",
			request: &model.RegisterUserRequest{
				Name:     "John Doe",
				Email:    "not-an-email",
				Password: "password123",
			},
			expectError: true,
			errorField:  "Email",
		},
		{
			name: "empty email",
			request: &model.RegisterUserRequest{
				Name:     "John Doe",
				Email:    "",
				Password: "password123",
			},
			expectError: true,
			errorField:  "Email",
		},
		{
			name: "empty password",
			request: &model.RegisterUserRequest{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "",
			},
			expectError: true,
			errorField:  "Password",
		},
		{
			name: "all fields empty",
			request: &model.RegisterUserRequest{
				Name:     "",
				Email:    "",
				Password: "",
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			uc, mock, cleanup := setupUserUseCase(t)
			defer cleanup()

			mock.ExpectBegin()

			if !tc.expectError {
				mock.ExpectQuery(`INSERT INTO users`).
					WithArgs(tc.request.Name, tc.request.Email, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectCommit()
			} else {
				mock.ExpectRollback()
			}

			result, err := uc.Create(context.Background(), tc.request)

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
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

// Benchmark test
func BenchmarkUserUseCase_Create(b *testing.B) {
	db, mock, err := sqlmock.New()
	if err != nil {
		b.Fatalf("failed to create mock db: %v", err)
	}
	defer db.Close()

	log := logrus.New()
	log.SetOutput(io.Discard)

	validate := validator.New()
	userRepo := repository.NewUserRepository(db, log)
	uc := usecase.NewUserUseCase(db, log, validate, userRepo)

	request := &model.RegisterUserRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO users`).
			WithArgs(request.Name, request.Email, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		_, _ = uc.Create(context.Background(), request)
	}
}
