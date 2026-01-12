package transport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/conbanwa/todo/internal/dao/cache/api"
	"github.com/conbanwa/todo/internal/dao/db"
	"github.com/conbanwa/todo/internal/model"
	"github.com/gin-gonic/gin"
)

// setupIntegrationTestDB creates a temporary database for integration tests
// Note: These tests require CGO_ENABLED=1 to work with SQLite
func setupIntegrationTestDB(t *testing.T) (*db.SQLiteStore, func()) {
	t.Helper()
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "integration_test.db")

	store, err := db.NewSQLiteStore(dbPath)
	if err != nil {
		// Skip test if CGO is not enabled or SQLite driver is not available
		if err.Error() != "" && (contains(err.Error(), "cgo") || contains(err.Error(), "CGO_ENABLED")) {
			t.Skip("Skipping SQLite integration test: CGO not enabled or SQLite driver not available")
		}
		t.Fatalf("failed to create test database: %v", err)
	}

	cleanup := func() {
		store.Close()
		os.Remove(dbPath)
	}

	return store, cleanup
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || indexOf(s, substr) >= 0))
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// TestIntegration_CompleteWorkflow tests a complete CRUD workflow through the API
func TestIntegration_CompleteWorkflow(t *testing.T) {
	store, cleanup := setupIntegrationTestDB(t)
	defer cleanup()

	svc := api.NewService(store)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	RegisterRoutes(r, svc)

	// Step 1: Create a api
	todo := model.Todo{
		Name:        "Integration Test Todo",
		Description: "Testing complete workflow",
		DueDate:     time.Date(2026, 12, 31, 12, 0, 0, 0, time.UTC),
		Status:      model.NotStarted,
		Priority:    5,
		Tags:        []string{"test", "integration"},
	}

	createBody, _ := json.Marshal(todo)
	createReq := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewReader(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	r.ServeHTTP(createW, createReq)

	if createW.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d: %s", createW.Code, createW.Body.String())
	}

	var createdTodo model.Todo
	if err := json.NewDecoder(createW.Body).Decode(&createdTodo); err != nil {
		t.Fatalf("failed to decode created api: %v", err)
	}

	if createdTodo.ID == 0 {
		t.Fatal("expected non-zero ID")
	}
	if createdTodo.Name != todo.Name {
		t.Errorf("expected name %q, got %q", todo.Name, createdTodo.Name)
	}

	// Step 2: Get the created api
	getReq := httptest.NewRequest(http.MethodGet, "/todos/"+fmt.Sprintf("%d", createdTodo.ID), nil)
	getW := httptest.NewRecorder()
	r.ServeHTTP(getW, getReq)

	if getW.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", getW.Code)
	}

	var retrievedTodo model.Todo
	if err := json.NewDecoder(getW.Body).Decode(&retrievedTodo); err != nil {
		t.Fatalf("failed to decode retrieved api: %v", err)
	}

	if retrievedTodo.ID != createdTodo.ID {
		t.Errorf("expected ID %d, got %d", createdTodo.ID, retrievedTodo.ID)
	}
	if retrievedTodo.Name != todo.Name {
		t.Errorf("expected name %q, got %q", todo.Name, retrievedTodo.Name)
	}

	// Step 3: Update the api
	updatedTodo := model.Todo{
		Name:        "Updated Integration Test Todo",
		Description: "Updated description",
		Status:      model.InProgress,
		Priority:    10,
		Tags:        []string{"test", "integration", "updated"},
	}

	updateBody, _ := json.Marshal(updatedTodo)
	updateReq := httptest.NewRequest(http.MethodPut, "/todos/"+fmt.Sprintf("%d", createdTodo.ID), bytes.NewReader(updateBody))
	updateReq.Header.Set("Content-Type", "application/json")
	updateW := httptest.NewRecorder()
	r.ServeHTTP(updateW, updateReq)

	if updateW.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", updateW.Code, updateW.Body.String())
	}

	var updatedResponse model.Todo
	if err := json.NewDecoder(updateW.Body).Decode(&updatedResponse); err != nil {
		t.Fatalf("failed to decode updated api: %v", err)
	}

	if updatedResponse.Name != updatedTodo.Name {
		t.Errorf("expected name %q, got %q", updatedTodo.Name, updatedResponse.Name)
	}
	if updatedResponse.Status != model.InProgress {
		t.Errorf("expected status %q, got %q", model.InProgress, updatedResponse.Status)
	}

	// Step 4: List todos (should include our api)
	listReq := httptest.NewRequest(http.MethodGet, "/todos", nil)
	listW := httptest.NewRecorder()
	r.ServeHTTP(listW, listReq)

	if listW.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", listW.Code)
	}

	var todos []model.Todo
	if err := json.NewDecoder(listW.Body).Decode(&todos); err != nil {
		t.Fatalf("failed to decode todos list: %v", err)
	}

	if len(todos) == 0 {
		t.Fatal("expected at least one api in list")
	}

	found := false
	for _, todo := range todos {
		if todo.ID == createdTodo.ID {
			found = true
			if todo.Name != updatedTodo.Name {
				t.Errorf("expected updated name %q in list, got %q", updatedTodo.Name, todo.Name)
			}
			break
		}
	}
	if !found {
		t.Fatal("created api not found in list")
	}

	// Step 5: Delete the api
	deleteReq := httptest.NewRequest(http.MethodDelete, "/todos/"+fmt.Sprintf("%d", createdTodo.ID), nil)
	deleteW := httptest.NewRecorder()
	r.ServeHTTP(deleteW, deleteReq)

	if deleteW.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", deleteW.Code)
	}

	// Step 6: Verify deletion by trying to get it again
	getDeletedReq := httptest.NewRequest(http.MethodGet, "/todos/"+fmt.Sprintf("%d", createdTodo.ID), nil)
	getDeletedW := httptest.NewRecorder()
	r.ServeHTTP(getDeletedW, getDeletedReq)

	if getDeletedW.Code != http.StatusNotFound {
		t.Errorf("expected status 404 after deletion, got %d", getDeletedW.Code)
	}
}

// TestIntegration_MultipleTodos tests operations with multiple todos
func TestIntegration_MultipleTodos(t *testing.T) {
	store, cleanup := setupIntegrationTestDB(t)
	defer cleanup()

	svc := api.NewService(store)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	RegisterRoutes(r, svc)

	// Create multiple todos with different statuses
	todos := []model.Todo{
		{Name: "Todo 1", Status: model.NotStarted},
		{Name: "Todo 2", Status: model.InProgress},
		{Name: "Todo 3", Status: model.Completed},
		{Name: "Todo 4", Status: model.NotStarted},
	}

	var createdIDs []int64
	for _, todo := range todos {
		body, _ := json.Marshal(todo)
		req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("failed to create api: got status %d", w.Code)
		}

		var created model.Todo
		json.NewDecoder(w.Body).Decode(&created)
		createdIDs = append(createdIDs, created.ID)
	}

	// Test filtering by status
	filterReq := httptest.NewRequest(http.MethodGet, "/todos?status=not_started", nil)
	filterW := httptest.NewRecorder()
	r.ServeHTTP(filterW, filterReq)

	if filterW.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", filterW.Code)
	}

	var filteredTodos []model.Todo
	json.NewDecoder(filterW.Body).Decode(&filteredTodos)

	if len(filteredTodos) != 2 {
		t.Errorf("expected 2 todos with NotStarted status, got %d", len(filteredTodos))
	}

	for _, todo := range filteredTodos {
		if todo.Status != model.NotStarted {
			t.Errorf("expected status NotStarted, got %q", todo.Status)
		}
	}

	// Test sorting
	sortReq := httptest.NewRequest(http.MethodGet, "/todos?sort_by=name&order=asc", nil)
	sortW := httptest.NewRecorder()
	r.ServeHTTP(sortW, sortReq)

	if sortW.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", sortW.Code)
	}

	var sortedTodos []model.Todo
	json.NewDecoder(sortW.Body).Decode(&sortedTodos)

	if len(sortedTodos) != 4 {
		t.Fatalf("expected 4 todos, got %d", len(sortedTodos))
	}

	// Verify sorted order
	if sortedTodos[0].Name != "Todo 1" || sortedTodos[1].Name != "Todo 2" ||
		sortedTodos[2].Name != "Todo 3" || sortedTodos[3].Name != "Todo 4" {
		t.Errorf("todos not sorted correctly: %v", getTodoNames(sortedTodos))
	}
}

// TestIntegration_ErrorHandling tests error scenarios
func TestIntegration_ErrorHandling(t *testing.T) {
	store, cleanup := setupIntegrationTestDB(t)
	defer cleanup()

	svc := api.NewService(store)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	RegisterRoutes(r, svc)

	// Test creating api with empty name
	emptyBody := []byte(`{"name":""}`)
	req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewReader(emptyBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for empty name, got %d", w.Code)
	}

	// Test getting non-existent api
	getReq := httptest.NewRequest(http.MethodGet, "/todos/999", nil)
	getW := httptest.NewRecorder()
	r.ServeHTTP(getW, getReq)

	if getW.Code != http.StatusNotFound {
		t.Errorf("expected status 404 for non-existent api, got %d", getW.Code)
	}

	// Test updating non-existent api
	updateBody := []byte(`{"name":"Updated"}`)
	updateReq := httptest.NewRequest(http.MethodPut, "/todos/999", bytes.NewReader(updateBody))
	updateReq.Header.Set("Content-Type", "application/json")
	updateW := httptest.NewRecorder()
	r.ServeHTTP(updateW, updateReq)

	if updateW.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for updating non-existent api, got %d", updateW.Code)
	}

	// Test deleting non-existent api
	deleteReq := httptest.NewRequest(http.MethodDelete, "/todos/999", nil)
	deleteW := httptest.NewRecorder()
	r.ServeHTTP(deleteW, deleteReq)

	if deleteW.Code != http.StatusNotFound {
		t.Errorf("expected status 404 for deleting non-existent api, got %d", deleteW.Code)
	}

	// Test invalid JSON
	invalidReq := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewReader([]byte("invalid json")))
	invalidReq.Header.Set("Content-Type", "application/json")
	invalidW := httptest.NewRecorder()
	r.ServeHTTP(invalidW, invalidReq)

	if invalidW.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for invalid JSON, got %d", invalidW.Code)
	}
}

// Helper functions
func getTodoNames(todos []model.Todo) []string {
	names := make([]string, len(todos))
	for i, t := range todos {
		names[i] = t.Name
	}
	return names
}
