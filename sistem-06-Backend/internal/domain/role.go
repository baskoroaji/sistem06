package domain

type Role struct {
	ID   int
	Name string
}

type RolesWithPermissions struct {
	Role
	Permissions []string
}
