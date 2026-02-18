package ws

import (
	"log"

	"github.com/gofiber/contrib/websocket"
)

// Client represents a single WebSocket connection from a client.
// It holds the connection, a reference to the Hub, and a buffered send channel.
type Client struct {
	// Hub is the central message broker this client is connected to
	Hub *Hub

	// Conn is the underlying WebSocket connection
	Conn *websocket.Conn

	// Send is a buffered channel of outbound messages
	Send chan []byte
}

// WritePump pumps messages from the Hub to the WebSocket connection.
// A goroutine running WritePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) WritePump() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		message, ok := <-c.Send
		if !ok {
			// The Hub closed the channel, send close message
			if err := c.Conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
				log.Printf("Error sending close message: %v", err)
			}
			return
		}

		// Write the message to the WebSocket connection
		if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("Error writing message to client: %v", err)
			return
		}
	}
}

// Close gracefully closes the client connection and cleans up resources.
func (c *Client) Close() {
	if c.Conn != nil {
		c.Conn.Close()
	}
}
