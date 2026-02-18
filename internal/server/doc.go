// Package server provides HTTP and WebSocket server functionality for the macro-analyst API.
//
// The server package integrates Fiber web framework with WebSocket support for
// real-time market data streaming.
//
// # Architecture
//
// The FiberServer embeds a Fiber application and maintains a reference to a
// WebSocket Hub for message broadcasting:
//
//	┌──────────────────┐
//	│  FiberServer     │
//	│                  │
//	│  ┌────────────┐  │
//	│  │ Fiber App  │  │  ← HTTP routes and middleware
//	│  └────────────┘  │
//	│                  │
//	│  ┌────────────┐  │
//	│  │  WS Hub    │  │  ← WebSocket message broker
//	│  └────────────┘  │
//	└──────────────────┘
//
// # Endpoints
//
// HTTP Endpoints:
//   - GET /        - Hello World (API info)
//   - GET /health  - Health check with active client count
//
// WebSocket Endpoints:
//   - GET /ws/prices - Real-time price updates
//
// # Usage
//
// Basic server setup:
//
//	hub := ws.NewHub()
//	go hub.Run()
//
//	srv := server.New(hub)
//	srv.RegisterFiberRoutes()
//
//	port := 8080
//	log.Fatal(srv.Listen(fmt.Sprintf(":%d", port)))
//
// # Custom Configuration
//
// The server can be configured with custom settings:
//
//	srv := server.New(hub, server.Config{
//	    ServerHeader: "my-app",
//	    AppName:      "my-app-v1",
//	})
//
// # Middleware
//
// The server includes CORS middleware by default, configured to:
//   - Allow all origins (*)
//   - Support common HTTP methods
//   - Cache preflight requests for 5 minutes
//
// # WebSocket Handling
//
// The WebSocket endpoint handles:
//   - Connection upgrade and client registration
//   - Message broadcasting via WritePump goroutine
//   - Graceful connection cleanup on close
//   - Error handling for unexpected disconnects
//
// # Thread Safety
//
// All routes are safe for concurrent access. The WebSocket handler
// properly manages goroutines and channel cleanup.
//
// # Testing
//
// The package includes comprehensive tests for all endpoints.
// See routes_test.go for examples.
package server
