package entity

type User struct {
	ID        int
	Name      string
	Email     string
	Password  string
	CreatedAt int64
	UpdatedAt int64
}

// type UserWithRole struct {
// 	User
// 	Roles []RolesWithPermissions
// }

func (u *UserWithRole) MapRole() error {

}
