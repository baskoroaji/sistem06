package ports

import (
	"context"
	"sistem-06-Backend/internal/domain/entity"
)

type RolesRepository interface {
	GetRolesByUserID(ctx context.Context, id int) (*entity.Role, error)
	GetRolesWithPermissionsByUserID(ctx context.Context, id int) (*entity.Role, error)
	AssignRoleToUser(ctx context.Context, userId int, rolesId int) error
	RemoveRoleFromUser(ctx context.Context, userId int, rolesId int) error
}
