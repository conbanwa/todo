package api

import "github.com/conbanwa/todo/internal/model"

// TeamStore defines the interface for team storage operations
type TeamStore interface {
	Create(team *model.Team) (int64, error)
	GetByID(id int64) (*model.Team, error)
	GetUserTeams(userID int64) ([]model.Team, error)
	AddUserToTeam(userID, teamID int64, role model.Role) error
	RemoveUserFromTeam(userID, teamID int64) error
	GetTeamMembers(teamID int64) ([]model.UserTeam, error)
	GetUserRoleInTeam(userID, teamID int64) (model.Role, error)
	Update(team *model.Team) error
}
