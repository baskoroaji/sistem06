package usecase

import (
	"context"
	"database/sql"
	"time"

	"backend-sistem06.com/internal/entity"
	"backend-sistem06.com/internal/model"
	"backend-sistem06.com/internal/repository"
	"backend-sistem06.com/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase struct {
	DB              *sql.DB
	Log             *logrus.Logger
	Validate        *validator.Validate
	UserRepository  *repository.UserRepository
	TokenRepository *repository.TokenRepository
}

func NewAuthUseCase(db *sql.DB, log *logrus.Logger, validate *validator.Validate, userRepository *repository.UserRepository, token *repository.TokenRepository) *AuthUseCase {
	return &AuthUseCase{
		DB:              db,
		Log:             log,
		Validate:        validate,
		UserRepository:  userRepository,
		TokenRepository: token,
	}
}

func (c *AuthUseCase) Login(ctx context.Context, request *model.LoginUserRequest) (*model.LoginResponse, error) {

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

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		c.Log.Warn("Invalid password")
		return nil, fiber.ErrUnauthorized
	}

	tokenValue := utils.GenerateToken()
	expiration := time.Now().Add(24 * time.Hour).Unix()

	token := &entity.PersonalAccessToken{
		UserID:    user.ID,
		Token:     tokenValue,
		CreatedAt: time.Now().Unix(),
		ExpiredAt: expiration,
	}
	if err := c.TokenRepository.CreateToken(tx, token); err != nil {
		c.Log.Warnf("Failed to create token: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit(); err != nil {
		c.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return &model.LoginResponse{
		Token: token.Token,
	}, nil
}

func (c *AuthUseCase) Verify(ctx *context.Context, request *model.VerifyUserRequest) *model
