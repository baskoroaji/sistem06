package user

type UserEntity struct {
	ID        int
	Name      string
	Email     string
	Password  string
	CreatedAt int64
	UpdatedAt int64
}

type UserWithRole struct {
	ID       int
	Email    string
	Name     string
	Password string
	Roles    []RolesWithPermissions
}

type Role struct {
	ID   int
	Name string
}

type Permissions struct {
	ID   int
	Name string
}

type RolesPermissions struct {
	ID            int
	RolesID       int
	PermissionsID int
}

type UserRoles struct {
	ID      int
	UserID  int
	RolesID int
}

type RolesWithPermissions struct {
	ID          int
	Name        string
	Permissions []string
}
