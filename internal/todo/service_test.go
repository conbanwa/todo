package todo

import (
	"testing"
	"time"
)

func TestService_Create_Validation(t *testing.T) {
	s := NewService(NewInMemoryStore())
	_, err := s.Create(&Todo{})
	if err == nil {
		t.Fatalf("expected error for missing name")
	}

	id, err := s.Create(&Todo{Name: "ok", DueDate: time.Now().Add(48 * time.Hour)})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if id == 0 {
		t.Fatalf("expected id > 0")
	}
}
