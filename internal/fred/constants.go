package fred

// Ticker represents a FRED data series identifier.
type Ticker string

const (
	// WALCL - Federal Reserve Total Assets (Money Printing/QT)
	TickerWALCL Ticker = "WALCL"

	// TGA - Treasury General Account (US Treasury Cash Balance)
	TickerTGA Ticker = "WTREGEN"

	// RRPONTSYD - Overnight Reverse Repurchase Agreements
	// (Money withdrawn from the system)
	TickerRRPONTSYD Ticker = "RRPONTSYD"

	// FEDFUNDS - Federal Funds Effective Rate
	TickerFEDFUNDS Ticker = "FEDFUNDS"

	// CPIAUCSL - Consumer Price Index for All Urban Consumers (Inflation)
	TickerCPIAUCSL Ticker = "CPIAUCSL"

	// DTWEXBGS - Trade Weighted U.S. Dollar Index: Broad, Goods and Services
	TickerDTWEXBGS Ticker = "DTWEXBGS"
)

// AllTickers returns all supported macro tickers.
func AllTickers() []Ticker {
	return []Ticker{
		TickerWALCL,
		TickerTGA,
		TickerRRPONTSYD,
		TickerFEDFUNDS,
		TickerCPIAUCSL,
		TickerDTWEXBGS,
	}
}

// String returns the string representation of a Ticker.
func (t Ticker) String() string {
	return string(t)
}

// Description returns a human-readable description of the ticker.
func (t Ticker) Description() string {
	descriptions := map[Ticker]string{
		TickerWALCL:     "Federal Reserve Total Assets",
		TickerTGA:       "Treasury General Account",
		TickerRRPONTSYD: "Overnight Reverse Repo",
		TickerFEDFUNDS:  "Federal Funds Rate",
		TickerCPIAUCSL:  "Consumer Price Index (CPI)",
		TickerDTWEXBGS:  "US Dollar Index",
	}
	return descriptions[t]
}
