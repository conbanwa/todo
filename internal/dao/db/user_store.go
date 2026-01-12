package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/conbanwa/todo/internal/dao/cache"
	"github.com/conbanwa/todo/internal/model"
)

// CreateUser inserts a new user into the database
func (s *SQLiteStore) CreateUser(user *model.User) (int64, error) {
	if user.Username == "" {
		return 0, fmt.Errorf("username is required")
	}
	if user.Email == "" {
		return 0, fmt.Errorf("email is required")
	}
	if user.PasswordHash == "" {
		return 0, fmt.Errorf("password_hash is required")
	}

	query := `
	INSERT INTO users (username, email, password_hash, created_at, updated_at)
	VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`
	result, err := s.db.Exec(query, user.Username, user.Email, user.PasswordHash)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}

	user.ID = id
	return id, nil
}

// GetUserByID retrieves a user by ID
func (s *SQLiteStore) GetUserByID(id int64) (*model.User, error) {
	query := `
	SELECT id, username, email, password_hash, created_at, updated_at
	FROM users
	WHERE id = ?
	`
	row := s.db.QueryRow(query, id)

	var user model.User
	var createdAtStr, updatedAtStr string

	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &createdAtStr, &updatedAtStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, cache.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Parse timestamps
	if createdAtStr != "" {
		if t, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
			user.CreatedAt = t
		} else if t, err := time.Parse("2006-01-02 15:04:05", createdAtStr); err == nil {
			user.CreatedAt = t
		}
	}
	if updatedAtStr != "" {
		if t, err := time.Parse(time.RFC3339, updatedAtStr); err == nil {
			user.UpdatedAt = t
		} else if t, err := time.Parse("2006-01-02 15:04:05", updatedAtStr); err == nil {
			user.UpdatedAt = t
		}
	}

	return &user, nil
}

// GetUserByEmail retrieves a user by email
func (s *SQLiteStore) GetUserByEmail(email string) (*model.User, error) {
	query := `
	SELECT id, username, email, password_hash, created_at, updated_at
	FROM users
	WHERE email = ?
	`
	row := s.db.QueryRow(query, email)

	var user model.User
	var createdAtStr, updatedAtStr string

	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &createdAtStr, &updatedAtStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, cache.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	// Parse timestamps
	if createdAtStr != "" {
		if t, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
			user.CreatedAt = t
		} else if t, err := time.Parse("2006-01-02 15:04:05", createdAtStr); err == nil {
			user.CreatedAt = t
		}
	}
	if updatedAtStr != "" {
		if t, err := time.Parse(time.RFC3339, updatedAtStr); err == nil {
			user.UpdatedAt = t
		} else if t, err := time.Parse("2006-01-02 15:04:05", updatedAtStr); err == nil {
			user.UpdatedAt = t
		}
	}

	return &user, nil
}

// GetUserByUsername retrieves a user by username
func (s *SQLiteStore) GetUserByUsername(username string) (*model.User, error) {
	query := `
	SELECT id, username, email, password_hash, created_at, updated_at
	FROM users
	WHERE username = ?
	`
	row := s.db.QueryRow(query, username)

	var user model.User
	var createdAtStr, updatedAtStr string

	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &createdAtStr, &updatedAtStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, cache.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	// Parse timestamps
	if createdAtStr != "" {
		if t, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
			user.CreatedAt = t
		} else if t, err := time.Parse("2006-01-02 15:04:05", createdAtStr); err == nil {
			user.CreatedAt = t
		}
	}
	if updatedAtStr != "" {
		if t, err := time.Parse(time.RFC3339, updatedAtStr); err == nil {
			user.UpdatedAt = t
		} else if t, err := time.Parse("2006-01-02 15:04:05", updatedAtStr); err == nil {
			user.UpdatedAt = t
		}
	}

	return &user, nil
}

// UpdateUser updates an existing user
func (s *SQLiteStore) UpdateUser(user *model.User) error {
	if user.ID == 0 {
		return fmt.Errorf("id is required")
	}

	// Check if user exists
	_, err := s.GetUserByID(user.ID)
	if err != nil {
		return err
	}

	query := `
	UPDATE users
	SET username = ?, email = ?, password_hash = ?, updated_at = CURRENT_TIMESTAMP
	WHERE id = ?
	`
	_, err = s.db.Exec(query, user.Username, user.Email, user.PasswordHash, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}
