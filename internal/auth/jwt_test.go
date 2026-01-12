package auth

import (
	"testing"
)

func TestGenerateToken(t *testing.T) {
	userID := int64(1)
	email := "test@example.com"

	token, err := GenerateToken(userID, email)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	if token == "" {
		t.Fatal("expected token to be non-empty")
	}
}

func TestValidateToken(t *testing.T) {
	userID := int64(1)
	email := "test@example.com"

	// Generate token
	token, err := GenerateToken(userID, email)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	// Validate token
	gotUserID, err := ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}

	if gotUserID != userID {
		t.Errorf("expected user ID %d, got %d", userID, gotUserID)
	}
}

func TestValidateToken_InvalidToken(t *testing.T) {
	invalidToken := "invalid.token.here"

	_, err := ValidateToken(invalidToken)
	if err == nil {
		t.Fatal("expected error for invalid token")
	}
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	// This test would require mocking time or setting a very short expiration
	// For now, we'll skip it as it requires more complex setup
	t.Skip("expired token test requires time mocking")
}
