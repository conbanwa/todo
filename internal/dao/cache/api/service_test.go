package api

import (
	"testing"
	"time"

	"github.com/conbanwa/todo/internal/dao/cache"
	"github.com/conbanwa/todo/internal/model"
)

func TestService_Create_Validation(t *testing.T) {
	s := NewService(cache.NewInMemoryStore())
	_, err := s.Create(&model.Todo{})
	if err == nil {
		t.Fatalf("expected error for missing name")
	}

	id, err := s.Create(&model.Todo{Name: "ok", DueDate: time.Now().Add(48 * time.Hour)})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if id == 0 {
		t.Fatalf("expected id > 0")
	}
}
