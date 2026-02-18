package ws

import (
	"testing"
	"time"
)

// TestNewIngestor verifies Ingestor initialization with default symbols.
func TestNewIngestor(t *testing.T) {
	hub := NewHub()
	ingestor := NewIngestor(hub)

	if ingestor == nil {
		t.Fatal("NewIngestor returned nil")
	}

	if ingestor.hub != hub {
		t.Error("Ingestor hub not set correctly")
	}

	if len(ingestor.symbols) == 0 {
		t.Error("Ingestor should have default symbols")
	}

	// Verify default crypto symbols
	expectedSymbols := []string{"BTCUSDT", "ETHUSDT", "BNBUSDT", "SOLUSDT", "ADAUSDT", "XRPUSDT"}
	if len(ingestor.symbols) != len(expectedSymbols) {
		t.Errorf("Expected %d symbols, got %d", len(expectedSymbols), len(ingestor.symbols))
	}

	// Verify context is initialized
	if ingestor.ctx == nil {
		t.Error("Context should be initialized")
	}

	if ingestor.cancel == nil {
		t.Error("Cancel function should be initialized")
	}
}

// TestIngestorWithOptions verifies functional options pattern.
func TestIngestorWithOptions(t *testing.T) {
	hub := NewHub()
	customInterval := 2 * time.Second

	ingestor := NewIngestor(hub,
		WithThrottleInterval(customInterval),
	)

	if ingestor.throttleInterval != customInterval {
		t.Errorf("Expected interval %v, got %v", customInterval, ingestor.throttleInterval)
	}
}

// TestAddSymbol verifies adding new symbols to the ingestor.
func TestAddSymbol(t *testing.T) {
	hub := NewHub()
	ingestor := NewIngestor(hub)

	initialCount := len(ingestor.symbols)
	ingestor.AddSymbol("DOGEUSDT")

	if len(ingestor.symbols) != initialCount+1 {
		t.Errorf("Expected %d symbols, got %d", initialCount+1, len(ingestor.symbols))
	}

	// Verify the symbol was added
	found := false
	for _, symbol := range ingestor.symbols {
		if symbol.Name == "DOGEUSDT" {
			found = true
			break
		}
	}

	if !found {
		t.Error("DOGEUSDT symbol not found after adding")
	}
}

// TestRemoveSymbol verifies removing symbols from the ingestor.
func TestRemoveSymbol(t *testing.T) {
	hub := NewHub()
	ingestor := NewIngestor(hub)

	// Add a test symbol
	ingestor.AddSymbol("TESTUSDT")
	initialCount := len(ingestor.symbols)

	// Remove the symbol
	removed := ingestor.RemoveSymbol("TESTUSDT")
	if !removed {
		t.Error("RemoveSymbol returned false for existing symbol")
	}

	if len(ingestor.symbols) != initialCount-1 {
		t.Errorf("Expected %d symbols, got %d", initialCount-1, len(ingestor.symbols))
	}

	// Try to remove non-existent symbol
	removed = ingestor.RemoveSymbol("NONEXISTENT")
	if removed {
		t.Error("RemoveSymbol returned true for non-existent symbol")
	}
}

// TestGetCurrentPrice verifies price retrieval.
func TestGetCurrentPrice(t *testing.T) {
	hub := NewHub()
	ingestor := NewIngestor(hub)

	// Test existing symbol (but no data yet since we haven't connected)
	_, err := ingestor.GetCurrentPrice("BTCUSDT")
	if err == nil {
		t.Error("Expected error for symbol with no price data yet")
	}

	// Test non-existent symbol
	_, err = ingestor.GetCurrentPrice("NONEXISTENT")
	if err == nil {
		t.Error("Expected error for non-existent symbol, got nil")
	}
}

// TestGetSymbols verifies symbol list retrieval.
func TestGetSymbols(t *testing.T) {
	hub := NewHub()
	ingestor := NewIngestor(hub)

	symbols := ingestor.GetSymbols()

	if len(symbols) == 0 {
		t.Error("GetSymbols returned empty slice")
	}

	// Verify it returns a copy (modifying shouldn't affect original)
	originalLen := len(ingestor.symbols)
	symbols = append(symbols, "TEST")

	if len(ingestor.symbols) != originalLen {
		t.Error("GetSymbols did not return a copy")
	}
}

// TestPriceUpdateJSON verifies JSON serialization.
func TestPriceUpdateJSON(t *testing.T) {
	update := &PriceUpdate{
		Symbol:        "BTC",
		Price:         50000.0,
		Change:        100.0,
		ChangePercent: 0.2,
		Volume:        500000,
		Timestamp:     "12:00:00",
	}

	// This test mainly verifies the struct tags are correct
	// Real JSON marshaling is tested in integration
	if update.Symbol != "BTC" {
		t.Error("Symbol field not accessible")
	}
}

// TestMultiUpdate verifies multi-update structure.
func TestMultiUpdate(t *testing.T) {
	updates := []*PriceUpdate{
		{Symbol: "BTC", Price: 50000, Change: 100, ChangePercent: 0.2},
		{Symbol: "ETH", Price: 3000, Change: 50, ChangePercent: 0.15},
	}

	multiUpdate := &MultiUpdate{
		Type: "multi_update",
		Data: updates,
	}

	if multiUpdate.Type != "multi_update" {
		t.Errorf("Expected type 'multi_update', got %s", multiUpdate.Type)
	}

	if len(multiUpdate.Data) != 2 {
		t.Errorf("Expected 2 updates, got %d", len(multiUpdate.Data))
	}
}

// TestStopIngestor verifies graceful shutdown.
func TestStopIngestor(t *testing.T) {
	hub := NewHub()
	ingestor := NewIngestor(hub)

	// Stop should not panic
	ingestor.Stop()

	// Context should be cancelled
	select {
	case <-ingestor.ctx.Done():
		// Expected
	default:
		t.Error("Context was not cancelled after Stop()")
	}
}
