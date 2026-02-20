package fred

import "testing"

// TestTickerString verifies Ticker string conversion.
func TestTickerString(t *testing.T) {
	tests := []struct {
		ticker   Ticker
		expected string
	}{
		{TickerWALCL, "WALCL"},
		{TickerTGA, "WTREGEN"},
		{TickerRRPONTSYD, "RRPONTSYD"},
		{TickerFEDFUNDS, "FEDFUNDS"},
		{TickerCPIAUCSL, "CPIAUCSL"},
		{TickerDTWEXBGS, "DTWEXBGS"},
	}

	for _, tt := range tests {
		result := tt.ticker.String()
		if result != tt.expected {
			t.Errorf("Ticker %v: expected %s, got %s", tt.ticker, tt.expected, result)
		}
	}
}

// TestTickerDescription verifies descriptions are not empty.
func TestTickerDescription(t *testing.T) {
	tests := []struct {
		ticker              Ticker
		expectedDescription string
	}{
		{TickerWALCL, "Federal Reserve Total Assets"},
		{TickerTGA, "Treasury General Account"},
		{TickerRRPONTSYD, "Overnight Reverse Repo"},
		{TickerFEDFUNDS, "Federal Funds Rate"},
		{TickerCPIAUCSL, "Consumer Price Index (CPI)"},
		{TickerDTWEXBGS, "US Dollar Index"},
	}

	for _, tt := range tests {
		result := tt.ticker.Description()
		if result != tt.expectedDescription {
			t.Errorf("Ticker %v: expected description '%s', got '%s'", 
				tt.ticker, tt.expectedDescription, result)
		}
	}
}

// TestAllTickers verifies all tickers are returned.
func TestAllTickers(t *testing.T) {
	tickers := AllTickers()

	if len(tickers) != 6 {
		t.Errorf("Expected 6 tickers, got %d", len(tickers))
	}

	expectedTickers := map[Ticker]bool{
		TickerWALCL:     false,
		TickerTGA:       false,
		TickerRRPONTSYD: false,
		TickerFEDFUNDS:  false,
		TickerCPIAUCSL:  false,
		TickerDTWEXBGS:  false,
	}

	for _, ticker := range tickers {
		if _, exists := expectedTickers[ticker]; !exists {
			t.Errorf("Unexpected ticker: %s", ticker)
		}
		expectedTickers[ticker] = true
	}

	for ticker, found := range expectedTickers {
		if !found {
			t.Errorf("Missing ticker: %s", ticker)
		}
	}
}

// TestTickerConstants verifies ticker constants are unique.
func TestTickerConstants(t *testing.T) {
	tickers := AllTickers()
	seen := make(map[string]bool)

	for _, ticker := range tickers {
		tickerStr := ticker.String()
		if seen[tickerStr] {
			t.Errorf("Duplicate ticker found: %s", tickerStr)
		}
		seen[tickerStr] = true
	}
}

// TestTickerDescriptionsNotEmpty verifies no empty descriptions.
func TestTickerDescriptionsNotEmpty(t *testing.T) {
	tickers := AllTickers()

	for _, ticker := range tickers {
		desc := ticker.Description()
		if desc == "" {
			t.Errorf("Ticker %s has empty description", ticker)
		}
	}
}
