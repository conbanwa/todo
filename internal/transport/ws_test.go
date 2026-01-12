package transport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/conbanwa/todo/internal/dao/cache/api"
	"github.com/conbanwa/todo/internal/dao/db"
	"github.com/conbanwa/todo/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func TestFullWorkflow_RESTAndWebSocket(t *testing.T) {
	// Setup test database
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "integration_test.db")
	store, err := db.NewSQLiteStore(dbPath)
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}
	defer store.Close()

	svc := api.NewService(store)
	hub := NewHub()
	defer hub.Close()
	go hub.Run()

	gin.SetMode(gin.TestMode)
	r := gin.New()
	RegisterRoutesWithHub(r, svc, hub)
	r.GET("/ws", func(c *gin.Context) {
		HandleWebSocket(c, hub)
	})

	s := httptest.NewServer(r)
	defer s.Close()

	wsURL := "ws" + s.URL[4:] + "/ws"

	// Connect WebSocket client
	wsConn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("failed to connect WebSocket: %v", err)
	}
	defer wsConn.Close()

	time.Sleep(50 * time.Millisecond)

	// Step 1: Create api via REST API
	todo1 := model.Todo{Name: "Integration Test 1", Status: model.NotStarted, Description: "Test description"}
	body1, _ := json.Marshal(todo1)

	resp1, err := http.Post(s.URL+"/todos", "application/json", bytes.NewReader(body1))
	if err != nil {
		t.Fatalf("failed to create api: %v", err)
	}
	defer resp1.Body.Close()

	if resp1.StatusCode != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", resp1.StatusCode)
	}

	var created1 model.Todo
	json.NewDecoder(resp1.Body).Decode(&created1)

	// Step 2: Verify WebSocket broadcast received
	wsConn.SetReadDeadline(time.Now().Add(1 * time.Second))
	var wsMsg1 WSMessage
	if err := wsConn.ReadJSON(&wsMsg1); err != nil {
		t.Fatalf("failed to read WebSocket message: %v", err)
	}

	if wsMsg1.Type != "create" {
		t.Errorf("expected type 'create', got %q", wsMsg1.Type)
	}
	if wsMsg1.Payload.ID != created1.ID {
		t.Errorf("expected ID %d, got %d", created1.ID, wsMsg1.Payload.ID)
	}

	// Step 3: Update api via REST API
	todo2 := model.Todo{Name: "Updated Integration Test", Status: model.InProgress}
	body2, _ := json.Marshal(todo2)

	req2, _ := http.NewRequest(http.MethodPut, s.URL+"/todos/"+int64ToString(created1.ID), bytes.NewReader(body2))
	req2.Header.Set("Content-Type", "application/json")
	resp2, err := http.DefaultClient.Do(req2)
	if err != nil {
		t.Fatalf("failed to update api: %v", err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp2.StatusCode)
	}

	// Step 4: Verify WebSocket broadcast for update
	wsConn.SetReadDeadline(time.Now().Add(1 * time.Second))
	var wsMsg2 WSMessage
	if err := wsConn.ReadJSON(&wsMsg2); err != nil {
		t.Fatalf("failed to read update WebSocket message: %v", err)
	}

	if wsMsg2.Type != "update" {
		t.Errorf("expected type 'update', got %q", wsMsg2.Type)
	}
	if wsMsg2.Payload.ID != created1.ID {
		t.Errorf("expected ID %d, got %d", created1.ID, wsMsg2.Payload.ID)
	}

	// Step 5: Delete api via REST API
	req3, _ := http.NewRequest(http.MethodDelete, s.URL+"/todos/"+int64ToString(created1.ID), nil)
	resp3, err := http.DefaultClient.Do(req3)
	if err != nil {
		t.Fatalf("failed to delete api: %v", err)
	}
	defer resp3.Body.Close()

	if resp3.StatusCode != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", resp3.StatusCode)
	}

	// Step 6: Verify WebSocket broadcast for delete
	wsConn.SetReadDeadline(time.Now().Add(1 * time.Second))
	var wsMsg3 WSMessage
	if err := wsConn.ReadJSON(&wsMsg3); err != nil {
		t.Fatalf("failed to read delete WebSocket message: %v", err)
	}

	if wsMsg3.Type != "delete" {
		t.Errorf("expected type 'delete', got %q", wsMsg3.Type)
	}
	if wsMsg3.Payload.ID != created1.ID {
		t.Errorf("expected ID %d, got %d", created1.ID, wsMsg3.Payload.ID)
	}
}

func TestMultipleClients_RealTimeSync(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "sync_test.db")
	store, err := db.NewSQLiteStore(dbPath)
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}
	defer store.Close()

	svc := api.NewService(store)
	hub := NewHub()
	defer hub.Close()
	go hub.Run()

	gin.SetMode(gin.TestMode)
	r := gin.New()
	RegisterRoutesWithHub(r, svc, hub)
	r.GET("/ws", func(c *gin.Context) {
		HandleWebSocket(c, hub)
	})

	s := httptest.NewServer(r)
	defer s.Close()

	wsURL := "ws" + s.URL[4:] + "/ws"

	// Connect multiple WebSocket clients
	const numClients = 3
	var conns []*websocket.Conn

	for i := 0; i < numClients; i++ {
		conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		if err != nil {
			t.Fatalf("failed to connect client %d: %v", i+1, err)
		}
		conns = append(conns, conn)
		defer conn.Close()
	}

	time.Sleep(100 * time.Millisecond)

	// Create api via REST API
	todo := model.Todo{Name: "Multi-Client Sync Test", Status: model.NotStarted}
	body, _ := json.Marshal(todo)

	resp, err := http.Post(s.URL+"/todos", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("failed to create api: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", resp.StatusCode)
	}

	// All clients should receive the broadcast
	for i, conn := range conns {
		conn.SetReadDeadline(time.Now().Add(1 * time.Second))
		var msg WSMessage
		if err := conn.ReadJSON(&msg); err != nil {
			t.Errorf("client %d failed to read message: %v", i+1, err)
			continue
		}

		if msg.Type != "create" {
			t.Errorf("client %d: expected type 'create', got %q", i+1, msg.Type)
		}
		if msg.Payload.Name != "Multi-Client Sync Test" {
			t.Errorf("client %d: expected name 'Multi-Client Sync Test', got %q", i+1, msg.Payload.Name)
		}
	}
}

func TestWebSocketReconnect(t *testing.T) {
	hub := NewHub()
	defer hub.Close()
	go hub.Run()

	svc := api.NewService(&mockStore{todos: make(map[int64]*model.Todo)})
	gin.SetMode(gin.TestMode)
	r := gin.New()
	RegisterRoutesWithHub(r, svc, hub)
	r.GET("/ws", func(c *gin.Context) {
		HandleWebSocket(c, hub)
	})

	s := httptest.NewServer(r)
	defer s.Close()

	wsURL := "ws" + s.URL[4:] + "/ws"

	// Connect first time
	conn1, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}

	time.Sleep(50 * time.Millisecond)

	// Verify client registered
	hub.mu.Lock()
	count1 := len(hub.clients)
	hub.mu.Unlock()

	if count1 != 1 {
		t.Errorf("expected 1 client, got %d", count1)
	}

	// Close connection
	conn1.Close()
	time.Sleep(100 * time.Millisecond)

	// Verify client unregistered
	hub.mu.Lock()
	count2 := len(hub.clients)
	hub.mu.Unlock()

	if count2 != 0 {
		t.Errorf("expected 0 clients after close, got %d", count2)
	}

	// Reconnect
	conn2, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("failed to reconnect: %v", err)
	}
	defer conn2.Close()

	time.Sleep(50 * time.Millisecond)

	// Verify reconnected client registered
	hub.mu.Lock()
	count3 := len(hub.clients)
	hub.mu.Unlock()

	if count3 != 1 {
		t.Errorf("expected 1 client after reconnect, got %d", count3)
	}
}

func TestHTMLFrontendIntegration(t *testing.T) {
	tmpDir := t.TempDir()
	staticDir := filepath.Join(tmpDir, "static")
	os.MkdirAll(staticDir, 0755)

	indexPath := filepath.Join(staticDir, "index.html")
	testHTML := `<!DOCTYPE html><html><head><title>Todo App</title></head><body><h1>Todo List</h1></body></html>`
	os.WriteFile(indexPath, []byte(testHTML), 0644)

	store, err := db.NewSQLiteStore(filepath.Join(tmpDir, "test.db"))
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}
	defer store.Close()

	svc := api.NewService(store)
	hub := NewHub()
	defer hub.Close()
	go hub.Run()

	gin.SetMode(gin.TestMode)
	r := gin.New()
	RegisterRoutesWithHub(r, svc, hub)
	r.GET("/ws", func(c *gin.Context) {
		HandleWebSocket(c, hub)
	})
	r.Static("/static", staticDir)
	r.GET("/", func(c *gin.Context) {
		c.File(indexPath)
	})

	s := httptest.NewServer(r)
	defer s.Close()

	// Test root route serves HTML
	resp, err := http.Get(s.URL + "/")
	if err != nil {
		t.Fatalf("failed to fetch root: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if string(body) != testHTML {
		t.Errorf("expected HTML content, got different content")
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "text/html; charset=utf-8" && contentType != "text/html" {
		t.Errorf("expected Content-Type 'text/html', got %q", contentType)
	}

	// Test static file route
	resp2, err := http.Get(s.URL + "/static/index.html")
	if err != nil {
		t.Fatalf("failed to fetch static file: %v", err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 for static file, got %d", resp2.StatusCode)
	}

	// Test WebSocket endpoint
	wsURL := "ws" + s.URL[4:] + "/ws"
	conn, resp3, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("failed to connect WebSocket: %v", err)
	}
	defer conn.Close()

	if resp3.StatusCode != http.StatusSwitchingProtocols {
		t.Errorf("expected status 101, got %d", resp3.StatusCode)
	}

	// Test REST API still works
	resp4, err := http.Get(s.URL + "/todos")
	if err != nil {
		t.Fatalf("failed to fetch todos: %v", err)
	}
	defer resp4.Body.Close()

	if resp4.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 for /todos, got %d", resp4.StatusCode)
	}
}

// Helper function
func int64ToString(n int64) string {
	return fmt.Sprintf("%d", n)
}
