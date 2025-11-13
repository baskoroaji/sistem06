package usecase

import (
	"database/sql"

	"backend-sistem06.com/internal/model"
	"backend-sistem06.com/internal/repository"
	"backend-sistem06.com/utils"
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

func (c *AuthUseCase) Login(ctx *fiber.Ctx, request *model.LoginUserRequest) (*model.LoginResponse, error) {
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

	return &model.LoginResponse{
		Message: sess.ID(),
	}, nil
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
