package api

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/conbanwa/todo/internal/model"
)

// RegisterRequest represents a user registration request
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest represents a user login request
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthService handles authentication operations
type AuthService struct {
	userStore UserStore
	jwtSecret []byte
}

// NewAuthService creates a new AuthService
func NewAuthService(userStore UserStore) *AuthService {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default-secret-key-change-in-production"
	}
	return &AuthService{
		userStore: userStore,
		jwtSecret: []byte(secret),
	}
}

// Register registers a new user
func (a *AuthService) Register(req *RegisterRequest) (*model.User, error) {
	if req.Username == "" {
		return nil, ErrInvalid("username is required")
	}
	if req.Email == "" {
		return nil, ErrInvalid("email is required")
	}
	if req.Password == "" {
		return nil, ErrInvalid("password is required")
	}

	// Check if user already exists
	_, err := a.userStore.GetByEmail(req.Email)
	if err == nil {
		return nil, ErrInvalid("email already registered")
	}

	// Check if username already exists
	_, err = a.userStore.GetByUsername(req.Username)
	if err == nil {
		return nil, ErrInvalid("username already taken")
	}

	// Hash password
	hashedPassword, err := a.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
	}

	id, err := a.userStore.Create(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	user.ID = id
	return user, nil
}

// Login authenticates a user and returns a JWT token
func (a *AuthService) Login(req *LoginRequest) (*model.User, string, error) {
	if req.Email == "" {
		return nil, "", ErrInvalid("email is required")
	}
	if req.Password == "" {
		return nil, "", ErrInvalid("password is required")
	}

	// Get user by email
	user, err := a.userStore.GetByEmail(req.Email)
	if err != nil {
		return nil, "", ErrInvalid("invalid email or password")
	}

	// Verify password
	err = a.VerifyPassword(user.PasswordHash, req.Password)
	if err != nil {
		return nil, "", ErrInvalid("invalid email or password")
	}

	// Generate JWT token
	token, err := a.GenerateToken(user)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	return user, token, nil
}

// HashPassword hashes a password using bcrypt
func (a *AuthService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// VerifyPassword verifies a password against a hash
func (a *AuthService) VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// GenerateToken generates a JWT token for a user
func (a *AuthService) GenerateToken(user *model.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(a.jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the user ID
func (a *AuthService) ValidateToken(tokenString string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return a.jwtSecret, nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			return 0, errors.New("invalid token claims")
		}
		return int64(userIDFloat), nil
	}

	return 0, errors.New("invalid token")
}
