package cache

import (
	"errors"
	"sync"

	"github.com/conbanwa/todo/internal/model"
)

var ErrNotFound = errors.New("api not found")

type ListOptions struct {
	Status    model.Status
	SortBy    string // due_date, status, name
	SortOrder string // asc, desc
}

type InMemoryStore struct {
	mu    sync.RWMutex
	next  int64
	items map[int64]*model.Todo
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{items: make(map[int64]*model.Todo), next: 1}
}

func (s *InMemoryStore) Create(t *model.Todo) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	t.ID = s.next
	s.next++
	if t.Status == "" {
		t.Status = model.NotStarted
	}
	// copy
	c := *t
	s.items[t.ID] = &c
	return t.ID, nil
}

func (s *InMemoryStore) Get(id int64) (*model.Todo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if v, ok := s.items[id]; ok {
		c := *v
		return &c, nil
	}
	return nil, ErrNotFound
}

func (s *InMemoryStore) Update(t *model.Todo) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[t.ID]; !ok {
		return ErrNotFound
	}
	c := *t
	s.items[t.ID] = &c
	return nil
}

func (s *InMemoryStore) Delete(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[id]; !ok {
		return ErrNotFound
	}
	delete(s.items, id)
	return nil
}

func (s *InMemoryStore) List(opts ListOptions) ([]model.Todo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	// copy items into a slice
	out := make([]model.Todo, 0, len(s.items))
	for _, v := range s.items {
		out = append(out, *v)
	}

	// delegate filtering/sorting to helper
	res := FilterAndSort(out, opts)
	return res, nil
}
