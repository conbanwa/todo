package transport

import (
	"testing"
	"time"

	"github.com/conbanwa/todo/internal/model"
)

func BenchmarkHubBroadcast_1Client(b *testing.B) {
	hub := NewHub()
	go hub.Run()

	client := &Client{
		hub:  hub,
		send: make(chan WSMessage, 256), // match real buffer size
	}
	hub.register <- client

	defer func() { hub.unregister <- client }()

	msg := WSMessage{
		Type: "create",
		Payload: model.Todo{
			Name:        "Benchmark Todo",
			Description: "Performance test",
			Status:      model.NotStarted,
			Priority:    5,
			Tags:        []string{"bench"},
		},
		Timestamp: time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hub.Broadcast(msg)        // Use public method
		select {
		case <-client.send:       // drain
		default:
		}
	}
}

func BenchmarkHubBroadcast_100Clients(b *testing.B) {
	hub := NewHub()
	go hub.Run()

	const numClients = 100
	clients := make([]*Client, numClients)
	for i := range clients {
		c := &Client{
			hub:  hub,
			send: make(chan WSMessage, 256),
		}
		clients[i] = c
		hub.register <- c
	}
	defer func() {
		for _, c := range clients {
			hub.unregister <- c
		}
	}()

	msg := WSMessage{
		Type: "update",
		Payload: model.Todo{
			ID:     1,
			Name:   "Updated Todo",
			Status: model.Completed,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hub.Broadcast(msg)
		for _, c := range clients {
			select {
			case <-c.send:
			default:
			}
		}
	}
}

func BenchmarkHubBroadcast_1000Clients(b *testing.B) {
	hub := NewHub()
	go hub.Run()

	const numClients = 1000
	clients := make([]*Client, numClients)
	for i := range clients {
		c := &Client{
			hub:  hub,
			send: make(chan WSMessage, 256),
		}
		clients[i] = c
		hub.register <- c
	}
	defer func() {
		for _, c := range clients {
			hub.unregister <- c
		}
	}()

	msg := WSMessage{
		Type:    "delete",
		Payload: model.Todo{ID: 42},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hub.Broadcast(msg)
		for _, c := range clients {
			select {
			case <-c.send:
			default:
			}
		}
	}
}