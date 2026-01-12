package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/conbanwa/todo/internal/dao/cache"
	"github.com/conbanwa/todo/internal/dao/cache/api"
	"github.com/conbanwa/todo/internal/model"
	_ "github.com/glebarez/sqlite"
)

// SQLiteStore implements the Store interface using SQLite database
type SQLiteStore struct {
	db *sql.DB
}

// NewSQLiteStore creates a new SQLite store instance
// If dbPath is empty, it defaults to "todos.db" in the current directory
func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	if dbPath == "" {
		dbPath = "todos.db"
	}

	// Ensure the directory exists
	dir := filepath.Dir(dbPath)
	if dir != "." && dir != "" {
		// Note: In production, you might want to explicitly create the directory
		// os.MkdirAll(dir, 0755)
	}

	// Use WAL mode for better concurrency
	dsn := fmt.Sprintf("file:%s?_journal_mode=WAL", dbPath)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	store := &SQLiteStore{db: db}

	if err := store.initSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return store, nil
}

// initSchema creates all tables if they don't exist
func (s *SQLiteStore) initSchema() error {
	// Create users table
	usersQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
	)
	`
	_, err := s.db.Exec(usersQuery)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	// Create teams table
	teamsQuery := `
	CREATE TABLE IF NOT EXISTS teams (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT,
		created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
	)
	`
	_, err = s.db.Exec(teamsQuery)
	if err != nil {
		return fmt.Errorf("failed to create teams table: %w", err)
	}

	// Create user_teams junction table
	userTeamsQuery := `
	CREATE TABLE IF NOT EXISTS user_teams (
		user_id INTEGER NOT NULL,
		team_id INTEGER NOT NULL,
		role TEXT NOT NULL DEFAULT 'member',
		created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (user_id, team_id),
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE
	)
	`
	_, err = s.db.Exec(userTeamsQuery)
	if err != nil {
		return fmt.Errorf("failed to create user_teams table: %w", err)
	}

	// Create todos table (or update if exists)
	todosQuery := `
	CREATE TABLE IF NOT EXISTS todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT,
		due_date TEXT,
		status TEXT NOT NULL DEFAULT 'not_started',
		priority INTEGER DEFAULT 0,
		tags TEXT,
		team_id INTEGER,
		created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE
	)
	`
	_, err = s.db.Exec(todosQuery)
	if err != nil {
		return fmt.Errorf("failed to create todos table: %w", err)
	}

	// Add team_id column to todos if it doesn't exist (migration for existing databases)
	migrationQuery := `
	SELECT COUNT(*) FROM pragma_table_info('todos') WHERE name='team_id'
	`
	var count int
	err = s.db.QueryRow(migrationQuery).Scan(&count)
	if err == nil && count == 0 {
		alterQuery := `ALTER TABLE todos ADD COLUMN team_id INTEGER REFERENCES teams(id) ON DELETE CASCADE`
		_, err = s.db.Exec(alterQuery)
		if err != nil {
			return fmt.Errorf("failed to add team_id column: %w", err)
		}
	}

	// Create indexes
	indexQuery := `
	CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
	CREATE INDEX IF NOT EXISTS idx_user_teams_user_id ON user_teams(user_id);
	CREATE INDEX IF NOT EXISTS idx_user_teams_team_id ON user_teams(team_id);
	CREATE INDEX IF NOT EXISTS idx_todos_status ON todos(status);
	CREATE INDEX IF NOT EXISTS idx_todos_due_date ON todos(due_date);
	CREATE INDEX IF NOT EXISTS idx_todos_team_id ON todos(team_id);
	`
	_, err = s.db.Exec(indexQuery)
	if err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	return nil
}

// Close closes the database connection
func (s *SQLiteStore) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// Create inserts a new api into the database
func (s *SQLiteStore) Create(t *model.Todo) (int64, error) {
	if t.Name == "" {
		return 0, api.ErrInvalid("name is required")
	}

	if t.Status == "" {
		t.Status = model.NotStarted
	}

	// Serialize tags as JSON
	tagsJSON, err := json.Marshal(t.Tags)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal tags: %w", err)
	}

	var dueDateStr sql.NullString
	if !t.DueDate.IsZero() {
		dueDateStr = sql.NullString{String: t.DueDate.Format(time.RFC3339), Valid: true}
	}

	var teamID sql.NullInt64
	if t.TeamID != 0 {
		teamID = sql.NullInt64{Int64: t.TeamID, Valid: true}
	}

	query := `
	INSERT INTO todos (name, description, due_date, status, priority, tags, team_id)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	result, err := s.db.Exec(query, t.Name, t.Description, dueDateStr, string(t.Status), t.Priority, string(tagsJSON), teamID)
	if err != nil {
		return 0, fmt.Errorf("failed to create api: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return id, nil
}

// Get retrieves a api by ID
func (s *SQLiteStore) Get(id int64) (*model.Todo, error) {
	query := `
	SELECT id, name, description, due_date, status, priority, tags, team_id
	FROM todos
	WHERE id = ?
	`
	row := s.db.QueryRow(query, id)

	var t model.Todo
	var dueDateStr sql.NullString
	var tagsJSON string
	var statusStr string
	var teamID sql.NullInt64

	err := row.Scan(&t.ID, &t.Name, &t.Description, &dueDateStr, &statusStr, &t.Priority, &tagsJSON, &teamID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, cache.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get api: %w", err)
	}

	// Parse due_date
	if dueDateStr.Valid && dueDateStr.String != "" {
		dueDate, err := time.Parse(time.RFC3339, dueDateStr.String)
		if err != nil {
			// Try alternative format
			dueDate, err = time.Parse("2006-01-02 15:04:05", dueDateStr.String)
			if err != nil {
				return nil, fmt.Errorf("failed to parse due_date: %w", err)
			}
		}
		t.DueDate = dueDate
	}

	t.Status = model.Status(statusStr)

	// Parse tags JSON
	if tagsJSON != "" {
		if err := json.Unmarshal([]byte(tagsJSON), &t.Tags); err != nil {
			return nil, fmt.Errorf("failed to unmarshal tags: %w", err)
		}
	}
	// Ensure tags is never nil (always at least an empty slice)
	if t.Tags == nil {
		t.Tags = []string{}
	}

	// Parse team_id
	if teamID.Valid {
		t.TeamID = teamID.Int64
	}

	return &t, nil
}

// Update updates an existing api
func (s *SQLiteStore) Update(t *model.Todo) error {
	if t.ID == 0 {
		return api.ErrInvalid("id is required")
	}

	// Check if api exists
	_, err := s.Get(t.ID)
	if err != nil {
		return err
	}

	// Serialize tags as JSON
	tagsJSON, err := json.Marshal(t.Tags)
	if err != nil {
		return fmt.Errorf("failed to marshal tags: %w", err)
	}

	var dueDateStr sql.NullString
	if !t.DueDate.IsZero() {
		dueDateStr = sql.NullString{String: t.DueDate.Format(time.RFC3339), Valid: true}
	}

	if t.Status == "" {
		t.Status = model.NotStarted
	}

	var teamID sql.NullInt64
	if t.TeamID != 0 {
		teamID = sql.NullInt64{Int64: t.TeamID, Valid: true}
	}

	query := `
	UPDATE todos
	SET name = ?, description = ?, due_date = ?, status = ?, priority = ?, tags = ?, team_id = ?, updated_at = CURRENT_TIMESTAMP
	WHERE id = ?
	`
	_, err = s.db.Exec(query, t.Name, t.Description, dueDateStr, string(t.Status), t.Priority, string(tagsJSON), teamID, t.ID)
	if err != nil {
		return fmt.Errorf("failed to update api: %w", err)
	}

	return nil
}

// Delete removes a api from the database
func (s *SQLiteStore) Delete(id int64) error {
	query := `DELETE FROM todos WHERE id = ?`
	result, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete api: %w", err)
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

// List retrieves all todos with optional filtering and sorting
func (s *SQLiteStore) List(opts cache.ListOptions) ([]model.Todo, error) {
	var conditions []string
	var args []interface{}

	// Apply status filter if specified
	if string(opts.Status) != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, string(opts.Status))
	}

	// Build WHERE clause
	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Always order by id for consistent results, then FilterAndSort will handle final sorting
	query := fmt.Sprintf(`
	SELECT id, name, description, due_date, status, priority, tags, team_id
	FROM todos
	%s
	ORDER BY id ASC
	`, whereClause)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list todos: %w", err)
	}
	defer rows.Close()

	var todos []model.Todo
	for rows.Next() {
		var t model.Todo
		var dueDateStr sql.NullString
		var tagsJSON string
		var statusStr string
		var teamID sql.NullInt64

		err := rows.Scan(&t.ID, &t.Name, &t.Description, &dueDateStr, &statusStr, &t.Priority, &tagsJSON, &teamID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan api: %w", err)
		}

		// Parse due_date
		if dueDateStr.Valid && dueDateStr.String != "" {
			dueDate, err := time.Parse(time.RFC3339, dueDateStr.String)
			if err != nil {
				// Try alternative format
				dueDate, err = time.Parse("2006-01-02 15:04:05", dueDateStr.String)
				if err != nil {
					return nil, fmt.Errorf("failed to parse due_date: %w", err)
				}
			}
			t.DueDate = dueDate
		}

		t.Status = model.Status(statusStr)

		// Parse tags JSON
		if tagsJSON != "" {
			if err := json.Unmarshal([]byte(tagsJSON), &t.Tags); err != nil {
				return nil, fmt.Errorf("failed to unmarshal tags: %w", err)
			}
		}
		// Ensure tags is never nil (always at least an empty slice)
		if t.Tags == nil {
			t.Tags = []string{}
		}

		// Parse team_id
		if teamID.Valid {
			t.TeamID = teamID.Int64
		}

		todos = append(todos, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	// Apply FilterAndSort for consistent behavior with in-memory store
	// This handles any remaining filtering and ensures sorting is consistent
	return cache.FilterAndSort(todos, opts), nil
}
