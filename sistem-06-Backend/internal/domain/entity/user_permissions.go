package entity

func (u *User) HasPermission(p Permissions) bool {
	for _, r := range u.Role {
		for _, perm := range r.Permission {
			if perm == p {
				return true
			}
		}
	}
	return false
}
