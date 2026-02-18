package ws

import (
	"log"
	"sync"
)

const (
	// BroadcastBufferSize is the buffer size for the broadcast channel
	BroadcastBufferSize = 256
)

// Hub maintains the set of active clients and broadcasts messages to them.
// It acts as the central message broker using Go channels for concurrent communication.
type Hub struct {
	// clients holds all currently connected clients
	clients map[*Client]bool

	// broadcast is the channel for inbound messages from data sources
	broadcast chan []byte

	// register is the channel for requests to register new clients
	register chan *Client

	// unregister is the channel for requests to unregister clients
	unregister chan *Client

	// mu protects concurrent access to the clients map
	mu sync.RWMutex
}

// NewHub creates and initializes a new Hub instance.
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, BroadcastBufferSize),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts the hub's main loop to handle client registration, unregistration,
// and message broadcasting. This should be run in a separate goroutine.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case message := <-h.broadcast:
			h.broadcastMessage(message)
		}
	}
}

// registerClient adds a new client to the hub.
func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	h.clients[client] = true
	clientCount := len(h.clients)
	h.mu.Unlock()

	log.Printf("New client connected! Total active clients: %d", clientCount)
}

// unregisterClient removes a client from the hub and closes its send channel.
func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	if _, exists := h.clients[client]; exists {
		delete(h.clients, client)
		close(client.Send)
		clientCount := len(h.clients)
		h.mu.Unlock()
		log.Printf("Client disconnected! Remaining clients: %d", clientCount)
	} else {
		h.mu.Unlock()
	}
}

// broadcastMessage sends a message to all connected clients.
// If a client's send channel is full, the client is removed.
func (h *Hub) broadcastMessage(message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		select {
		case client.Send <- message:
			// Message sent successfully
		default:
			// Client's send channel is full, likely disconnected
			// Schedule for removal by closing channel
			go func(c *Client) {
				h.unregister <- c
			}(client)
		}
	}
}

// GetClientCount returns the number of currently connected clients.
// This method is safe for concurrent use.
func (h *Hub) GetClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// Broadcast returns the broadcast channel for sending messages to all clients.
// External data sources can write to this channel.
func (h *Hub) Broadcast() chan<- []byte {
	return h.broadcast
}

// Register returns the register channel for adding new clients.
func (h *Hub) Register() chan<- *Client {
	return h.register
}

// Unregister returns the unregister channel for removing clients.
func (h *Hub) Unregister() chan<- *Client {
	return h.unregister
}
