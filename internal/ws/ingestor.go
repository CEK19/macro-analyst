package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2"
)

const (
	// DefaultThrottleInterval is the default time between price broadcasts
	// to prevent overwhelming the clients
	DefaultThrottleInterval = 500 * time.Millisecond

	// MaxUpdatesPerSecond limits the number of updates sent to clients
	MaxUpdatesPerSecond = 10
)

// PriceUpdate represents a single price update for a financial instrument.
type PriceUpdate struct {
	Symbol        string  `json:"symbol"`        // Trading symbol (e.g., "BTCUSDT")
	Price         float64 `json:"price"`         // Current price
	Change        float64 `json:"change"`        // Absolute price change
	ChangePercent float64 `json:"changePercent"` // Percentage change
	Volume        int64   `json:"volume"`        // Trading volume
	Timestamp     string  `json:"timestamp"`     // Update timestamp
}

// MultiUpdate represents a batch of price updates for multiple symbols.
type MultiUpdate struct {
	Type string         `json:"type"` // Always "multi_update"
	Data []*PriceUpdate `json:"data"` // Array of price updates
}

// Symbol represents a trading symbol being tracked.
type Symbol struct {
	Name         string
	LastPrice    string
	LastChange   string
	LastVolume   string
	LastUpdateAt time.Time
}

// Ingestor connects to Binance WebSocket and streams real-time market data
// to connected clients via the Hub. It implements throttling to prevent
// overwhelming clients with too many updates.
type Ingestor struct {
	hub              *Hub
	symbols          []*Symbol
	throttleInterval time.Duration
	ctx              context.Context
	cancel           context.CancelFunc
	doneChannels     []chan struct{} // Track all WebSocket connections
}

// IngestorOption is a functional option for configuring the Ingestor.
type IngestorOption func(*Ingestor)

// WithThrottleInterval sets the minimum interval between broadcasts.
func WithThrottleInterval(interval time.Duration) IngestorOption {
	return func(i *Ingestor) {
		i.throttleInterval = interval
	}
}

// NewIngestor creates a new Ingestor with default crypto symbols.
func NewIngestor(hub *Hub, opts ...IngestorOption) *Ingestor {
	ctx, cancel := context.WithCancel(context.Background())

	// Initialize with popular crypto trading symbols
	symbols := []*Symbol{
		{Name: "BTCUSDT"},
		{Name: "ETHUSDT"},
		{Name: "BNBUSDT"},
		{Name: "SOLUSDT"},
		{Name: "ADAUSDT"},
		{Name: "XRPUSDT"},
	}

	ingestor := &Ingestor{
		hub:              hub,
		symbols:          symbols,
		throttleInterval: DefaultThrottleInterval,
		ctx:              ctx,
		cancel:           cancel,
		doneChannels:     make([]chan struct{}, 0),
	}

	// Apply options
	for _, opt := range opts {
		opt(ingestor)
	}

	return ingestor
}

// Start begins streaming real-time data from Binance WebSocket.
// It connects to Binance's Combined Ticker Stream for multiple symbols
// and broadcasts updates with throttling to prevent client overload.
func (i *Ingestor) Start() {
	log.Printf("Price Ingestor started - connecting to Binance WebSocket")
	log.Printf("Tracking symbols: %v", i.GetSymbols())

	// Start the multi-symbol stream
	i.StartMultiSymbol()
}

// StartMultiSymbol connects to Binance WebSocket for multiple symbols.
// It uses CombinedSymbolTickerServe to get all symbols in one connection.
func (i *Ingestor) StartMultiSymbol() {
	symbols := i.GetSymbols()
	if len(symbols) == 0 {
		log.Println("No symbols to track")
		return
	}

	log.Printf("Connecting to Binance for %d symbols...", len(symbols))

	// Create a throttle ticker to control broadcast rate
	throttleTicker := time.NewTicker(i.throttleInterval)
	defer throttleTicker.Stop()

	var pendingUpdate *MultiUpdate

	// Handler for WebSocket events
	wsHandler := func(event *binance.WsMarketStatEvent) {
		// Update the symbol's cached data
		i.updateSymbolData(event)

		// Prepare the update (but don't send yet - wait for throttle)
		priceUpdate := i.convertEventToPriceUpdate(event)

		// Queue this update
		if pendingUpdate == nil {
			pendingUpdate = &MultiUpdate{
				Type: "multi_update",
				Data: []*PriceUpdate{priceUpdate},
			}
		} else {
			// Check if this symbol already exists in pending updates
			found := false
			for idx, existing := range pendingUpdate.Data {
				if existing.Symbol == priceUpdate.Symbol {
					// Replace with newer data
					pendingUpdate.Data[idx] = priceUpdate
					found = true
					break
				}
			}
			if !found {
				pendingUpdate.Data = append(pendingUpdate.Data, priceUpdate)
			}
		}
	}

	errHandler := func(err error) {
		log.Printf("Binance WebSocket error: %v", err)
	}

	// Connect to Binance WebSocket (this is non-blocking)
	doneC, _, err := binance.WsCombinedMarketStatServe(symbols, wsHandler, errHandler)
	if err != nil {
		log.Printf("Failed to connect to Binance: %v", err)
		return
	}

	i.doneChannels = append(i.doneChannels, doneC)

	// Throttle loop - broadcasts updates at controlled rate
	go func() {
		for {
			select {
			case <-i.ctx.Done():
				log.Println("Ingestor stopped")
				return
			case <-throttleTicker.C:
				if pendingUpdate != nil && len(pendingUpdate.Data) > 0 {
					// Broadcast the batched updates
					jsonData, err := json.Marshal(pendingUpdate)
					if err != nil {
						log.Printf("Error marshaling update: %v", err)
						continue
					}

					select {
					case i.hub.broadcast <- jsonData:
						log.Printf("✓ Broadcasted %d symbol updates", len(pendingUpdate.Data))
					default:
						log.Println("⚠ Broadcast channel full, skipping update")
					}

					// Reset pending updates
					pendingUpdate = nil
				}
			}
		}
	}()

	// Wait for WebSocket to close or context cancellation
	select {
	case <-doneC:
		log.Println("Binance WebSocket connection closed")
	case <-i.ctx.Done():
		log.Println("Ingestor context cancelled")
	}
}

// Stop gracefully stops the ingestor and closes all WebSocket connections.
func (i *Ingestor) Stop() {
	log.Println("Stopping Price Ingestor...")
	i.cancel()

	// Close all WebSocket connections
	for _, doneC := range i.doneChannels {
		close(doneC)
	}
}

// updateSymbolData updates the cached symbol data from a Binance event.
func (i *Ingestor) updateSymbolData(event *binance.WsMarketStatEvent) {
	for _, symbol := range i.symbols {
		if symbol.Name == event.Symbol {
			symbol.LastPrice = event.LastPrice
			symbol.LastChange = event.PriceChangePercent
			symbol.LastVolume = event.BaseVolume
			symbol.LastUpdateAt = time.Now()
			break
		}
	}
}

// convertEventToPriceUpdate converts a Binance event to our PriceUpdate format.
func (i *Ingestor) convertEventToPriceUpdate(event *binance.WsMarketStatEvent) *PriceUpdate {
	price, _ := strconv.ParseFloat(event.LastPrice, 64)
	change, _ := strconv.ParseFloat(event.PriceChange, 64)
	changePercent, _ := strconv.ParseFloat(event.PriceChangePercent, 64)
	volume, _ := strconv.ParseFloat(event.BaseVolume, 64)

	return &PriceUpdate{
		Symbol:        event.Symbol,
		Price:         price,
		Change:        change,
		ChangePercent: changePercent,
		Volume:        int64(volume),
		Timestamp:     time.Now().Format("15:04:05.000"),
	}
}

// AddSymbol adds a new trading symbol to the ingestor's watchlist.
// Note: You'll need to restart the ingestor for this to take effect.
func (i *Ingestor) AddSymbol(name string) {
	symbol := &Symbol{
		Name: name,
	}
	i.symbols = append(i.symbols, symbol)
	log.Printf("Added symbol: %s (restart required)", name)
}

// RemoveSymbol removes a symbol from the ingestor's watchlist.
// Note: You'll need to restart the ingestor for this to take effect.
func (i *Ingestor) RemoveSymbol(name string) bool {
	for idx, symbol := range i.symbols {
		if symbol.Name == name {
			// Remove symbol by swapping with last element and truncating
			i.symbols[idx] = i.symbols[len(i.symbols)-1]
			i.symbols = i.symbols[:len(i.symbols)-1]
			log.Printf("Removed symbol: %s (restart required)", name)
			return true
		}
	}
	return false
}

// GetCurrentPrice returns the last known price of a symbol.
func (i *Ingestor) GetCurrentPrice(name string) (string, error) {
	for _, symbol := range i.symbols {
		if symbol.Name == name {
			if symbol.LastPrice == "" {
				return "", fmt.Errorf("no price data yet for: %s", name)
			}
			return symbol.LastPrice, nil
		}
	}
	return "", fmt.Errorf("symbol not found: %s", name)
}

// GetSymbols returns a copy of all tracked symbols.
func (i *Ingestor) GetSymbols() []string {
	symbols := make([]string, len(i.symbols))
	for i, symbol := range i.symbols {
		symbols[i] = symbol.Name
	}
	return symbols
}
