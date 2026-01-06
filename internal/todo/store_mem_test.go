package todo

import (
	"testing"
	"time"
)

func TestInMemoryStore_CreateGetUpdateDelete_List(t *testing.T) {
	s := NewInMemoryStore()

	id, err := s.Create(&Todo{Name: "task1", DueDate: time.Now().Add(24 * time.Hour)})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	got, err := s.Get(id)
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if got.Name != "task1" {
		t.Fatalf("unexpected name: %s", got.Name)
	}

	got.Description = "updated"
	if err := s.Update(got); err != nil {
		t.Fatalf("update failed: %v", err)
	}

	if err := s.Delete(id); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	// list
	_, _ = s.List(ListOptions{})
}
