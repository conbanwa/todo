package todo

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

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

// initSchema creates the todos table if it doesn't exist
func (s *SQLiteStore) initSchema() error {
	query := `
	CREATE TABLE IF NOT EXISTS todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT,
		due_date TEXT,
		status TEXT NOT NULL DEFAULT 'not_started',
		priority INTEGER DEFAULT 0,
		tags TEXT,
		created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
	)
	`
	_, err := s.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create todos table: %w", err)
	}

	// Create index for common queries
	indexQuery := `
	CREATE INDEX IF NOT EXISTS idx_todos_status ON todos(status);
	CREATE INDEX IF NOT EXISTS idx_todos_due_date ON todos(due_date);
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

// Create inserts a new todo into the database
func (s *SQLiteStore) Create(t *Todo) (int64, error) {
	if t.Name == "" {
		return 0, ErrInvalid("name is required")
	}

	if t.Status == "" {
		t.Status = NotStarted
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

	query := `
	INSERT INTO todos (name, description, due_date, status, priority, tags)
	VALUES (?, ?, ?, ?, ?, ?)
	`
	result, err := s.db.Exec(query, t.Name, t.Description, dueDateStr, string(t.Status), t.Priority, string(tagsJSON))
	if err != nil {
		return 0, fmt.Errorf("failed to create todo: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return id, nil
}

// Get retrieves a todo by ID
func (s *SQLiteStore) Get(id int64) (*Todo, error) {
	query := `
	SELECT id, name, description, due_date, status, priority, tags
	FROM todos
	WHERE id = ?
	`
	row := s.db.QueryRow(query, id)

	var t Todo
	var dueDateStr sql.NullString
	var tagsJSON string
	var statusStr string

	err := row.Scan(&t.ID, &t.Name, &t.Description, &dueDateStr, &statusStr, &t.Priority, &tagsJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get todo: %w", err)
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

	t.Status = Status(statusStr)

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

	return &t, nil
}

// Update updates an existing todo
func (s *SQLiteStore) Update(t *Todo) error {
	if t.ID == 0 {
		return ErrInvalid("id is required")
	}

	// Check if todo exists
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
		t.Status = NotStarted
	}

	query := `
	UPDATE todos
	SET name = ?, description = ?, due_date = ?, status = ?, priority = ?, tags = ?, updated_at = CURRENT_TIMESTAMP
	WHERE id = ?
	`
	_, err = s.db.Exec(query, t.Name, t.Description, dueDateStr, string(t.Status), t.Priority, string(tagsJSON), t.ID)
	if err != nil {
		return fmt.Errorf("failed to update todo: %w", err)
	}

	return nil
}

// Delete removes a todo from the database
func (s *SQLiteStore) Delete(id int64) error {
	query := `DELETE FROM todos WHERE id = ?`
	result, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// List retrieves all todos with optional filtering and sorting
func (s *SQLiteStore) List(opts ListOptions) ([]Todo, error) {
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
	SELECT id, name, description, due_date, status, priority, tags
	FROM todos
	%s
	ORDER BY id ASC
	`, whereClause)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list todos: %w", err)
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var t Todo
		var dueDateStr sql.NullString
		var tagsJSON string
		var statusStr string

		err := rows.Scan(&t.ID, &t.Name, &t.Description, &dueDateStr, &statusStr, &t.Priority, &tagsJSON)
		if err != nil {
			return nil, fmt.Errorf("failed to scan todo: %w", err)
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

		t.Status = Status(statusStr)

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

		todos = append(todos, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	// Apply FilterAndSort for consistent behavior with in-memory store
	// This handles any remaining filtering and ensures sorting is consistent
	return FilterAndSort(todos, opts), nil
}
