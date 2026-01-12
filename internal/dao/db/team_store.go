package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/conbanwa/todo/internal/dao/cache"
	"github.com/conbanwa/todo/internal/model"
)

// CreateTeam inserts a new team into the database
func (s *SQLiteStore) CreateTeam(team *model.Team) (int64, error) {
	if team.Name == "" {
		return 0, fmt.Errorf("name is required")
	}

	query := `
	INSERT INTO teams (name, description, created_at, updated_at)
	VALUES (?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`
	result, err := s.db.Exec(query, team.Name, team.Description)
	if err != nil {
		return 0, fmt.Errorf("failed to create team: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}

	team.ID = id
	return id, nil
}

// GetTeamByID retrieves a team by ID
func (s *SQLiteStore) GetTeamByID(id int64) (*model.Team, error) {
	query := `
	SELECT id, name, description, created_at, updated_at
	FROM teams
	WHERE id = ?
	`
	row := s.db.QueryRow(query, id)

	var team model.Team
	var createdAtStr, updatedAtStr string

	err := row.Scan(&team.ID, &team.Name, &team.Description, &createdAtStr, &updatedAtStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, cache.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get team: %w", err)
	}

	// Parse timestamps
	if createdAtStr != "" {
		if t, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
			team.CreatedAt = t
		} else if t, err := time.Parse("2006-01-02 15:04:05", createdAtStr); err == nil {
			team.CreatedAt = t
		}
	}
	if updatedAtStr != "" {
		if t, err := time.Parse(time.RFC3339, updatedAtStr); err == nil {
			team.UpdatedAt = t
		} else if t, err := time.Parse("2006-01-02 15:04:05", updatedAtStr); err == nil {
			team.UpdatedAt = t
		}
	}

	return &team, nil
}

// AddUserToTeam adds a user to a team with a specific role
func (s *SQLiteStore) AddUserToTeam(userID, teamID int64, role model.Role) error {
	if !role.IsValid() {
		return fmt.Errorf("invalid role: %v", role)
	}

	query := `
	INSERT OR REPLACE INTO user_teams (user_id, team_id, role, created_at)
	VALUES (?, ?, ?, CURRENT_TIMESTAMP)
	`
	_, err := s.db.Exec(query, userID, teamID, string(role))
	if err != nil {
		return fmt.Errorf("failed to add user to team: %w", err)
	}

	return nil
}

// RemoveUserFromTeam removes a user from a team
func (s *SQLiteStore) RemoveUserFromTeam(userID, teamID int64) error {
	query := `DELETE FROM user_teams WHERE user_id = ? AND team_id = ?`
	result, err := s.db.Exec(query, userID, teamID)
	if err != nil {
		return fmt.Errorf("failed to remove user from team: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return cache.ErrNotFound
	}

	return nil
}

// GetUserTeams retrieves all teams that a user belongs to
func (s *SQLiteStore) GetUserTeams(userID int64) ([]model.Team, error) {
	query := `
	SELECT t.id, t.name, t.description, t.created_at, t.updated_at
	FROM teams t
	INNER JOIN user_teams ut ON t.id = ut.team_id
	WHERE ut.user_id = ?
	ORDER BY t.id ASC
	`
	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user teams: %w", err)
	}
	defer rows.Close()

	var teams []model.Team
	for rows.Next() {
		var team model.Team
		var createdAtStr, updatedAtStr string

		err := rows.Scan(&team.ID, &team.Name, &team.Description, &createdAtStr, &updatedAtStr)
		if err != nil {
			return nil, fmt.Errorf("failed to scan team: %w", err)
		}

		// Parse timestamps
		if createdAtStr != "" {
			if t, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
				team.CreatedAt = t
			} else if t, err := time.Parse("2006-01-02 15:04:05", createdAtStr); err == nil {
				team.CreatedAt = t
			}
		}
		if updatedAtStr != "" {
			if t, err := time.Parse(time.RFC3339, updatedAtStr); err == nil {
				team.UpdatedAt = t
			} else if t, err := time.Parse("2006-01-02 15:04:05", updatedAtStr); err == nil {
				team.UpdatedAt = t
			}
		}

		teams = append(teams, team)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return teams, nil
}

// GetTeamMembers retrieves all users in a team with their roles
func (s *SQLiteStore) GetTeamMembers(teamID int64) ([]model.UserTeam, error) {
	query := `
	SELECT user_id, team_id, role
	FROM user_teams
	WHERE team_id = ?
	ORDER BY user_id ASC
	`
	rows, err := s.db.Query(query, teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to get team members: %w", err)
	}
	defer rows.Close()

	var members []model.UserTeam
	for rows.Next() {
		var ut model.UserTeam
		var roleStr string

		err := rows.Scan(&ut.UserID, &ut.TeamID, &roleStr)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user_team: %w", err)
		}

		ut.Role = model.Role(roleStr)
		members = append(members, ut)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return members, nil
}

// GetUserRoleInTeam retrieves a user's role in a team
func (s *SQLiteStore) GetUserRoleInTeam(userID, teamID int64) (model.Role, error) {
	query := `
	SELECT role
	FROM user_teams
	WHERE user_id = ? AND team_id = ?
	`
	row := s.db.QueryRow(query, userID, teamID)

	var roleStr string
	err := row.Scan(&roleStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", cache.ErrNotFound
		}
		return "", fmt.Errorf("failed to get user role in team: %w", err)
	}

	return model.Role(roleStr), nil
}

// UpdateTeam updates an existing team
func (s *SQLiteStore) UpdateTeam(team *model.Team) error {
	if team.ID == 0 {
		return fmt.Errorf("id is required")
	}

	// Check if team exists
	_, err := s.GetTeamByID(team.ID)
	if err != nil {
		return err
	}

	query := `
	UPDATE teams
	SET name = ?, description = ?, updated_at = CURRENT_TIMESTAMP
	WHERE id = ?
	`
	_, err = s.db.Exec(query, team.Name, team.Description, team.ID)
	if err != nil {
		return fmt.Errorf("failed to update team: %w", err)
	}

	return nil
}
