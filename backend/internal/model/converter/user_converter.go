package converter

import (
	"backend-sistem06.com/internal/entity"
	"backend-sistem06.com/internal/model"
)

func UserToResponse(user *entity.UserEntity) *model.UserResponse {
	return &model.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
