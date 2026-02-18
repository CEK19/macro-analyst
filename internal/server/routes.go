package server

import (
	"log"

	"macro-analyst/internal/ws"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

const (
	// ClientSendBufferSize is the buffer size for client send channels
	ClientSendBufferSize = 256
)

// RegisterFiberRoutes registers all HTTP and WebSocket routes for the application.
func (s *FiberServer) RegisterFiberRoutes() {
	s.setupMiddleware()
	s.setupHTTPRoutes()
	s.setupWebSocketRoutes()
}

// setupMiddleware configures global middleware for the application.
func (s *FiberServer) setupMiddleware() {
	// CORS middleware for cross-origin requests
	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type",
		AllowCredentials: false,
		MaxAge:           300,
	}))
}

// setupHTTPRoutes registers all HTTP routes.
func (s *FiberServer) setupHTTPRoutes() {
	s.App.Get("/", s.HelloWorldHandler)
	s.App.Get("/health", s.HealthHandler)
}

// setupWebSocketRoutes registers all WebSocket routes.
func (s *FiberServer) setupWebSocketRoutes() {
	// WebSocket upgrade endpoint for real-time price updates
	s.App.Get("/ws/prices", websocket.New(s.handleWebSocket))
}

// handleWebSocket handles WebSocket connections for real-time price streaming.
func (s *FiberServer) handleWebSocket(c *websocket.Conn) {
	// Create a new client for this connection
	client := &ws.Client{
		Hub:  s.Hub,
		Conn: c,
		Send: make(chan []byte, ClientSendBufferSize),
	}

	// Register the client with the Hub
	s.Hub.Register() <- client

	// Ensure cleanup on connection close
	defer func() {
		s.Hub.Unregister() <- client
		client.Close()
	}()

	// Start the write pump in a goroutine to send messages to the client
	go client.WritePump()

	// Keep the connection open and read messages from the client
	// This allows clients to send commands (e.g., subscribe to specific symbols)
	s.readLoop(c)
}

// readLoop continuously reads messages from the WebSocket connection.
// This keeps the connection alive and allows clients to send commands.
func (s *FiberServer) readLoop(c *websocket.Conn) {
	for {
		messageType, message, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket unexpected close error: %v", err)
			}
			break
		}

		// Handle client messages (e.g., subscription requests)
		// TODO: Implement message handling for client commands
		log.Printf("Received message type %d: %s", messageType, string(message))
	}
}

// HelloWorldHandler handles the root endpoint.
func (s *FiberServer) HelloWorldHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Hello World",
		"service": "macro-analyst",
		"version": "1.0.0",
	})
}

// HealthHandler handles the health check endpoint.
// Returns server status and the number of active WebSocket clients.
func (s *FiberServer) HealthHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":         "ok",
		"active_clients": s.Hub.GetClientCount(),
	})
}
