package entity

type Role struct {
	ID   int
	Name string
}

type RolesWithPermissions struct {
	Roles       string
	Permissions []string
}
