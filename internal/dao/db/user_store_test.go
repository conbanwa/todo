package db

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/conbanwa/todo/internal/dao/cache"
	"github.com/conbanwa/todo/internal/model"
)

// setupTestUserDB creates a temporary database for user testing
func setupTestUserDB(t *testing.T) (*SQLiteStore, func()) {
	t.Helper()

	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test_users.db")

	store, err := NewSQLiteStore(dbPath)
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}

	cleanup := func() {
		store.Close()
		os.Remove(dbPath)
	}

	return store, cleanup
}

func TestUserStore_Create(t *testing.T) {
	store, cleanup := setupTestUserDB(t)
	defer cleanup()

	user := &model.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
	}

	id, err := store.CreateUser(user)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}
	if id == 0 {
		t.Fatal("expected non-zero ID")
	}
	if user.ID != id {
		t.Errorf("expected user.ID to be set to %d, got %d", id, user.ID)
	}
}

func TestUserStore_GetByID(t *testing.T) {
	store, cleanup := setupTestUserDB(t)
	defer cleanup()

	user := &model.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
	}

	id, err := store.CreateUser(user)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	got, err := store.GetUserByID(id)
	if err != nil {
		t.Fatalf("GetUserByID failed: %v", err)
	}

	if got.ID != id {
		t.Errorf("expected ID %d, got %d", id, got.ID)
	}
	if got.Username != "testuser" {
		t.Errorf("expected username 'testuser', got %q", got.Username)
	}
	if got.Email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got %q", got.Email)
	}
	if got.PasswordHash != "hashed_password" {
		t.Errorf("expected password hash 'hashed_password', got %q", got.PasswordHash)
	}
}

func TestUserStore_GetByEmail(t *testing.T) {
	store, cleanup := setupTestUserDB(t)
	defer cleanup()

	user := &model.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
	}

	_, err := store.CreateUser(user)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	got, err := store.GetUserByEmail("test@example.com")
	if err != nil {
		t.Fatalf("GetUserByEmail failed: %v", err)
	}

	if got.Email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got %q", got.Email)
	}
}

func TestUserStore_GetByUsername(t *testing.T) {
	store, cleanup := setupTestUserDB(t)
	defer cleanup()

	user := &model.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
	}

	_, err := store.CreateUser(user)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	got, err := store.GetUserByUsername("testuser")
	if err != nil {
		t.Fatalf("GetUserByUsername failed: %v", err)
	}

	if got.Username != "testuser" {
		t.Errorf("expected username 'testuser', got %q", got.Username)
	}
}

func TestUserStore_GetByID_NotFound(t *testing.T) {
	store, cleanup := setupTestUserDB(t)
	defer cleanup()

	_, err := store.GetUserByID(999)
	if err != cache.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestUserStore_Create_DuplicateEmail(t *testing.T) {
	store, cleanup := setupTestUserDB(t)
	defer cleanup()

	user1 := &model.User{
		Username:     "user1",
		Email:        "test@example.com",
		PasswordHash: "hash1",
	}

	_, err := store.CreateUser(user1)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	user2 := &model.User{
		Username:     "user2",
		Email:        "test@example.com",
		PasswordHash: "hash2",
	}

	_, err = store.CreateUser(user2)
	if err == nil {
		t.Fatal("expected error for duplicate email")
	}
}

func TestUserStore_Create_DuplicateUsername(t *testing.T) {
	store, cleanup := setupTestUserDB(t)
	defer cleanup()

	user1 := &model.User{
		Username:     "testuser",
		Email:        "test1@example.com",
		PasswordHash: "hash1",
	}

	_, err := store.CreateUser(user1)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	user2 := &model.User{
		Username:     "testuser",
		Email:        "test2@example.com",
		PasswordHash: "hash2",
	}

	_, err = store.CreateUser(user2)
	if err == nil {
		t.Fatal("expected error for duplicate username")
	}
}

func TestUserStore_Update(t *testing.T) {
	store, cleanup := setupTestUserDB(t)
	defer cleanup()

	user := &model.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
	}

	id, err := store.CreateUser(user)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	user.ID = id
	user.Username = "updateduser"
	user.Email = "updated@example.com"

	err = store.UpdateUser(user)
	if err != nil {
		t.Fatalf("UpdateUser failed: %v", err)
	}

	got, err := store.GetUserByID(id)
	if err != nil {
		t.Fatalf("GetUserByID failed: %v", err)
	}

	if got.Username != "updateduser" {
		t.Errorf("expected username 'updateduser', got %q", got.Username)
	}
	if got.Email != "updated@example.com" {
		t.Errorf("expected email 'updated@example.com', got %q", got.Email)
	}
}
