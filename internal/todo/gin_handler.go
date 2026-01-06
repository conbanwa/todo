package todo

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers todo REST routes on the provided Gin engine.
func RegisterRoutes(r *gin.Engine, svc *Service) {
	g := r.Group("/todos")
	g.GET("", func(c *gin.Context) { handleList(c, svc) })
	g.POST("", func(c *gin.Context) { handleCreate(c, svc) })
	g.GET(":id", func(c *gin.Context) { handleGet(c, svc) })
	g.PUT(":id", func(c *gin.Context) { handleUpdate(c, svc) })
	g.DELETE(":id", func(c *gin.Context) { handleDelete(c, svc) })
}

// @Summary List todos
// @Description Get a list of todos
// @Tags todos
// @Accept json
// @Produce json
// @Param sort_by query string false "sort field"
// @Param order query string false "sort order"
// @Success 200 {array} Todo
// @Router /todos [get]
func handleList(c *gin.Context, svc *Service) {
	q := c.Request.URL.Query()
	opts := ListOptions{SortBy: q.Get("sort_by"), SortOrder: q.Get("order")}
	if s := q.Get("status"); s != "" {
		opts.Status = Status(s)
	}
	items, _ := svc.List(opts)
	c.JSON(http.StatusOK, items)
}

// @Summary Create todo
// @Description Create a new todo
// @Tags todos
// @Accept json
// @Produce json
// @Param todo body Todo true "Todo to create"
// @Success 201 {object} Todo
// @Router /todos [post]
func handleCreate(c *gin.Context, svc *Service) {
	var t Todo
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}
	id, err := svc.Create(&t)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	t.ID = id
	c.JSON(http.StatusCreated, t)
}

// @Summary Get todo
// @Description Get a todo by ID
// @Tags todos
// @Accept json
// @Produce json
// @Param id path int true "Todo ID"
// @Success 200 {object} Todo
// @Failure 404 {object} map[string]string
// @Router /todos/{id} [get]
func handleGet(c *gin.Context, svc *Service) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	t, err := svc.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, t)
}

// @Summary Update todo
// @Description Update a todo by ID
// @Tags todos
// @Accept json
// @Produce json
// @Param id path int true "Todo ID"
// @Param todo body Todo true "Todo to update"
// @Success 200 {object} Todo
// @Failure 400 {object} map[string]string
// @Router /todos/{id} [put]
func handleUpdate(c *gin.Context, svc *Service) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var t Todo
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}
	t.ID = id
	if err := svc.Update(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, t)
}

// @Summary Delete todo
// @Description Delete a todo by ID
// @Tags todos
// @Param id path int true "Todo ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /todos/{id} [delete]
func handleDelete(c *gin.Context, svc *Service) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := svc.Delete(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
