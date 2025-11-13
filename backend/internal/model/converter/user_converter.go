package converter

import (
	"backend-sistem06.com/internal/entity"
	"backend-sistem06.com/internal/model"
)

func UserToResponse(user *entity.UserEntity) *model.UserResponse {
	return &model.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func UserWithRoleToResponse(user *entity.UserWithRole) *model.UserResponse {
	if user == nil {
		return nil
	}

	var roles []string
	for _, r := range user.RoleName {
		roles = append(roles, r.Name)
	}

	return &model.UserResponse{
		ID:    int(user.ID),
		Name:  user.Name,
		Email: user.Email,
		Roles: roles,
	}
}
