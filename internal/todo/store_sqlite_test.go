package todo

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// setupTestDB creates a temporary database for testing
func setupTestDB(t *testing.T) (*SQLiteStore, func()) {
	t.Helper()

	// Create a temporary database file
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test_todos.db")

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

// TestSQLiteStore_CreateGetUpdateDelete_List tests basic CRUD operations
func TestSQLiteStore_CreateGetUpdateDelete_List(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	// Create
	todo := &Todo{
		Name:        "task1",
		Description: "First task",
		DueDate:     time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC),
		Status:      NotStarted,
		Priority:    5,
		Tags:        []string{"work", "urgent"},
	}

	id, err := store.Create(todo)
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if id == 0 {
		t.Fatal("expected non-zero ID")
	}

	// Get
	got, err := store.Get(id)
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if got.ID != id {
		t.Errorf("expected ID %d, got %d", id, got.ID)
	}
	if got.Name != "task1" {
		t.Errorf("expected name 'task1', got %q", got.Name)
	}
	if got.Description != "First task" {
		t.Errorf("expected description 'First task', got %q", got.Description)
	}
	if !got.DueDate.Equal(time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)) {
		t.Errorf("expected due date 2026-01-15 10:00:00, got %v", got.DueDate)
	}
	if got.Status != NotStarted {
		t.Errorf("expected status NotStarted, got %q", got.Status)
	}
	if got.Priority != 5 {
		t.Errorf("expected priority 5, got %d", got.Priority)
	}
	if len(got.Tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(got.Tags))
	}
	if got.Tags[0] != "work" || got.Tags[1] != "urgent" {
		t.Errorf("expected tags [work urgent], got %v", got.Tags)
	}

	// Update
	got.Description = "updated"
	got.Status = InProgress
	if err := store.Update(got); err != nil {
		t.Fatalf("update failed: %v", err)
	}

	// Verify update
	updated, err := store.Get(id)
	if err != nil {
		t.Fatalf("get after update failed: %v", err)
	}
	if updated.Description != "updated" {
		t.Errorf("expected updated description 'updated', got %q", updated.Description)
	}
	if updated.Status != InProgress {
		t.Errorf("expected updated status InProgress, got %q", updated.Status)
	}

	// Delete
	if err := store.Delete(id); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	// Verify deletion
	_, err = store.Get(id)
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}

	// List (should be empty)
	list, err := store.List(ListOptions{})
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if len(list) != 0 {
		t.Errorf("expected empty list after deletion, got %d items", len(list))
	}
}

// TestSQLiteStore_ListOptions_SortAndFilter tests filtering and sorting
func TestSQLiteStore_ListOptions_SortAndFilter(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	// deterministic times
	t1 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)
	t3 := time.Date(2026, 1, 3, 0, 0, 0, 0, time.UTC)

	_, _ = store.Create(&Todo{Name: "alpha", DueDate: t2, Status: InProgress})
	id2, _ := store.Create(&Todo{Name: "bravo", DueDate: t1, Status: NotStarted})
	_, _ = store.Create(&Todo{Name: "charlie", DueDate: t3, Status: Completed})

	// filter by status
	list, err := store.List(ListOptions{Status: NotStarted})
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 todo with NotStarted status, got %d", len(list))
	}
	if list[0].ID != id2 || list[0].Name != "bravo" {
		t.Fatalf("expected only bravo (id=%d), got: id=%d name=%s", id2, list[0].ID, list[0].Name)
	}

	// sort by due_date asc
	list, err = store.List(ListOptions{SortBy: "due_date", SortOrder: "asc"})
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if len(list) != 3 {
		t.Fatalf("expected 3 todos, got %d", len(list))
	}
	if list[0].Name != "bravo" || list[1].Name != "alpha" || list[2].Name != "charlie" {
		t.Fatalf("unexpected due_date asc order: %v", []string{list[0].Name, list[1].Name, list[2].Name})
	}

	// sort by due_date desc
	list, err = store.List(ListOptions{SortBy: "due_date", SortOrder: "desc"})
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if list[0].Name != "charlie" || list[2].Name != "bravo" {
		t.Fatalf("unexpected due_date desc order: %v", []string{list[0].Name, list[1].Name, list[2].Name})
	}

	// sort by name asc
	list, err = store.List(ListOptions{SortBy: "name", SortOrder: "asc"})
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if list[0].Name != "alpha" || list[1].Name != "bravo" || list[2].Name != "charlie" {
		t.Fatalf("unexpected name asc order: %v", []string{list[0].Name, list[1].Name, list[2].Name})
	}

	// sort by status asc
	list, err = store.List(ListOptions{SortBy: "status", SortOrder: "asc"})
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if len(list) != 3 {
		t.Fatalf("expected 3 items for status sort, got %d", len(list))
	}
	// Status should be sorted: Completed, InProgress, NotStarted (alphabetically)
	if list[0].Status != Completed || list[1].Status != InProgress || list[2].Status != NotStarted {
		t.Errorf("unexpected status sort order: %v", []Status{list[0].Status, list[1].Status, list[2].Status})
	}
}

// TestSQLiteStore_Create_EdgeCases tests edge cases for Create
func TestSQLiteStore_Create_EdgeCases(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	t.Run("creates todo with default status when empty", func(t *testing.T) {
		todo := &Todo{Name: "Test Todo"}
		id, err := store.Create(todo)
		if err != nil {
			t.Fatalf("create failed: %v", err)
		}

		got, err := store.Get(id)
		if err != nil {
			t.Fatalf("get failed: %v", err)
		}

		if got.Status != NotStarted {
			t.Errorf("expected default status NotStarted, got %q", got.Status)
		}
	})

	t.Run("creates todo without due date", func(t *testing.T) {
		todo := &Todo{Name: "No Due Date"}
		id, err := store.Create(todo)
		if err != nil {
			t.Fatalf("create failed: %v", err)
		}

		got, err := store.Get(id)
		if err != nil {
			t.Fatalf("get failed: %v", err)
		}

		if !got.DueDate.IsZero() {
			t.Errorf("expected zero due date, got %v", got.DueDate)
		}
	})

	t.Run("creates todo with empty tags", func(t *testing.T) {
		todo := &Todo{Name: "Empty Tags", Tags: []string{}}
		id, err := store.Create(todo)
		if err != nil {
			t.Fatalf("create failed: %v", err)
		}

		got, err := store.Get(id)
		if err != nil {
			t.Fatalf("get failed: %v", err)
		}

		if got.Tags == nil || len(got.Tags) != 0 {
			t.Errorf("expected empty tags slice, got %v", got.Tags)
		}
	})

	t.Run("creates todo with nil tags", func(t *testing.T) {
		todo := &Todo{Name: "Nil Tags", Tags: nil}
		id, err := store.Create(todo)
		if err != nil {
			t.Fatalf("create failed: %v", err)
		}

		got, err := store.Get(id)
		if err != nil {
			t.Fatalf("get failed: %v", err)
		}

		// Tags should be unmarshalled as empty slice, not nil
		if got.Tags == nil {
			t.Error("expected tags to be empty slice, got nil")
		}
	})

	t.Run("creates multiple todos with sequential IDs", func(t *testing.T) {
		id1, err := store.Create(&Todo{Name: "Todo 1"})
		if err != nil {
			t.Fatalf("create 1 failed: %v", err)
		}

		id2, err := store.Create(&Todo{Name: "Todo 2"})
		if err != nil {
			t.Fatalf("create 2 failed: %v", err)
		}

		id3, err := store.Create(&Todo{Name: "Todo 3"})
		if err != nil {
			t.Fatalf("create 3 failed: %v", err)
		}

		if id2 != id1+1 {
			t.Errorf("expected sequential IDs, got id1=%d, id2=%d", id1, id2)
		}
		if id3 != id2+1 {
			t.Errorf("expected sequential IDs, got id2=%d, id3=%d", id2, id3)
		}
	})
}

// TestSQLiteStore_Get_EdgeCases tests edge cases for Get
func TestSQLiteStore_Get_EdgeCases(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	t.Run("returns ErrNotFound for non-existent todo", func(t *testing.T) {
		_, err := store.Get(999)
		if err != ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("gets todo with complex tags", func(t *testing.T) {
		tags := []string{"tag1", "tag with spaces", "tag-with-dashes", "tag_with_underscores"}
		todo := &Todo{
			Name: "Complex Tags",
			Tags: tags,
		}

		id, err := store.Create(todo)
		if err != nil {
			t.Fatalf("create failed: %v", err)
		}

		got, err := store.Get(id)
		if err != nil {
			t.Fatalf("get failed: %v", err)
		}

		if len(got.Tags) != len(tags) {
			t.Fatalf("expected %d tags, got %d", len(tags), len(got.Tags))
		}
		for i, expectedTag := range tags {
			if got.Tags[i] != expectedTag {
				t.Errorf("expected tag %q at index %d, got %q", expectedTag, i, got.Tags[i])
			}
		}
	})
}

// TestSQLiteStore_Update_EdgeCases tests edge cases for Update
func TestSQLiteStore_Update_EdgeCases(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	t.Run("updates all fields", func(t *testing.T) {
		original := &Todo{
			Name:        "Original",
			Description: "Original Description",
			Status:      NotStarted,
			Priority:    1,
			Tags:        []string{"old"},
		}

		id, err := store.Create(original)
		if err != nil {
			t.Fatalf("create failed: %v", err)
		}

		updated := &Todo{
			ID:          id,
			Name:        "Updated",
			Description: "Updated Description",
			DueDate:     time.Date(2026, 12, 31, 23, 59, 0, 0, time.UTC),
			Status:      Completed,
			Priority:    10,
			Tags:        []string{"new", "updated"},
		}

		err = store.Update(updated)
		if err != nil {
			t.Fatalf("update failed: %v", err)
		}

		got, err := store.Get(id)
		if err != nil {
			t.Fatalf("get failed: %v", err)
		}

		if got.Name != updated.Name {
			t.Errorf("expected name %q, got %q", updated.Name, got.Name)
		}
		if got.Description != updated.Description {
			t.Errorf("expected description %q, got %q", updated.Description, got.Description)
		}
		if !got.DueDate.Equal(updated.DueDate) {
			t.Errorf("expected due date %v, got %v", updated.DueDate, got.DueDate)
		}
		if got.Status != updated.Status {
			t.Errorf("expected status %q, got %q", updated.Status, got.Status)
		}
		if got.Priority != updated.Priority {
			t.Errorf("expected priority %d, got %d", updated.Priority, got.Priority)
		}
		if len(got.Tags) != len(updated.Tags) {
			t.Fatalf("expected %d tags, got %d", len(updated.Tags), len(got.Tags))
		}
	})

	t.Run("returns ErrNotFound for non-existent todo", func(t *testing.T) {
		todo := &Todo{ID: 999, Name: "Non-existent"}
		err := store.Update(todo)
		if err != ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("updates todo to have empty due date", func(t *testing.T) {
		dueDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
		todo := &Todo{Name: "With Date", DueDate: dueDate}

		id, err := store.Create(todo)
		if err != nil {
			t.Fatalf("create failed: %v", err)
		}

		// Update to remove due date
		updated := &Todo{ID: id, Name: "Without Date", DueDate: time.Time{}}
		err = store.Update(updated)
		if err != nil {
			t.Fatalf("update failed: %v", err)
		}

		got, err := store.Get(id)
		if err != nil {
			t.Fatalf("get failed: %v", err)
		}

		// Due date should be zero/empty after update
		if !got.DueDate.IsZero() {
			t.Errorf("expected zero due date after update, got %v", got.DueDate)
		}
	})
}

// TestSQLiteStore_Delete_EdgeCases tests edge cases for Delete
func TestSQLiteStore_Delete_EdgeCases(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	t.Run("returns ErrNotFound for non-existent todo", func(t *testing.T) {
		err := store.Delete(999)
		if err != ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("can delete multiple todos", func(t *testing.T) {
		id1, _ := store.Create(&Todo{Name: "Todo 1"})
		id2, _ := store.Create(&Todo{Name: "Todo 2"})
		id3, _ := store.Create(&Todo{Name: "Todo 3"})

		// Delete middle one
		if err := store.Delete(id2); err != nil {
			t.Fatalf("delete failed: %v", err)
		}

		// Verify others still exist
		if _, err := store.Get(id1); err != nil {
			t.Errorf("todo 1 should still exist: %v", err)
		}
		if _, err := store.Get(id3); err != nil {
			t.Errorf("todo 3 should still exist: %v", err)
		}

		// Verify deleted one is gone
		if _, err := store.Get(id2); err != ErrNotFound {
			t.Errorf("todo 2 should be deleted, got: %v", err)
		}
	})
}

// TestSQLiteStore_Persistence tests that data persists across store instances
func TestSQLiteStore_Persistence(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "persist_test.db")

	// Create first store and add data
	store1, err := NewSQLiteStore(dbPath)
	if err != nil {
		t.Fatalf("failed to create first store: %v", err)
	}

	id, err := store1.Create(&Todo{
		Name:        "Persistent Todo",
		Description: "This should persist",
		Status:      Completed,
	})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	store1.Close()

	// Create second store with same database path
	store2, err := NewSQLiteStore(dbPath)
	if err != nil {
		t.Fatalf("failed to create second store: %v", err)
	}
	defer store2.Close()

	// Verify data persisted
	got, err := store2.Get(id)
	if err != nil {
		t.Fatalf("get failed after reopening: %v", err)
	}

	if got.Name != "Persistent Todo" {
		t.Errorf("expected name %q, got %q", "Persistent Todo", got.Name)
	}
	if got.Description != "This should persist" {
		t.Errorf("expected description %q, got %q", "This should persist", got.Description)
	}
	if got.Status != Completed {
		t.Errorf("expected status %q, got %q", Completed, got.Status)
	}
}

// TestSQLiteStore_Close tests the Close function
func TestSQLiteStore_Close(t *testing.T) {
	t.Run("closes database connection", func(t *testing.T) {
		store, _ := setupTestDB(t)

		err := store.Close()
		if err != nil {
			t.Fatalf("close failed: %v", err)
		}

		// Attempting to use closed store should fail
		_, err = store.Create(&Todo{Name: "Should Fail"})
		if err == nil {
			t.Fatal("expected error when using closed store")
		}
	})

	t.Run("can close multiple times safely", func(t *testing.T) {
		store, _ := setupTestDB(t)

		err1 := store.Close()
		err2 := store.Close()

		if err1 != nil {
			t.Fatalf("first close failed: %v", err1)
		}
		if err2 != nil {
			t.Fatalf("second close failed: %v", err2)
		}
	})
}

// TestNewSQLiteStore tests database initialization
func TestNewSQLiteStore(t *testing.T) {
	t.Run("creates database with default path", func(t *testing.T) {
		store, err := NewSQLiteStore("")
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		defer store.Close()
		defer os.Remove("todos.db")

		if store.db == nil {
			t.Fatal("expected db connection to be initialized")
		}
	})

	t.Run("creates database with custom path", func(t *testing.T) {
		tmpDir := t.TempDir()
		dbPath := filepath.Join(tmpDir, "custom.db")

		store, err := NewSQLiteStore(dbPath)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		defer store.Close()

		if _, err := os.Stat(dbPath); os.IsNotExist(err) {
			t.Fatalf("database file should exist at %s", dbPath)
		}
	})

	t.Run("creates schema on initialization", func(t *testing.T) {
		store, cleanup := setupTestDB(t)
		defer cleanup()

		// Try to query the table structure
		rows, err := store.db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name='todos'")
		if err != nil {
			t.Fatalf("failed to query schema: %v", err)
		}
		defer rows.Close()

		if !rows.Next() {
			t.Fatal("expected todos table to exist")
		}
	})
}
