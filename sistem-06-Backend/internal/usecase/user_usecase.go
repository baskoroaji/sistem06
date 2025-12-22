package usecase

import (
	"context"
	"database/sql"
	"strings"

	"sistem-06-Backend/internal/domain/entity"
	"sistem-06-Backend/internal/domain/ports"
	"sistem-06-Backend/internal/dto"
	"sistem-06-Backend/internal/pkg/errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type userUseCase struct {
	DB       *sql.DB
	Log      *logrus.Logger
	Validate *validator.Validate
	Repo     ports.UserRepository
}

func NewUserUseCase(db *sql.DB, log *logrus.Logger, Validate *validator.Validate, repo ports.UserRepository) *userUseCase {
	return &userUseCase{
		Log:      log,
		Validate: Validate,
		Repo:     repo,
	}
}

func (c *userUseCase) Create(ctx context.Context, request *dto.RegisterUserRequest) (*dto.UserResponse, error) {
	if err := c.Validate.Struct(request); err != nil {
		validationErrors := errors.ValidationError(err, errors.UserErrorMessages)
		c.Log.Warnf("Validation failed: %+v", validationErrors)

		// Return response error yang bisa dibaca frontend
		return nil, fiber.NewError(fiber.StatusBadRequest, errors.FormatValidationErrors(validationErrors))
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Log.Warnf("failed to generate bcrypt hash: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	user := &entity.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: string(password),
	}

	if err := c.UserRepository.CreateUser(tx, user); err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return nil, fiber.NewError(fiber.StatusConflict, "email or name already exist")
		}
		c.Log.Warnf("Database insert error: %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	if err := tx.Commit(); err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	return converter.UserToResponse(user), nil
}
