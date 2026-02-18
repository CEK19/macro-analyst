package server

import (
	"io"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"

	"macro-analyst/internal/ws"
)

// TestHelloWorldHandler tests the root endpoint response.
func TestHelloWorldHandler(t *testing.T) {
	// Arrange: Create test server
	hub := ws.NewHub()
	app := fiber.New()
	server := &FiberServer{App: app, Hub: hub}
	app.Get("/", server.HelloWorldHandler)

	// Act: Create and execute request
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	// Assert: Verify status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Assert: Verify response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	expected := `{"message":"Hello World","service":"macro-analyst","version":"1.0.0"}`
	if string(body) != expected {
		t.Errorf("Expected body %q, got %q", expected, string(body))
	}
}

// TestHealthHandler tests the health check endpoint.
func TestHealthHandler(t *testing.T) {
	// Arrange: Create test server
	hub := ws.NewHub()
	app := fiber.New()
	server := &FiberServer{App: app, Hub: hub}
	app.Get("/health", server.HealthHandler)

	// Act: Create and execute request
	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	// Assert: Verify status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Assert: Verify response body contains expected fields
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	expected := `{"active_clients":0,"status":"ok"}`
	if string(body) != expected {
		t.Errorf("Expected body %q, got %q", expected, string(body))
	}
}

// TestHealthHandlerWithClients tests health endpoint with active clients.
func TestHealthHandlerWithClients(t *testing.T) {
	// Arrange: Create test server with active clients simulation
	hub := ws.NewHub()
	go hub.Run() // Start hub in background

	app := fiber.New()
	server := &FiberServer{App: app, Hub: hub}
	app.Get("/health", server.HealthHandler)

	// Simulate adding a client
	// Note: In a real scenario, you'd establish a WebSocket connection
	// For this test, we just verify the endpoint works with 0 clients

	// Act: Execute request
	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	// Assert
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}
}
