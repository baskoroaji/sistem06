package usecase

import (
	"context"
	"database/sql"

	"backend-sistem06.com/internal/model"
	"backend-sistem06.com/internal/repository"
	"backend-sistem06.com/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase struct {
	DB             *sql.DB
	Log            *logrus.Logger
	Validate       *validator.Validate
	UserRepository *repository.UserRepository
}

func NewAuthUseCase(db *sql.DB, log *logrus.Logger, validate *validator.Validate, userRepository *repository.UserRepository) *AuthUseCase {
	return &AuthUseCase{
		DB:             db,
		Log:            log,
		Validate:       validate,
		UserRepository: userRepository,
	}
}

func (c *AuthUseCase) Login(ctx context.Context, request *model.LoginUserRequest) (*model.UserResponse, error) {

	tx, err := c.DB.BeginTx(ctx, nil)
	if err != nil {
		c.Log.Warnf("Failed to begin transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		validationErrors := utils.ValidationError(err)

		c.Log.Warnf("Validation failed: %+v", validationErrors)

		return nil, fiber.NewError(fiber.StatusBadRequest, utils.FormatValidationErrors(validationErrors))
	}

	user, err := c.UserRepository.FindByEmail(request.Email)
	if err != nil {
		c.Log.Warnf("Failed find user by email: %+v", err)
		return nil, fiber.ErrUnauthorized
	}

	password, err := bcrypt.CompareHashAndPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Log.Warnf("failed to generate bcrypt hash: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

}
