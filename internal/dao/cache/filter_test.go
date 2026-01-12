package cache

import (
	"testing"
	"time"

	"github.com/conbanwa/todo/internal/model"
)

func TestFilterAndSort_StatusAndSort(t *testing.T) {
	// deterministic times
	t1 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)
	t3 := time.Date(2026, 1, 3, 0, 0, 0, 0, time.UTC)

	items := []model.Todo{
		{ID: 1, Name: "alpha", DueDate: t2, Status: model.InProgress},
		{ID: 2, Name: "bravo", DueDate: t1, Status: model.NotStarted},
		{ID: 3, Name: "charlie", DueDate: t3, Status: model.Completed},
	}

	// filter by status
	got := FilterAndSort(items, ListOptions{Status: model.NotStarted})
	if len(got) != 1 || got[0].Name != "bravo" {
		t.Fatalf("expected only bravo, got: %v", got)
	}

	// due_date asc
	got = FilterAndSort(items, ListOptions{SortBy: "due_date", SortOrder: "asc"})
	if got[0].Name != "bravo" || got[1].Name != "alpha" || got[2].Name != "charlie" {
		t.Fatalf("unexpected due_date asc: %v", []string{got[0].Name, got[1].Name, got[2].Name})
	}

	// name desc
	got = FilterAndSort(items, ListOptions{SortBy: "name", SortOrder: "desc"})
	if got[0].Name != "charlie" || got[2].Name != "alpha" {
		t.Fatalf("unexpected name desc: %v", []string{got[0].Name, got[1].Name, got[2].Name})
	}
}
