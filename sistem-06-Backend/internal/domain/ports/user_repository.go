package domain

import "context"

type UserRepository interface {
	CreateUser(ctx context.Context, user *UserEntity) error
	CountById(ctx context.Context, id int) (int, error)
	CountByName(ctx context.Context, name string) (int, error)
	FindByEmail(ctx context.Context, email string) (*UserEntity, error)
	FindByID(ctx context.Context, id int) (*UserEntity, error)
}
