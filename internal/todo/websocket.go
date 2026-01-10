package todo

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// WSMessage represents a WebSocket message
type WSMessage struct {
	Type      string    `json:"type"` // "create", "update", "delete"
	Payload   Todo      `json:"payload"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}

// Client represents a WebSocket connection
type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan WSMessage
	mu     sync.Mutex
	closed bool
}

// Hub maintains the set of active clients and broadcasts messages to them
type Hub struct {
	// Registered clients
	clients map[*Client]bool

	// Inbound messages from clients (for future use)
	broadcast chan WSMessage

	// Register requests from clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	mu sync.RWMutex
}

// NewHub creates a new Hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan WSMessage, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts the hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("Client registered. Total clients: %d", len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
			log.Printf("Client unregistered. Total clients: %d", len(h.clients))

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					// Client send buffer is full, remove client
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// Broadcast sends a message to all connected clients
func (h *Hub) Broadcast(msg WSMessage) {
	if msg.Timestamp.IsZero() {
		msg.Timestamp = time.Now()
	}
	select {
	case h.broadcast <- msg:
	default:
		log.Println("Hub broadcast channel full, dropping message")
	}
}

// Close closes the hub and all client connections
func (h *Hub) Close() {
	h.mu.Lock()
	defer h.mu.Unlock()
	for client := range h.clients {
		client.conn.Close()
		close(client.send)
		delete(h.clients, client)
	}
}

// readPump pumps messages from the websocket connection to the hub
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}
		// For now, we only handle incoming pings/pongs
		// Future: could handle incoming messages from client
	}
}

// writePump pumps messages from the hub to the websocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				// Hub closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteJSON(message); err != nil {
				log.Printf("WebSocket write error: %v", err)
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// HandleWebSocket handles websocket requests from clients
func HandleWebSocket(c *gin.Context, hub *Hub) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &Client{
		hub:  hub,
		conn: conn,
		send: make(chan WSMessage, 256),
	}

	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in new goroutines
	go client.writePump()
	go client.readPump()
}

// BroadcastCreate broadcasts a create event
func (h *Hub) BroadcastCreate(todo *Todo) {
	h.Broadcast(WSMessage{
		Type:    "create",
		Payload: *todo,
	})
}

// BroadcastUpdate broadcasts an update event
func (h *Hub) BroadcastUpdate(todo *Todo) {
	h.Broadcast(WSMessage{
		Type:    "update",
		Payload: *todo,
	})
}

// BroadcastDelete broadcasts a delete event
func (h *Hub) BroadcastDelete(id int64) {
	h.Broadcast(WSMessage{
		Type:    "delete",
		Payload: Todo{ID: id},
	})
}
