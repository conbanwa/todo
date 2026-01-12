package auth

import (
	"net/http"
	"strings"

	"github.com/conbanwa/todo/internal/dao/cache/api"
	"github.com/gin-gonic/gin"
)

const (
	UserIDKey = "user_id"
)

// AuthMiddleware validates JWT tokens and sets user context
func AuthMiddleware(authService *api.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		token := parts[1]

		// Validate token
		userID, err := authService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		// Set user ID in context
		c.Set(UserIDKey, userID)
		c.Next()
	}
}

// ExtractToken extracts JWT token from request (from Authorization header or query parameter)
func ExtractToken(c *gin.Context) string {
	// Try Authorization header first
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1]
		}
	}

	// Try query parameter
	if token := c.Query("token"); token != "" {
		return token
	}

	return ""
}

// GetUserID gets the user ID from the context
func GetUserID(c *gin.Context) (int64, bool) {
	userID, exists := c.Get(UserIDKey)
	if !exists {
		return 0, false
	}

	id, ok := userID.(int64)
	return id, ok
}
