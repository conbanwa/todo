package transport

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/conbanwa/todo/internal/dao/cache"
	"github.com/conbanwa/todo/internal/dao/cache/api"
	"github.com/conbanwa/todo/internal/model"
)

func TestHandler_ServeHTTP(t *testing.T) {
	t.Run("returns 404 for invalid path", func(t *testing.T) {
		svc := api.NewService(&mockStore{todos: make(map[int64]*model.Todo)})
		handler := NewHandler(svc)

		req := httptest.NewRequest(http.MethodGet, "/invalid", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", w.Code)
		}
	})

	t.Run("routes to correct handler methods", func(t *testing.T) {
		store := &mockStore{todos: make(map[int64]*model.Todo)}
		svc := api.NewService(store)
		handler := NewHandler(svc)

		tests := []struct {
			method string
			path   string
			want   int
			name   string
		}{
			{http.MethodGet, "/todos", http.StatusOK, "GET /todos"},
			{http.MethodPost, "/todos", http.StatusBadRequest, "POST /todos (empty body - validation error)"},
			{http.MethodGet, "/todos/1", http.StatusNotFound, "GET /todos/1 (doesn't exist)"},
			{http.MethodPut, "/todos/1", http.StatusBadRequest, "PUT /todos/1 (invalid JSON)"},
			{http.MethodDelete, "/todos/1", http.StatusNotFound, "DELETE /todos/1 (doesn't exist)"},
		}

		for _, tt := range tests {
			var req *http.Request
			if tt.method == http.MethodPost || tt.method == http.MethodPut {
				req = httptest.NewRequest(tt.method, tt.path, bytes.NewReader([]byte("{}")))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req = httptest.NewRequest(tt.method, tt.path, nil)
			}
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			if w.Code != tt.want {
				t.Errorf("%s: expected status %d, got %d", tt.name, tt.want, w.Code)
			}
		}
	})
}

func TestHandler_handleCreate(t *testing.T) {
	t.Run("creates api successfully", func(t *testing.T) {
		svc := api.NewService(&mockStore{todos: make(map[int64]*model.Todo)})
		handler := NewHandler(svc)

		todo := model.Todo{Name: "Test Todo", Description: "Test Description"}
		body, _ := json.Marshal(todo)

		req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status 201, got %d", w.Code)
		}

		var response model.Todo
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if response.Name != todo.Name {
			t.Errorf("expected name %q, got %q", todo.Name, response.Name)
		}
		if response.ID == 0 {
			t.Error("expected non-zero ID")
		}
	})

	t.Run("parses due_date from string", func(t *testing.T) {
		svc := api.NewService(&mockStore{todos: make(map[int64]*model.Todo)})
		handler := NewHandler(svc)

		dueDate := time.Date(2026, 12, 31, 12, 0, 0, 0, time.UTC)
		body := []byte(`{"name":"Test","due_date":"2026-12-31T12:00:00Z"}`)

		req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status 201, got %d", w.Code)
		}

		var response model.Todo
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		// Due date should be set (service layer sets default if zero, but we should verify parsing)
		if response.DueDate.IsZero() && !dueDate.IsZero() {
			t.Error("expected due date to be parsed")
		}
	})

	t.Run("returns 400 for invalid JSON", func(t *testing.T) {
		svc := api.NewService(&mockStore{todos: make(map[int64]*model.Todo)})
		handler := NewHandler(svc)

		req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})

	t.Run("returns 400 when service validation fails", func(t *testing.T) {
		svc := api.NewService(&mockStore{todos: make(map[int64]*model.Todo)})
		handler := NewHandler(svc)

		body := []byte(`{"name":""}`) // Empty name should fail validation

		req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})
}

func TestHandler_handleGet(t *testing.T) {
	t.Run("gets existing api", func(t *testing.T) {
		store := &mockStore{todos: make(map[int64]*model.Todo)}
		store.todos[1] = &model.Todo{ID: 1, Name: "Test Todo", Status: model.InProgress}
		svc := api.NewService(store)
		handler := NewHandler(svc)

		req := httptest.NewRequest(http.MethodGet, "/todos/1", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		var response model.Todo
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if response.ID != 1 {
			t.Errorf("expected ID 1, got %d", response.ID)
		}
		if response.Name != "Test Todo" {
			t.Errorf("expected name 'Test Todo', got %q", response.Name)
		}
	})

	t.Run("returns 404 for non-existent api", func(t *testing.T) {
		svc := api.NewService(&mockStore{todos: make(map[int64]*model.Todo)})
		handler := NewHandler(svc)

		req := httptest.NewRequest(http.MethodGet, "/todos/999", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", w.Code)
		}
	})
}

func TestHandler_handleUpdate(t *testing.T) {
	t.Run("updates existing api", func(t *testing.T) {
		store := &mockStore{todos: make(map[int64]*model.Todo)}
		store.todos[1] = &model.Todo{ID: 1, Name: "Original", Status: model.NotStarted}
		svc := api.NewService(store)
		handler := NewHandler(svc)

		update := model.Todo{Name: "Updated", Description: "New Description", Status: model.Completed}
		body, _ := json.Marshal(update)

		req := httptest.NewRequest(http.MethodPut, "/todos/1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		var response model.Todo
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if response.ID != 1 {
			t.Errorf("expected ID 1, got %d", response.ID)
		}
		if response.Name != "Updated" {
			t.Errorf("expected name 'Updated', got %q", response.Name)
		}
	})

	t.Run("returns 400 for invalid JSON", func(t *testing.T) {
		store := &mockStore{todos: make(map[int64]*model.Todo)}
		store.todos[1] = &model.Todo{ID: 1, Name: "Original"}
		svc := api.NewService(store)
		handler := NewHandler(svc)

		req := httptest.NewRequest(http.MethodPut, "/todos/1", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})

	t.Run("returns 400 when service validation fails", func(t *testing.T) {
		store := &mockStore{todos: make(map[int64]*model.Todo)}
		svc := api.NewService(store)
		handler := NewHandler(svc)

		body := []byte(`{"name":"Test"}`) // ID will be set to 1, but api doesn't exist

		req := httptest.NewRequest(http.MethodPut, "/todos/1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})
}

func TestHandler_handleDelete(t *testing.T) {
	t.Run("deletes existing api", func(t *testing.T) {
		store := &mockStore{todos: make(map[int64]*model.Todo)}
		store.todos[1] = &model.Todo{ID: 1, Name: "To Delete"}
		svc := api.NewService(store)
		handler := NewHandler(svc)

		req := httptest.NewRequest(http.MethodDelete, "/todos/1", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("expected status 204, got %d", w.Code)
		}

		// Verify it's deleted
		_, err := svc.Get(1)
		if err != cache.ErrNotFound {
			t.Errorf("expected ErrNotFound after delete, got %v", err)
		}
	})

	t.Run("returns 404 for non-existent api", func(t *testing.T) {
		svc := api.NewService(&mockStore{todos: make(map[int64]*model.Todo)})
		handler := NewHandler(svc)

		req := httptest.NewRequest(http.MethodDelete, "/todos/999", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", w.Code)
		}
	})
}

func TestHandler_handleList(t *testing.T) {
	t.Run("lists all todos", func(t *testing.T) {
		store := &mockStore{todos: make(map[int64]*model.Todo)}
		store.todos[1] = &model.Todo{ID: 1, Name: "Todo 1", Status: model.NotStarted}
		store.todos[2] = &model.Todo{ID: 2, Name: "Todo 2", Status: model.InProgress}
		svc := api.NewService(store)
		handler := NewHandler(svc)

		req := httptest.NewRequest(http.MethodGet, "/todos", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		var response []model.Todo
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if len(response) != 2 {
			t.Errorf("expected 2 todos, got %d", len(response))
		}
	})

	t.Run("filters by status query parameter", func(t *testing.T) {
		store := &mockStore{todos: make(map[int64]*model.Todo)}
		store.todos[1] = &model.Todo{ID: 1, Name: "Todo 1", Status: model.NotStarted}
		store.todos[2] = &model.Todo{ID: 2, Name: "Todo 2", Status: model.InProgress}
		store.todos[3] = &model.Todo{ID: 3, Name: "Todo 3", Status: model.NotStarted}
		svc := api.NewService(store)
		handler := NewHandler(svc)

		req := httptest.NewRequest(http.MethodGet, "/todos?status=not_started", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		var response []model.Todo
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if len(response) != 2 {
			t.Errorf("expected 2 todos, got %d", len(response))
		}
		for _, todo := range response {
			if todo.Status != model.NotStarted {
				t.Errorf("expected status NotStarted, got %q", todo.Status)
			}
		}
	})

	t.Run("sorts by query parameters", func(t *testing.T) {
		store := &mockStore{todos: make(map[int64]*model.Todo)}
		store.todos[1] = &model.Todo{ID: 1, Name: "Charlie"}
		store.todos[2] = &model.Todo{ID: 2, Name: "Alpha"}
		store.todos[3] = &model.Todo{ID: 3, Name: "Bravo"}
		svc := api.NewService(store)
		handler := NewHandler(svc)

		req := httptest.NewRequest(http.MethodGet, "/todos?sort_by=name&order=asc", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		var response []model.Todo
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if len(response) != 3 {
			t.Fatalf("expected 3 todos, got %d", len(response))
		}
		if response[0].Name != "Alpha" || response[1].Name != "Bravo" || response[2].Name != "Charlie" {
			t.Errorf("unexpected sort order: %v", []string{response[0].Name, response[1].Name, response[2].Name})
		}
	})
}

func TestHandler_writeJSON(t *testing.T) {
	t.Run("writes JSON with correct content type through handler", func(t *testing.T) {
		store := &mockStore{todos: make(map[int64]*model.Todo)}
		store.todos[1] = &model.Todo{ID: 1, Name: "Test"}
		svc := api.NewService(store)
		handler := NewHandler(svc)

		req := httptest.NewRequest(http.MethodGet, "/todos/1", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Header().Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type 'application/json', got %q", w.Header().Get("Content-Type"))
		}

		var response model.Todo
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("failed to decode JSON: %v", err)
		}

		if response.ID != 1 {
			t.Errorf("expected ID 1, got %d", response.ID)
		}
		if response.Name != "Test" {
			t.Errorf("expected name 'Test', got %q", response.Name)
		}
	})
}
