package api

import (
	"testing"

	"github.com/conbanwa/todo/internal/model"
)

// MockUserStore for testing
type mockUserStore struct {
	users   map[int64]*model.User
	byEmail map[string]*model.User
	byUser  map[string]*model.User
	nextID  int64
}

func newMockUserStore() *mockUserStore {
	return &mockUserStore{
		users:   make(map[int64]*model.User),
		byEmail: make(map[string]*model.User),
		byUser:  make(map[string]*model.User),
		nextID:  1,
	}
}

func (m *mockUserStore) Create(user *model.User) (int64, error) {
	user.ID = m.nextID
	m.nextID++
	m.users[user.ID] = user
	m.byEmail[user.Email] = user
	m.byUser[user.Username] = user
	return user.ID, nil
}

func (m *mockUserStore) GetByID(id int64) (*model.User, error) {
	user, ok := m.users[id]
	if !ok {
		return nil, ErrNotFound
	}
	return user, nil
}

func (m *mockUserStore) GetByEmail(email string) (*model.User, error) {
	user, ok := m.byEmail[email]
	if !ok {
		return nil, ErrNotFound
	}
	return user, nil
}

func (m *mockUserStore) GetByUsername(username string) (*model.User, error) {
	user, ok := m.byUser[username]
	if !ok {
		return nil, ErrNotFound
	}
	return user, nil
}

func (m *mockUserStore) Update(user *model.User) error {
	if _, ok := m.users[user.ID]; !ok {
		return ErrNotFound
	}
	m.users[user.ID] = user
	m.byEmail[user.Email] = user
	m.byUser[user.Username] = user
	return nil
}

var ErrNotFound = &notFoundError{}

type notFoundError struct{}

func (e *notFoundError) Error() string { return "not found" }

func TestAuthService_Register(t *testing.T) {
	userStore := newMockUserStore()
	authService := NewAuthService(userStore)

	req := RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	user, err := authService.Register(&req)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	if user.Username != "testuser" {
		t.Errorf("expected username 'testuser', got %q", user.Username)
	}
	if user.Email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got %q", user.Email)
	}
	if user.PasswordHash == "" {
		t.Error("expected password hash to be set")
	}
	if user.PasswordHash == "password123" {
		t.Error("password should be hashed, not stored as plain text")
	}
}

func TestAuthService_Register_DuplicateEmail(t *testing.T) {
	userStore := newMockUserStore()
	authService := NewAuthService(userStore)

	req1 := RegisterRequest{
		Username: "user1",
		Email:    "test@example.com",
		Password: "password123",
	}

	_, err := authService.Register(&req1)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	req2 := RegisterRequest{
		Username: "user2",
		Email:    "test@example.com",
		Password: "password123",
	}

	_, err = authService.Register(&req2)
	if err == nil {
		t.Fatal("expected error for duplicate email")
	}
}

func TestAuthService_Login(t *testing.T) {
	userStore := newMockUserStore()
	authService := NewAuthService(userStore)

	// Register user
	req := RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	registeredUser, err := authService.Register(&req)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	// Login
	loginReq := LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	user, token, err := authService.Login(&loginReq)
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	if user.ID != registeredUser.ID {
		t.Errorf("expected user ID %d, got %d", registeredUser.ID, user.ID)
	}
	if token == "" {
		t.Error("expected token to be returned")
	}
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	userStore := newMockUserStore()
	authService := NewAuthService(userStore)

	// Register user
	req := RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	_, err := authService.Register(&req)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	// Login with wrong password
	loginReq := LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	_, _, err = authService.Login(&loginReq)
	if err == nil {
		t.Fatal("expected error for invalid password")
	}
}

func TestAuthService_HashPassword(t *testing.T) {
	authService := NewAuthService(nil)

	password := "testpassword"
	hash, err := authService.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if hash == "" {
		t.Error("expected hash to be non-empty")
	}
	if hash == password {
		t.Error("hash should be different from password")
	}

	// Verify password
	err = authService.VerifyPassword(hash, password)
	if err != nil {
		t.Errorf("VerifyPassword failed: %v", err)
	}
}
