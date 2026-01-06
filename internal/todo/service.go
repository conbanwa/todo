package todo

import (
	"time"
)

type Store interface {
	Create(*Todo) (int64, error)
	Get(int64) (*Todo, error)
	Update(*Todo) error
	Delete(int64) error
	List(ListOptions) ([]Todo, error)
}

type Service struct {
	store Store
}

func NewService(s Store) *Service { return &Service{store: s} }

func (s *Service) Create(t *Todo) (int64, error) {
	if t.Name == "" {
		return 0, ErrInvalid("name is required")
	}
	if t.DueDate.IsZero() {
		t.DueDate = time.Now().Add(24 * time.Hour)
	}
	return s.store.Create(t)
}

func (s *Service) Get(id int64) (*Todo, error) { return s.store.Get(id) }

func (s *Service) Update(t *Todo) error {
	if t.ID == 0 {
		return ErrInvalid("id is required")
	}
	return s.store.Update(t)
}

func (s *Service) Delete(id int64) error { return s.store.Delete(id) }

func (s *Service) List(opts ListOptions) ([]Todo, error) { return s.store.List(opts) }

type ErrInvalid string

func (e ErrInvalid) Error() string { return string(e) }
