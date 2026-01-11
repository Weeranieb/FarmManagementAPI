package constants

const (
	// UserLevelNormal - Normal user with limited permissions within their client
	UserLevelNormal = iota + 1

	// UserLevelClientAdmin - Client admin with full permissions within their client
	UserLevelClientAdmin

	// UserLevelSuperAdmin - Super admin with full permissions across all clients
	UserLevelSuperAdmin
)

// IsSuperAdmin checks if user level is super admin
func IsSuperAdmin(userLevel int) bool {
	return userLevel >= UserLevelSuperAdmin
}

// IsClientAdminOrAbove checks if user is client admin or super admin
func IsClientAdminOrAbove(userLevel int) bool {
	return userLevel >= UserLevelClientAdmin
}

// CanAccessClient checks if user can access a specific client
// Super admins can access all clients, others only their own
func CanAccessClient(userLevel int, userClientId int, targetClientId int) bool {
	if IsSuperAdmin(userLevel) {
		return true
	}
	return userClientId == targetClientId
}
