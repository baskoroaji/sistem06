package entity

type UserEntity struct {
	ID        int
	Name      string
	Email     string
	Password  string
	CreatedAt int64
	UpdatedAt int64
}

type UserWithRole struct {
	ID       int64
	Email    string
	Password string
	RoleID   int64
	RoleName []Role
}

type Role struct {
	ID          int64
	Name        string
	Permissions []string
}
