package ws

import (
	"testing"
	"time"
)

// TestNewHub verifies Hub initialization.
func TestNewHub(t *testing.T) {
	hub := NewHub()

	if hub == nil {
		t.Fatal("NewHub returned nil")
	}

	if hub.clients == nil {
		t.Error("Hub clients map not initialized")
	}

	if hub.broadcast == nil {
		t.Error("Hub broadcast channel not initialized")
	}

	if hub.register == nil {
		t.Error("Hub register channel not initialized")
	}

	if hub.unregister == nil {
		t.Error("Hub unregister channel not initialized")
	}
}

// TestHubGetClientCount verifies client count tracking.
func TestHubGetClientCount(t *testing.T) {
	hub := NewHub()

	// Initially should be 0
	if count := hub.GetClientCount(); count != 0 {
		t.Errorf("Expected 0 clients, got %d", count)
	}

	// Start hub in background
	go hub.Run()

	// Give hub time to start
	time.Sleep(10 * time.Millisecond)

	// Verify count is still 0
	if count := hub.GetClientCount(); count != 0 {
		t.Errorf("Expected 0 clients after Run(), got %d", count)
	}
}

// TestHubChannelAccess verifies channel accessor methods.
func TestHubChannelAccess(t *testing.T) {
	hub := NewHub()

	if hub.Broadcast() == nil {
		t.Error("Broadcast() returned nil channel")
	}

	if hub.Register() == nil {
		t.Error("Register() returned nil channel")
	}

	if hub.Unregister() == nil {
		t.Error("Unregister() returned nil channel")
	}
}

// TestHubConcurrentAccess tests thread-safety of GetClientCount.
func TestHubConcurrentAccess(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	// Give hub time to start
	time.Sleep(10 * time.Millisecond)

	// Spawn multiple goroutines reading client count
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				hub.GetClientCount()
			}
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

// TestHubRegisterClient verifies client registration.
func TestHubRegisterClient(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	time.Sleep(10 * time.Millisecond)

	client := &Client{
		Hub:  hub,
		Send: make(chan []byte, 256),
	}

	// Register client
	hub.register <- client
	time.Sleep(10 * time.Millisecond)

	if count := hub.GetClientCount(); count != 1 {
		t.Errorf("Expected 1 client after registration, got %d", count)
	}
}

// TestHubUnregisterClient verifies client unregistration.
func TestHubUnregisterClient(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	time.Sleep(10 * time.Millisecond)

	client := &Client{
		Hub:  hub,
		Send: make(chan []byte, 256),
	}

	// Register then unregister
	hub.register <- client
	time.Sleep(10 * time.Millisecond)

	hub.unregister <- client
	time.Sleep(10 * time.Millisecond)

	if count := hub.GetClientCount(); count != 0 {
		t.Errorf("Expected 0 clients after unregistration, got %d", count)
	}

	// Verify send channel is closed
	_, ok := <-client.Send
	if ok {
		t.Error("Client send channel should be closed after unregistration")
	}
}

// TestHubBroadcast verifies message broadcasting to clients.
func TestHubBroadcast(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	time.Sleep(10 * time.Millisecond)

	client := &Client{
		Hub:  hub,
		Send: make(chan []byte, 256),
	}

	hub.register <- client
	time.Sleep(10 * time.Millisecond)

	// Broadcast a message
	testMessage := []byte("test message")
	hub.broadcast <- testMessage

	// Wait for message to be delivered
	select {
	case msg := <-client.Send:
		if string(msg) != string(testMessage) {
			t.Errorf("Expected message %s, got %s", testMessage, msg)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Timeout waiting for broadcast message")
	}
}

// TestHubBroadcastToMultipleClients verifies broadcasting to multiple clients.
func TestHubBroadcastToMultipleClients(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	time.Sleep(10 * time.Millisecond)

	clients := []*Client{
		{Hub: hub, Send: make(chan []byte, 256)},
		{Hub: hub, Send: make(chan []byte, 256)},
		{Hub: hub, Send: make(chan []byte, 256)},
	}

	// Register all clients
	for _, client := range clients {
		hub.register <- client
	}
	time.Sleep(10 * time.Millisecond)

	if count := hub.GetClientCount(); count != 3 {
		t.Errorf("Expected 3 clients, got %d", count)
	}

	// Broadcast a message
	testMessage := []byte("broadcast to all")
	hub.broadcast <- testMessage

	// Verify all clients received the message
	for i, client := range clients {
		select {
		case msg := <-client.Send:
			if string(msg) != string(testMessage) {
				t.Errorf("Client %d: expected message %s, got %s", i, testMessage, msg)
			}
		case <-time.After(100 * time.Millisecond):
			t.Errorf("Client %d: timeout waiting for broadcast message", i)
		}
	}
}

// TestHubUnregisterNonExistentClient verifies handling of non-existent client.
func TestHubUnregisterNonExistentClient(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	time.Sleep(10 * time.Millisecond)

	client := &Client{
		Hub:  hub,
		Send: make(chan []byte, 256),
	}

	// Try to unregister without registering first (should not panic)
	hub.unregister <- client
	time.Sleep(10 * time.Millisecond)

	if count := hub.GetClientCount(); count != 0 {
		t.Errorf("Expected 0 clients, got %d", count)
	}
}

// TestHubBroadcastWithFullChannel verifies behavior when client channel is full.
func TestHubBroadcastWithFullChannel(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	time.Sleep(10 * time.Millisecond)

	// Create a client with a small buffer
	client := &Client{
		Hub:  hub,
		Send: make(chan []byte, 1),
	}

	hub.register <- client
	time.Sleep(10 * time.Millisecond)

	// Fill the channel
	client.Send <- []byte("filling")

	// Try to broadcast (should trigger removal of client)
	hub.broadcast <- []byte("test")
	time.Sleep(50 * time.Millisecond)

	// Client should be removed
	if count := hub.GetClientCount(); count != 0 {
		t.Errorf("Expected 0 clients after channel full, got %d", count)
	}
}

// TestMultipleRegisterSameClient verifies registering the same client twice.
func TestMultipleRegisterSameClient(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	time.Sleep(10 * time.Millisecond)

	client := &Client{
		Hub:  hub,
		Send: make(chan []byte, 256),
	}

	// Register same client twice
	hub.register <- client
	hub.register <- client
	time.Sleep(10 * time.Millisecond)

	// Should still only count as 1 client (map behavior)
	if count := hub.GetClientCount(); count != 1 {
		t.Errorf("Expected 1 client (same client registered twice), got %d", count)
	}
}

