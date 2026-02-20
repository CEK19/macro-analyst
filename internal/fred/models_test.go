package fred

import (
	"encoding/json"
	"testing"
	"time"
)

// TestObservationJSON verifies Observation JSON serialization.
func TestObservationJSON(t *testing.T) {
	obs := Observation{
		Date:  "2024-01-15",
		Value: "100.5",
	}

	data, err := json.Marshal(obs)
	if err != nil {
		t.Fatalf("Failed to marshal Observation: %v", err)
	}

	var decoded Observation
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal Observation: %v", err)
	}

	if decoded.Date != obs.Date {
		t.Errorf("Expected date %s, got %s", obs.Date, decoded.Date)
	}

	if decoded.Value != obs.Value {
		t.Errorf("Expected value %s, got %s", obs.Value, decoded.Value)
	}
}

// TestSeriesDataJSON verifies SeriesData JSON serialization.
func TestSeriesDataJSON(t *testing.T) {
	now := time.Now()
	seriesData := SeriesData{
		Ticker:      TickerWALCL,
		Description: "Federal Reserve Total Assets",
		Title:       "Assets: Total Assets: Total Assets: Wednesday Level",
		Observations: []Observation{
			{Date: "2024-01-15", Value: "50000.5"},
		},
		Units:       "Millions of Dollars",
		UnitsShort:  "Mil. of $",
		Frequency:   "Weekly, As of Wednesday",
		Notes:       "Assets: Total Assets data from Federal Reserve",
		LastUpdated: now,
	}

	data, err := json.Marshal(seriesData)
	if err != nil {
		t.Fatalf("Failed to marshal SeriesData: %v", err)
	}

	var decoded SeriesData
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal SeriesData: %v", err)
	}

	if decoded.Ticker != seriesData.Ticker {
		t.Errorf("Expected ticker %s, got %s", seriesData.Ticker, decoded.Ticker)
	}

	if len(decoded.Observations) != 1 {
		t.Errorf("Expected 1 observation, got %d", len(decoded.Observations))
	}

	if decoded.Title != seriesData.Title {
		t.Errorf("Expected title %s, got %s", seriesData.Title, decoded.Title)
	}

	if decoded.Units != seriesData.Units {
		t.Errorf("Expected units %s, got %s", seriesData.Units, decoded.Units)
	}

	if decoded.UnitsShort != seriesData.UnitsShort {
		t.Errorf("Expected units_short %s, got %s", seriesData.UnitsShort, decoded.UnitsShort)
	}

	if decoded.Frequency != seriesData.Frequency {
		t.Errorf("Expected frequency %s, got %s", seriesData.Frequency, decoded.Frequency)
	}

	if decoded.Notes != seriesData.Notes {
		t.Errorf("Expected notes %s, got %s", seriesData.Notes, decoded.Notes)
	}
}

// TestLatestValueJSON verifies LatestValue JSON serialization.
func TestLatestValueJSON(t *testing.T) {
	now := time.Now()
	latest := LatestValue{
		Ticker:      TickerCPIAUCSL,
		Description: "Consumer Price Index",
		Value:       "310.5",
		Date:        "2024-01-15",
		UpdatedAt:   now,
	}

	data, err := json.Marshal(latest)
	if err != nil {
		t.Fatalf("Failed to marshal LatestValue: %v", err)
	}

	var decoded LatestValue
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal LatestValue: %v", err)
	}

	if decoded.Ticker != latest.Ticker {
		t.Errorf("Expected ticker %s, got %s", latest.Ticker, decoded.Ticker)
	}

	if decoded.Value != latest.Value {
		t.Errorf("Expected value %s, got %s", latest.Value, decoded.Value)
	}
}

// TestMultiTickerResponseJSON verifies MultiTickerResponse JSON serialization.
func TestMultiTickerResponseJSON(t *testing.T) {
	now := time.Now()
	multiResp := MultiTickerResponse{
		Data: []LatestValue{
			{
				Ticker:      TickerWALCL,
				Description: "Fed Assets",
				Value:       "50000",
				Date:        "2024-01-15",
				UpdatedAt:   now,
			},
			{
				Ticker:      TickerFEDFUNDS,
				Description: "Fed Funds Rate",
				Value:       "5.25",
				Date:        "2024-01-15",
				UpdatedAt:   now,
			},
		},
		Timestamp: now,
	}

	data, err := json.Marshal(multiResp)
	if err != nil {
		t.Fatalf("Failed to marshal MultiTickerResponse: %v", err)
	}

	var decoded MultiTickerResponse
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal MultiTickerResponse: %v", err)
	}

	if len(decoded.Data) != 2 {
		t.Errorf("Expected 2 data items, got %d", len(decoded.Data))
	}
}

// TestFREDAPIResponseJSON verifies FREDAPIResponse JSON serialization.
func TestFREDAPIResponseJSON(t *testing.T) {
	fredResp := FREDAPIResponse{
		RealtimeStart:  "2024-01-01",
		RealtimeEnd:    "2024-01-31",
		Observations: []Observation{
			{Date: "2024-01-15", Value: "100.5"},
		},
		Count:          1,
		Offset:         0,
		Limit:          100,
		Units:          "Millions of Dollars",
		UnitsShort:     "Mil. of $",
		Frequency:      "Weekly",
		FrequencyShort: "W",
		OrderBy:        "observation_date",
		SortOrder:      "desc",
	}

	data, err := json.Marshal(fredResp)
	if err != nil {
		t.Fatalf("Failed to marshal FREDAPIResponse: %v", err)
	}

	var decoded FREDAPIResponse
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal FREDAPIResponse: %v", err)
	}

	if decoded.Count != fredResp.Count {
		t.Errorf("Expected count %d, got %d", fredResp.Count, decoded.Count)
	}

	if decoded.RealtimeStart != fredResp.RealtimeStart {
		t.Errorf("Expected realtime_start %s, got %s", fredResp.RealtimeStart, decoded.RealtimeStart)
	}

	if decoded.Units != fredResp.Units {
		t.Errorf("Expected units %s, got %s", fredResp.Units, decoded.Units)
	}

	if decoded.UnitsShort != fredResp.UnitsShort {
		t.Errorf("Expected units_short %s, got %s", fredResp.UnitsShort, decoded.UnitsShort)
	}
}

// TestFREDSeriesInfoJSON verifies FREDSeriesInfo JSON serialization.
func TestFREDSeriesInfoJSON(t *testing.T) {
	seriesInfo := FREDSeriesInfo{
		ID:                      "WALCL",
		Title:                   "Federal Reserve Total Assets",
		ObservationStart:        "2002-12-18",
		ObservationEnd:          "2024-02-14",
		Frequency:               "Weekly, As of Wednesday",
		FrequencyShort:          "W",
		Units:                   "Millions of Dollars",
		UnitsShort:              "Mil. of $",
		SeasonalAdjustment:      "Not Seasonally Adjusted",
		SeasonalAdjustmentShort: "NSA",
		LastUpdated:             "2024-02-15 16:17:03-06",
		Popularity:              100,
		Notes:                   "Assets: Total Assets: Total Assets: Wednesday Level",
	}

	data, err := json.Marshal(seriesInfo)
	if err != nil {
		t.Fatalf("Failed to marshal FREDSeriesInfo: %v", err)
	}

	var decoded FREDSeriesInfo
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal FREDSeriesInfo: %v", err)
	}

	if decoded.ID != seriesInfo.ID {
		t.Errorf("Expected ID %s, got %s", seriesInfo.ID, decoded.ID)
	}

	if decoded.Title != seriesInfo.Title {
		t.Errorf("Expected title %s, got %s", seriesInfo.Title, decoded.Title)
	}

	if decoded.Units != seriesInfo.Units {
		t.Errorf("Expected units %s, got %s", seriesInfo.Units, decoded.Units)
	}

	if decoded.UnitsShort != seriesInfo.UnitsShort {
		t.Errorf("Expected units_short %s, got %s", seriesInfo.UnitsShort, decoded.UnitsShort)
	}

	if decoded.Frequency != seriesInfo.Frequency {
		t.Errorf("Expected frequency %s, got %s", seriesInfo.Frequency, decoded.Frequency)
	}

	if decoded.Notes != seriesInfo.Notes {
		t.Errorf("Expected notes %s, got %s", seriesInfo.Notes, decoded.Notes)
	}
}

// TestEmptyObservations verifies handling of empty observations.
func TestEmptyObservations(t *testing.T) {
	seriesData := SeriesData{
		Ticker:       TickerWALCL,
		Description:  "Test",
		Observations: []Observation{},
		LastUpdated:  time.Now(),
	}

	if len(seriesData.Observations) != 0 {
		t.Error("Expected empty observations")
	}
}

// TestMultipleObservations verifies multiple observations handling.
func TestMultipleObservations(t *testing.T) {
	observations := []Observation{
		{Date: "2024-01-15", Value: "100.5"},
		{Date: "2024-01-14", Value: "99.8"},
		{Date: "2024-01-13", Value: "98.2"},
	}

	seriesData := SeriesData{
		Ticker:       TickerWALCL,
		Description:  "Test",
		Observations: observations,
		LastUpdated:  time.Now(),
	}

	if len(seriesData.Observations) != 3 {
		t.Errorf("Expected 3 observations, got %d", len(seriesData.Observations))
	}

	for i, obs := range seriesData.Observations {
		if obs.Date != observations[i].Date {
			t.Errorf("Observation %d: expected date %s, got %s", 
				i, observations[i].Date, obs.Date)
		}
	}
}
