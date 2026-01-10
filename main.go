// @title Todo API
// @version 0.1.0
// @description Minimal Todo API generated with swag
// @termsOfService http://example.com/terms/
// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	docs "github.com/conbanwa/todo/docs"
	"github.com/conbanwa/todo/internal/todo"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	addr := ":8080"
	if v := os.Getenv("PORT"); v != "" {
		addr = ":" + v
	}

	// Get database path from environment or use default
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "todos.db"
	}

	// Initialize SQLite store
	store, err := todo.NewSQLiteStore(dbPath)
	if err != nil {
		log.Fatalf("failed to initialize SQLite store: %v", err)
	}
	defer func() {
		if err := store.Close(); err != nil {
			log.Printf("error closing database: %v", err)
		}
	}()

	svc := todo.NewService(store)

	r := gin.Default()

	// update swagger host to match runtime
	docs.SwaggerInfo.Host = "localhost:8080"

	// serve static swagger.json
	r.StaticFile("/docs/swagger.json", "docs/swagger.json")

	// swagger UI at /swagger/index.html and /swagger/*any
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/docs/swagger.json")))

	// register API routes
	todo.RegisterRoutes(r, svc)

	// Graceful shutdown handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
	log.Printf("starting server on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server exited: %v", err)
	}
	}()

	// Wait for interrupt signal
	<-sigChan
	log.Println("shutting down server...")
}
