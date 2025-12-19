package domain

import "context"

type PermissionRepository interface {
	GetPermissionsByRoleID(ctx context.Context, id int) (*UserWithRole, error)
	GetPermissionsByUserID(ctx context.Context, id int) (*UserWithRole, error)
}
