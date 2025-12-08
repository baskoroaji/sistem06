package domain

type UserEntity struct {
	ID        int
	Name      string
	Email     string
	Password  string
	CreatedAt int64
	UpdatedAt int64
}

type UserWithRole struct {
	UserEntity
	Roles []RolesWithPermissions
}
