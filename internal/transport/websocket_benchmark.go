package transport

import (
	"testing"
)

func BenchmarkHubBroadcast_1Client(b *testing.B) {
	hub := NewHub()
	go hub.Run()
	client := &Client{hub: hub, send: make(chan []byte, 100000)}
	hub.register <- client
	defer func() { hub.unregister <- client }()

	message := []byte(`{"action":"create","todo":{"id":1,"name":"test"}}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hub.broadcast <- message
		<-client.send // Drain to keep buffer free
	}
}

func BenchmarkHubBroadcast_100Clients(b *testing.B) {
	hub := NewHub()
	go hub.Run()

	const numClients = 100
	clients := make([]*Client, numClients)
	for i := range clients {
		c := &Client{hub: hub, send: make(chan []byte, 100000)}
		clients[i] = c
		hub.register <- c
	}
	defer func() {
		for _, c := range clients {
			hub.unregister <- c
		}
	}()

	message := []byte(`{"action":"update","todo":{"id":1,"status":"completed"}}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hub.broadcast <- message
		// Drain all (non-blocking)
		for _, c := range clients {
			select {
			case <-c.send:
			default:
			}
		}
	}
}

func BenchmarkHubBroadcast_1000Clients(b *testing.B) {
	// Similar to above, but with 1000 clients — tests scalability
	// (Large buffer to prevent blocking)
	hub := NewHub()
	go hub.Run()

	const numClients = 1000
	clients := make([]*Client, numClients)
	for i := range clients {
		c := &Client{hub: hub, send: make(chan []byte, 100000)}
		clients[i] = c
		hub.register <- c
	}
	defer func() {
		for _, c := range clients {
			hub.unregister <- c
		}
	}()

	message := []byte(`{"action":"delete","id":42}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hub.broadcast <- message
		// Drain all client channels (simplified — in real bench, drain in parallel if needed)
		for _, c := range clients {
			select {
			case <-c.send:
			default:
			}
		}
	}
}