package usecase

import (
	"context"
	"database/sql"

	"backend-sistem06.com/internal/model"
	"backend-sistem06.com/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
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

func (l *AuthUseCase) Login(ctx context.Context, request *model.LoginUserRequest) (*model.UserResponse, error) {

	/* TODO create auth
	 */
}
