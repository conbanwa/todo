package model

import (
	"testing"
	"time"
)

func TestUser_Fields(t *testing.T) {
	now := time.Now()
	user := User{
		ID:           1,
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if user.ID != 1 {
		t.Errorf("User.ID = %v, want 1", user.ID)
	}
	if user.Username != "testuser" {
		t.Errorf("User.Username = %v, want testuser", user.Username)
	}
	if user.Email != "test@example.com" {
		t.Errorf("User.Email = %v, want test@example.com", user.Email)
	}
	if user.PasswordHash != "hashed_password" {
		t.Errorf("User.PasswordHash = %v, want hashed_password", user.PasswordHash)
	}
}

func TestUser_ZeroValue(t *testing.T) {
	var user User
	if user.ID != 0 {
		t.Errorf("User.ID should be 0 for zero value, got %v", user.ID)
	}
	if user.Username != "" {
		t.Errorf("User.Username should be empty for zero value, got %v", user.Username)
	}
	if user.Email != "" {
		t.Errorf("User.Email should be empty for zero value, got %v", user.Email)
	}
	if !user.CreatedAt.IsZero() {
		t.Errorf("User.CreatedAt should be zero for zero value")
	}
	if !user.UpdatedAt.IsZero() {
		t.Errorf("User.UpdatedAt should be zero for zero value")
	}
}
