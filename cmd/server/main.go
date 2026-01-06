package main

import (
	"log"
	"net/http"
	"os"

	"github.com/conbanwa/todo/internal/todo"
)

func main() {
	addr := ":8080"
	if v := os.Getenv("PORT"); v != "" {
		addr = ":" + v
	}

	svc := todo.NewService(todo.NewInMemoryStore())
	mux := http.NewServeMux()
	mux.Handle("/todos", todo.NewHandler(svc))
	mux.Handle("/todos/", todo.NewHandler(svc))

	log.Printf("starting server on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server exited: %v", err)
	}
}
