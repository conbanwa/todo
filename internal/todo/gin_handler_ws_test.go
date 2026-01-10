package todo

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func TestGinHandler_WebSocketIntegration_Create(t *testing.T) {
	hub := NewHub()
	defer hub.Close()
	go hub.Run()

	svc := NewService(&mockStore{todos: make(map[int64]*Todo)})
	gin.SetMode(gin.TestMode)
	r := gin.New()
	RegisterRoutesWithHub(r, svc, hub)

	// Setup WebSocket route
	r.GET("/ws", func(c *gin.Context) {
		HandleWebSocket(c, hub)
	})

	s := httptest.NewServer(r)
	defer s.Close()

	// Connect WebSocket client
	wsURL := "ws" + s.URL[4:] + "/ws"
	wsConn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("failed to connect WebSocket: %v", err)
	}
	defer wsConn.Close()

	time.Sleep(50 * time.Millisecond)

	// Create todo via REST API
	todo := Todo{Name: "WebSocket Test", Description: "Testing broadcast"}
	body, _ := json.Marshal(todo)

	req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d: %s", w.Code, w.Body.String())
	}

	// Read WebSocket message
	wsConn.SetReadDeadline(time.Now().Add(1 * time.Second))
	var wsMsg WSMessage
	if err := wsConn.ReadJSON(&wsMsg); err != nil {
		t.Fatalf("failed to read WebSocket message: %v", err)
	}

	if wsMsg.Type != "create" {
		t.Errorf("expected type 'create', got %q", wsMsg.Type)
	}
	if wsMsg.Payload.Name != "WebSocket Test" {
		t.Errorf("expected name 'WebSocket Test', got %q", wsMsg.Payload.Name)
	}
}

func TestGinHandler_WebSocketIntegration_Update(t *testing.T) {
	hub := NewHub()
	defer hub.Close()
	go hub.Run()

	store := &mockStore{todos: make(map[int64]*Todo)}
	store.todos[1] = &Todo{ID: 1, Name: "Original", Status: NotStarted}
	svc := NewService(store)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	RegisterRoutesWithHub(r, svc, hub)

	r.GET("/ws", func(c *gin.Context) {
		HandleWebSocket(c, hub)
	})

	s := httptest.NewServer(r)
	defer s.Close()

	wsURL := "ws" + s.URL[4:] + "/ws"
	wsConn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("failed to connect WebSocket: %v", err)
	}
	defer wsConn.Close()

	time.Sleep(50 * time.Millisecond)

	// Update todo via REST API
	update := Todo{Name: "Updated", Status: Completed}
	body, _ := json.Marshal(update)

	req := httptest.NewRequest(http.MethodPut, "/todos/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	// Read WebSocket message
	wsConn.SetReadDeadline(time.Now().Add(1 * time.Second))
	var wsMsg WSMessage
	if err := wsConn.ReadJSON(&wsMsg); err != nil {
		t.Fatalf("failed to read WebSocket message: %v", err)
	}

	if wsMsg.Type != "update" {
		t.Errorf("expected type 'update', got %q", wsMsg.Type)
	}
	if wsMsg.Payload.ID != 1 {
		t.Errorf("expected ID 1, got %d", wsMsg.Payload.ID)
	}
	if wsMsg.Payload.Name != "Updated" {
		t.Errorf("expected name 'Updated', got %q", wsMsg.Payload.Name)
	}
}

func TestGinHandler_WebSocketIntegration_Delete(t *testing.T) {
	hub := NewHub()
	defer hub.Close()
	go hub.Run()

	store := &mockStore{todos: make(map[int64]*Todo)}
	store.todos[1] = &Todo{ID: 1, Name: "To Delete"}
	svc := NewService(store)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	RegisterRoutesWithHub(r, svc, hub)

	r.GET("/ws", func(c *gin.Context) {
		HandleWebSocket(c, hub)
	})

	s := httptest.NewServer(r)
	defer s.Close()

	wsURL := "ws" + s.URL[4:] + "/ws"
	wsConn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("failed to connect WebSocket: %v", err)
	}
	defer wsConn.Close()

	time.Sleep(50 * time.Millisecond)

	// Delete todo via REST API
	req := httptest.NewRequest(http.MethodDelete, "/todos/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", w.Code)
	}

	// Read WebSocket message
	wsConn.SetReadDeadline(time.Now().Add(1 * time.Second))
	var wsMsg WSMessage
	if err := wsConn.ReadJSON(&wsMsg); err != nil {
		t.Fatalf("failed to read WebSocket message: %v", err)
	}

	if wsMsg.Type != "delete" {
		t.Errorf("expected type 'delete', got %q", wsMsg.Type)
	}
	if wsMsg.Payload.ID != 1 {
		t.Errorf("expected ID 1, got %d", wsMsg.Payload.ID)
	}
}

func TestGinHandler_WebSocketIntegration_MultipleClients(t *testing.T) {
	hub := NewHub()
	defer hub.Close()
	go hub.Run()

	svc := NewService(&mockStore{todos: make(map[int64]*Todo)})
	gin.SetMode(gin.TestMode)
	r := gin.New()
	RegisterRoutesWithHub(r, svc, hub)

	r.GET("/ws", func(c *gin.Context) {
		HandleWebSocket(c, hub)
	})

	s := httptest.NewServer(r)
	defer s.Close()

	wsURL := "ws" + s.URL[4:] + "/ws"

	// Connect two WebSocket clients
	wsConn1, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("failed to connect WebSocket client 1: %v", err)
	}
	defer wsConn1.Close()

	wsConn2, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("failed to connect WebSocket client 2: %v", err)
	}
	defer wsConn2.Close()

	time.Sleep(50 * time.Millisecond)

	// Create todo via REST API
	todo := Todo{Name: "Multi-Client Test"}
	body, _ := json.Marshal(todo)

	req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}

	// Both clients should receive the broadcast
	wsConn1.SetReadDeadline(time.Now().Add(1 * time.Second))
	wsConn2.SetReadDeadline(time.Now().Add(1 * time.Second))

	var msg1, msg2 WSMessage
	if err := wsConn1.ReadJSON(&msg1); err != nil {
		t.Fatalf("client 1 failed to read: %v", err)
	}
	if err := wsConn2.ReadJSON(&msg2); err != nil {
		t.Fatalf("client 2 failed to read: %v", err)
	}

	if msg1.Type != "create" || msg2.Type != "create" {
		t.Errorf("both clients should receive 'create' message: client1=%q, client2=%q", msg1.Type, msg2.Type)
	}
	if msg1.Payload.Name != "Multi-Client Test" || msg2.Payload.Name != "Multi-Client Test" {
		t.Errorf("both clients should receive correct payload")
	}
}
