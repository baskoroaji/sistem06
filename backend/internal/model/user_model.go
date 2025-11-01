package model

type UserResponse struct {
	ID        int    `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	CreatedAt int64  `json:"created_at,omitempty"`
	UpdatedAt int64  `json:"updated_at,omitempty"`
}
type RegisterUserRequest struct {
	Name     string `json:"name" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginUserRequest struct {
	email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8 ,max=100"`
}
