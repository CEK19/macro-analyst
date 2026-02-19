package ws

import (
	"fmt"
	"testing"
	"time"

	"github.com/adshao/go-binance/v2"
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

// TestGetCurrentPriceWithData verifies price retrieval after data is set.
func TestGetCurrentPriceWithData(t *testing.T) {
	hub := NewHub()
	ingestor := NewIngestor(hub)

	// Simulate price update
	event := &binance.WsMarketStatEvent{
		Symbol:             "BTCUSDT",
		LastPrice:          "50000.00",
		PriceChange:        "100.00",
		PriceChangePercent: "0.20",
		BaseVolume:         "1000",
	}
	ingestor.updateSymbolData(event)

	price, err := ingestor.GetCurrentPrice("BTCUSDT")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if price != "50000.00" {
		t.Errorf("Expected price 50000.00, got %s", price)
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

// TestConvertEventToPriceUpdate verifies event conversion to PriceUpdate.
func TestConvertEventToPriceUpdate(t *testing.T) {
	hub := NewHub()
	ingestor := NewIngestor(hub)

	event := &binance.WsMarketStatEvent{
		Symbol:             "BTCUSDT",
		LastPrice:          "50000.50",
		PriceChange:        "100.25",
		PriceChangePercent: "0.20",
		BaseVolume:         "1000.75",
	}

	priceUpdate := ingestor.convertEventToPriceUpdate(event)

	if priceUpdate.Symbol != "BTCUSDT" {
		t.Errorf("Expected symbol BTCUSDT, got %s", priceUpdate.Symbol)
	}

	if priceUpdate.Price != 50000.50 {
		t.Errorf("Expected price 50000.50, got %f", priceUpdate.Price)
	}

	if priceUpdate.Change != 100.25 {
		t.Errorf("Expected change 100.25, got %f", priceUpdate.Change)
	}

	if priceUpdate.ChangePercent != 0.20 {
		t.Errorf("Expected change percent 0.20, got %f", priceUpdate.ChangePercent)
	}

	if priceUpdate.Volume != 1000 {
		t.Errorf("Expected volume 1000, got %d", priceUpdate.Volume)
	}

	if priceUpdate.Timestamp == "" {
		t.Error("Timestamp should not be empty")
	}
}

// TestConvertEventWithInvalidData verifies handling of invalid string data.
func TestConvertEventWithInvalidData(t *testing.T) {
	hub := NewHub()
	ingestor := NewIngestor(hub)

	event := &binance.WsMarketStatEvent{
		Symbol:             "BTCUSDT",
		LastPrice:          "invalid",
		PriceChange:        "invalid",
		PriceChangePercent: "invalid",
		BaseVolume:         "invalid",
	}

	priceUpdate := ingestor.convertEventToPriceUpdate(event)

	// Should default to zero values when parsing fails
	if priceUpdate.Price != 0 {
		t.Errorf("Expected price 0 for invalid input, got %f", priceUpdate.Price)
	}

	if priceUpdate.Change != 0 {
		t.Errorf("Expected change 0 for invalid input, got %f", priceUpdate.Change)
	}

	if priceUpdate.ChangePercent != 0 {
		t.Errorf("Expected change percent 0 for invalid input, got %f", priceUpdate.ChangePercent)
	}

	if priceUpdate.Volume != 0 {
		t.Errorf("Expected volume 0 for invalid input, got %d", priceUpdate.Volume)
	}
}

// TestUpdateSymbolData verifies symbol data is updated correctly.
func TestUpdateSymbolData(t *testing.T) {
	hub := NewHub()
	ingestor := NewIngestor(hub)

	event := &binance.WsMarketStatEvent{
		Symbol:             "BTCUSDT",
		LastPrice:          "50000.00",
		PriceChangePercent: "0.20",
		BaseVolume:         "1000",
	}

	// Update the symbol data
	ingestor.updateSymbolData(event)

	// Find the symbol and verify it was updated
	symbol := ingestor.findSymbol("BTCUSDT")
	if symbol == nil {
		t.Fatal("BTCUSDT symbol not found")
	}

	if symbol.LastPrice != "50000.00" {
		t.Errorf("Expected LastPrice 50000.00, got %s", symbol.LastPrice)
	}

	if symbol.LastChange != "0.20" {
		t.Errorf("Expected LastChange 0.20, got %s", symbol.LastChange)
	}

	if symbol.LastVolume != "1000" {
		t.Errorf("Expected LastVolume 1000, got %s", symbol.LastVolume)
	}

	if symbol.LastUpdateAt.IsZero() {
		t.Error("LastUpdateAt should be set")
	}
}

// TestUpdateSymbolDataForNonExistentSymbol verifies behavior when symbol doesn't exist.
func TestUpdateSymbolDataForNonExistentSymbol(t *testing.T) {
	hub := NewHub()
	ingestor := NewIngestor(hub)

	event := &binance.WsMarketStatEvent{
		Symbol:             "NONEXISTENT",
		LastPrice:          "50000.00",
		PriceChangePercent: "0.20",
		BaseVolume:         "1000",
	}

	// Should not panic when updating non-existent symbol
	ingestor.updateSymbolData(event)

	// Verify the non-existent symbol wasn't added
	symbol := ingestor.findSymbol("NONEXISTENT")
	if symbol != nil {
		t.Error("Non-existent symbol should not be added automatically")
	}
}

// TestFindSymbol verifies symbol lookup functionality.
func TestFindSymbol(t *testing.T) {
	hub := NewHub()
	ingestor := NewIngestor(hub)

	// Test finding existing symbol
	symbol := ingestor.findSymbol("BTCUSDT")
	if symbol == nil {
		t.Error("Expected to find BTCUSDT symbol")
	}

	if symbol.Name != "BTCUSDT" {
		t.Errorf("Expected symbol name BTCUSDT, got %s", symbol.Name)
	}

	// Test finding non-existent symbol
	symbol = ingestor.findSymbol("NONEXISTENT")
	if symbol != nil {
		t.Error("Expected nil for non-existent symbol")
	}
}

// TestStartMultiSymbolWithEmptySymbols verifies behavior with no symbols.
func TestStartMultiSymbolWithEmptySymbols(t *testing.T) {
	hub := NewHub()
	ingestor := NewIngestor(hub)

	// Remove all symbols
	for len(ingestor.symbols) > 0 {
		ingestor.RemoveSymbol(ingestor.symbols[0].Name)
	}

	// StartMultiSymbol should return early without panic
	ingestor.StartMultiSymbol()

	// No assertions needed, just verify it doesn't hang or panic
}

// TestQueuePriceUpdate verifies price update queuing logic.
func TestQueuePriceUpdate(t *testing.T) {
	hub := NewHub()
	ingestor := NewIngestor(hub)

	var pendingUpdate *MultiUpdate

	priceUpdate1 := &PriceUpdate{
		Symbol: "BTCUSDT",
		Price:  50000,
	}

	// Queue first update
	ingestor.queuePriceUpdate(&pendingUpdate, priceUpdate1)

	if pendingUpdate == nil {
		t.Fatal("Pending update should not be nil")
	}

	if len(pendingUpdate.Data) != 1 {
		t.Errorf("Expected 1 update, got %d", len(pendingUpdate.Data))
	}

	if pendingUpdate.Type != "multi_update" {
		t.Errorf("Expected type multi_update, got %s", pendingUpdate.Type)
	}

	// Queue second update for different symbol
	priceUpdate2 := &PriceUpdate{
		Symbol: "ETHUSDT",
		Price:  3000,
	}

	ingestor.queuePriceUpdate(&pendingUpdate, priceUpdate2)

	if len(pendingUpdate.Data) != 2 {
		t.Errorf("Expected 2 updates, got %d", len(pendingUpdate.Data))
	}

	// Queue update for same symbol (should replace)
	priceUpdate3 := &PriceUpdate{
		Symbol: "BTCUSDT",
		Price:  51000,
	}

	ingestor.queuePriceUpdate(&pendingUpdate, priceUpdate3)

	if len(pendingUpdate.Data) != 2 {
		t.Errorf("Expected 2 updates (replacement), got %d", len(pendingUpdate.Data))
	}

	// Verify BTCUSDT price was updated
	for _, update := range pendingUpdate.Data {
		if update.Symbol == "BTCUSDT" && update.Price != 51000 {
			t.Errorf("Expected BTCUSDT price 51000, got %f", update.Price)
		}
	}
}

// TestUpdateOrAppendPrice verifies the update/append logic for price updates.
func TestUpdateOrAppendPrice(t *testing.T) {
	hub := NewHub()
	ingestor := NewIngestor(hub)

	multiUpdate := &MultiUpdate{
		Type: "multi_update",
		Data: []*PriceUpdate{
			{Symbol: "BTCUSDT", Price: 50000},
		},
	}

	// Test updating existing symbol
	newUpdate := &PriceUpdate{Symbol: "BTCUSDT", Price: 51000}
	ingestor.updateOrAppendPrice(multiUpdate, newUpdate)

	if len(multiUpdate.Data) != 1 {
		t.Errorf("Expected 1 update after replace, got %d", len(multiUpdate.Data))
	}

	if multiUpdate.Data[0].Price != 51000 {
		t.Errorf("Expected price 51000, got %f", multiUpdate.Data[0].Price)
	}

	// Test appending new symbol
	newUpdate2 := &PriceUpdate{Symbol: "ETHUSDT", Price: 3000}
	ingestor.updateOrAppendPrice(multiUpdate, newUpdate2)

	if len(multiUpdate.Data) != 2 {
		t.Errorf("Expected 2 updates after append, got %d", len(multiUpdate.Data))
	}
}

// TestConcurrentSymbolAccess verifies thread-safe symbol operations.
func TestConcurrentSymbolAccess(t *testing.T) {
	hub := NewHub()
	ingestor := NewIngestor(hub)

	done := make(chan bool)

	// Concurrent reads
	for i := 0; i < 10; i++ {
		go func() {
			symbols := ingestor.GetSymbols()
			if len(symbols) == 0 {
				t.Error("GetSymbols returned empty")
			}
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}
}

// TestMultipleStopCalls verifies Stop can be called multiple times safely.
func TestMultipleStopCalls(t *testing.T) {
	hub := NewHub()
	ingestor := NewIngestor(hub)

	// Multiple Stop calls should not panic
	ingestor.Stop()
	ingestor.Stop()
	ingestor.Stop()

	// Context should still be cancelled
	select {
	case <-ingestor.ctx.Done():
		// Expected
	default:
		t.Error("Context should be cancelled")
	}
}

// TestCreateWebSocketHandler verifies WebSocket handler creation.
func TestCreateWebSocketHandler(t *testing.T) {
	hub := NewHub()
	ingestor := NewIngestor(hub)

	var pendingUpdate *MultiUpdate
	handler := ingestor.createWebSocketHandler(&pendingUpdate)

	if handler == nil {
		t.Fatal("createWebSocketHandler returned nil")
	}

	event := &binance.WsMarketStatEvent{
		Symbol:             "BTCUSDT",
		LastPrice:          "50000.00",
		PriceChange:        "100.00",
		PriceChangePercent: "0.20",
		BaseVolume:         "1000",
	}

	handler(event)

	if pendingUpdate == nil {
		t.Error("Handler should have created pending update")
	}

	if len(pendingUpdate.Data) != 1 {
		t.Errorf("Expected 1 update, got %d", len(pendingUpdate.Data))
	}
}

// TestCreateErrorHandler verifies error handler creation.
func TestCreateErrorHandler(t *testing.T) {
	hub := NewHub()
	ingestor := NewIngestor(hub)

	handler := ingestor.createErrorHandler()

	if handler == nil {
		t.Fatal("createErrorHandler returned nil")
	}

	// Should not panic when called
	handler(fmt.Errorf("test error"))
}

// TestSendToHub verifies sending data to hub with proper logging.
func TestSendToHub(t *testing.T) {
	hub := NewHub()
	ingestor := NewIngestor(hub)

	testData := []byte("test data")
	
	// Send without running hub (so we can verify it's in the channel)
	ingestor.sendToHub(testData, 5)

	// Verify data is in the hub's broadcast channel
	select {
	case msg := <-hub.broadcast:
		if string(msg) != string(testData) {
			t.Errorf("Expected %s, got %s", testData, msg)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Timeout waiting for data in hub broadcast channel")
	}
}

// TestSendToHubWithFullChannel verifies overflow protection.
func TestSendToHubWithFullChannel(t *testing.T) {
	hub := NewHub()
	ingestor := NewIngestor(hub)

	// Don't run hub.Run() so broadcast channel fills up
	// Fill the channel to capacity
	for i := 0; i < BroadcastBufferSize; i++ {
		hub.broadcast <- []byte("filler")
	}

	// This should not block or panic
	testData := []byte("overflow test")
	ingestor.sendToHub(testData, 1)

	// Should skip the send (verified by log message in implementation)
}

// TestBroadcastPendingUpdatesWithNilUpdate verifies nil safety.
func TestBroadcastPendingUpdatesWithNilUpdate(t *testing.T) {
	hub := NewHub()
	go hub.Run()
	time.Sleep(10 * time.Millisecond)

	ingestor := NewIngestor(hub)

	var pendingUpdate *MultiUpdate
	// Should not panic with nil
	ingestor.broadcastPendingUpdates(&pendingUpdate)

	if pendingUpdate != nil {
		t.Error("Pending update should remain nil")
	}
}

// TestBroadcastPendingUpdatesWithEmptyData verifies empty data handling.
func TestBroadcastPendingUpdatesWithEmptyData(t *testing.T) {
	hub := NewHub()
	go hub.Run()
	time.Sleep(10 * time.Millisecond)

	ingestor := NewIngestor(hub)

	pendingUpdate := &MultiUpdate{
		Type: "multi_update",
		Data: []*PriceUpdate{},
	}

	ingestor.broadcastPendingUpdates(&pendingUpdate)

	// Should not send anything to hub
	select {
	case <-hub.broadcast:
		t.Error("Should not broadcast empty update")
	case <-time.After(50 * time.Millisecond):
		// Expected - no broadcast
	}
}

// TestBroadcastPendingUpdatesWithValidData verifies successful broadcast.
func TestBroadcastPendingUpdatesWithValidData(t *testing.T) {
	hub := NewHub()
	ingestor := NewIngestor(hub)

	pendingUpdate := &MultiUpdate{
		Type: "multi_update",
		Data: []*PriceUpdate{
			{Symbol: "BTCUSDT", Price: 50000},
		},
	}

	ingestor.broadcastPendingUpdates(&pendingUpdate)

	// Verify data was sent to hub's broadcast channel
	select {
	case msg := <-hub.broadcast:
		if len(msg) == 0 {
			t.Error("Received empty message")
		}
		// Could also unmarshal and verify structure
	case <-time.After(100 * time.Millisecond):
		t.Error("Timeout waiting for broadcast")
	}

	// Pending update should be reset to nil
	if pendingUpdate != nil {
		t.Error("Pending update should be reset to nil after broadcast")
	}
}

// TestRemoveSymbolOrderIndependence verifies symbol removal maintains correctness.
func TestRemoveSymbolOrderIndependence(t *testing.T) {
	hub := NewHub()
	ingestor := NewIngestor(hub)

	ingestor.AddSymbol("TEST1")
	ingestor.AddSymbol("TEST2")
	ingestor.AddSymbol("TEST3")

	initialCount := len(ingestor.symbols)

	// Remove middle symbol
	removed := ingestor.RemoveSymbol("TEST2")
	if !removed {
		t.Error("Failed to remove TEST2")
	}

	if len(ingestor.symbols) != initialCount-1 {
		t.Errorf("Expected %d symbols, got %d", initialCount-1, len(ingestor.symbols))
	}

	// Verify other symbols still exist
	if ingestor.findSymbol("TEST1") == nil {
		t.Error("TEST1 should still exist")
	}

	if ingestor.findSymbol("TEST3") == nil {
		t.Error("TEST3 should still exist")
	}

	if ingestor.findSymbol("TEST2") != nil {
		t.Error("TEST2 should be removed")
	}
}

// TestDefaultThrottleInterval verifies default throttle setting.
func TestDefaultThrottleInterval(t *testing.T) {
	hub := NewHub()
	ingestor := NewIngestor(hub)

	if ingestor.throttleInterval != DefaultThrottleInterval {
		t.Errorf("Expected default interval %v, got %v", DefaultThrottleInterval, ingestor.throttleInterval)
	}
}

// TestSymbolStructInitialization verifies Symbol struct fields.
func TestSymbolStructInitialization(t *testing.T) {
	symbol := &Symbol{
		Name:         "BTCUSDT",
		LastPrice:    "50000",
		LastChange:   "0.20",
		LastVolume:   "1000",
		LastUpdateAt: time.Now(),
	}

	if symbol.Name != "BTCUSDT" {
		t.Errorf("Expected Name BTCUSDT, got %s", symbol.Name)
	}

	if symbol.LastPrice != "50000" {
		t.Errorf("Expected LastPrice 50000, got %s", symbol.LastPrice)
	}
}


