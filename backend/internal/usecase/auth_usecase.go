package usecase

import (
	"context"
	"database/sql"

	"backend-sistem06.com/internal/model"
	"backend-sistem06.com/internal/model/converter"
	"backend-sistem06.com/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase struct {
	DB             *sql.DB
	Log            *logrus.Logger
	Validate       *validator.Validate
	UserRepository repository.UserRepositoryInterface
	Session        *session.Store
}

func NewAuthUseCase(db *sql.DB, log *logrus.Logger, validate *validator.Validate, userRepository repository.UserRepositoryInterface, session *session.Store) *AuthUseCase {
	return &AuthUseCase{
		DB:             db,
		Log:            log,
		Validate:       validate,
		UserRepository: userRepository,
		Session:        session,
	}
}

func (c *AuthUseCase) Login(ctx context.Context, request *model.LoginUserRequest) (*model.UserResponse, error) {
	if err := c.Validate.Struct(request); err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid request")
	}

	user, err := c.UserRepository.FindByEmail(ctx, request.Email)
	if err != nil {
		return nil, fiber.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return nil, fiber.ErrUnauthorized
	}

	userWithRoles, err := c.UserRepository.FindWithRoles(ctx, user.ID)
	if err != nil {
		c.Log.Errorf("Failed to load roles: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.UserWithRolesToResponse(userWithRoles), nil
}

// func (c *AuthUseCase) Verify(ctx context.Context, tokenID int) (*model.Auth, error) {
// 	token, err := c.TokenRepository.FindTokenById(tokenID)
// 	if err != nil {
// 		c.Log.Warnf("Failed to find token: %+v", err)
// 		return nil, fiber.ErrUnauthorized
// 	}

// 	if time.Now().Unix() > token.ExpiredAt {
// 		c.Log.Warn("Token expired")
// 		return nil, fiber.ErrUnauthorized
// 	}

// 	user, err := c.UserRepository.FindByID(token.UserID)
// 	if err != nil {
// 		c.Log.Warnf("Failed to find user by token: %+v", err)
// 		return nil, fiber.ErrUnauthorized
// 	}

// 	return &model.Auth{
// 		ID: user.ID,
// 	}, nil
// }
