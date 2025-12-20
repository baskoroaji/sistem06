package entity

func (r *Role) HasPermission(p Permissions) bool {
	for _, perm := range r.Permission {
		if perm == p {
			return true
		}
	}
	return false
}
