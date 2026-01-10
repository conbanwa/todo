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

	// Initialize WebSocket hub
	hub := todo.NewHub()
	go hub.Run()

	r := gin.Default()

	// update swagger host to match runtime
	docs.SwaggerInfo.Host = "localhost:8080"

	// serve static swagger.json
	r.StaticFile("/docs/swagger.json", "docs/swagger.json")

	// swagger UI at /swagger/index.html and /swagger/*any
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/docs/swagger.json")))

	// serve static files (HTML, CSS, JS)
	r.Static("/static", "./static")

	// serve index.html at root
	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	// register WebSocket route
	r.GET("/ws", func(c *gin.Context) {
		todo.HandleWebSocket(c, hub)
	})

	// register API routes with WebSocket broadcasting
	todo.RegisterRoutesWithHub(r, svc, hub)

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

	// Close WebSocket hub gracefully
	hub.Close()
}
