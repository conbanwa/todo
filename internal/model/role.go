package model

// Role represents a user's role in a team
type Role string

const (
	RoleAdmin  Role = "admin"
	RoleMember Role = "member"
	RoleViewer Role = "viewer"
)

// Permission represents what actions a user can perform
type Permission string

const (
	PermissionCreate Permission = "create"
	PermissionRead   Permission = "read"
	PermissionUpdate Permission = "update"
	PermissionDelete Permission = "delete"
)

// IsValid checks if the role is a valid role
func (r Role) IsValid() bool {
	return r == RoleAdmin || r == RoleMember || r == RoleViewer
}

// HasPermission checks if the role has the specified permission
// Admin has all permissions
// Member can create, read, update (but not delete)
// Viewer can only read
func (r Role) HasPermission(perm Permission) bool {
	if !r.IsValid() {
		return false
	}

	switch r {
	case RoleAdmin:
		return true // Admin has all permissions
	case RoleMember:
		return perm == PermissionCreate || perm == PermissionRead || perm == PermissionUpdate
	case RoleViewer:
		return perm == PermissionRead
	default:
		return false
	}
}

// UserTeam represents the many-to-many relationship between users and teams
// with a role assignment
type UserTeam struct {
	UserID int64 `json:"user_id"`
	TeamID int64 `json:"team_id"`
	Role   Role  `json:"role"`
}
