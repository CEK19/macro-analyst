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
