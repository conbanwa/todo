package auth

import (
	"fmt"
	"net/http"

	"github.com/conbanwa/todo/internal/dao/cache/api"
	"github.com/conbanwa/todo/internal/model"
	"github.com/gin-gonic/gin"
)

const (
	TeamIDKey = "team_id"
	RoleKey   = "role"
)

// RequireRole creates middleware that requires the user to have one of the specified roles
func RequireRole(teamStore api.TeamStore, roles ...model.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context (set by AuthMiddleware)
		userID, exists := GetUserID(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			c.Abort()
			return
		}

		// Get team ID from context or route parameter
		teamID, err := GetTeamID(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "team_id required"})
			c.Abort()
			return
		}

		// Get user's role in team
		userRole, err := teamStore.GetUserRoleInTeam(userID, teamID)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "user is not a member of this team"})
			c.Abort()
			return
		}

		// Check if user has one of the required roles
		hasRole := false
		for _, requiredRole := range roles {
			if userRole == requiredRole {
				hasRole = true
				break
			}
		}

		// Admin has all permissions
		if userRole == model.RoleAdmin {
			hasRole = true
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			c.Abort()
			return
		}

		// Set role and team ID in context
		c.Set(RoleKey, userRole)
		c.Set(TeamIDKey, teamID)
		c.Next()
	}
}

// RequireTeam creates middleware that requires the user to be a member of a team
func RequireTeam(teamStore api.TeamStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context (set by AuthMiddleware)
		userID, exists := GetUserID(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			c.Abort()
			return
		}

		// Get team ID from context or route parameter
		teamID, err := GetTeamID(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "team_id required"})
			c.Abort()
			return
		}

		// Verify user is a member of the team
		_, err = teamStore.GetUserRoleInTeam(userID, teamID)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "user is not a member of this team"})
			c.Abort()
			return
		}

		// Set team ID in context
		c.Set(TeamIDKey, teamID)
		c.Next()
	}
}

// GetTeamID gets the team ID from context or route parameter
func GetTeamID(c *gin.Context) (int64, error) {
	// Try to get from context first
	if teamID, exists := c.Get(TeamIDKey); exists {
		if id, ok := teamID.(int64); ok {
			return id, nil
		}
	}

	// Try route parameter
	if teamIDStr := c.Param("team_id"); teamIDStr != "" {
		var teamID int64
		_, err := fmt.Sscanf(teamIDStr, "%d", &teamID)
		if err == nil {
			return teamID, nil
		}
	}

	// Try query parameter
	if teamIDStr := c.Query("team_id"); teamIDStr != "" {
		var teamID int64
		_, err := fmt.Sscanf(teamIDStr, "%d", &teamID)
		if err == nil {
			return teamID, nil
		}
	}

	return 0, fmt.Errorf("team_id not found")
}

// GetRole gets the role from the context
func GetRole(c *gin.Context) (model.Role, bool) {
	role, exists := c.Get(RoleKey)
	if !exists {
		return "", false
	}

	r, ok := role.(model.Role)
	return r, ok
}
