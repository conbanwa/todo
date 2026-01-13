package api

import (
	"time"

	"github.com/conbanwa/todo/internal/dao/cache"
	"github.com/conbanwa/todo/internal/model"
)

type Store interface {
	Create(*model.Todo) (int64, error)
	Get(int64) (*model.Todo, error)
	Update(*model.Todo) error
	Delete(int64) error
	List(cache.ListOptions) ([]model.Todo, error)
}

type Service struct {
	store Store
	UserStore
	TeamStore TeamStore
}

func NewService(s Store) *Service { return &Service{store: s} }

func (s *Service) Create(t *model.Todo) (int64, error) {
	if t.Name == "" {
		return 0, ErrInvalid("name is required")
	}
	if t.DueDate.IsZero() {
		t.DueDate = time.Now().Add(24 * time.Hour)
	}
	return s.store.Create(t)
}

func (s *Service) Get(id int64) (*model.Todo, error) { return s.store.Get(id) }

func (s *Service) Update(t *model.Todo) error {
	if t.ID == 0 {
		return ErrInvalid("id is required")
	}
	return s.store.Update(t)
}

func (s *Service) Delete(id int64) error { return s.store.Delete(id) }

func (s *Service) List(opts cache.ListOptions) ([]model.Todo, error) { return s.store.List(opts) }

type ErrInvalid string

func (e ErrInvalid) Error() string { return string(e) }
