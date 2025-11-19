package converter

import (
	"backend-sistem06.com/internal/entity"
	"backend-sistem06.com/internal/model"
)

func UserToResponse(user *entity.UserEntity) *model.UserResponse {
	return &model.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,

		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func UserWithRolesToResponse(user *entity.UserWithRole) *model.UserResponse {
	if user == nil {
		return nil
	}

	// Convert roles
	roles := make([]model.RoleResponse, len(user.Roles))
	permissionSet := make(map[string]bool)

	for i, role := range user.Roles {
		roles[i] = model.RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			Permissions: role.Permissions,
		}

		// Collect all unique permissions
		for _, perm := range role.Permissions {
			permissionSet[perm] = true
		}
	}

	return &model.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Roles: roles,
	}
}
