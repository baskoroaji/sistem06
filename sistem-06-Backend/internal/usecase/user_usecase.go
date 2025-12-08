package user

import (
	"context"
	"database/sql"
	"strings"

	"backend-sistem06.com/internal/entity"
	"backend-sistem06.com/internal/model"
	"backend-sistem06.com/internal/model/converter"
	"backend-sistem06.com/internal/repository"
	"backend-sistem06.com/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	DB             *sql.DB
	Log            *logrus.Logger
	validate       *validator.Validate
	UserRepository repository.UserRepositoryInterface
}

func NewUserUseCase(db *sql.DB, log *logrus.Logger, validate *validator.Validate, userRepository repository.UserRepositoryInterface) *UserUseCase {
	return &UserUseCase{
		DB:             db,
		Log:            log,
		validate:       validate,
		UserRepository: userRepository,
	}
}

func (c *UserUseCase) Create(ctx context.Context, request *model.RegisterUserRequest) (*model.UserResponse, error) {
	tx, err := c.DB.BeginTx(ctx, nil)
	if err != nil {
		c.Log.Warnf("Failed to begin transaction: %+v", err)

		if err == context.Canceled || err == context.DeadlineExceeded {
			return nil, fiber.NewError(fiber.StatusRequestTimeout, "request timeout or canceled")
		}

		return nil, fiber.ErrInternalServerError
	}
	defer tx.Rollback()

	if err := c.validate.Struct(request); err != nil {
		validationErrors := utils.ValidationError(err)
		c.Log.Warnf("Validation failed: %+v", validationErrors)

		// Return response error yang bisa dibaca frontend
		return nil, fiber.NewError(fiber.StatusBadRequest, utils.FormatValidationErrors(validationErrors))
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Log.Warnf("failed to generate bcrypt hash: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	user := &entity.UserEntity{
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
