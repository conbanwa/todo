package cache

import (
	"testing"
	"time"

	"github.com/conbanwa/todo/internal/model"
)

func TestInMemoryStore_CreateGetUpdateDelete_List(t *testing.T) {
	s := NewInMemoryStore()

	id, err := s.Create(&model.Todo{Name: "task1", DueDate: time.Now().Add(24 * time.Hour)})
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

func TestInMemoryStore_ListOptions_SortAndFilter(t *testing.T) {
	s := NewInMemoryStore()

	// deterministic times
	t1 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)
	t3 := time.Date(2026, 1, 3, 0, 0, 0, 0, time.UTC)

	_, _ = s.Create(&model.Todo{Name: "alpha", DueDate: t2, Status: model.InProgress})
	_, _ = s.Create(&model.Todo{Name: "bravo", DueDate: t1, Status: model.NotStarted})
	_, _ = s.Create(&model.Todo{Name: "charlie", DueDate: t3, Status: model.Completed})

	// filter by status
	list, err := s.List(ListOptions{Status: model.NotStarted})
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if len(list) != 1 || list[0].Name != "bravo" {
		t.Fatalf("expected only bravo, got: %v", list)
	}

	// sort by due_date asc
	list, _ = s.List(ListOptions{SortBy: "due_date", SortOrder: "asc"})
	if len(list) != 3 || list[0].Name != "bravo" || list[1].Name != "alpha" || list[2].Name != "charlie" {
		t.Fatalf("unexpected due_date asc order: %v", []string{list[0].Name, list[1].Name, list[2].Name})
	}

	// sort by due_date desc
	list, _ = s.List(ListOptions{SortBy: "due_date", SortOrder: "desc"})
	if list[0].Name != "charlie" || list[2].Name != "bravo" {
		t.Fatalf("unexpected due_date desc order: %v", []string{list[0].Name, list[1].Name, list[2].Name})
	}

	// sort by name asc
	list, _ = s.List(ListOptions{SortBy: "name", SortOrder: "asc"})
	if list[0].Name != "alpha" || list[1].Name != "bravo" || list[2].Name != "charlie" {
		t.Fatalf("unexpected name asc order: %v", []string{list[0].Name, list[1].Name, list[2].Name})
	}

	// sort by status asc
	list, _ = s.List(ListOptions{SortBy: "status", SortOrder: "asc"})
	if len(list) != 3 {
		t.Fatalf("expected 3 items for status sort, got %d", len(list))
	}
}
