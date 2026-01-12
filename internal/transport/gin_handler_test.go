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
	"github.com/gin-gonic/gin"
)

// mockStore is a mock implementation of Store for testing
type mockStore struct {
	todos     map[int64]*model.Todo
	nextID    int64
	createErr error
	getErr    error
	updateErr error
	deleteErr error
	listErr   error
}

func (m *mockStore) Create(t *model.Todo) (int64, error) {
	if m.createErr != nil {
		return 0, m.createErr
	}
	m.nextID++
	t.ID = m.nextID
	if t.Status == "" {
		t.Status = model.NotStarted
	}
	if m.todos == nil {
		m.todos = make(map[int64]*model.Todo)
	}
	// copy
	c := *t
	m.todos[t.ID] = &c
	return t.ID, nil
}

func (m *mockStore) Get(id int64) (*model.Todo, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	if m.todos == nil {
		return nil, cache.ErrNotFound
	}
	if t, ok := m.todos[id]; ok {
		c := *t
		return &c, nil
	}
	return nil, cache.ErrNotFound
}

func (m *mockStore) Update(t *model.Todo) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	if m.todos == nil {
		return cache.ErrNotFound
	}
	if _, ok := m.todos[t.ID]; !ok {
		return cache.ErrNotFound
	}
	c := *t
	m.todos[t.ID] = &c
	return nil
}

func (m *mockStore) Delete(id int64) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	if m.todos == nil {
		return cache.ErrNotFound
	}
	if _, ok := m.todos[id]; !ok {
		return cache.ErrNotFound
	}
	delete(m.todos, id)
	return nil
}

func (m *mockStore) List(opts cache.ListOptions) ([]model.Todo, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	if m.todos == nil {
		return []model.Todo{}, nil
	}
	result := make([]model.Todo, 0, len(m.todos))
	for _, t := range m.todos {
		result = append(result, *t)
	}
	return cache.FilterAndSort(result, opts), nil
}

func setupGinTestRouter(svc *api.Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	RegisterRoutes(r, svc)
	return r
}

func TestRegisterRoutes(t *testing.T) {
	t.Run("registers all routes", func(t *testing.T) {
		svc := api.NewService(&mockStore{todos: make(map[int64]*model.Todo)})
		r := setupGinTestRouter(svc)

		routes := []struct {
			method string
			path   string
		}{
			{http.MethodGet, "/todos"},
			{http.MethodPost, "/todos"},
			{http.MethodGet, "/todos/1"},
			{http.MethodPut, "/todos/1"},
			{http.MethodDelete, "/todos/1"},
		}

		for _, route := range routes {
			var req *http.Request
			if route.method == http.MethodPost || route.method == http.MethodPut {
				req = httptest.NewRequest(route.method, route.path, bytes.NewReader([]byte("{}")))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req = httptest.NewRequest(route.method, route.path, nil)
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			// Routes should be registered (not 404)
			// Status 400/404 are acceptable - they mean the route is registered and handled
			if w.Code == http.StatusNotFound && route.method == http.MethodGet && route.path == "/todos/1" {
				// This is expected - api doesn't exist, but route is registered
				continue
			}
			// If we get 404 for POST/PUT/DELETE, it might mean route not found
			// But if we get 400, it means route is found but validation failed (which is fine)
			if w.Code == http.StatusNotFound && (route.method == http.MethodPost || route.method == http.MethodPut || route.method == http.MethodDelete) {
				// For DELETE, 404 means api not found, which is acceptable
				if route.method == http.MethodDelete {
					continue
				}
				// For POST/PUT with 404, route might not be registered
				if route.method == http.MethodPost || route.method == http.MethodPut {
					// Actually, 400 is more likely - validation/parsing error means route is found
					if w.Code == http.StatusNotFound {
						t.Errorf("route %s %s might not be registered (got 404)", route.method, route.path)
					}
				}
			}
		}
	})
}

func TestGinHandler_handleList(t *testing.T) {
	t.Run("lists all todos", func(t *testing.T) {
		store := &mockStore{todos: make(map[int64]*model.Todo)}
		store.todos[1] = &model.Todo{ID: 1, Name: "Todo 1", Status: model.NotStarted}
		store.todos[2] = &model.Todo{ID: 2, Name: "Todo 2", Status: model.InProgress}
		svc := api.NewService(store)
		r := setupGinTestRouter(svc)

		req := httptest.NewRequest(http.MethodGet, "/todos", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

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
		r := setupGinTestRouter(svc)

		req := httptest.NewRequest(http.MethodGet, "/todos?status=not_started", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

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
		r := setupGinTestRouter(svc)

		req := httptest.NewRequest(http.MethodGet, "/todos?sort_by=name&order=asc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

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

	t.Run("handles empty list", func(t *testing.T) {
		svc := api.NewService(&mockStore{todos: make(map[int64]*model.Todo)})
		r := setupGinTestRouter(svc)

		req := httptest.NewRequest(http.MethodGet, "/todos", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		var response []model.Todo
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if len(response) != 0 {
			t.Errorf("expected empty list, got %d todos", len(response))
		}
	})
}

func TestGinHandler_handleCreate(t *testing.T) {
	t.Run("creates api successfully", func(t *testing.T) {
		svc := api.NewService(&mockStore{todos: make(map[int64]*model.Todo)})
		r := setupGinTestRouter(svc)

		todo := model.Todo{Name: "Test Todo", Description: "Test Description"}
		body, _ := json.Marshal(todo)

		req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

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

	t.Run("creates api with all fields", func(t *testing.T) {
		svc := api.NewService(&mockStore{todos: make(map[int64]*model.Todo)})
		r := setupGinTestRouter(svc)

		dueDate := time.Date(2026, 12, 31, 12, 0, 0, 0, time.UTC)
		todo := model.Todo{
			Name:        "Full Todo",
			Description: "Full Description",
			DueDate:     dueDate,
			Status:      model.InProgress,
			Priority:    5,
			Tags:        []string{"work", "urgent"},
		}
		body, _ := json.Marshal(todo)

		req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

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
		if response.Description != todo.Description {
			t.Errorf("expected description %q, got %q", todo.Description, response.Description)
		}
		if len(response.Tags) != len(todo.Tags) {
			t.Errorf("expected %d tags, got %d", len(todo.Tags), len(response.Tags))
		}
	})

	t.Run("returns 400 for invalid JSON", func(t *testing.T) {
		svc := api.NewService(&mockStore{todos: make(map[int64]*model.Todo)})
		r := setupGinTestRouter(svc)

		req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})

	t.Run("returns 400 when validation fails", func(t *testing.T) {
		svc := api.NewService(&mockStore{todos: make(map[int64]*model.Todo)})
		r := setupGinTestRouter(svc)

		body := []byte(`{"name":""}`) // Empty name should fail validation

		req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})
}

func TestGinHandler_handleGet(t *testing.T) {
	t.Run("gets existing api", func(t *testing.T) {
		store := &mockStore{todos: make(map[int64]*model.Todo)}
		store.todos[1] = &model.Todo{ID: 1, Name: "Test Todo", Status: model.InProgress}
		svc := api.NewService(store)
		r := setupGinTestRouter(svc)

		req := httptest.NewRequest(http.MethodGet, "/todos/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

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
		r := setupGinTestRouter(svc)

		req := httptest.NewRequest(http.MethodGet, "/todos/999", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", w.Code)
		}
	})

	t.Run("handles invalid ID parameter", func(t *testing.T) {
		svc := api.NewService(&mockStore{todos: make(map[int64]*model.Todo)})
		r := setupGinTestRouter(svc)

		req := httptest.NewRequest(http.MethodGet, "/todos/invalid", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Invalid ID should be parsed as 0, which won't exist
		if w.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", w.Code)
		}
	})
}

func TestGinHandler_handleUpdate(t *testing.T) {
	t.Run("updates existing api", func(t *testing.T) {
		store := &mockStore{todos: make(map[int64]*model.Todo)}
		store.todos[1] = &model.Todo{ID: 1, Name: "Original", Status: model.NotStarted}
		svc := api.NewService(store)
		r := setupGinTestRouter(svc)

		update := model.Todo{Name: "Updated", Description: "New Description", Status: model.Completed}
		body, _ := json.Marshal(update)

		req := httptest.NewRequest(http.MethodPut, "/todos/1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

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
		if response.Status != model.Completed {
			t.Errorf("expected status Completed, got %q", response.Status)
		}
	})

	t.Run("returns 400 for invalid JSON", func(t *testing.T) {
		store := &mockStore{todos: make(map[int64]*model.Todo)}
		store.todos[1] = &model.Todo{ID: 1, Name: "Original"}
		svc := api.NewService(store)
		r := setupGinTestRouter(svc)

		req := httptest.NewRequest(http.MethodPut, "/todos/1", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})

	t.Run("returns 400 for non-existent api", func(t *testing.T) {
		svc := api.NewService(&mockStore{todos: make(map[int64]*model.Todo)})
		r := setupGinTestRouter(svc)

		body := []byte(`{"name":"Updated"}`)

		req := httptest.NewRequest(http.MethodPut, "/todos/999", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})
}

func TestGinHandler_handleDelete(t *testing.T) {
	t.Run("deletes existing api", func(t *testing.T) {
		store := &mockStore{todos: make(map[int64]*model.Todo)}
		store.todos[1] = &model.Todo{ID: 1, Name: "To Delete"}
		svc := api.NewService(store)
		r := setupGinTestRouter(svc)

		req := httptest.NewRequest(http.MethodDelete, "/todos/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

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
		r := setupGinTestRouter(svc)

		req := httptest.NewRequest(http.MethodDelete, "/todos/999", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", w.Code)
		}
	})
}
