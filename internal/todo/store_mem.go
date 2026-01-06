package todo

import (
	"errors"
	"sort"
	"sync"
)

var ErrNotFound = errors.New("todo not found")

type ListOptions struct {
	Status    Status
	SortBy    string // due_date, status, name
	SortOrder string // asc, desc
}

type InMemoryStore struct {
	mu    sync.RWMutex
	next  int64
	items map[int64]*Todo
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{items: make(map[int64]*Todo), next: 1}
}

func (s *InMemoryStore) Create(t *Todo) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	t.ID = s.next
	s.next++
	if t.Status == "" {
		t.Status = NotStarted
	}
	// copy
	c := *t
	s.items[t.ID] = &c
	return t.ID, nil
}

func (s *InMemoryStore) Get(id int64) (*Todo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if v, ok := s.items[id]; ok {
		c := *v
		return &c, nil
	}
	return nil, ErrNotFound
}

func (s *InMemoryStore) Update(t *Todo) error {
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

func (s *InMemoryStore) List(opts ListOptions) ([]Todo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Todo, 0, len(s.items))
	for _, v := range s.items {
		if opts.Status != "" && v.Status != opts.Status {
			continue
		}
		out = append(out, *v)
	}

	// sort
	cmp := func(i, j int) bool { return out[i].ID < out[j].ID }
	switch opts.SortBy {
	case "due_date":
		cmp = func(i, j int) bool { return out[i].DueDate.Before(out[j].DueDate) }
	case "status":
		cmp = func(i, j int) bool { return out[i].Status < out[j].Status }
	case "name":
		cmp = func(i, j int) bool { return out[i].Name < out[j].Name }
	}
	sort.SliceStable(out, func(i, j int) bool {
		if opts.SortOrder == "desc" {
			return !cmp(i, j)
		}
		return cmp(i, j)
	})

	return out, nil
}
