package security

type Permission string

const (
	PermissionCreateUser Permission = "create-user"
	PermissionEditUser   Permission = "edit-user"
	PermissionDeleteUser Permission = "delete-user"

	PermissionCreateRole     Permission = "create-role"
	PermissionEditRole       Permission = "edit-role"
	PermissionDeleteRole     Permission = "delete-role"
	PermissionAssignUserRole Permission = "assign-user-role"
	PermissionRemoveUserRole Permission = "remove-user-role"
)

func (p Permission) String() string {
	return string(p)
}

type Role string

const (
	RoleAdmin Role = "Admin"
)

func (r Role) String() string {
	return string(r)
}
