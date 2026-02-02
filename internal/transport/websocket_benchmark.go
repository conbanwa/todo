package transport

import (
	"testing"
)

func BenchmarkHubBroadcast_1Client(b *testing.B) {
	hub := NewHub()
	go hub.Run()
	client := &Client{hub: hub, send: make(chan []byte, 10000)} // Large buffer to avoid blocking
	hub.register <- client
	defer func() { hub.unregister <- client }()

	message := []byte(`{"type":"broadcast","data":"test"}`)

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
	for i := 0; i < numClients; i++ {
		c := &Client{hub: hub, send: make(chan []byte, 10000)}
		clients[i] = c
		hub.register <- c
	}
	defer func() {
		for _, c := range clients {
			hub.unregister <- c
		}
	}()

	message := []byte(`{"type":"broadcast","data":"test"}`)

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