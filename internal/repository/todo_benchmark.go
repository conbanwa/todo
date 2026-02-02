package repository

import (
	"database/sql"
	"testing"

	"github.com/conbanwa/todo/internal/model"
	_ "github.com/glebarez/sqlite" // Your SQLite driver
)

func setupBenchDB(b *testing.B) *sql.DB {
	// In-memory DB for fast, isolated benchmarks
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		b.Fatal(err)
	}

	// Create table — adapt this EXACTLY to your schema
	_, err = db.Exec(`
		CREATE TABLE todos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			description TEXT,
			completed BOOLEAN DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME
		);
	`)
	if err != nil {
		b.Fatal(err)
	}

	return db
}

func BenchmarkCreateTodo(b *testing.B) {
	db := setupBenchDB(b)
	defer db.Close()
	repo := NewTodoRepository(db) // Your constructor

	todo := &model.Todo{
		Title:       "Benchmark Todo",
		Description: "Testing performance",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := repo.Create(todo); err != nil { // Adapt method name/signature
			b.Fatal(err)
		}
	}
}

func BenchmarkGetAllTodos_Empty(b *testing.B) {
	db := setupBenchDB(b)
	defer db.Close()
	repo := NewTodoRepository(db)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := repo.GetAll(); err != nil { // Adapt method name
			b.Fatal(err)
		}
	}
}

func BenchmarkGetAllTodos_1000(b *testing.B) {
	db := setupBenchDB(b)
	defer db.Close()
	repo := NewTodoRepository(db)

	// Pre-populate with 1000 todos
	for i := 0; i < 1000; i++ {
		todo := &model.Todo{Title: "Preload"}
		if err := repo.Create(todo); err != nil {
			b.Fatal(err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := repo.GetAll(); err != nil {
			b.Fatal(err)
		}
	}
}

// Add more if you have Update/Delete methods