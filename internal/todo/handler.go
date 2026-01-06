package todo

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) http.Handler { return &Handler{svc: svc} }

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/todos")
	if path == "" || path == "/" {
		switch r.Method {
		case http.MethodGet:
			h.handleList(w, r)
			return
		case http.MethodPost:
			h.handleCreate(w, r)
			return
		}
	}
	// /todos/{id}
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 1 {
		id, _ := strconv.ParseInt(parts[0], 10, 64)
		switch r.Method {
		case http.MethodGet:
			h.handleGet(w, r, id)
			return
		case http.MethodPut:
			h.handleUpdate(w, r, id)
			return
		case http.MethodDelete:
			h.handleDelete(w, r, id)
			return
		}
	}
	http.NotFound(w, r)
}

func (h *Handler) writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func (h *Handler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var t Todo
	body, _ := io.ReadAll(r.Body)
	if err := json.Unmarshal(body, &t); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if t.DueDate.IsZero() {
		// try to parse due_date from body if provided as string
		var tmp map[string]interface{}
		_ = json.Unmarshal(body, &tmp)
		if v, ok := tmp["due_date"].(string); ok {
			if parsed, err := time.Parse(time.RFC3339, v); err == nil {
				t.DueDate = parsed
			}
		}
	}
	id, err := h.svc.Create(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	t.ID = id
	w.WriteHeader(http.StatusCreated)
	h.writeJSON(w, t)
}

func (h *Handler) handleGet(w http.ResponseWriter, r *http.Request, id int64) {
	t, err := h.svc.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	h.writeJSON(w, t)
}

func (h *Handler) handleUpdate(w http.ResponseWriter, r *http.Request, id int64) {
	var t Todo
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	t.ID = id
	if err := h.svc.Update(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.writeJSON(w, t)
}

func (h *Handler) handleDelete(w http.ResponseWriter, r *http.Request, id int64) {
	if err := h.svc.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) handleList(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	opts := ListOptions{SortBy: q.Get("sort_by"), SortOrder: q.Get("order")}
	if s := q.Get("status"); s != "" {
		opts.Status = Status(strings.ToLower(s))
	}
	items, _ := h.svc.List(opts)
	h.writeJSON(w, items)
}
