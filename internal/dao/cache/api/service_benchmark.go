package api_test

import (
	"testing"
	"time"

	"github.com/conbanwa/todo/internal/dao/cache"
	"github.com/conbanwa/todo/internal/dao/cache/api"
	"github.com/conbanwa/todo/internal/dao/db"
	"github.com/conbanwa/todo/internal/model"
)

func setupBenchService(b *testing.B) *api.Service {
	store, err := db.NewSQLiteStore(":memory:")
	if err != nil {
		b.Fatal(err)
	}
	return api.NewService(store)
}

func BenchmarkService_Create(b *testing.B) {
	svc := setupBenchService(b)
	todo := &model.Todo{
		Name:        "Benchmark Todo",
		Description: "Performance test",
		Status:      model.NotStarted,
		Priority:    5,
		Tags:        []string{"bench", "test"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := svc.Create(todo)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkService_List_Empty(b *testing.B) {
	svc := setupBenchService(b)
	opts := cache.ListOptions{} // Zero value = default (full list, no filters)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := svc.List(opts)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkService_List_1000Items(b *testing.B) {
	svc := setupBenchService(b)
	opts := cache.ListOptions{}

	// Pre-populate 1000 todos
	for i := 0; i < 1000; i++ {
		todo := &model.Todo{
			Name:   "Preload Todo",
			Status: model.InProgress,
		}
		_, err := svc.Create(todo)
		if err != nil {
			b.Fatal(err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := svc.List(opts)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkService_Get(b *testing.B) {
	svc := setupBenchService(b)

	// Create one todo to fetch
	todo := &model.Todo{Name: "Test Todo", Status: model.Completed}
	id, err := svc.Create(todo)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := svc.Get(id)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkService_Update(b *testing.B) {
	svc := setupBenchService(b)

	// Create one todo to update
	todo := &model.Todo{Name: "Updatable Todo", Status: model.NotStarted}
	id, err := svc.Create(todo)
	if err != nil {
		b.Fatal(err)
	}
	updateTodo := &model.Todo{
		ID:          id,
		Name:        "Updated Todo",
		Status:      model.Completed,
		Description: "Updated description",
		DueDate:     time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := svc.Update(updateTodo)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkService_Delete(b *testing.B) {
	svc := setupBenchService(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		// Create a fresh todo each iteration (delete requires existing ID)
		todo := &model.Todo{Name: "Deletable Todo", Status: model.NotStarted}
		id, err := svc.Create(todo)
		if err != nil {
			b.Fatal(err)
		}
		b.StartTimer()

		err = svc.Delete(id)
		if err != nil {
			b.Fatal(err)
		}
	}
}