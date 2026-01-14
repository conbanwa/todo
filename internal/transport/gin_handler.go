package transport

import (
	"net/http"
	"strconv"

	"github.com/conbanwa/todo/internal/dao/cache"
	"github.com/conbanwa/todo/internal/dao/cache/api"
	"github.com/conbanwa/todo/internal/model"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers api REST routes on the provided Gin engine.
func RegisterRoutes(r *gin.Engine, svc *api.Service) {
	RegisterRoutesWithHub(r, svc, nil)
}

// RegisterRoutesWithHub registers api REST routes with WebSocket hub for broadcasting.
// If hub is nil, routes work without WebSocket broadcasting (backward compatible).
func RegisterRoutesWithHub(r *gin.Engine, svc *api.Service, hub *Hub) {
	g := r.Group("/todos")
	g.GET("", func(c *gin.Context) { handleList(c, svc) })
	g.POST("", func(c *gin.Context) { handleCreateWithBroadcast(c, svc, hub) })
	g.GET(":id", func(c *gin.Context) { handleGet(c, svc) })
	g.PUT(":id", func(c *gin.Context) { handleUpdateWithBroadcast(c, svc, hub) })
	g.DELETE(":id", func(c *gin.Context) { handleDeleteWithBroadcast(c, svc, hub) })
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
func handleList(c *gin.Context, svc *api.Service) {
	q := c.Request.URL.Query()
	opts := cache.ListOptions{SortBy: q.Get("sort_by"), SortOrder: q.Get("order")}
	if s := q.Get("status"); s != "" {
		opts.Status = model.Status(s)
	}
	items, _ := svc.List(opts)
	c.JSON(http.StatusOK, items)
}

// @Summary Create api
// @Description Create a new api
// @Tags todos
// @Accept json
// @Produce json
// @Param api body Todo true "Todo to create"
// @Success 201 {object} Todo
// @Router /todos [post]
func handleCreate(c *gin.Context, svc *api.Service) {
	handleCreateWithBroadcast(c, svc, nil)
}

func handleCreateWithBroadcast(c *gin.Context, svc *api.Service, hub *Hub) {
	var t model.Todo
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

	// Broadcast create event if hub is available
	if hub != nil {
		hub.BroadcastCreate(&t)
	}

	c.JSON(http.StatusCreated, t)
}

// @Summary Get api
// @Description Get a api by ID
// @Tags todos
// @Accept json
// @Produce json
// @Param id path int true "Todo ID"
// @Success 200 {object} Todo
// @Failure 404 {object} map[string]string
// @Router /todos/{id} [get]
func handleGet(c *gin.Context, svc *api.Service) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	t, err := svc.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, t)
}

// @Summary Update api
// @Description Update a api by ID
// @Tags todos
// @Accept json
// @Produce json
// @Param id path int true "Todo ID"
// @Param api body Todo true "Todo to update"
// @Success 200 {object} Todo
// @Failure 400 {object} map[string]string
// @Router /todos/{id} [put]
func handleUpdate(c *gin.Context, svc *api.Service) {
	handleUpdateWithBroadcast(c, svc, nil)
}

func handleUpdateWithBroadcast(c *gin.Context, svc *api.Service, hub *Hub) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var t model.Todo
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}
	t.ID = id
	if err := svc.Update(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get updated api to broadcast complete state and return accurate data
	updated, err := svc.Get(id)
	if err == nil {
		if hub != nil {
			hub.BroadcastUpdate(updated)
		}
		c.JSON(http.StatusOK, updated)
	} else {
		c.JSON(http.StatusOK, t)
	}
}

// @Summary Delete api
// @Description Delete a api by ID
// @Tags todos
// @Param id path int true "Todo ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /todos/{id} [delete]
func handleDelete(c *gin.Context, svc *api.Service) {
	handleDeleteWithBroadcast(c, svc, nil)
}

func handleDeleteWithBroadcast(c *gin.Context, svc *api.Service, hub *Hub) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := svc.Delete(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Broadcast delete event if hub is available
	if hub != nil {
		hub.BroadcastDelete(id)
	}

	c.Status(http.StatusNoContent)
}
