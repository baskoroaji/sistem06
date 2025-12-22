package ports

import (
	"context"
	"sistem-06-Backend/internal/domain/entity"
)

type PermissionRepository interface {
	GetPermissionsByRoleID(ctx context.Context, id int) (*entity.Role, error)
	GetPermissionsByUserID(ctx context.Context, id int) (*entity.User, error)
}
