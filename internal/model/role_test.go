package model

import "testing"

func TestRole_String(t *testing.T) {
	tests := []struct {
		name string
		role Role
		want string
	}{
		{"Admin role", RoleAdmin, "admin"},
		{"Member role", RoleMember, "member"},
		{"Viewer role", RoleViewer, "viewer"},
		{"Empty role", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := string(tt.role); got != tt.want {
				t.Errorf("Role.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRole_IsValid(t *testing.T) {
	tests := []struct {
		name string
		role Role
		want bool
	}{
		{"Valid admin", RoleAdmin, true},
		{"Valid member", RoleMember, true},
		{"Valid viewer", RoleViewer, true},
		{"Invalid role", Role("invalid"), false},
		{"Empty role", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.role.IsValid(); got != tt.want {
				t.Errorf("Role.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRole_HasPermission(t *testing.T) {
	tests := []struct {
		name       string
		role       Role
		permission Permission
		want       bool
	}{
		{"Admin can create", RoleAdmin, PermissionCreate, true},
		{"Admin can read", RoleAdmin, PermissionRead, true},
		{"Admin can update", RoleAdmin, PermissionUpdate, true},
		{"Admin can delete", RoleAdmin, PermissionDelete, true},
		{"Member can create", RoleMember, PermissionCreate, true},
		{"Member can read", RoleMember, PermissionRead, true},
		{"Member can update", RoleMember, PermissionUpdate, true},
		{"Member cannot delete", RoleMember, PermissionDelete, false},
		{"Viewer can read", RoleViewer, PermissionRead, true},
		{"Viewer cannot create", RoleViewer, PermissionCreate, false},
		{"Viewer cannot update", RoleViewer, PermissionUpdate, false},
		{"Viewer cannot delete", RoleViewer, PermissionDelete, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.role.HasPermission(tt.permission); got != tt.want {
				t.Errorf("Role.HasPermission() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserTeam_Role(t *testing.T) {
	ut := UserTeam{
		UserID: 1,
		TeamID: 1,
		Role:   RoleAdmin,
	}

	if ut.Role != RoleAdmin {
		t.Errorf("UserTeam.Role = %v, want %v", ut.Role, RoleAdmin)
	}
}
