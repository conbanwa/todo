package transport

import (
	"net/http"

	"github.com/conbanwa/todo/internal/auth"
	"github.com/conbanwa/todo/internal/dao/cache/api"
	"github.com/conbanwa/todo/internal/model"
	"github.com/gin-gonic/gin"
)

// RegisterAuthRoutes registers authentication routes
func RegisterAuthRoutes(r *gin.Engine, authService *api.AuthService) {
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", handleRegister(authService))
		authGroup.POST("/login", handleLogin(authService))
		authGroup.GET("/me", auth.AuthMiddleware(authService), handleGetMe(authService.UserStore))
	}
}

// RegisterRequest represents a user registration request
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest represents a user login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	Token string     `json:"token"`
	User  model.User `json:"user"`
}

// @Summary Register new user
// @Description Register a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration request"
// @Success 201 {object} model.User
// @Failure 400 {object} map[string]string
// @Router /auth/register [post]
func handleRegister(authService *api.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		registerReq := &api.RegisterRequest{
			Username: req.Username,
			Email:    req.Email,
			Password: req.Password,
		}

		user, err := authService.Register(registerReq)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, user)
	}
}

// @Summary Login user
// @Description Login with email and password, returns JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login request"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/login [post]
func handleLogin(authService *api.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		loginReq := &api.LoginRequest{
			Email:    req.Email,
			Password: req.Password,
		}

		user, token, err := authService.Login(loginReq)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, LoginResponse{
			Token: token,
			User:  *user,
		})
	}
}

// @Summary Get current user
// @Description Get the current authenticated user's information
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.User
// @Failure 401 {object} map[string]string
// @Router /auth/me [get]
func handleGetMe(userStore api.UserStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := auth.GetUserID(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			return
		}

		user, err := userStore.GetByID(userID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
