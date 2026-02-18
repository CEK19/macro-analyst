package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"macro-analyst/internal/server"
	"macro-analyst/internal/ws"
)

const (
	// DefaultPort is used if PORT environment variable is not set
	DefaultPort = 8080

	// ShutdownTimeout is the maximum time to wait for graceful shutdown
	ShutdownTimeout = 5 * time.Second
)

func main() {
	// Initialize the WebSocket Hub
	hub := ws.NewHub()
	go hub.Run()
	log.Println("WebSocket Hub started")

	// Initialize the Price Ingestor with custom throttle interval
	ingestor := ws.NewIngestor(hub,
		ws.WithThrottleInterval(500*time.Millisecond),
	)

	// Start the ingestor - connects to Binance WebSocket
	go ingestor.Start()
	log.Println("Price Ingestor started - connecting to Binance for real-time data")

	// Initialize the HTTP/WebSocket server
	srv := server.New(hub)
	srv.RegisterFiberRoutes()

	// Start the server in a goroutine
	port := getPort()
	go startServer(srv, port)

	// Wait for shutdown signal and perform graceful shutdown
	waitForShutdown(srv, ingestor)
}

// getPort retrieves the port number from environment variable or returns default.
func getPort() int {
	portStr := os.Getenv("PORT")
	if portStr == "" {
		return DefaultPort
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Printf("Invalid PORT value '%s', using default %d", portStr, DefaultPort)
		return DefaultPort
	}

	return port
}

// startServer starts the HTTP/WebSocket server on the specified port.
func startServer(srv *server.FiberServer, port int) {
	log.Printf("Server starting on port %d", port)
	log.Printf("WebSocket endpoint: ws://localhost:%d/ws/prices", port)
	log.Printf("Health check: http://localhost:%d/health", port)

	addr := fmt.Sprintf(":%d", port)
	if err := srv.Listen(addr); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

// waitForShutdown blocks until an interrupt signal is received,
// then performs a graceful shutdown of the server.
func waitForShutdown(srv *server.FiberServer, ingestor *ws.Ingestor) {
	// Create a channel to listen for interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive a signal
	sig := <-quit
	log.Printf("Received signal %v, shutting down gracefully...", sig)

	// Stop the ingestor first
	if ingestor != nil {
		ingestor.Stop()
	}

	// Create a context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.ShutdownWithContext(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	} else {
		log.Println("Server shutdown completed successfully")
	}
}
