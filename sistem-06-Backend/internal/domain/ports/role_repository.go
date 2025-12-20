package domain

import "context"

type RolesRepository interface {
	GetRolesByUserID(ctx context.Context, id int) (*Role, error)
	GetRolesWithPermissionsByUserID(ctx context.Context, id int) (*UserWithRole, error)
	AssignRoleToUser(ctx context.Context, userId int, rolesId int) error
	RemoveRoleFromUser(ctx context.Context, userId int, rolesId int) error
}
