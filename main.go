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

	svc := todo.NewService(todo.NewInMemoryStore())

	r := gin.Default()

	// update swagger host to match runtime
	docs.SwaggerInfo.Host = "localhost:8080"

	// serve static swagger.json
	r.StaticFile("/docs/swagger.json", "docs/swagger.json")

	// swagger UI at /swagger/index.html and /swagger/*any
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/docs/swagger.json")))

	// register API routes
	todo.RegisterRoutes(r, svc)

	log.Printf("starting server on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server exited: %v", err)
	}
}
