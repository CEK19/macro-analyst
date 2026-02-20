package fred

import "time"

// Observation represents a single data point from FRED.
type Observation struct {
	Date  string `json:"date"`
	Value string `json:"value"`
}

// SeriesData represents the complete response for a series query.
type SeriesData struct {
	Ticker       Ticker        `json:"ticker"`
	Description  string        `json:"description"`
	Title        string        `json:"title"`
	Observations []Observation `json:"observations"`
	Units        string        `json:"units"`
	UnitsShort   string        `json:"units_short"`
	Frequency    string        `json:"frequency"`
	Notes        string        `json:"notes,omitempty"`
	LastUpdated  time.Time     `json:"last_updated"`
}

// FREDAPIResponse represents the raw response from FRED API observations endpoint.
type FREDAPIResponse struct {
	RealtimeStart  string        `json:"realtime_start"`
	RealtimeEnd    string        `json:"realtime_end"`
	Observations   []Observation `json:"observations"`
	Count          int           `json:"count"`
	Offset         int           `json:"offset"`
	Limit          int           `json:"limit"`
	Units          string        `json:"units,omitempty"`
	UnitsShort     string        `json:"units_short,omitempty"`
	Frequency      string        `json:"frequency,omitempty"`
	FrequencyShort string        `json:"frequency_short,omitempty"`
	OrderBy        string        `json:"order_by,omitempty"`
	SortOrder      string        `json:"sort_order,omitempty"`
}

// FREDSeriesResponse represents the response from FRED API series endpoint.
type FREDSeriesResponse struct {
	Seriess []FREDSeriesInfo `json:"seriess"`
}

// FREDSeriesInfo represents metadata about a FRED series.
type FREDSeriesInfo struct {
	ID                      string `json:"id"`
	Title                   string `json:"title"`
	ObservationStart        string `json:"observation_start"`
	ObservationEnd          string `json:"observation_end"`
	Frequency               string `json:"frequency"`
	FrequencyShort          string `json:"frequency_short"`
	Units                   string `json:"units"`
	UnitsShort              string `json:"units_short"`
	SeasonalAdjustment      string `json:"seasonal_adjustment"`
	SeasonalAdjustmentShort string `json:"seasonal_adjustment_short"`
	LastUpdated             string `json:"last_updated"`
	Popularity              int    `json:"popularity"`
	Notes                   string `json:"notes"`
}

// LatestValue represents the most recent data point for a ticker.
type LatestValue struct {
	Ticker      Ticker    `json:"ticker"`
	Description string    `json:"description"`
	Value       string    `json:"value"`
	Date        string    `json:"date"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// MultiTickerResponse represents data for multiple tickers.
type MultiTickerResponse struct {
	Data      []LatestValue `json:"data"`
	Timestamp time.Time     `json:"timestamp"`
}
