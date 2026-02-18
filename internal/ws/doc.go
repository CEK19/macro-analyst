// Package ws provides WebSocket support for real-time market data streaming.
//
// # Architecture
//
// The package implements a Hub-and-Spoke pattern for efficient message broadcasting:
//
//	┌─────────────┐
//	│  Binance    │  ← Real-time data from Binance WebSocket API
//	│  WebSocket  │
//	└──────┬──────┘
//	       │
//	       ▼
//	┌─────────────┐
//	│  Ingestor   │  ← Data ingestion with throttling
//	└──────┬──────┘
//	       │
//	       ▼
//	┌─────────────┐
//	│     Hub     │  ← Central message broker (goroutine-safe)
//	└──────┬──────┘
//	       │
//	       ├───────┬───────┬───────┐
//	       ▼       ▼       ▼       ▼
//	   Client  Client  Client  Client  ← WebSocket connections
//
// # Components
//
// Hub: Central message broker that manages client connections and broadcasts
// messages using Go channels. Thread-safe for concurrent access.
//
// Client: Represents a single WebSocket connection with a send buffer.
// Each client runs a WritePump goroutine to handle outbound messages.
//
// Ingestor: Connects to Binance WebSocket API and streams real-time market data.
// Implements throttling to prevent overwhelming clients with high-frequency updates.
// Uses adshao/go-binance SDK for reliable WebSocket connections with auto-reconnect.
//
// # Usage
//
// Basic setup:
//
//	// Create and start the Hub
//	hub := ws.NewHub()
//	go hub.Run()
//
//	// Create and start the Ingestor (connects to Binance)
//	ingestor := ws.NewIngestor(hub)
//	go ingestor.Start()
//
//	// Register WebSocket endpoint
//	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
//	    client := &ws.Client{
//	        Hub:  hub,
//	        Conn: c,
//	        Send: make(chan []byte, 256),
//	    }
//	    hub.Register() <- client
//	    defer func() {
//	        hub.Unregister() <- client
//	        client.Close()
//	    }()
//	    go client.WritePump()
//	    // ... handle client reads
//	}))
//
// # Customization
//
// The Ingestor supports throttle interval configuration:
//
//	ingestor := ws.NewIngestor(hub,
//	    ws.WithThrottleInterval(500 * time.Millisecond),
//	)
//
// Default symbols tracked: BTC, ETH, BNB, SOL, ADA, XRP (all vs USDT)
//
// # Data Throttling
//
// Binance sends market updates at very high frequency (multiple times per second).
// The Ingestor implements throttling to:
//   - Batch multiple symbol updates together
//   - Control broadcast rate (default: 500ms intervals)
//   - Prevent React/frontend from excessive re-renders
//   - Reduce network bandwidth usage
//
// # Thread Safety
//
// All components are designed for concurrent use:
//   - Hub uses sync.RWMutex for client map protection
//   - Channels are used for goroutine communication
//   - Client write operations are serialized via WritePump
//   - Context-based cancellation for graceful shutdown
//
// # Performance
//
// The implementation is optimized for high throughput:
//   - Buffered channels prevent blocking on slow clients
//   - Non-blocking sends with default cases
//   - Efficient batching of multi-symbol updates
//   - Auto-reconnection on WebSocket disconnection
//
// # Production Features
//
// Production-ready features:
//   - Real-time data from Binance via adshao/go-binance SDK
//   - Automatic reconnection with backoff strategy
//   - Graceful shutdown with context cancellation
//   - Configurable throttling to control data flow
//   - Multi-symbol support via combined streams
//
// # Future Enhancements
//
// Potential improvements:
//   - Add authentication middleware for WebSocket connections
//   - Implement per-client subscription filtering
//   - Use Redis Pub/Sub for horizontal scaling
//   - Add connection rate limiting and compression
//   - Support more exchanges beyond Binance
package ws
