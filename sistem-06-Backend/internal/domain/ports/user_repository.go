package ports

import (
	"context"
	"sistem-06-Backend/internal/domain/entity"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) error
	CountById(ctx context.Context, id int) (int, error)
	CountByName(ctx context.Context, name string) (int, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByID(ctx context.Context, id int) (*entity.User, error)
}
